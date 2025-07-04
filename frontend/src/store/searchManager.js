
import { authState } from './auth';

const state = {
  searchQuery: '',
  useRegex: false,
  searchResults: [],
  isLoading: false,
  error: null,
};

const mutations = {
  SET_SEARCH_QUERY(state, query) {
    state.searchQuery = query;
  },
  SET_USE_REGEX(state, useRegex) {
    state.useRegex = useRegex;
  },
  SET_SEARCH_RESULTS(state, results) {
    state.searchResults = results;
    state.error = null;
  },
  SET_IS_LOADING(state, isLoading) {
    state.isLoading = isLoading;
  },
  SET_ERROR(state, error) {
    state.error = error;
    state.searchResults = [];
  },
  CLEAR_SEARCH(state) {
    state.searchQuery = '';
    state.useRegex = false;
    state.searchResults = [];
    state.isLoading = false;
    state.error = null;
  },
  CLEAR_SEARCH_RESULTS(state) {
    state.searchResults = [];
    state.isLoading = false;
    state.error = null;
  },
};

const actions = {
  async performSearch({ commit, state, rootState }) {
    if (!state.searchQuery.trim()) {
      commit('SET_SEARCH_RESULTS', []);
      return;
    }
    commit('SET_IS_LOADING', true);
    try {
      const results = await authState.apiClient.searchFiles(state.searchQuery, state.useRegex);
      commit('SET_SEARCH_RESULTS', results);
    } catch (error) {
      commit('SET_ERROR', error.message || 'An unknown error occurred');
    } finally {
      commit('SET_IS_LOADING', false);
    }
  },
  setSearchQuery({ commit }, query) {
    commit('SET_SEARCH_QUERY', query);
  },
  setUseRegex({ commit }, useRegex) {
    commit('SET_USE_REGEX', useRegex);
  },
  clearSearch({commit}) {
    commit('CLEAR_SEARCH');
  },
  clearSearchResults({ commit }) {
    commit('CLEAR_SEARCH_RESULTS');
  }
};

const getters = {
  searchResults: state => state.searchResults,
  isSearchLoading: state => state.isLoading,
  searchError: state => state.error,
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
};
