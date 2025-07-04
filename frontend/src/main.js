/** @format */

if (process.env.NODE_ENV === 'production') {
  console.log = () => {}
}

import Vue from 'vue'
import './global.js'
import App from './App.vue'
import store from './store' // Import the Vuex store
import './registerServiceWorker'
import ElementUI from 'element-ui' // Import Element UI
import 'element-ui/lib/theme-chalk/index.css' // Import Element UI CSS
import { authService } from '@/store/auth'
import i18n from './locales'

authService.initAuth()

Vue.use(ElementUI) // Use Element UI globally

const updateCanonical = (url) => {
  let link = document.querySelector("link[rel='canonical']")
  if (!link) {
    link = document.createElement('link')
    link.setAttribute('rel', 'canonical')
    document.head.appendChild(link)
  }
  link.setAttribute('href', url)
}


new Vue({
  store, // Add the store to the Vue instance
  i18n,
  render: (h) => h(App),
}).$mount('#app')
