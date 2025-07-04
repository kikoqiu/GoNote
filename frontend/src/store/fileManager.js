import { authState } from '@/store/auth'
import defaultText from '@config/default'
import { MessageBox } from 'element-ui'
import i18n from '@/locales'

async function _fetchTreeData(apiClient, path) {
  // Now we can fetch the entire directory tree recursively with a single API call.
  const files = await apiClient.listDirectory(path, true)

  // The API client now returns a nested structure when recursive=true.
  // We need to transform it into the format our application expects.
  const transformNode = (node, currentPath = '') => {
    const newPath = currentPath ? `${currentPath}/${node.name}` : node.name
    const transformed = {
      name: node.name,
      path: newPath,
      type: node.is_dir ? 'directory' : 'file',
      size: node.size,
      mod_time: node.mod_time,
      is_dir: node.is_dir,
    }

    if (node.is_dir) {
      transformed.children = node.children
        ? node.children.map((child) => transformNode(child, newPath))
        : []
      // Sort children: directories first, then by name. Special case for 'Recycle'.
      transformed.children.sort((a, b) => {
        if (a.name === 'Recycle') return 1
        if (b.name === 'Recycle') return -1
        if (a.type === 'directory' && b.type !== 'directory') return -1
        if (a.type !== 'directory' && b.type === 'directory') return 1
        return a.name.localeCompare(b.name)
      })
    } else {
      transformed.isLeaf = true
    }

    return transformed
  }

  const fullTreeChildren = files.map((node) => transformNode(node, path))

  // Sort root level items
  fullTreeChildren.sort((a, b) => {
    if (a.name === 'Recycle') return 1
    if (b.name === 'Recycle') return -1
    if (a.type === 'directory' && b.type !== 'directory') return -1
    if (a.type !== 'directory' && b.type === 'directory') return 1
    return a.name.localeCompare(b.name)
  })

  return {
    fullTreeChildren,
  }
}
// Helper function to find a node by path in a tree
function findNodeByPath(nodes, path) {
  for (const node of nodes) {
    if (node.path === path) {
      return node
    }
    if (node.type === 'directory' && node.children) {
      const found = findNodeByPath(node.children, path)
      if (found) {
        return found
      }
    }
  }
  return null
}

// New helper function for handling Markdown file cache
async function handleMdFileCache(file) {
  const cacheKey = `filecache_${file.path}`;
  const cachedDataJSON = localStorage.getItem(cacheKey);

  if (cachedDataJSON) {
    const cachedData = JSON.parse(cachedDataJSON);
    const formattedDate = new Date(cachedData.timestamp).toLocaleString();

    // 使用一个无限循环，直到用户做出最终决定（恢复、确认放弃）或发生错误
    while (true) {
      try {
        // --- 步骤 1: 显示第一个选择框 (恢复 vs. 放弃) ---
        await MessageBox.confirm(
          i18n.t('fileManager.cacheFoundMessage', { formattedDate }),
          i18n.t('fileManager.restoreCacheTitle'),
          {
            confirmButtonText: i18n.t('fileManager.restoreButton'),
            cancelButtonText: i18n.t('fileManager.discardButton'),
            type: 'warning',
          }
        );

        // --- 用户选择“恢复” ---
        // 这是用户的最终决定，执行恢复操作并退出函数
        try {
          await authState.apiClient.writeFile(file.path, cachedData.content);
          localStorage.removeItem(cacheKey); // 成功后清除缓存
          console.log(`Restored and saved cached content for ${file.path}`);
          return; // 成功恢复，退出循环和函数
        } catch (writeError) {
          console.error(`Error saving restored cache for ${file.path}:`, writeError);
          //throw writeError; // 恢复失败，向上抛出错误
          continue; // 恢复失败，继续下一次循环
        }
      } catch (error) {
        // --- 用户在第一个弹窗选择“放弃”，进入二次确认流程 ---
        try {
          await MessageBox.confirm(
            i18n.t('fileManager.confirmDiscardMessage'),
            i18n.t('fileManager.discardCacheTitle'),
            {
              confirmButtonText: i18n.t('fileManager.confirmDiscardButton'),
              cancelButtonText: i18n.t('common.cancel'),
              type: 'warning',
            }
          );
          
          // --- 用户“确认放弃” ---
          // 这是用户的最终决定，清除缓存并退出函数
          localStorage.removeItem(cacheKey); // 清除缓存
          console.log(`Discarded cache for ${file.path}`);
          return; // 确认放弃，退出循环和函数

        } catch (discardError) {
          // --- 用户在第二个弹窗选择“取消” ---
          // 用户的意图是返回上一步，而不是做出最终决定
          console.log(`Cancelled discard for ${file.path}. Re-prompting.`);
          // 使用 'continue' 关键字跳过本次循环的剩余部分，
          // 直接开始下一次循环，从而重新显示第一个选择框。
          continue; 
        }
      }
    }
  }
}

const state = {
  isLoading: false,
  filesInSelectedFolder: [],
  directoryTree: [], // For el-tree, dirs only
  fullFileSystemTree: [], // Complete tree with files
}

const mutations = {
  SET_LOADING(state, status) {
    state.isLoading = status
  },
  SET_FILES_IN_SELECTED_FOLDER(state, files) {
    state.filesInSelectedFolder = files
  },
}

const actions = {
  // New action for loading the complete tree structure
  async loadCompleteDirectoryTree({ commit, state, dispatch }, path = '') {
    commit('SET_LOADING', true)

    // Helper function to build directory-only tree from the full tree.
    const _buildDirectoryTreeFromFull = (fullTreeNodes) => {
      const sorter = (a, b) => {
        if (a.name === 'Recycle') return 1
        if (b.name === 'Recycle') return -1
        return a.name.localeCompare(b.name)
      }

      const calculateTree = (nodes) => {
        let mdFileCount = 0
        const dirOnlyChildren = []

        if (nodes) {
          for (const node of nodes) {
            if (node.type === 'directory') {
              const subTreeResult = calculateTree(node.children || [])
              mdFileCount += subTreeResult.totalMdFileCount
              dirOnlyChildren.push({
                label: node.name,
                path: node.path,
                name: node.name,
                type: 'directory',
                children: subTreeResult.dirOnlyChildren,
                mdFileCount: subTreeResult.totalMdFileCount,
              })
            } else if (node.name.endsWith('.md')) {
              mdFileCount++
            }
          }
        }

        dirOnlyChildren.sort(sorter)
        return { dirOnlyChildren, totalMdFileCount: mdFileCount }
      }

      return calculateTree(fullTreeNodes).dirOnlyChildren
    }

    try {
      if (path === '') {
        await dispatch('ensureRecycleDirectory')
      }

      // Since _fetchTreeData now gets the whole tree, we always call it from the root.
      const { fullTreeChildren } = await _fetchTreeData(authState.apiClient, '')

      // 1. Replace the entire fullFileSystemTree with the new tree.
      state.fullFileSystemTree = fullTreeChildren

      // 2. Recalculate directoryTree from the updated fullFileSystemTree.
      const newDirectoryTree = _buildDirectoryTreeFromFull(state.fullFileSystemTree)

      // 3. Replace the entire directoryTree.
      state.directoryTree = newDirectoryTree
    } catch (error) {
      console.error(`Error loading directory tree for path "${path}":`, error)
    } finally {
      commit('SET_LOADING', false)
    }
  },

  // --- Existing actions below ---
  async ensureRecycleDirectory({ commit }) {
    try {
      const recyclePath = 'Recycle'
      const files = await authState.apiClient.listDirectory('')
      const recycleExists = files.some((file) => file.is_dir && file.name === recyclePath)

      if (!recycleExists) {
        console.log(`Creating Recycle directory: ${recyclePath}`)
        await authState.apiClient.createDirectory(recyclePath)
      }
    } catch (error) {
      console.error('Error ensuring Recycle directory exists:', error)
    }
  },

  async handleFolderClick({ commit, dispatch, state, rootState }, data) {
    if (data === null) {
      dispatch('ui/selectFolder', null, { root: true });
      dispatch('ui/selectFile', null, { root: true });
      return;
    }
    if (data.type === 'directory') {
      commit('SET_LOADING', true)
      try {
        const selectedFolderNode = findNodeByPath(state.fullFileSystemTree, data.path)

        let mdFiles = []
        if (selectedFolderNode && selectedFolderNode.children) {
          mdFiles = selectedFolderNode.children
            .filter((file) => file.type === 'file' && file.name.endsWith('.md'))
            .map((file) => ({
              name: file.name,
              path: file.path,
            }))
        }

        commit('SET_FILES_IN_SELECTED_FOLDER', mdFiles)
        dispatch('ui/selectFolder', data, { root: true })

        // If the newly selected folder is not the ui selectedFile's folder, clear the selectedFile
        if (rootState.ui.selectedFile) {
          const selectedFileFolderPath = rootState.ui.selectedFile.path.substring(0, rootState.ui.selectedFile.path.lastIndexOf('/'));
          if (selectedFileFolderPath !== data.path) {
            dispatch('ui/selectFile', null, { root: true });
          }
        }
      } catch (error) {
        console.error(`Error loading files for folder ${data.path}:`, error)
        commit('SET_FILES_IN_SELECTED_FOLDER', [])
        // You might want to dispatch a global error message here
      } finally {
        commit('SET_LOADING', false)
      }
    }
  },

  async handleFileSelect({ commit, dispatch, rootState }, file) {
    if (file === null) {
      dispatch('ui/selectFile', null, { root: true });
      return;
    }
    commit('SET_LOADING', true)
    try {
      const fileExtension = file.name.split('.').pop().toLowerCase()
      const imageExtensions = ['png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp', 'svg']

      let fileToSelect = null;

      if (imageExtensions.includes(fileExtension)) {
        fileToSelect = { ...file, imagePreviewUrl: authState.apiClient.getDownloadUrl(file.path) };
      } else if (fileExtension === 'md') {
        await handleMdFileCache(file);
        const fileContent = await authState.apiClient.readFile(file.path);
        fileToSelect = { ...file, content: fileContent.Content };
      } else {
        try {
          await MessageBox.confirm(
            i18n.t('fileManager.downloadConfirm', { name: file.name }),
            i18n.t('fileManager.downloadTitle'),
            {
              confirmButtonText: i18n.t('fileManager.downloadButton'),
              cancelButtonText: i18n.t('common.cancel'),
              type: 'info',
            }
          );
          const url = authState.apiClient.getDownloadUrl(file.path);
          const a = document.createElement('a');
          a.href = url;
          a.download = file.name;
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          commit('SET_LOADING', false);
          return file; // No content to return for download, just the file object
        } catch (error) {
          if (error !== 'cancel') {
            console.error(`Failed to download file: ${error.message}`);
            throw error;
          }
          commit('SET_LOADING', false);
          return file; // Return file even if download is cancelled
        }
      }

      // Check if the selected file's folder is different from the currently selected folder
      const fileFolderPath = file.path.substring(0, file.path.lastIndexOf('/'));
      if (!rootState.ui.selectedFolder || rootState.ui.selectedFolder.path !== fileFolderPath) {
        // If different, automatically change the selected folder to the file's folder
        // We need to find the folder object from the fullFileSystemTree
        const newSelectedFolder = findNodeByPath(rootState.fileManager.fullFileSystemTree, fileFolderPath);
        if (newSelectedFolder) {
          // Call handleFolderClick to update the folder and file list
          await dispatch('handleFolderClick', newSelectedFolder);
        }
      }

      // Set the selected file in the UI store
      dispatch('ui/selectFile', fileToSelect, { root: true });

      commit('SET_LOADING', false);
      return fileToSelect;
    } catch (error) {
      console.error(`Error handling file click for ${file.path}:`, error);
      commit('SET_LOADING', false);
      throw error; // Re-throw the error to be handled by the component
    }
  },

  async saveFileContent({ commit, rootGetters }, { path, content }) {
    commit('SET_LOADING', true)
    const cacheKey = `filecache_${path}`
    try {
      // Cache the content before saving
      const cacheData = {
        content,
        timestamp: new Date().toISOString(),
      }
      localStorage.setItem(cacheKey, JSON.stringify(cacheData))

      await authState.apiClient.writeFile(path, content)

      // Clear the cache on successful save
      localStorage.removeItem(cacheKey)
      console.log(`File saved and cache cleared: ${path}`)
    } catch (error) {
      console.error(`Error saving file ${path}:`, error)
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async handleDelete({ commit, dispatch, rootState }, { targetPath, isFolder, isFromRecycle }) {
    commit('SET_LOADING', true)
    try {
      if (isFromRecycle) {
        // Permanent delete
        if (isFolder) {
          await authState.apiClient.deleteDirectory(targetPath)
        } else {
          await authState.apiClient.deleteFile(targetPath)
        }
      } else {
        // Move to Recycle
        const targetName = targetPath.substring(targetPath.lastIndexOf('/') + 1)
        const recyclePath = 'Recycle'
        let newRecyclePath
        if (isFolder) {
          newRecyclePath = `${recyclePath}/${targetName}_${Date.now()}`
        } else {
          const dotIndex = targetName.lastIndexOf('.')
          if (dotIndex !== -1) {
            const name = targetName.substring(0, dotIndex)
            const ext = targetName.substring(dotIndex)
            newRecyclePath = `${recyclePath}/${name}_${Date.now()}${ext}`
          } else {
            newRecyclePath = `${recyclePath}/${targetName}_${Date.now()}`
          }
        }

        if (isFolder) {
          await authState.apiClient.renameDirectory(targetPath, newRecyclePath)
        } else {
          await authState.apiClient.renameFile(targetPath, newRecyclePath)
        }
      }
      if (isFolder) {
          dispatch('ui/selectFolder', null, { root: true });
          dispatch('ui/selectFile', null, { root: true });
      }else{
          dispatch('ui/selectFile', null, { root: true });
      }
      return { success: true, targetPath, isFolder, isFromRecycle }
    } catch (error) {
      console.error(`Error deleting ${isFolder ? 'folder' : 'file'} ${targetPath}:`, error)
      throw error // Re-throw to be caught by the component
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async performFileManagerAction(
    { commit, dispatch, rootState },
    { actionType, inputValue, selectedFolder, selectedFile, createInRoot }
  ) {
    commit('SET_LOADING', true)
    let newPath = ''
    try {
      const currentFolder = selectedFolder
      const currentFile = selectedFile
      const currentPath = currentFolder ? currentFolder.path : ''

      if (actionType === 'newFolder' && createInRoot) {
        newPath = inputValue
      } else {
        newPath = currentPath ? `${currentPath}/${inputValue}` : inputValue
      }

      switch (actionType) {
        case 'newFolder':
          await authState.apiClient.createDirectory(newPath)
          break
        case 'newFile':
          await authState.apiClient.writeFile(newPath, '')
          break
        case 'renameFolder':
          if (currentFolder) {
            const oldPath = currentFolder.path
            newPath = oldPath.substring(0, oldPath.lastIndexOf('/')) + '/' + inputValue
            await authState.apiClient.renameDirectory(oldPath, newPath)
          }
          break
        case 'renameFile':
          if (currentFile) {
            const oldPath = currentFile.path
            newPath = oldPath.substring(0, oldPath.lastIndexOf('/')) + '/' + inputValue
            await authState.apiClient.renameFile(oldPath, newPath)
          }
          break
      }
      return { success: true, actionType, newPath, selectedFolder, selectedFile }
    } catch (error) {
      console.error(`Error performing ${actionType}:`, error)
      throw error // Re-throw to be caught by the component
    } finally {
      commit('SET_LOADING', false)
    }
  },
}

const getters = {
  isLoading: (state) => state.isLoading,
  filesInSelectedFolder: (state) => state.filesInSelectedFolder,
  directoryTree: (state) => state.directoryTree,
  fullFileSystemTree: (state) => state.fullFileSystemTree,
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
}
;('')
