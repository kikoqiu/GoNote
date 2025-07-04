
<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :fullscreen="true"
    custom-class="search-dialog"
    :before-close="handleClose"
    :title="$t('fileManager.searchFiles')"
  >
    <div class="search-container">
      <div class="search-bar">
        <el-input
          :placeholder="$t('fileManager.searchPlaceholder')"
          v-model="searchQuery"
          @keyup.enter.native="performSearch"
          clearable
        >
          <el-checkbox v-model="useRegex" slot="append">{{ $t('fileManager.useRegex') }}</el-checkbox>
        </el-input>
        <el-button type="primary" @click="performSearch" :loading="isSearchLoading">{{ $t('common.search') }}</el-button>
      </div>

      <div class="search-results">
        <el-table v-if="!isMobile" :data="searchResults" height="calc(100vh - 200px)" style="width: 100%">
          <el-table-column prop="path" :label="$t('fileManager.filePath')" width="250"></el-table-column>
          <el-table-column prop="context" :label="$t('fileManager.contentPreview')">
            <template slot-scope="scope">
              <div v-for="(line, index) in scope.row.context" :key="index" v-html="highlightMatches(line)" class="file-context-desktop"></div>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.actions')" width="80">
            <template slot-scope="scope">
              <el-button size="mini" @click="selectFile(scope.row)">{{ $t('common.select') }}</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div v-else class="mobile-search-results">
          <div v-for="(result, index) in searchResults" :key="index" class="search-result-item">
            <div class="path-actions">
              <span class="file-path">{{ result.path }}</span>
              <el-button size="mini" @click="selectFile(result)">{{ $t('common.select') }}</el-button>
            </div>
            <div class="file-context">
              <div v-for="(line, lineIndex) in result.context" :key="lineIndex" v-html="highlightMatches(line)"></div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="searchError" class="error-message">
        {{ searchError }}
      </div>
    </div>
  </el-dialog>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import { Dialog, Input, Button, Checkbox, Table, TableColumn } from 'element-ui';

export default {
  name: 'FileSearchDialog',
  components: {
    ElDialog: Dialog,
    ElInput: Input,
    ElButton: Button,
    ElCheckbox: Checkbox,
    ElTable: Table,
    ElTableColumn: TableColumn,
  },
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    isMobile: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      dialogVisible: this.visible,
    };
  },
  computed: {
    ...mapState('searchManager', {
      searchResults: state => state.searchResults,
      isSearchLoading: state => state.isLoading,
      searchError: state => state.error,
    }),
    searchQuery: {
      get() {
        return this.$store.state.searchManager.searchQuery;
      },
      set(value) {
        this.$store.dispatch('searchManager/setSearchQuery', value);
      },
    },
    useRegex: {
      get() {
        return this.$store.state.searchManager.useRegex;
      },
      set(value) {
        this.$store.dispatch('searchManager/setUseRegex', value);
      },
    },
  },
  watch: {
    visible(newVal) {
      this.dialogVisible = newVal;
      if (!newVal) {
        this.$store.dispatch('searchManager/clearSearchResults');
      }
    },
  },
  methods: {
    ...mapActions('searchManager', ['performSearch']),
    ...mapActions('fileManager', ['handleFileSelect']),
    handleClose() {
      this.$emit('close');
    },
    selectFile(file) {
      // Extract the file name from the path
      const fileName = file.path.substring(file.path.lastIndexOf('/') + 1);
      // Create a new file object with the required 'name' property
      const fileToSelect = {
        name: fileName,
        path: file.path,
      };
      this.handleFileSelect(fileToSelect);
      this.handleClose();
    },
    highlightMatches(content) {
      const prefixMatch = content.match(/^(\d+):\s*(.*)$/);
      let prefix = '';
      let textToProcess = content;

      if (prefixMatch) {
        const lineNumber = String(prefixMatch[1]).padStart(4, ' ');
        prefix = `<span class="line-number-prefix">${lineNumber}:</span> `;
        textToProcess = prefixMatch[2];
      }

      if (!this.searchQuery) {
        return prefix + this.escapeHtml(textToProcess);
      }

      try {
        const escapedTextToProcess = this.escapeHtml(textToProcess);
        const regex = new RegExp(this.searchQuery, 'gi');
        const highlightedText = escapedTextToProcess.replace(regex, match => `<span class="highlight">${match}</span>`);
        return prefix + highlightedText;
      } catch (e) {
        return prefix + this.escapeHtml(textToProcess); // Invalid regex, return original content escaped
      }
    },
    escapeHtml(text) {
      const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
      };
      return text.replace(/[&<>"']/g, function(m) { return map[m]; });
    }
  },
};
</script>

<style lang="less">
.search-dialog {
  .el-dialog__header {
    border-bottom: 1px solid #eee;
  }
  .search-container {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
  .search-bar {
    display: flex;
    margin-bottom: 20px;
    .el-input {
      margin-right: 10px;
    }
  }
  .search-results {
    flex-grow: 1;
    overflow-y: auto;
  }
  .highlight {
    background-color: yellow;
    font-weight: bold;
  }
  .error-message {
    color: red;
    margin-top: 10px;
  }

  .line-number-prefix {
    display: inline-block;
    /* width: 4em; Fixed width for 4 digits + colon */
    text-align: right;
    margin-right: 0.5em;
    font-weight: bold;
    color: #333;
  }

  .mobile-search-results {
    .search-result-item {
      border: 1px solid #eee;
      border-radius: 5px;
      padding: 10px;
      margin-bottom: 10px;
      .path-actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 5px;
        .file-path {
          font-weight: bold;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          flex-grow: 1;
          margin-right: 10px;
          text-align: left;
        }
      }
      .file-context {
        font-size: 0.9em;
        color: #666;
        white-space: pre-wrap; /* Preserve whitespace and wrap text */
        word-break: break-all; /* Break long words */
        text-align: left;
      }
    }
  }

  .file-context-desktop {
    font-size: 0.9em;
    color: #666;
    white-space: pre-wrap; /* Preserve whitespace and wrap text */
    word-break: break-all; /* Break long words */
    text-align: left;
  }
}

@media (max-width: 768px) {
  .search-bar {
    flex-direction: column;
    .el-input {
      margin-right: 0;
      margin-bottom: 10px;
    }
  }
}
</style>
