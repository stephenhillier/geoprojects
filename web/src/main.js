import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import BootstrapVue from 'bootstrap-vue'
import './registerServiceWorker'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import axios from 'axios'
import VueNoty from 'vuejs-noty'
import VueApexCharts from 'vue-apexcharts'

// Auth0 auth service
import AuthService from '@/components/common/AuthService.js'

// Icons
import { library } from '@fortawesome/fontawesome-svg-core'
import { faSpinner, faInfoCircle, faPrint, faLink } from '@fortawesome/free-solid-svg-icons'
import { faTrashAlt, faPlusSquare, faTimesCircle, faFile, faFileAlt } from '@fortawesome/free-regular-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// custom helpers (forms, displaying data etc)
import FormInput from '@/components/common/FormInput.vue'
import '@/components/common/readableFilter.js'

// 3rd party library css
import 'leaflet/dist/leaflet.css'
import './ag-grid.scss'
import 'vuejs-noty/dist/vuejs-noty.css'

Vue.use(VueNoty, {
  timeout: 2000,
  layout: 'bottomRight'
})

// font awesome icons
library.add(faSpinner, faInfoCircle, faTrashAlt, faPlusSquare, faTimesCircle, faPrint, faLink, faFile, faFileAlt)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('form-input', FormInput)

// apexcharts chart component
Vue.use(VueApexCharts)
Vue.component('apexchart', VueApexCharts)

Vue.use(BootstrapVue)

// create an axios client for accessing the API server
const axiosClient = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8000/api/v1/'
})

// file client for accessing PDF generator service
const fileClient = axios.create({
  baseURL: process.env.VUE_APP_BASE_URL || 'http://localhost:8081/'
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
Vue.prototype.$file = fileClient
Vue.prototype.$auth = AuthService

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
