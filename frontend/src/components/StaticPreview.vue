<!-- src/components/StaticPreview.vue -->
<template>
  <div class="static-preview-wrapper" v-loading="isLoading" element-loading-text="内容生成中...">
    <div
      v-show="!isLoading"
      ref="preview"
      id="capture-target"
      class="preview-content-area vditor-preview"
    >
      <!-- Vditor 将会把渲染后的 HTML 内容注入到这里 -->
    </div>
  </div>
</template>

<script>
import Vditor from 'vditor';
import 'vditor/dist/index.css';

export default {
  name: 'StaticPreview',
  props: {
    pdata: {
      type: String,
      required: true,
      default: '',
    },
  },
  data() {
    return {
      isLoading: true,
    };
  },
  mounted() {
    this.renderMarkdown();
  },
  methods: {
    renderMarkdown() {
      const previewElement = this.$refs.preview;
      if (!previewElement) {
        console.error('StaticPreview: Preview element not found.');
        this.isLoading = false;
        return;
      }
      try {
        Vditor.preview(previewElement, this.pdata, {
          mode: 'light',
          hljs: { enable: true, style: 'native', lineNumber: true },
          markdown: { toc: true },
          preview: { actions: [] },
        });
      } catch (error) {
        console.error('Vditor preview rendering failed:', error);
      } finally {
        this.$nextTick(() => {
          this.isLoading = false;
        });
      }
    },
  },
};
</script>

<style lang="less" scoped>
.static-preview-wrapper {
  width: 100%;
  min-height: 200px;
}

.preview-content-area {
  /* --- 关键改动：精确模拟A4纸张 --- */
  /* A4 纸张的标准宽度 */
  width: 210mm; 
  /* 设置一个最小高度，使其看起来像一页纸 */
  min-height: 297mm;
  /* 关键：确保 padding 和 border 被计算在 width/height 内 */
  box-sizing: border-box;

  padding: 15.0mm 25.4mm;

  /* 页面在屏幕上居中显示，并带有阴影，模拟真实效果 */
  margin: 0 auto;
  background-color: #ffffff;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.1);
  
  /* 强制换行规则依然重要 */
  text-align: left;
  word-wrap: break-word;
  overflow-wrap: break-word;
  word-break: break-word;

  :deep(.vditor-reset) {
    h1 {
      text-align: center;
    }
  }
}

/* 响应式设计：在小屏幕上，移除A4模拟，恢复为占满宽度 */
@media (max-width: 768px) {
  .preview-content-area {
    /* 覆盖A4尺寸设置 */
    width: 100% !important;
    min-height: auto; /* 移除最小高度 */
    margin: 0 !important;
    padding: 15px;
    box-shadow: none;
  }
}
</style>