// frontend/src/store/ui.js
import { appTitle } from '../config/constant'

const state = {
  isShowingAbout: false,
  displayState: 1, // 1: show folder/file list, 2: show file list full width, 3: show editor
  selectedFolder: null, // Stores the selected folder object { name, path, childrenCount }
  selectedFile: null, // Stores the selected file object { name, path }
}

const mutations = {
  SET_DISPLAY_STATE(state, newState) {
    state.displayState = newState
  },
  SET_SELECTED_FOLDER(state, folder) {
    state.selectedFolder = folder
  },
  SET_SELECTED_FILE(state, file) {
    state.selectedFile = file
  },
  setAboutShowing(state, isShowing) {
    state.isShowingAbout = isShowing;
  },
}

const actions = {
  moveBack({ commit, state }) {
    if (state.displayState == 1) {
      commit('SET_DISPLAY_STATE', 1)
    } else {
      commit('SET_DISPLAY_STATE', state.displayState - 1)
    }
    //commit('SET_SELECTED_FOLDER', null);
    //commit('SET_SELECTED_FILE', null);
  },

  // Action to set display state and title based on logo click
  resetDisplay({ commit }) {
    commit('SET_DISPLAY_STATE', 1)
    //commit('SET_SELECTED_FOLDER', null);
    commit('SET_SELECTED_FILE', null)
  },

  // Action when a folder is selected
  selectFolder({ commit }, folder) {
    if(folder!=null){
      commit('SET_DISPLAY_STATE', 2)
      commit('SET_SELECTED_FOLDER', folder)
    }
    //commit('SET_SELECTED_FILE', null) // Clear selected file when a new folder is selected
  },

  // Action when a file is selected
  selectFile({ commit }, file) {
    if(file){
      commit('SET_DISPLAY_STATE', 3)
    }
    commit('SET_SELECTED_FILE', file)
  },

  // Action to go back from editor to file list
  backToFileList({ commit, state }) {
    commit('SET_DISPLAY_STATE', 2)
    commit('SET_SELECTED_FILE', null)
  },

  // Action to go back from file list to folder list
  backToFolderList({ commit }) {
    commit('SET_DISPLAY_STATE', 1)
    //commit('SET_SELECTED_FOLDER', null);
    commit('SET_SELECTED_FILE', null)
  },
  showAbout({ commit }) {
    commit('setAboutShowing', true);
  },
  hideAbout({ commit }) {
    commit('setAboutShowing', false);
  },
}

const getters = {
  currentDisplayTitle(state) {
    if (state.displayState === 1) {
      return appTitle // Application title
    } else if (state.displayState === 2 && state.selectedFolder) {
      return state.selectedFolder.name // Selected directory name
    } else if (state.displayState === 3 && state.selectedFile) {
      return state.selectedFile.name // Selected file name
    }
    return ''
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
}
