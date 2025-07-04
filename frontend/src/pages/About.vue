<!-- @format -->

<template>
  <div class="about-page">
    <div class="button-group">
      <el-button round @click="onBackToMainPage">{{ $t('aboutPage.backToMainPage') }}</el-button>
    </div>
    <PreviewVditor :pdata="pdata" @callback="handleCallBack" />
  </div>
</template>

<script>
import PreviewVditor from '@components/PreviewVditor'
import { mapActions } from 'vuex';
import About from '@config/about'
export default {
  name: 'about',

  data() {
    return {
      isLoading: true,
      pdata: About,
    }
  },

  created() {},

  components: {
    PreviewVditor,
  },

  mounted() {},

  methods: {
    handleCallBack() {
      this.isLoading = false
    },
    ...mapActions('ui', ['hideAbout']),
    onBackToMainPage() {
      this.hideAbout();
    },
  },
}
</script>

<style lang="less">
.about-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  /*min-height: 100vh;  Ensure it takes full viewport height */
  padding: 20px;
  box-sizing: border-box;
  width: 100%;

  .button-group {
    margin-bottom: 20px;
  }

  .vditor {
    max-width: 800px; /* Limit width for better readability */
    width: 100%;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    border-radius: 8px;
    overflow: hidden;

    @media (max-width: 960px) {
      overflow: auto; /* Allow scrolling on mobile */
    }
  }
}
</style>
