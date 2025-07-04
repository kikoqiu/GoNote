<template>
  <el-dialog
    :visible.sync="dialogVisible"
    fullscreen
    :before-close="handleClose"
    custom-class="file-history-dialog"
  >
    <div slot="title" class="dialog-title">
      <i class="el-icon-time"></i>
      <span>{{ $t('fileHistory.fileHistory') }}: {{ currentFilePath }}</span>
    </div>

    <el-row :gutter="20" class="dialog-content">
      <el-col :span="8" class="history-list-container">
        <el-card class="box-card" shadow="never">
          <div slot="header" class="clearfix">
            <span>{{ $t('fileHistory.versionList') }}</span>
            <el-button
              style="float: right; padding: 3px 0"
              type="text"
              icon="el-icon-refresh"
              @click="refreshHistory"
              :loading="loadingHistory"
              >{{ $t('fileHistory.refresh') }}</el-button
            >
          </div>
          <el-table
            :data="history"
            style="width: 100%"
            highlight-current-row
            @current-change="handleHistorySelect"
            v-loading="loadingHistory"
            :element-loading-text="$t('fileHistory.loadingHistory')"
            element-loading-spinner="el-icon-loading"
            element-loading-background="rgba(0, 0, 0, 0.8)"
            :empty-text="$t('fileHistory.noVersionSelected')"
          >
            <el-table-column :label="$t('fileHistory.date')" width="auto">
              <template slot-scope="scope">
                {{ formatDate(scope.row.timestamp) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('fileHistory.actions')" width="100">
              <template slot-scope="scope">
                <el-tooltip :content="$t('fileHistory.viewContent')" placement="top">
                  <el-button
                    size="mini"
                    icon="el-icon-document"
                    circle
                    @click="viewContent(scope.row.id)"
                    :loading="loadingVersion && currentVersionId === scope.row.id"
                  ></el-button>
                </el-tooltip>
                <el-tooltip :content="$t('fileHistory.viewDiff')" placement="top">
                  <el-button
                    size="mini"
                    type="primary"
                    icon="el-icon-files"
                    circle
                    @click="viewDiff(scope.row.id)"
                    :loading="loadingVersion && currentVersionId === scope.row.id"
                  ></el-button>
                </el-tooltip>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="16" class="content-display-container">
        <el-card class="box-card" shadow="never">
          <div slot="header" class="clearfix">
            <span>{{ $t('fileHistory.contentPreviewDiff') }}</span>
            <el-button
              v-if="currentVersionContent"
              style="float: right; padding: 3px 0"
              type="text"
              icon="el-icon-document-copy"
              @click="useThisVersion"
              >{{ $t('fileHistory.useThisVersion') }}</el-button
            >
          </div>
          <div v-loading="loadingVersion" :element-loading-text="$t('fileHistory.loadingVersionContent')" class="content-area">
            
            <!-- 【修改】将原来的 diff-viewer 替换为 diff2html 的渲染容器 -->
            <div v-if="diffHtml" v-html="diffHtml" class="diff-container"></div>
            
            <div v-else-if="currentVersionContent">
              <h3>{{ $t('fileHistory.versionContent', { id: currentVersionId }) }}</h3>
              <pre class="markdown-content">{{ currentVersionContent }}</pre>
            </div>
            <div v-else class="no-selection-placeholder">
              <p>{{ $t('fileHistory.noVersionSelected') }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </el-dialog>
</template>

<script>
import { mapState, mapActions, mapMutations } from 'vuex'
import { Message, MessageBox } from 'element-ui'

// 【新增】引入 diff 和 diff2html
import * as Diff from 'diff'
//import { Diff2Html } from 'diff2html'
import * as Diff2Html from 'diff2html';
import 'diff2html/bundles/css/diff2html.min.css' // 【新增】引入 diff2html 的核心样式

export default {
  name: 'FileHistoryDialog',
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    filePath: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      dialogVisible: this.visible,
    }
  },
  computed: {
    ...mapState('fileHistoryManager', [
      'currentFilePath',
      'history',
      'loadingHistory',
      'loadingVersion',
      'currentVersionContent',
      'currentVersionId',
      'diffContent',
    ]),

    // 【新增】计算属性，用于生成 diff HTML
    diffHtml() {
      // 只有在 diffContent 有效时才进行计算
      if (!this.diffContent || !this.diffContent.current || !this.diffContent.historical) {
        return ''
      }

      const { historical, current } = this.diffContent
      
      // 1. 使用 js-diff 生成统一的 diff 补丁
      const patch = Diff.createTwoFilesPatch(
        'historical-version', // 旧文件名（可自定义）
        'current-version',    // 新文件名（可自定义）
        historical,           // 旧文件内容
        current,              // 新文件内容
        '',                   // 旧文件头信息（可选）
        '',                   // 新文件头信息（可选）
        { context: 5 }        // 上下文显示的行数
      )

      // 2. 使用 diff2html 将补丁解析为 JSON 格式
      const diffJson = Diff2Html.parse(patch)
      
      // 3. 使用 diff2html 将 JSON 渲染为 HTML
      const diffHtml = Diff2Html.html(diffJson, {
        drawFileList: false,       // 不显示文件列表头
        matching: 'lines',         // 'lines', 'words', or 'none'
        outputFormat: 'side-by-side', // 默认使用并排视图 'side-by-side' 或 'line-by-line'
        renderNothingWhenEmpty: true,
      })

      return diffHtml
    },
  },
  watch: {
    visible(newVal) {
      this.dialogVisible = newVal
      if (newVal && this.filePath) {
        this.SET_CURRENT_FILE_PATH(this.filePath)
        this.fetchFileHistory(this.filePath)
      } else if (!newVal) {
        this.CLEAR_HISTORY_STATE()
      }
    },
    dialogVisible(newVal) {
      this.$emit('update:visible', newVal)
    },
  },
  methods: {
    ...mapActions('fileHistoryManager', [
      'fetchFileHistory',
      'fetchFileVersion',
      'applyFileVersion',
      'generateDiff',
    ]),
    ...mapMutations('fileHistoryManager', [
      'SET_CURRENT_FILE_PATH',
      'SET_CURRENT_VERSION_ID',
      'SET_CURRENT_VERSION_CONTENT',
      'SET_DIFF_CONTENT',
      'CLEAR_HISTORY_STATE',
    ]),
    handleClose(done) {
      this.CLEAR_HISTORY_STATE()
      done()
    },
    formatDate(timestamp) {
      if (!timestamp) return ''
      const date = new Date(timestamp)
      return date.toLocaleString()
    },
    refreshHistory() {
      if (this.currentFilePath) {
        this.fetchFileHistory(this.currentFilePath)
      }
    },
    handleHistorySelect(row) {
      if (row) {
        this.SET_CURRENT_VERSION_ID(row.id)
        this.SET_CURRENT_VERSION_CONTENT(null)
        this.SET_DIFF_CONTENT(null)
      }
    },
    viewContent(versionId) {
      if (this.currentFilePath && versionId) {
        this.fetchFileVersion({ filePath: this.currentFilePath, versionId })
        this.SET_DIFF_CONTENT(null)
      }
    },
    viewDiff(versionId) {
      if (this.currentFilePath && versionId) {
        // 确保 Vuex action 'generateDiff' 会获取两个版本的内容
        // 并将其存入 state 的 diffContent: { historical: '...', current: '...' }
        this.generateDiff({ filePath: this.currentFilePath, versionId })
        this.SET_CURRENT_VERSION_CONTENT(null)
      }
    },
    useThisVersion() {
      if (!this.currentVersionContent) {
        Message.warning(this.$t('fileHistory.selectVersionToLoad'))
        return
      }

      MessageBox.confirm(
        this.$t('fileHistory.confirmRevert'),
        this.$t('fileHistory.warning'),
        {
          confirmButtonText: this.$t('common.confirm'),
          cancelButtonText: this.$t('common.cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          this.applyFileVersion({
            filePath: this.currentFilePath,
            versionContent: this.currentVersionContent,
            comment: `Reverted to version ID: ${this.currentVersionId}`,
          }).then(() => {
            this.$emit('file-reverted')
            this.dialogVisible = false
          })
        })
        .catch(() => {
          Message.info(this.$t('fileHistory.revertCancelled'))
        })
    },
  },
}
</script>

<style lang="less">
.file-history-dialog {
  .el-dialog__header {
    padding: 20px;
    border-bottom: 1px solid #eee;
    .dialog-title {
      font-size: 20px;
      font-weight: bold;
      color: #333;
      .el-icon-time {
        margin-right: 10px;
      }
    }
  }
  .el-dialog__body {
    padding: 20px;
    height: calc(100vh - 70px);
    display: flex;
    flex-direction: column;
  }

  .dialog-content {
    flex: 1;
    display: flex;
    height: 100%;
  }

  .history-list-container,
  .content-display-container {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .box-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    .el-card__header {
      padding: 10px 20px;
      border-bottom: 1px solid #ebeef5;
    }
    .el-card__body {
      flex: 1;
      overflow-y: auto;
      padding: 10px 20px;
    }
  }

  .content-area {
    min-height: 200px;
    position: relative;
  }

  .markdown-content {
    white-space: pre-wrap;
    word-wrap: break-word;
    font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
    background-color: #f5f5f5;
    padding: 15px;
    border-radius: 4px;
    border: 1px solid #ddd;
    max-height: calc(100vh - 250px);
    overflow-y: auto;
    text-align: left;
  }

  .no-selection-placeholder {
    text-align: center;
    color: #909399;
    padding: 50px;
  }

  // 【新增】diff-container 的样式和响应式调整
  .diff-container {
    .d2h-file-header {
      display: none; // 隐藏 diff2html 自带的文件头，因为我们已有标题
    }
    .d2h-file-wrapper {
      border: none; // 移除外层边框，使其与 el-card 融合
      border-radius: 0;
      margin-bottom: 0;
    }
    // 调整字体大小使其与UI更协调
    .d2h-code-line, .d2h-code-side-line {
        font-size: 14px;
    }
  }
  
  /* 【新增】对原始 diff-viewer 的清理，如果您已删除模板中的该部分，这行可以不加 */
  .diff-viewer { display: none; }

  /* 响应式调整 */
  @media (max-width: 768px) {
    .dialog-content {
      flex-direction: column;
    }
    .el-col {
      width: 100%;
    }
    .history-list-container {
      margin-bottom: 20px;
      // 移动端限制历史列表高度，防止过长
      height: 30vh;
      min-height: 200px;
    }
    .markdown-content, .diff-pane {
      max-height: 300px;
    }

    // 【新增】移动端 diff 视图的样式调整
    .diff-container {
        .d2h-files-diff {
            // 将并排视图的两列垂直堆叠
            flex-direction: column;
        }
        .d2h-file-side-diff {
            // 让每一列都占满宽度
            width: 100%;
            margin-bottom: 1em; // 在两列之间增加一些间距
        }
    }
  }
}
</style>