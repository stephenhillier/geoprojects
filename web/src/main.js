import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import BootstrapVue from 'bootstrap-vue'
import './registerServiceWorker'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import axios from 'axios'
// import qs from 'qs'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faSpinner, faInfoCircle, faPrint, faLink } from '@fortawesome/free-solid-svg-icons'
import { faTrashAlt, faPlusSquare, faTimesCircle, faFile, faFileAlt } from '@fortawesome/free-regular-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import FormInput from '@/components/common/FormInput.vue'
import '@/components/common/readableFilter.js'

import 'leaflet/dist/leaflet.css'
import './ag-grid.scss'

library.add(faSpinner, faInfoCircle, faTrashAlt, faPlusSquare, faTimesCircle, faPrint, faLink, faFile, faFileAlt)

Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('form-input', FormInput)

Vue.use(BootstrapVue)

const axiosClient = axios.create({
  baseURL: process.env.VUE_APP_BASE_URL || 'http://localhost:8000/api/v1/'
})

// axiosClient.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded'

axiosClient.interceptors.request.use(function (config) {
  // if (config.data) {
  //   config.data = qs.stringify(config.data)
  // }
  return config
}, function (error) {
  console.log(error)
  return Promise.reject(error)
})

axiosClient.interceptors.response.use(function (response) {
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
