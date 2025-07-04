<template>
  <div class="export-page">
    <header v-if="exportType !== 'ppt'" class="export-header">
      <div class="button-group">
        <el-button round @click="onBackToMainPage">{{ $t('export.backToMainPage') }}</el-button>
        <el-button round @click="onExportBtnClick" type="primary" :loading="isSaving">
          {{ isSaving ? $t('export.exporting') : $t('export.generateExport') }}
        </el-button>
      </div>
    </header>

    <main class="export-body">
      <div v-if="exportType !== 'ppt'" class="preview-container">
        <StaticPreview :pdata="pdata" />
      </div>

      <div v-if="exportType === 'ppt'" class="export-ppt">
        <div class="reveal">
          <div class="slides">
            <section data-markdown data-separator="---" data-separator-vertical="--">
              <textarea data-template>{{ savedMdContent }}</textarea>
            </section>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import { generateScreenshot } from '@helper/export';
import html2pdf from 'html2pdf.js';
import Reveal from 'reveal.js/js/reveal';
import 'reveal.js/css/reset.css';
import 'reveal.js/css/reveal.css';
import 'reveal.js/css/theme/beige.css';
import StaticPreview from '@components/StaticPreview';

export default {
  name: 'Export',
  components: {
    StaticPreview,
  },
  data() {
    return {
      pdata: localStorage.getItem('vditorvditor') || '',
      savedMdContent: localStorage.getItem('vditorvditor') || '',
      isSaving: false,
    };
  },
  computed: {
    ...mapState('exportManager', ['isExporting', 'exportType']),
  },
  mounted() {
    if (this.exportType === 'ppt') {
      this.$nextTick(() => this.initReveal());
    }
  },
  methods: {
    ...mapActions('exportManager', ['endExport']),

    onBackToMainPage() {
      this.endExport();
    },

    onExportBtnClick() {
      if (this.isSaving) return;
      this.isSaving = true;

      const captureElement = document.querySelector('#capture-target');
      if (!captureElement) {
        this.$message.error('无法找到导出目标元素');
        this.isSaving = false;
        return;
      }
      
      const originalShadow = captureElement.style.boxShadow;
      captureElement.style.boxShadow = 'none';

      const exportPromise = new Promise((resolve, reject) => {
        this.$nextTick(() => {
          if (this.exportType === 'png') {
            this.exportAndDownloadImg().then(resolve).catch(reject);
          } else if (this.exportType === 'pdf') {
            this.exportAndDownloadPdf().then(resolve).catch(reject);
          } else {
            reject(new Error('Unknown export type'));
          }
        });
      });

      exportPromise.finally(() => {
        captureElement.style.boxShadow = originalShadow;
        this.isSaving = false;
        this.endExport();
      });
    },

    exportAndDownloadImg() {
      // PNG导出逻辑保持不变，它工作正常
      return new Promise(async (resolve, reject) => {
        try {
          const element = document.querySelector('#capture-target');
          if (!element) throw new Error('Capture target element not found');

          const filename = this.$utils.getExportFileName() + '.png';
          const canvas = await generateScreenshot(element);
          const link = document.createElement('a');
          link.download = filename;
          link.href = canvas.toDataURL('image/png');
          link.click();
          resolve();
        } catch (error) {
          console.error('导出图片失败:', error);
          this.$message.error(this.$t('export.imageExportFailed'));
          reject(error);
        }
      });
    },

    exportAndDownloadPdf() {
      return new Promise((resolve, reject) => {
        const element = document.querySelector('#capture-target');
        if (!element) {
          this.$message.error(this.$t('export.pdfExportFailed'));
          return reject(new Error('Capture target element not found'));
        }

        const filename = this.$utils.getExportFileName() + '.pdf';
        const opt = {
          margin: 0, // 0.5in on all sides
          filename,
          image: { type: 'jpeg', quality: 0.98 },
          jsPDF: { unit: 'in', format: 'a4', orientation: 'portrait' },
          html2canvas: {
            scale: 2, // 提高分辨率
            useCORS: true,
            logging: false,
            backgroundColor: '#ffffff',            
          },
        };
        
        html2pdf()
          .set(opt)
          .from(element)
          .save()
          .then(() => {
            this.$message.success(this.$t('export.pdfExportSuccess'));
            resolve();
          })
          .catch((error) => {
            console.error('PDF导出失败:', error);
            this.$message.error(this.$t('export.pdfExportFailed'));
            reject(error);
          });
      });
    },

    initReveal() {
      window.Reveal = Reveal;
      const revealSourcePath = `https://cdn.jsdelivr.net/npm/reveal.js@3.8.0`;
      Reveal.initialize({
        controls: true,
        progress: true,
        center: true,
        hash: true,
        transition: 'slide',
        display: 'block',
        dependencies: [
          {
            src: `${revealSourcePath}/plugin/markdown/marked.js`,
            condition: function () {
              return !!document.querySelector('[data-markdown]');
            },
          },
          {
            src: `${revealSourcePath}/plugin/markdown/markdown.js`,
            condition: function () {
              return !!document.querySelector('[data-markdown]');
            },
          },
          { src: `${revealSourcePath}/plugin/highlight/highlight.js`, async: true },
          { src: `${revealSourcePath}/plugin/search/search.js`, async: true },
          { src: `${revealSourcePath}/plugin/zoom-js/zoom.js`, async: true },
          { src: `${revealSourcePath}/plugin/notes/notes.js`, async: true },
        ],
      });
    },
  },
};
</script>

<style lang="less" scoped>
/* 所有样式都保持不变，因为问题不在于此 */
.export-page {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  background-color: #f0f2f5;
}
.export-header {
  flex-shrink: 0;
  padding: 12px 24px;
  background-color: #ffffff;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  text-align: center;
  z-index: 10;
}
.button-group {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
}
.export-body {
  flex-grow: 1;
  overflow-y: auto;
  padding: 24px;
}
.preview-container {
  width: 100%;
}
.export-ppt {
  width: 100%;
  .reveal {
    font-size: 2em;
    background-color: #ffffff;
    height: calc(100vh - 60px);
    h1 {
      font-size: 2em !important;
    }
  }
}
@media (max-width: 768px) {
  .export-header {
    padding: 12px;
  }
  .button-group {
    flex-direction: column;
    width: 100%;
    gap: 10px;
    .el-button {
      width: 100%;
      margin-left: 0 !important;
    }
  }
  .export-body {
    padding: 12px;
  }
}
</style>