<template>
  <el-dialog
    :visible.sync="dialogVisible"
    fullscreen
    :show-close="false"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    custom-class="login-dialog-fullscreen"
  >
    <div class="login-container">
      <el-card class="login-card">
        <div slot="header" class="clearfix">
          <h2 class="login-title">{{ $t('login.title') }}</h2>
        </div>
        <el-form
          ref="loginForm"
          :model="loginForm"
          :rules="loginRules"
          label-position="left"
          label-width="80px"
          @submit.native.prevent="handleLogin"
        >
          <el-form-item :label="$t('login.apiUrl')" prop="apiUrl">
            <el-input v-model="loginForm.apiUrl" :placeholder="$t('login.apiUrlPlaceholder')"></el-input>
          </el-form-item>
          <el-form-item :label="$t('login.username')" prop="username">
            <el-input v-model="loginForm.username" :placeholder="$t('login.usernamePlaceholder')"></el-input>
          </el-form-item>
          <el-form-item :label="$t('login.password')" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              :placeholder="$t('login.passwordPlaceholder')"
              show-password
              @keyup.enter.native="handleLogin"
            ></el-input>
          </el-form-item>
          <el-form-item class="button-group-form-item">
            <el-button @click="handleCancel" class="cancel-button">{{ $t('login.cancelButton') }}</el-button>
            <el-button type="primary" @click="handleLogin" :loading="loading" class="login-button">{{ $t('login.loginButton') }}</el-button>
          </el-form-item>
          <el-alert v-if="error" :title="error" type="error" show-icon :closable="false"></el-alert>
        </el-form>
      </el-card>
    </div>
  </el-dialog>
</template>

<script>
import { authState, authService } from '@/store/auth';

export default {
  name: 'LoginDialog',
  data() {
    return {
      loginForm: {
        apiUrl: authState.apiUrl || '',
        username: authState.username || '',
        password: '',
      },
      loginRules: {
        apiUrl: [{ required: true, message: this.$t('login.validation.apiUrlRequired'), trigger: 'blur' }],
        username: [{ required: true, message: this.$t('login.validation.usernameRequired'), trigger: 'blur' }],
        password: [{ required: true, message: this.$t('login.validation.passwordRequired'), trigger: 'blur' }],
      },
      loading: false,
      error: null,
      wasInitiallyAuthenticated: authState.isAuthenticated, // Track if user was logged in when dialog opened
    };
  },
  computed: {
    dialogVisible: {
      get() {
        return authState.isLoginDialogVisible;
      },
      set(value) {
        authState.isLoginDialogVisible = value;
      },
    },
  },
  methods: {
    handleLogin() {
      this.$refs.loginForm.validate(async (valid) => {
        if (valid) {
          this.loading = true;
          this.error = null;
          const { apiUrl, username, password } = this.loginForm;
          const result = await authService.login(apiUrl, username, password);
          this.loading = false;
          if (result.success) {
            //this.$store.dispatch('fileManager/loadCompleteDirectoryTree');
          } else {
            this.error = result.message || this.$t('login.loginFailed');
          }
        }
      });
    },
    handleCancel() {
      authService.hideLoginDialog();
    },
  },
  watch: {
    'authState.isLoginDialogVisible'(newValue) {
      if (newValue) {
        // When dialog becomes visible, reset form and update initial auth state
        this.wasInitiallyAuthenticated = authState.isAuthenticated;
        this.loginForm.apiUrl = authState.apiUrl || '';
        this.loginForm.username = authState.username || '';
        this.loginForm.password = '';
        this.error = null;
      }
    },
  },
};
</script>

<style>
.login-dialog-fullscreen .el-dialog__header {
  display: none;
}
.login-dialog-fullscreen .el-dialog__body {
  padding: 0;
  background-color: #f5f7fa;
}
</style>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  border-radius: 12px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.login-title {
  text-align: center;
  font-size: 24px;
  color: #333;
  margin: 0;
}

.el-form-item {
  margin-bottom: 15px; /* Reduce vertical spacing between form items */
}

.el-form-item__label {
  line-height: 24px; /* Adjust label line height for compactness */
  padding-bottom: 5px; /* Reduce padding below label */
}

.button-group-form-item {
  margin-top: 20px; /* Add space above buttons */
}

.button-group-form-item ::v-deep .el-form-item__content {
  margin-left: 0 !important; /* Override default margin from label-width */
  display: flex; /* Use flexbox for horizontal layout */
  justify-content: space-between; /* Distribute space between buttons */
  gap: 10px; /* Add gap between buttons */
}

.login-button {
  flex-grow: 1; /* Allow button to grow and fill space */
  margin-left: 0 !important; /* Reset margin from previous attempts */
  margin-right: 0 !important; /* Reset margin from previous attempts */
}

.cancel-button {
  flex-grow: 1; /* Allow button to grow and fill space */
  margin-left: 0 !important; /* Reset margin from previous attempts */
  margin-right: 0 !important; /* Reset margin from previous attempts */
  margin-top: 0; /* Reset margin-top */
}
</style>
