<template>
  <header class="header-wrapper">
    <h1 class="header-area">
      <a href="/" class="header-link" target="_self" @click.prevent="$emit('logo-click')">
        <template v-if="isMobile">
          <template v-if="displayState == 1">
            <icon :name="'menu'" class="header-icon menu-icon" />
          </template>
          <template v-else>
            <icon :name="'chevron-left'" class="header-icon menu-icon" />
          </template>
          <span class="header-title-mobile">{{ displayTitle }}</span>
        </template>
        <template v-else>
          <img
            class="mark-markdown"
            src="@assets/images/markdown.png"
            :alt="$t('header.markdownEditor')"
          />
          <strong class="header-text">{{ displayTitle }}</strong>
        </template>
      </a>
      <nav class="button-group" v-if="!isMobile">
        <span class="hint--bottom" @click="onAboutClick" :aria-label="$t('header.about')">
          <icon class="header-icon" name="marker" />
        </span>
        <span class="hint--bottom" @click="onImportClick" :aria-label="$t('header.importFile')">
          <icon class="header-icon" name="upload" />
        </span>
        <el-dropdown trigger="click" @command="handleCommand">
          <span class="hint--bottom el-dropdown-link" :aria-label="$t('header.settings')">
            <icon class="header-icon" name="setting" />
          </span>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item disabled>
              <icon class="dropdown-icon" name="set-style" />
              {{
                $t('header.customStyle')
              }}
            </el-dropdown-item>
            <el-dropdown-item command="/export/ppt" divided>
              <icon class="dropdown-icon" name="preview" />
              {{ exportTextMap['/export/ppt'] }}
              
            </el-dropdown-item>
            <el-dropdown-item command="/export/png" divided>
              <icon class="dropdown-icon" name="download" />
             {{
                exportTextMap['/export/png']
              }}
            </el-dropdown-item>
            <el-dropdown-item command="/export/pdf">
              <icon class="dropdown-icon" name="download" />
              {{ exportTextMap['/export/pdf'] }}              
            </el-dropdown-item>
            <el-dropdown-item command="/export/html" disabled divided>
              <icon class="dropdown-icon" name="download" />{{
                $t('header.exportHtml')
              }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
        <el-dropdown trigger="click" @command="handleUserCommand">
          <span class="hint--bottom el-dropdown-link" :aria-label="$t('header.user')">
            <icon class="header-icon" name="user" />
          </span>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="login">{{ $t('header.login') }}</el-dropdown-item>
            <el-dropdown-item v-if="isAuthenticated" command="logout">{{
              $t('header.logout')
            }}</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
        <span
          class="hint--bottom full-screen"
          @click="onFullScreenClick"
          :aria-label="$t('header.fullScreen')"
        >
          <icon class="header-icon" name="full-screen" />
        </span>
      </nav>
      <el-dropdown trigger="click" @command="handleMobileCommand" v-else>
        <span class="el-dropdown-link">
          <icon :name="'menu'" class="header-icon" />
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="about">
            <icon class="dropdown-icon" name="marker" />{{ $t('header.about') }}
          </el-dropdown-item>
          <el-dropdown-item command="import" divided>
            <icon class="dropdown-icon" name="upload" />{{ $t('header.importFile') }}
          </el-dropdown-item>
          <el-dropdown-item command="/export/ppt">
            <icon class="dropdown-icon" name="preview" />{{ exportTextMap['/export/ppt'] }}
          </el-dropdown-item>
          <el-dropdown-item command="/export/png">
            <icon class="dropdown-icon" name="download" />{{ exportTextMap['/export/png'] }}
          </el-dropdown-item>
          <el-dropdown-item command="/export/pdf" divided>
            <icon class="dropdown-icon" name="download" />{{ exportTextMap['/export/pdf'] }}
          </el-dropdown-item>
          <el-dropdown-item command="login" divided>
            <icon class="dropdown-icon" name="user" />{{ $t('header.login') }}
          </el-dropdown-item>
          <el-dropdown-item command="logout" v-if="isAuthenticated">
            <icon class="dropdown-icon" name="user" />{{ $t('header.logout') }}
          </el-dropdown-item>
          <el-dropdown-item command="fullScreen" divided>
            <icon class="dropdown-icon" name="full-screen" />{{ $t('header.fullScreen') }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </h1>
  </header>
</template>

<script>
import 'hint.css'
import { exportTextMap } from '@config/constant'
import { authState, authService } from '@/store/auth';
import { mapActions } from 'vuex';
import { appTitle } from '../../config/constant'
import { mapState } from 'vuex'
export default {
  name: 'HeaderNav',
  props: {
    displayTitle: {
      type: String,
      default: appTitle,
    },
  },
  data() {
    return {
      isMobile: window.innerWidth <= 768,
      exportTextMap,
    }
  },

  computed: {
    ...mapState('ui', ['displayState', 'selectedFolder', 'selectedFile']),
    isAuthenticated() {
      return authState.isAuthenticated
    },
  },


  methods: {
    handleUserCommand(command) {
      if (command === 'logout') {
        authService.logout()
      } else if (command === 'login') {
        authService.showLoginDialog()
      }
    },
    launchFullScreen() {
      const element = document.getElementById('vditor')
      if (element.requestFullscreen) {
        element.requestFullscreen()
      } else if (element.msRequestFullscreen) {
        element.msRequestFullscreen()
      } else if (element.mozRequestFullScreen) {
        element.mozRequestFullScreen()
      } else if (element.webkitRequestFullscreen) {
        element.webkitRequestFullscreen(Element.ALLOW_KEYBOARD_INPUT)
      }
    },
    cancelFullScreen() {
      if (document.exitFullscreen) {
        document.exitFullscreen()
      } else if (document.msExitFullscreen) {
        document.msExitFullscreen()
      } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen()
      } else if (document.webkitExitFullscreen) {
        document.webkitExitFullscreen()
      }
    },
    ...mapActions('ui', ['showAbout']),
    onAboutClick() {
      this.showAbout();
    },
    onFullScreenClick() {
      const isFullScreen =
        document.fullscreenElement ||
        document.mozFullScreenElement ||
        document.msFullscreenElement ||
        document.webkitFullscreenElement
      isFullScreen ? this.cancelFullScreen() : this.launchFullScreen()
    },
    ...mapActions('exportManager', ['startExport']),
    handleCommand(command) {
      this.startExport(command.replace('/export/', ''));
    },
    handleMobileCommand(command) {
      switch (command) {
        case 'about':
          this.onAboutClick();
          break
        case 'import':
          this.onImportClick()
          break
        case 'fullScreen':
          this.onFullScreenClick()
          break
        case 'login':
        case 'logout':
          this.handleUserCommand(command)
          break
        default:
          // Handle export commands
          this.handleCommand(command)
          break
      }
    },
    onImportClick() {
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = '.md,.markdown,text/markdown'
      input.onchange = (e) => {
        const file = e.target.files[0]
        if (file) {
          const reader = new FileReader()
          reader.onload = (e) => {
            const content = e.target.result
            localStorage.setItem('vditorvditor', content)
            this.$root.$emit('reload-content')
          }
          reader.readAsText(file)
        }
      }
      input.click()
    },
  },
}
</script>

<style lang="less">
@import './../../assets/styles/style.less';

[class*='hint--']:after {
  border-radius: 3px;
}

.el-popper[x-placement^='bottom'] {
  margin-top: 10px;
}

.el-dropdown .el-dropdown-link {
  height: @header-height;
  .flex-box-center(column);
}

.hint--bottom {
  cursor: pointer;
  pointer-events: all;
}

.el-dropdown-menu {
  margin: 0;

  .dropdown-icon {
    fill: @deep-black;
    vertical-align: middle;
    margin-right: 10px;
  }

  .dropdown-text {
    vertical-align: middle;
  }
}

.header-wrapper {
  // 【必要改动】: position: fixed 是导致布局问题的根源，必须移除或注释掉。
  // position: fixed;
  // top: 0;
  
  // 【必要改动】: 增加 flex-shrink: 0，防止 Header 在 flex 布局中被压缩。
  flex-shrink: 0;
  
  width: 100%;
  height: @header-height;
  line-height: @header-height;
  z-index: @hint-css-zindex;
  background-color: #fff;
  box-shadow: 0 0 12px 2px rgba(0, 0, 0, 0.1);
  transition: border 0.5s cubic-bezier(0.455, 0.03, 0.515, 0.955),
    background 0.5s cubic-bezier(0.455, 0.03, 0.515, 0.955);

  .header-area {
    display: flex;
    height: @header-height;
    line-height: @header-height;
    padding: 0 10px; /* Add padding for spacing from edges */
    align-items: center; /* Vertically align items */
    justify-content: space-between; /* Distribute items to left and right */

    .header-link {
      display: flex; /* Use flex for better alignment of logo and text */
      height: @header-height;
      line-height: @header-height;
      align-items: center; /* Vertically align logo and text */

      .mark-markdown {
        width: @header-height; /* Adjust width to match new header height */
        vertical-align: middle;
      }

      .header-text {
        margin-left: 5px; /* Reduce margin-left */
        font-size: @font-medium;
        color: transparent;
        background-clip: text;
        background-image: linear-gradient(to right, #000000, #434343);
        vertical-align: middle;
        flex-grow: 1; /* Allow title to take available space */
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }

    .button-group {
      display: flex; /* Use flex for button group */
      align-items: center; /* Vertically align icons */

      .header-icon {
        margin: 0 5px; /* Reduce margin */
        fill: @deep-black;
      }

      .full-screen {
        margin-right: -5px; /* Adjust margin */
      }
    }
  }
}

@media (max-width: 960px) {
  .header-wrapper {
    .header-area {
      padding: 0 5px; /* Adjust padding for smaller screens */

      .header-link {
        .menu-icon {
          width: 24px;
          height: 24px;
          margin-right: 8px;
        }

        .header-title-mobile {
          font-size: 14px; /* Slightly smaller font size for mobile */
          color: @black;
          flex-grow: 1; /* Allow title to take available space */
          overflow: hidden;
          white-space: nowrap;
          text-overflow: ellipsis;
        }
      }
    }
  }
}
</style>