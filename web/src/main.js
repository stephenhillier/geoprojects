import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import BootstrapVue from 'bootstrap-vue'
import './registerServiceWorker'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import axios from 'axios'
import qs from 'qs'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faSpinner } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faSpinner)

Vue.component('font-awesome-icon', FontAwesomeIcon)

Vue.use(BootstrapVue)

const axiosClient = axios.create({
  baseURL: 'http://localhost:8000/'
})

axiosClient.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'

axiosClient.interceptors.request.use(function (config) {
  console.log(config)
  if (config.data) {
    config.data = qs.stringify(config.data)
  }
  return config
}, function (error) {
  console.log(error)
  return Promise.reject(error)
})

axiosClient.interceptors.response.use(function (response) {
  console.log(response)
  return response
}, function (error) {
  console.log(error)
  return Promise.reject(error)
})

Vue.prototype.$http = axiosClient

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
