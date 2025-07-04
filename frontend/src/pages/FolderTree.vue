<template>
  <div style="width:100%;height:100%">
    <div class="column-header folder-header">
      <h3>{{ $t('fileManager.directories') }}</h3>
      <div class="header-actions">
        <el-tooltip :content="$t('fileManager.searchFiles')" placement="top">
          <span class="action-icon" @click="openSearchDialog">
            <icon name="magnifying-glass" />
          </span>
        </el-tooltip>
        <el-tooltip :content="$t('fileManager.newFolderTooltip')" placement="top">
          <span class="action-icon" @click="openNewFolderDialog">
            <icon name="folder" />
          </span>
        </el-tooltip>
        <el-tooltip :content="$t('fileManager.renameTooltip')" placement="top">
          <span
            class="action-icon"
            @click="openRenameDialog"
            :class="{ disabled: !selectedFolder }"
          >
            <icon name="copy" />
          </span>
        </el-tooltip>
        <el-tooltip :content="$t('fileManager.deleteTooltip')" placement="top">
          <span class="action-icon" @click="confirmDelete" :class="{ disabled: !selectedFolder }">
            <icon name="trash-can" />
          </span>
        </el-tooltip>
      </div>
    </div>
    <el-tree
      ref="folderTree"
      :data="directoryTree"
      :props="defaultProps"
      accordion
      @node-click="onFolderClick"
      :render-content="renderContent"
      highlight-current
      node-key="path"
      :current-node-key="selectedFolder ? selectedFolder.path : ''"
      :key="treeKey"
      :default-expanded-keys="expandedFolderKeys"
      @node-expand="onNodeExpand"
      @node-collapse="onNodeCollapse"
      :expand-on-click-node="false"
    >
    </el-tree>
  </div>
</template>

<script>
import { mapState, mapActions, mapGetters } from 'vuex';
import Icon from '@/components/Icon';
import { Col, Tooltip, Tree, MessageBox } from 'element-ui';

export default {
  name: 'FolderTree',
  components: {
    Icon,
    ElCol: Col,
    ElTooltip: Tooltip,
    ElTree: Tree,
  },
  props: {
    isMobile: Boolean,
    displayState: Number,
    treeKey: Number,
  },
  data() {
    return {
      defaultProps: {
        children: 'children',
        label: 'label',
        isLeaf: 'isLeaf',
      },
      expandedFolderKeys: [],
      createInRoot: true, // New data property
    };
  },
  computed: {
    ...mapState('ui', ['selectedFolder']),
    ...mapGetters('fileManager', ['directoryTree']),
  },
  watch: {
    selectedFolder(newFolder) {
      if (this.$refs.folderTree) {
        if (newFolder) {
          this.$refs.folderTree.setCurrentKey(newFolder.path);
          const node = this.$refs.folderTree.getNode(newFolder.path);
          if (node) {
            let parent = node.parent;
            while (parent) {
              if (!parent.expanded) {
                parent.expanded = true;
              }
              parent = parent.parent;
            }
          }
        } else {
          this.$refs.folderTree.setCurrentKey(null);
        }
      }
    },
  },
  methods: {
    ...mapActions('fileManager', [
      'loadCompleteDirectoryTree',
      'handleFolderClick',
      'handleDelete',
    ]),
    openSearchDialog() {
      this.$emit('open-search-dialog');
    },
    onFolderClick(data) {
      this.handleFolderClick(data);
    },
    onNodeExpand(data) {
      if (!this.expandedFolderKeys.includes(data.path)) {
        this.expandedFolderKeys.push(data.path);
      }
    },
    onNodeCollapse(data) {
      const index = this.expandedFolderKeys.indexOf(data.path);
      if (index > -1) {
        this.expandedFolderKeys.splice(index, 1);
      }
    },
    renderContent(h, { node, data }) {
      let displayName = data.name;
      if (data.type === 'directory' && data.name === 'Recycle') {
        displayName = this.$t('common.recycle');
      }
      if (data.type === 'directory') {
        return (
          <span class="custom-tree-node">
            <span>{displayName}</span>
            {data.mdFileCount !== undefined && data.mdFileCount > 0 && (
              <span class="folder-children-count">({data.mdFileCount})</span>
            )}
          </span>
        );
      }
      return (
        <span class="custom-tree-node">
          <span>{displayName}</span>
        </span>
      );
    },

    // Action Handlers
    openNewFolderDialog() {
      this.$emit('open-folder-dialog', {
        actionType: 'newFolder',
        title: this.$t('fileManager.newFolderTitle'),
        placeholder: this.$t('fileManager.newFolderInputPlaceholder'),
        initialValue: '',
        createInRoot: this.createInRoot, // Pass the new flag
      });
    },
    openRenameDialog() {
      if (!this.selectedFolder) return;
      this.$emit('open-folder-dialog', {
        actionType: 'renameFolder',
        title: this.$t('fileManager.renameFolderTitle'),
        placeholder: this.$t('fileManager.renameFolderInputPlaceholder'),
        initialValue: this.selectedFolder.name,
      });
    },
    async confirmDelete() {
      if (!this.selectedFolder) return;

      const targetName = this.selectedFolder.name;
      const targetPath = this.selectedFolder.path;
      const isFromRecycle = targetPath.startsWith('Recycle/') || targetPath === 'Recycle';

      let confirmMessage = isFromRecycle
        ? this.$t('common.confirmPermanentDelete', { name: targetName })
        : this.$t('common.confirmDeleteFolderContents', { name: targetName });

      try {
        await MessageBox.confirm(confirmMessage, this.$t('common.confirmDeletionTitle'), {
          confirmButtonText: this.$t('common.ok'),
          cancelButtonText: this.$t('common.cancel'),
          type: 'warning',
          customClass: this.isMobile ? 'responsive-message-box' : '',
        });
        await this.handleDelete({ targetPath, isFolder: true, isFromRecycle });
        const parentPath = targetPath.substring(0, targetPath.lastIndexOf('/'));
        await this.loadCompleteDirectoryTree(parentPath);
      } catch (error) {
        if (error === 'cancel') {
          this.$message.info(this.$t('common.deletionCancelled'));
        } else {
          this.$message.error(`${this.$t('common.deletionFailed')} ${error.message || error}`);
        }
      }
    },
  },
};
</script>

<!-- Note: Styles are defined in the parent component (Main.vue) to maintain global scope and original behavior. -->
