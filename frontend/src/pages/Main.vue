<template>
  <div class="index-page" v-loading="isLoading">
    <LoginDialog />
    <HeaderNav :displayTitle="currentDisplayTitle" @logo-click="onLogoClick" />

    <!-- ======================================================================= -->
    <!-- ==================== MOBILE LAYOUT (v-if="isMobile") ==================== -->
    <!-- ======================================================================= -->
    <div v-if="isMobile" class="main-container mobile-layout">
      <!-- 1. Viewport: 宽度为 100vw，裁剪内部的宽容器 -->
      <div class="mobile-slider-viewport">
        <!-- 2. Slider Container: 宽度为 250vw，通过 transform 来移动 -->
        <div class="mobile-slider-container" :style="mobileSliderStyle">
          <!-- Column 1: Folder Tree (50vw) -->
          <div class="mobile-column folder-column-mobile">
            <FolderTree
              ref="folderTreeComponent"
              :is-mobile="isMobile"
              :display-state="displayState"
              :tree-key="treeKey"
              @open-folder-dialog="handleOpenFolderDialog"
              @open-search-dialog="openFileSearchDialog"
            />
          </div>

          <!-- Column 2: File List (100vw) -->
          <div class="mobile-column file-column-mobile">
            <div v-if="selectedFolder" class="column-header file-list-header">
              <div class="file-list-title">
                <h3>
                  {{ $t('fileManager.files') }} ({{
                    filesInSelectedFolder.length ? filesInSelectedFolder.length : '0'
                  }})
                </h3>
              </div>
              <div class="header-actions">
                <el-tooltip :content="$t('fileManager.newFileTooltip')" placement="top">
                  <span class="action-icon" @click="openNewFileDialog"
                    ><icon :name="'file'"
                  /></span>
                </el-tooltip>
              </div>
            </div>
            <ul class="file-list">
              <li
                v-for="file in filesInSelectedFolder"
                :key="file.path"
                @click="$store.dispatch('fileManager/handleFileSelect', file)"
                :class="{ 'is-selected': selectedFile && selectedFile.path === file.path }"
              >
                <el-tooltip :content="file.name" placement="top" :disabled="file.name.length < 20">
                  <span class="file-name-text">{{ file.name }}</span>
                </el-tooltip>
                <div class="file-actions">
                  <el-tooltip :content="$t('fileManager.renameTooltip')" placement="top">
                    <el-button
                      size="mini"
                      icon="el-icon-edit"
                      type="text"
                      @click.stop="openRenameDialog(file)"
                    ></el-button>
                  </el-tooltip>
                  <el-tooltip :content="$t('fileManager.deleteTooltip')" placement="top">
                    <el-button
                      size="mini"
                      icon="el-icon-delete"
                      type="text"
                      @click.stop="confirmDelete(file)"
                    ></el-button>
                  </el-tooltip>
                  <el-tooltip :content="$t('fileHistory.fileHistory')" placement="top">
                    <el-button
                      v-if="file.name.endsWith('.md')"
                      size="mini"
                      icon="el-icon-time"
                      type="text"
                      @click.stop="openFileHistory(file)"
                    ></el-button>
                  </el-tooltip>
                </div>
              </li>
            </ul>
          </div>

          <!-- Column 3: Vditor (100vw) -->
          <div class="mobile-column vditor-main-mobile">
            <div id="vditor" class="vditor" />
          </div>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- ================== DESKTOP LAYOUT (v-else) - UNCHANGED ================== -->
    <!-- ======================================================================= -->
    <el-container v-else class="main-container">
      <el-aside width="550px" class="file-browser-aside modern-sidebar">
        <el-row class="file-browser-row">
          <el-col :span="12" class="folder-column">
            <FolderTree
              ref="folderTreeComponent"
              :is-mobile="isMobile"
              :display-state="displayState"
              :tree-key="treeKey"
              @open-folder-dialog="handleOpenFolderDialog"
              @open-search-dialog="openFileSearchDialog"
            />
          </el-col>

          <el-col :span="12" class="file-column">
            <div v-if="selectedFolder" class="column-header file-list-header">
              <div class="file-list-title">
                <h3>
                  {{ $t('fileManager.files') }} ({{
                    filesInSelectedFolder.length ? filesInSelectedFolder.length : '0'
                  }})
                </h3>
              </div>
              <div class="header-actions">
                <el-tooltip :content="$t('fileManager.newFileTooltip')" placement="top">
                  <span class="action-icon" @click="openNewFileDialog"
                    ><icon :name="'file'"
                  /></span>
                </el-tooltip>
              </div>
            </div>
            <ul class="file-list">
              <li
                v-for="file in filesInSelectedFolder"
                :key="file.path"
                @click="$store.dispatch('fileManager/handleFileSelect', file)"
                :class="{ 'is-selected': selectedFile && selectedFile.path === file.path }"
              >
                <el-tooltip :content="file.name" placement="top" :disabled="file.name.length < 20">
                  <span class="file-name-text">{{ file.name }}</span>
                </el-tooltip>
                <div class="file-actions">
                  <el-tooltip :content="$t('fileManager.renameTooltip')" placement="top">
                    <el-button
                      size="mini"
                      icon="el-icon-edit"
                      type="text"
                      @click.stop="openRenameDialog(file)"
                    ></el-button>
                  </el-tooltip>
                  <el-tooltip :content="$t('fileManager.deleteTooltip')" placement="top">
                    <el-button
                      size="mini"
                      icon="el-icon-delete"
                      type="text"
                      @click.stop="confirmDelete(file)"
                    ></el-button>
                  </el-tooltip>
                  <el-tooltip :content="$t('fileHistory.fileHistory')" placement="top">
                    <el-button
                      v-if="file.name.endsWith('.md')"
                      size="mini"
                      icon="el-icon-time"
                      type="text"
                      @click.stop="openFileHistory(file)"
                    ></el-button>
                  </el-tooltip>
                </div>
              </li>
            </ul>
          </el-col>
        </el-row>
      </el-aside>
      <el-main class="vditor-main">
        <div id="vditor" class="vditor" />
      </el-main>
    </el-container>

    <!-- Image Preview Lightbox -->
    <div v-if="imagePreviewUrl" class="image-lightbox" @click="imagePreviewUrl = null">
      <img :src="imagePreviewUrl" />
    </div>

    <!-- Dialog for File Actions -->
    <el-dialog
      :title="dialogTitle"
      :visible.sync="dialogVisible"
      :width="isMobile ? '90%' : '30%'"
      :fullscreen="isMobile"
      @close="inputValue = ''"
    >
      <el-input
        v-model="inputValue"
        :placeholder="inputPlaceholder"
        @keyup.enter.native="handleDialogConfirm"
      ></el-input>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleDialogConfirm">{{
          $t('common.confirm')
        }}</el-button>
      </span>
    </el-dialog>

    <!-- Dialog for Folder Actions -->
    <el-dialog
      :title="folderDialog.title"
      :visible.sync="folderDialog.visible"
      :width="isMobile ? '90%' : '30%'"
      :fullscreen="isMobile"
      @close="folderDialog.inputValue = ''"
    >
      <el-input
        v-model="folderDialog.inputValue"
        :placeholder="folderDialog.placeholder"
        @keyup.enter.native="handleFolderDialogConfirm"
      ></el-input>
      <el-checkbox
        v-if="folderDialog.actionType === 'newFolder'"
        v-model="folderDialog.createInRoot"
        style="margin-top: 10px;"
      >{{ $t('fileManager.createInRoot') }}</el-checkbox>
      <span slot="footer" class="dialog-footer">
        <el-button @click="folderDialog.visible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleFolderDialogConfirm">{{
          $t('common.confirm')
        }}</el-button>
      </span>
    </el-dialog>

    <!-- File History Dialog -->
    <FileHistoryDialog
      :visible.sync="showFileHistoryDialog"
      :file-path="selectedFileForHistory"
    />

    <!-- File Search Dialog -->
    <FileSearchDialog
      :visible.sync="showFileSearchDialog"
      @close="showFileSearchDialog = false"
      :is-mobile="isMobile"
    />
  </div>
</template>

<script>
import Vditor from 'vditor'
import HeaderNav from './partials/HeaderNav'
import defaultText from '@config/default'
import LoginDialog from '@/components/LoginDialog.vue'
import FolderTree from './FolderTree.vue'
import floppyDisk from '@/assets/icons/floppyDisk.js'
import {
  Container,
  Aside,
  Main,
  Row,
  Col,
  Tooltip,
  Button,
  Dialog,
  Input,
  MessageBox,
} from 'element-ui'
import Icon from '@/components/Icon'
import { mapState, mapGetters, mapActions } from 'vuex'
import { authState } from '@/store/auth'
import Icons from '@assets/icons'
import FileHistoryDialog from '@/components/FileHistoryDialog.vue' // Import the new component
import FileSearchDialog from '@/components/FileSearchDialog.vue'

export default {
  name: 'Main',
  components: {
    HeaderNav,
    LoginDialog,
    FolderTree,
    ElContainer: Container,
    ElAside: Aside,
    ElMain: Main,
    ElRow: Row,
    ElCol: Col,
    ElTooltip: Tooltip,
    ElButton: Button,
    ElDialog: Dialog,
    ElInput: Input,
    Icon,
    FileHistoryDialog, // Register the new component
    FileSearchDialog, // Register the new component
  },
  data() {
    return {
      isMobile: window.innerWidth <= 960,
      vditor: null,
      dialogVisible: false,
      dialogTitle: '',
      inputPlaceholder: '',
      inputValue: '',
      actionType: '',
      treeKey: 0,
      contentUpdated: false,
      imagePreviewUrl: null,
      folderDialog: {
        visible: false,
        title: '',
        placeholder: '',
        inputValue: '',
        actionType: '',
        createInRoot: true, // New property
      },
      showFileHistoryDialog: false, // New data property
      selectedFileForHistory: null, // New data property
      showFileSearchDialog: false, // New data property for FileSearchDialog
      dialogParams: { // Stores parameters for dialog actions
        actionType: '',
        inputValue: '',
        selectedFolder: null,
        selectedFile: null,
        folderPathToReload: '',
      },
      autoSaveTimer: null,
    }
  },
  computed: {
    ...mapState('ui', ['displayState', 'selectedFolder', 'selectedFile']),
    ...mapState('fileManager', ['isLoading']),
    ...mapGetters('ui', ['currentDisplayTitle']),
    ...mapGetters('fileManager', ['fullFileSystemTree']),
    filesInSelectedFolder() {
      if (!this.selectedFolder || !this.fullFileSystemTree) return []
      const findNode = (nodes, path) => {
        for (const node of nodes) {
          if (node.path === path) return node
          if (node.children) {
            const found = findNode(node.children, path)
            if (found) return found
          }
        }
        return null
      }
      const folderNode = findNode(this.fullFileSystemTree, this.selectedFolder.path)
      return folderNode && folderNode.children
        ? folderNode.children.filter((child) => child.type === 'file')
        : []
    },

    // 3. NEW COMPUTED PROPERTY: Calculates the transform style for the mobile slider
    mobileSliderStyle() {
      if (!this.isMobile) return {}
      let offset = 0
      switch (this.displayState) {
        case 1: // Show folder list (left-aligned)
          offset = 0
          break
        case 2: // Show file list (slide by folder width)
          offset = -50 // Folder column is 50vw
          break
        case 3: // Show vditor (slide by folder + file width)
          offset = -150 // Folder (50vw) + File (100vw)
          break
        default:
          offset = 0
      }
      return {
        transform: `translateX(${offset}vw)`,
      }
    },
  },
  created() {
    if (process.env.NODE_ENV === 'production') console.log = () => {}
    this.$root.$on('reload-content', this.reloadContent)
    window.addEventListener('resize', this.handleResize)
  },
  watch: {
    displayState(newState, oldState) {
      if (this.isMobile && newState > oldState) {
        history.pushState({ displayState: newState }, '');
      }
      if(this.isMobile && newState==2 && oldState==3){
        this.autoSaveTrigger();
      }
    },
    '$store.state.ui.selectedFile': {
      handler(newFile, oldFile) {
        if (newFile !== oldFile) {
          this.handleFileSelected(newFile);
        }
      },
      deep: true,
      immediate: false, // Do not run on initial load
    },
  },
  mounted() {
    this.initVditor();
    window.addEventListener('popstate', this.moveBack);
    setInterval(()=>{
      this.autoSaveTrigger();
    },1000*60*10);
  },
  beforeDestroy() {
    this.$root.$off('reload-content', this.reloadContent);
    window.removeEventListener('resize', this.handleResize);
    window.removeEventListener('popstate', this.moveBack);
  },
  methods: {
    ...mapActions('ui', ['resetDisplay', 'backToFileList', 'backToFolderList', 'moveBack']),
    ...mapActions('fileManager', [
      'loadCompleteDirectoryTree',
      'handleFileSelect',
      'handleDelete',
      'performFileManagerAction',
      'saveFileContent',
    ]),

    autoSaveTrigger(){
      if(this.vditor ){
        const value=this.vditor.getValue();
        if(value.trim()!==""){
          this.autoSaveCheck(value)
        }     
      }
    },
    handleResize() {
      this.isMobile = window.innerWidth <= 960
    },
    // ... (rest of the methods are unchanged)
    openNewFileDialog() {
      if (!this.selectedFolder) {
        this.$message.info(this.$t('fileManager.selectFolderMessage'))
        return
      }
      this.dialogTitle = this.$t('fileManager.newFileTitle')
      this.inputPlaceholder = this.$t('fileManager.newFileInputPlaceholder')
      this.inputValue = '';
      this.dialogParams.actionType = 'newFile';
      this.dialogParams.selectedFolder = this.selectedFolder; // Store the current selected folder
      this.dialogParams.selectedFile = null;
      this.dialogParams.folderPathToReload = this.selectedFolder.path;
      this.dialogVisible = true
    },
    openRenameDialog(file) {
      const targetFile = file // Store the file being renamed
      if (!targetFile) return
      this.dialogTitle = this.$t('fileManager.renameFileTitle')
      this.inputPlaceholder = this.$t('fileManager.renameFileInputPlaceholder')
      this.inputValue = targetFile.name;
      this.dialogParams.actionType = 'renameFile';
      this.dialogParams.selectedFile = targetFile;
      this.dialogParams.selectedFolder = { path: targetFile.path.substring(0, targetFile.path.lastIndexOf('/')) };
      this.dialogParams.folderPathToReload = this.dialogParams.selectedFolder.path; // Reload the parent folder of the renamed file
      this.dialogVisible = true
    },
    async handleDialogConfirm() {
      this.dialogVisible = false
      try {
        let inputValue = this.inputValue
        if (this.dialogParams.actionType == 'newFile') {
          if (!inputValue.endsWith('.md')) {
            inputValue += '.md'
          }
        }

        const result = await this.performFileManagerAction({
          actionType: this.dialogParams.actionType,
          inputValue: inputValue,
          selectedFolder: this.dialogParams.selectedFolder,
          selectedFile: this.dialogParams.selectedFile,
        })

        let folderPathToReload = '';
        if (this.dialogParams.actionType === 'newFile'||this.dialogParams.actionType === 'renameFile') {
          folderPathToReload = this.dialogParams.selectedFolder.path; // Parent folder of the renamed file
        }
        await this.loadCompleteDirectoryTree(folderPathToReload)

        if (result.actionType === 'newFile') {
          this.$nextTick(async () => {
            const newFile = this.filesInSelectedFolder.find((f) => f.path === result.newPath)
            if (newFile) {
              await this.$store.dispatch('fileManager/handleFileSelect', newFile)
            }
          })
        }

        if (this.isMobile&&result.actionType === 'renameFile') {
          this.$nextTick(async () => {
            if(this.displayState==3){
              this.moveBack();
            }
          })
        }

      } catch (error) {
        this.$message.error(`Operation failed: ${error.message}`)
      }
    },
    handleOpenFolderDialog(payload) {
      this.folderDialog.title = payload.title
      this.folderDialog.placeholder = payload.placeholder
      this.folderDialog.inputValue = payload.initialValue
      this.folderDialog.actionType = payload.actionType
      this.folderDialog.createInRoot = payload.createInRoot || true; // Assign the new flag
      this.folderDialog.visible = true
    },

    async handleFolderDialogConfirm() {
      this.folderDialog.visible = false
      try {
        const result = await this.performFileManagerAction({
          actionType: this.folderDialog.actionType,
          inputValue: this.folderDialog.inputValue,
          selectedFolder: this.selectedFolder,
          selectedFile: null,
          createInRoot: this.folderDialog.createInRoot, // Pass the new flag
        })
        const parentPath = this.selectedFolder
          ? this.selectedFolder.path.substring(0, this.selectedFolder.path.lastIndexOf('/'))
          : ''
        await this.loadCompleteDirectoryTree(parentPath)
      } catch (error) {
        this.$message.error(`${this.$t('common.operationFailed')}: ${error.message}`)
      }
    },
    async confirmDelete(file) {
      const targetFile = file
      if (!targetFile) return

      const targetName = targetFile.name
      const targetPath = targetFile.path
      const isFromRecycle = targetPath.startsWith('Recycle/')

      const confirmMessage = isFromRecycle
        ? this.$t('common.confirmPermanentDelete', { name: targetName })
        : this.$t('common.confirmMoveToRecycle', { name: targetName })

      try {
        await MessageBox.confirm(confirmMessage, this.$t('common.confirmDeletionTitle'), {
          confirmButtonText: this.$t('common.ok'),
          cancelButtonText: this.$t('common.cancel'),
          type: 'warning',
          customClass: this.isMobile ? 'responsive-message-box' : ''
        })
        await this.handleDelete({ targetPath, isFolder: false, isFromRecycle })
        await this.loadCompleteDirectoryTree(this.selectedFolder.path)        
        if(this.isMobile){
          this.$nextTick(async () => {
            if(this.displayState==3){
              this.moveBack();
            }
          })
        }
      } catch (error) {
        if (error === 'cancel') {
          this.$message.info(this.$t('common.deletionCancelled'))
        } else {
          this.$message.error(`${this.$t('common.deletionFailed')} ${error.message || error}`)
        }
      }
    },
    handleFileSelected(file) {
      if (!file) {
        this.vditor.setValue(''); // Clear Vditor if not a markdown file
        return;
      }

      // If it's an image, set the image preview URL
      if (file.imagePreviewUrl) {
        this.imagePreviewUrl = file.imagePreviewUrl;
      } else {
        this.imagePreviewUrl = null; // Clear image preview if not an image
      }

      // If it's a markdown file, set Vditor content
      if (file.name.endsWith('.md') && this.vditor) {
        const mdPath = file.path;
        this.vditor.vditor.lute.SetLinkBase(authState.apiClient.getAttachmentBase(mdPath));
        this.vditor.setValue(file.content || ''); // Use content from store, or empty string
        //this.vditor.focus();
      } else if (this.vditor) {
        this.vditor.setValue(''); // Clear Vditor if not a markdown file
      }
    },
    async autoSaveCheck(value) {
          //console.log("blur0");
          let selectedFile = this.selectedFile
          //console.log(`${this.contentUpdated} && ${selectedFile.path} && ${value.trim() !== selectedFile.content.trim()}`)
          if (this.contentUpdated && selectedFile && value.trim() !== selectedFile.content.trim()) {
            let toSave={ path: selectedFile.path, content: value }
            if (this.autoSaveTimer) {
              clearTimeout(this.autoSaveTimer);
            }
            this.autoSaveTimer = setTimeout(async () => {
              //console.log("blur1");
              if(this.contentUpdated){
                try {              
                  await this.saveFileContent(toSave)
                  this.$message.success(this.$t('common.autoSaveSuccess'))
                } catch (error) {
                  this.$message.error(`${this.$t('common.autoSaveFailed')} ${error.message || error}`)
                }
              }
              this.contentUpdated = false
              this.autoSaveTimer = null;
            },500);
          }
        },
    initVditor() {
      const that = this
      const options = {
        // 【必要改动】1. JS部分: 添加此配置项，指示 Vditor 在滚动时固定工具栏。
        toolbarConfig: {
          pin: true,
        },
        cdn: '/vditor', // 这个路径就是你存放资源的目录
        width: '100%',
        tab: '\t',
        counter: '999999',
        typewriterMode: true,
        mode: 'ir',
        preview: { delay: 100, show: !this.isMobile },
        outline: true,
        upload: {
          max: 5 * 1024 * 1024,
          handler(files) {
            const mdPath = that.$store.state.ui.selectedFile
              ? that.$store.state.ui.selectedFile.path
              : ''
            if (!mdPath) {
              that.$message.error(that.$t('fileManager.noMarkdownSelected'))
              return
            }
            if (!authState.apiClient) {
              that.$message.error(that.$t('common.apiClientNotInitialized'))
              return
            }
            for (let i = 0; i < files.length; i++) {
              const file = files[i]
              authState.apiClient
                .uploadAttachment(mdPath, file)
                .then((response) => {
                  const imgMdStr = `![${file.name}](${response.attachPath})`
                  that.vditor.insertValue(imgMdStr)
                })
                .catch((error) => {
                  console.error('Upload failed:', error)
                  that.$message.error(`Upload failed: ${error.message}`)
                })
            }
          },
        },
        after: () => {
          //const content = defaultText
          const content = ``
          this.vditor.setValue(content)
          //this.vditor.focus()
        },
        input: (value) => {
          this.contentUpdated = true
        },
        keydown: (e) => {
          this.contentUpdated = true
        },
        blur: async (value) => {
          await this.autoSaveCheck(value)
        },
        placeholder: this.$t('fileManager.vditorPlaceholder'),
        toolbar: [
          {
            hotkey: '⌘S',
            name: 'Save',
            tipPosition: 's',
            tip: 'Save',
            className: 'Savebtn',
            icon: floppyDisk,
            click:() => {
              console.log("save");
              const forceSaveFile= async () => {
                this.contentUpdated = false
                let selectedFile = this.selectedFile
                if (selectedFile) {
                  try {
                    await this.saveFileContent({ path: selectedFile.path, content: this.vditor.getValue() })
                    this.$message.success(this.$t('common.saveSuccess'))
                  } catch (error) {
                    this.$message.error(`${this.$t('common.saveFailed')} ${error.message || error}`)
                  }
                }
              };
              forceSaveFile();    
            },
          },
          ...(this.isMobile?[
            /*{
              name: 'rename',
              tip: this.$t('fileManager.renameTooltip'),
              icon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><path fill="currentColor" d="M128 320v576h576V320H128zm-64-64h704v704H64V256z"/><path fill="currentColor" d="M832 64v576h64V64h-64zM320 128v64h576V128H320z"/><path fill="currentColor" d="M832 128v448h64V128h-64zM384 64v64h448V64H384z"/></svg>',
              click: () => {
                this.openRenameDialog()
              }
            },
            {
              name: 'delete',
              tip: this.$t('fileManager.deleteTooltip'),
              icon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><path fill="currentColor" d="M256 192h512v64H256zM192 256h640v64H192z"/><path fill="currentColor" d="M320 320h384v512H320z"/><path fill="currentColor" d="M128 832h768v64H128zM448 448v256h64V448zM512 448v256h64V448z"/></svg>',
              click: () => {
                this.confirmDelete()
              }
            },*/
            "upload",
            "undo",
            "redo",
            "edit-mode",
            {
                name: "more",
                toolbar: [
                    "record",
                    "table",
                    "fullscreen",
                    "both",
                    "code-theme",
                    "content-theme",
                    "export",
                    "outline",
                    "preview",
                    "devtools",
                    "info",
                    "help",
                ],
            }
          ]:[
            "emoji",
            "headings",
            "bold",
            "italic",
            "strike",
            "link",
            "|",
            "list",
            "ordered-list",
            "check",
            "outdent",
            "indent",
            "|",
            "quote",
            "line",
            "code",
            "inline-code",
            "insert-before",
            "insert-after",
            "|",
            "upload",
            "record",
            "table",
            "|",
            "undo",
            "redo",
            "|",
            "fullscreen",
            "edit-mode",
            {
                name: "more",
                toolbar: [
                    "both",
                    "code-theme",
                    "content-theme",
                    "export",
                    "outline",
                    "preview",
                    "devtools",
                    "info",
                    "help",
                ],
            },
          ]),
          ]
      }
      this.vditor = new Vditor('vditor', options)
    },
    reloadContent() {
      if (this.vditor && this.vditor.getValue && !this.selectedFile) {
        this.vditor.setValue(defaultText)
        //this.vditor.focus()
      }
    },
    onLogoClick() {
      this.moveBack()
      //this.vditor.setValue(defaultText)
      //this.reloadContent()
    },
    onBackClick() {
      if (this.displayState === 3) {
        this.backToFileList()
        this.vditor.setValue(defaultText)
      } else if (this.displayState === 2) {
        this.backToFolderList()
        this.vditor.setValue(defaultText)
      }
    },
    openFileHistory(file) {
      this.selectedFileForHistory = file.path
      this.showFileHistoryDialog = true
    },
    openFileSearchDialog() {
      this.showFileSearchDialog = true
    }
  },
}
</script>

<style lang="less">
// The entire original style block is kept here to apply globally, including to the child component.
@import './../assets/styles/style.less';

@file-list-header-height: 40px;

// STYLES FOR DESKTOP (and shared styles)
.index-page {
  width: 100%;
  height: 100%; // 【必要改动】2a. CSS: 使用 100vh 保证全屏，并设为 flex 容器。
  background-color: @white;
  display: flex; // 【必要改动】
  flex-direction: column; // 【必要改动】
  //.flex-box-center(column); // 此行被上面两行标准写法替代

  .image-lightbox {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.8);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    img {
      max-width: 90%;
      max-height: 90%;
    }
  }
  .custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
  }
  .folder-children-count {
    margin-left: 5px;
    color: #999;
    font-size: 12px;
  }
  .main-container {
    width: 100%;
    // 【必要改动】2b. CSS: 以下三行是问题的根源，必须替换。
    // height: calc(100vh - @header-height);
    // position: absolute;
    // top: @header-height;
    flex: 1; // 【必要改动】用此行替代上面三行，让容器自动填充剩余空间。
    position: relative; // absolute改为relative，更安全。
    overflow: hidden; // 增加此行防止容器自身滚动。
  }
  .file-browser-aside {
    background-color: @white;
    padding: 0;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  }
  .modern-sidebar {
    margin: 10px;
    height: calc(100% - 20px); // 此处用 100% 可以正确工作，因为父级已由flex定义高度
    border: 1px solid @border-grey;
    border-radius: 10px;
  }
  .file-browser-row {
    width: 100%;
    height: 100%;
  }
  .folder-column {
    /*border: 1px solid @border-grey; border-radius: 10px;*/
    height: 100%;
    padding-right: 5px;
    position: relative;
  }
  .file-column {
    position: relative;
    height: 100%;
    padding-left: 10px; /*border: 1px solid @border-grey;border-radius: 10px;*/
    border-left: 1px solid @border-grey;
  }
  .column-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 5px;
    border-bottom: 1px solid @border-grey;
    background-color: @white;
    position: sticky;
    top: 0;
    z-index: 10;
    height: @file-list-header-height;
    h3 {
      margin: 0;
      font-size: 16px;
      color: @black;
      flex-grow: 1;
      text-align: center;
    }
    .header-actions {
      display: flex;
      align-items: center;
      flex-shrink: 0;
    }
  }
  .folder-header {
    h3 {
      text-align: center;
    }
  }
  .action-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 4px;
    margin: 0 2px;
    cursor: pointer;
    transition: color 0.2s;
    &:hover {
      color: @brand;
      transform: scale(1.1);
    }
    &.disabled {
      color: @border-grey;
      cursor: not-allowed;
    }
    svg {
      width: @font-medium;
      height: @font-medium;
    }
    svg path {
      fill: @black;
      transition: fill 0.2s ease-in-out;
    }
    &:hover svg path {
      fill: @brand;
    }
  }
  .el-tree {
    background-color: transparent;
    padding-top: 10px;
    overflow-y: auto;
    height: calc(100% - @file-list-header-height);
  }
  .el-tree-node__content {
    height: 30px;
    line-height: 30px;
    padding-right: 8px;
    &:hover {
      background-color: @bg-grey;
    }
  }
  .el-tree-node.is-current > .el-tree-node__content {
    background-color: @brand !important;
    color: @white;
    border-radius: 4px;
    margin-right: 5px;
    .folder-children-count {
      color: @white;
    }
  }
  .file-list-header {
    margin-bottom: 0;
    .file-list-title {
      text-align: left;
    }
  }
  .file-list {
    list-style: none;
    padding: 10px 0 0 0;
    margin: 0;
    overflow-y: auto;
    height: calc(100% - @file-list-header-height);
    li {
      padding: 5px 5px;
      height: 30px;
      cursor: pointer;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      text-align: left;
      display: flex; /* Enable flexbox */
      justify-content: space-between; /* Space out content and button */
      align-items: center; /* Vertically align items */
      .file-name-text {
        flex-grow: 1;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        margin-right: 10px; /* Add some space between text and actions */
        line-height: 1.2;
      }
      &:hover {
        background-color: #f5f7fa;
      }
      &.is-selected {
        background-color: @brand;
        color: @white;
        border-radius: 4px;
      }
      .file-actions {
        display: flex;
        gap: 2px; /* 更小的间距 */
        .el-button--mini {
          margin-left: 0;
          padding: 4px 6px; /* 调整填充 */
          font-size: 12px; /* 调整字体大小 */
          border: none; /* 移除边框 */
          background: transparent; /* 透明背景 */
          color: #606266; /* 默认颜色 */
          &:hover {
            color: @brand; /* 悬停颜色 */
            background-color: rgba(64, 158, 255, 0.1); /* 悬停背景 */
          }
        }
      }
    }
  }
  .vditor-main {
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: flex-start;
    overflow: hidden;
    height:100%;
  }
  .vditor {
    width: 100%;
    height: calc(100% - 20px) !important; // 高度基于父容器，所以用100%
    margin: 10px auto;
    text-align: left;
    border-radius: 8px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    overflow: hidden;
    .vditor-toolbar {
      border-left: none;
      border-right: none;
      border-bottom: 1px solid @border-grey;
      background-color: @bg-grey;
    }
    .vditor-content {
      height: 100%;
      min-height: auto;
      border: none;
      border-top: none;
      overflow: auto;
    }
  }
  .vditor-reset {
    font-size: 14px;
  }
  .vditor-textarea {
    font-size: 14px;
    height: 100% !important;
  }
}

// =======================================================================
// ==================== MOBILE STYLES (@media max-width: 960px) ===================
// =======================================================================
@media (max-width: 960px) {
  .index-page {
    .main-container.mobile-layout {
      // height, position, top 已在 .main-container 中被修正，此处无需改动
      // top: calc(10px + @header-height); // 此行已失效
      left: 0;
      width: 100vw;
      // The main container itself is the viewport wrapper
      .mobile-slider-viewport {
        width: 100%;
        height: 100%;
        overflow: hidden;
        display: flex;
      }
      .mobile-slider-container {
        width: 250vw; // 50vw (folder) + 100vw (file) + 100vw (vditor)
        height: 100%;
        display: flex;
        transition: transform 0.1s ease-in;
      }
      .mobile-column {
        height: 100%;
        flex-shrink: 0; // Prevent columns from shrinking
        overflow-y: hidden; // Let inner content handle scroll
        display: flex;
        flex-direction: column;
      }
      // Specific column widths
      .folder-column-mobile {
        width: 50vw;
        border-right: 1px solid @border-grey; // Keep the visual separator
        padding-right: 5px; // Keep original style
      }
      .file-column-mobile {
        width: 100vw;
        padding-left: 10px; // Keep original style
      }
      .vditor-main-mobile {
        width: 100vw;
        // Styles for vditor wrapper
        padding: 0;
        display: flex;
        justify-content: center;
        align-items: flex-start;
      }
    }

    // Adjust vditor look for mobile
    .vditor {
      height: 100% !important;
      width: 100%;
      margin: 0;
      border-radius: 0;
      box-shadow: none;
      border: none;
      .vditor-content{
        overflow:auto;
      }
    }

    // Ensure headers and lists inside mobile columns scroll correctly
    .folder-column-mobile,
    .file-column-mobile {
      .el-tree,
      .file-list {
        overflow-y: auto;
        height: calc(100% - @file-list-header-height);
      }
    }
    .el-dialog.is-fullscreen {
      .el-dialog__header {
        padding: 15px;
      }
      .el-dialog__body {
        padding: 10px 15px;
      }
      .el-dialog__footer {
        padding: 10px 15px;
      }
    }
  }
  .responsive-message-box {
    width: 80vw;
  }
}
</style>