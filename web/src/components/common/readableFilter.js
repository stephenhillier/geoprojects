import Vue from 'vue'

Vue.filter('readable', value => {
  return value.replace(/_/g, ' ')
})
