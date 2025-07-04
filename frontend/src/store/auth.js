import { reactive } from 'vue';
import ApiClient from '@/apiClient';
import store from './index'; // Import the store

const AUTH_TOKEN_KEY = 'auth_token';
const API_URL_KEY = 'api_url';

export const authState = reactive({
  token: null,
  isAuthenticated: false,
  apiClient: null,
  apiUrl: null,
  username: null,
  isLoginDialogVisible: false, // Control dialog visibility
});

export const authService = {
  async initAuth() {
    const token = localStorage.getItem(AUTH_TOKEN_KEY);
    const apiUrl = localStorage.getItem(API_URL_KEY) || window.location.origin;
    authState.apiUrl = apiUrl;

    if (token) {
      try {
        const [username, password] = atob(token).split(':');
        authState.token = token;
        authState.username = username;
        authState.apiClient = new ApiClient(apiUrl, username, password);
        await authState.apiClient.listDirectory('');
        authState.isAuthenticated = true;
        store.dispatch('fileManager/loadCompleteDirectoryTree');
      } catch (e) {
        console.error('Failed to initialize auth from localStorage', e);
        //this.logout(); // Clear invalid state
        authState.apiClient = null;
        authState.isLoginDialogVisible = true;
      }
    } else {
      // Not authenticated, show login dialog
      authState.isLoginDialogVisible = true;
    }
  },

  async login(apiUrl, username, password) {
    const token = btoa(`${username}:${password}`);
    const apiClient = new ApiClient(apiUrl, username, password);

    try {
      // Test credentials
      await apiClient.listDirectory('');
      localStorage.setItem(AUTH_TOKEN_KEY, token);
      localStorage.setItem(API_URL_KEY, apiUrl);
      
      authState.token = token;
      authState.isAuthenticated = true;
      store.dispatch('fileManager/loadCompleteDirectoryTree');
      authState.apiClient = apiClient;
      authState.apiUrl = apiUrl;
      authState.username = username;
      authState.isLoginDialogVisible = false; // Hide dialog on success

      return { success: true };
    } catch (error) {
      console.error('Login failed:', error);
      return { success: false, message: error.message || 'An unknown error occurred during login.' };
    }
  },

  logout() {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    authState.token = null;
    authState.isAuthenticated = false;
    authState.apiClient = null;
    //authState.username = null;
    authState.isLoginDialogVisible = true; // Show login dialog after logout
  },

  showLoginDialog() {
    authState.isLoginDialogVisible = true;
  },

  hideLoginDialog() {
    authState.isLoginDialogVisible = false;
  },
};
