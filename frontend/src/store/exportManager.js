
/** @format */


export const exportState = {
  isExporting: false,
  exportType: null, // 'png', 'pdf', 'ppt'
};

export const exportMutations = {
  setExporting(state, isExporting) {
    state.isExporting = isExporting;
  },
  setExportType(state, exportType) {
    state.exportType = exportType;
  },
};

export const exportActions = {
  startExport({ commit }, exportType) {
    commit('setExportType', exportType);
    commit('setExporting', true);
  },
  endExport({ commit }) {
    commit('setExporting', false);
    commit('setExportType', null);
  },
};

export default {
  namespaced: true,
  state: exportState,
  mutations: exportMutations,
  actions: exportActions,
};
