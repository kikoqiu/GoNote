// frontend/src/store/fileHistoryManager.js
import { Message } from 'element-ui' // Assuming Element UI is used for messages
import { authState } from './auth'

const fileHistoryManager = {
  namespaced: true,
  state: {
    currentFilePath: null,
    history: [],
    loadingHistory: false,
    loadingVersion: false,
    currentVersionContent: null,
    currentVersionId: null,
    diffContent: null,
  },
  mutations: {
    SET_CURRENT_FILE_PATH(state, path) {
      state.currentFilePath = path
    },
    SET_HISTORY(state, history) {
      state.history = history
    },
    SET_LOADING_HISTORY(state, status) {
      state.loadingHistory = status
    },
    SET_LOADING_VERSION(state, status) {
      state.loadingVersion = status
    },
    SET_CURRENT_VERSION_CONTENT(state, content) {
      state.currentVersionContent = content
    },
    SET_CURRENT_VERSION_ID(state, id) {
      state.currentVersionId = id
    },
    SET_DIFF_CONTENT(state, diff) {
      state.diffContent = diff
    },
    CLEAR_HISTORY_STATE(state) {
      state.currentFilePath = null
      state.history = []
      state.loadingHistory = false
      state.loadingVersion = false
      state.currentVersionContent = null
      state.currentVersionId = null
      state.diffContent = null
    },
  },
  actions: {
    async fetchFileHistory({ commit, state }, filePath) {
      commit('SET_LOADING_HISTORY', true)
      commit('SET_CURRENT_FILE_PATH', filePath)
      try {
        const response = await authState.apiClient.getFileHistory(filePath)
        let historyData = Array.isArray(response) ? response : []
        // Sort history by timestamp descending (newest first)
        const sortedHistory = historyData.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
        commit('SET_HISTORY', sortedHistory)
      } catch (error) {
        Message.error(`Failed to fetch file history: ${error.message}`)
        console.error('Error fetching file history:', error)
      } finally {
        commit('SET_LOADING_HISTORY', false)
      }
    },

    async fetchFileVersion({ commit, state }, { filePath, versionId }) {
      commit('SET_LOADING_VERSION', true)
      commit('SET_CURRENT_VERSION_ID', versionId)
      try {
        const response = await authState.apiClient.getFileVersion(filePath, versionId)
        commit('SET_CURRENT_VERSION_CONTENT', response.content)
      } catch (error) {
        Message.error(`Failed to fetch file version: ${error.message}`)
        console.error('Error fetching file version:', error)
      } finally {
        commit('SET_LOADING_VERSION', false)
      }
    },

    async applyFileVersion({ commit, state, dispatch }, { filePath, versionContent, comment = 'Reverted to an older version' }) {
      try {
        await authState.apiClient.writeFile(filePath, versionContent, comment)
        Message.success('File successfully reverted to selected version.')
        // After reverting, update the UI's selected file to reflect the new content
        const updatedFile = { path: filePath, name: filePath.substring(filePath.lastIndexOf('/') + 1) };
        await dispatch('fileManager/handleFileSelect', updatedFile, { root: true });
      } catch (error) {
        Message.error(`Failed to apply file version: ${error.message}`)
        console.error('Error applying file version:', error)
      }
    },

    async fetchCurrentFileContent({ commit, state }, filePath) {
      try {
        const response = await authState.apiClient.readFile(filePath)
        return response.Content // Assuming 'Content' holds the markdown string
      } catch (error) {
        Message.error(`Failed to fetch current file content: ${error.message}`)
        console.error('Error fetching current file content:', error)
        return null
      }
    },

    async generateDiff({ commit, dispatch, state }, { filePath, versionId }) {
      commit('SET_LOADING_VERSION', true)
      try {
        const currentContent = await dispatch('fetchCurrentFileContent', filePath)
        const versionResponse = await authState.apiClient.getFileVersion(filePath, versionId)
        const historicalContent = versionResponse.content

        if (currentContent && historicalContent) {
          // This is a placeholder for actual diff generation.
          // In a real application, you would use a diff library here.
          // For now, we'll just store both contents for display.
          commit('SET_DIFF_CONTENT', {
            current: currentContent,
            historical: historicalContent,
            // You would typically generate a diff string here, e.g., using diff-match-patch
            // For example:
            // const dmp = new diff_match_patch();
            // const diff = dmp.diff_main(historicalContent, currentContent);
            // dmp.diff_cleanupSemantic(diff);
            // const diffHtml = dmp.diff_prettyHtml(diff);
            // commit('SET_DIFF_CONTENT', diffHtml);
          })
        } else {
          commit('SET_DIFF_CONTENT', null)
        }
      } catch (error) {
        Message.error(`Failed to generate diff: ${error.message}`)
        console.error('Error generating diff:', error)
        commit('SET_DIFF_CONTENT', null)
      } finally {
        commit('SET_LOADING_VERSION', false)
      }
    },
  },
}

export default fileHistoryManager
