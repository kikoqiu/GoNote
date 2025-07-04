import Vue from 'vue';
import Vuex from 'vuex';
import ui from './ui';
import fileManager from './fileManager';
import fileHistoryManager from './fileHistoryManager';

import exportManager from './exportManager';
import searchManager from './searchManager';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    ui,
    fileManager,
    fileHistoryManager,
    searchManager,
    exportManager,
    
    // auth module is already imported in main.js, so we don't need to import it here
    // but if it were a Vuex module, it would be added here.
  },
});
