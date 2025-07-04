<!-- @format -->

<template>
  <div id="app">
    <Main />
    <template v-if="isExporting">
      <div class="overlay">
        <Export />
      </div>
    </template>
    <template v-else-if="isShowingAbout">
      <div class="overlay">
        <About />
      </div>
    </template>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import Main from '@/pages/Main';
import About from '@/pages/About';
import Export from '@/pages/Export';

export default {
  name: 'app',

  components: {
    Main,
    About,
    Export,
  },

  computed: {
    ...mapState('exportManager', ['isExporting']),
    ...mapState('ui', ['isShowingAbout']),
  },
};
</script>

<style lang="less">
@import './assets/styles/element/index.css';
@import './assets/styles/style.less';
@import '~vditor/dist/index.css';

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  
  width:100%;
  height:100% !important;
  overflow: hidden;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1000; /* High z-index to ensure it's on top */
  background-color: white; /* Or any background color you prefer for the overlay */
  display: flex;
  justify-content: start;
  align-items: start;
  overflow: scroll;
}
</style>
