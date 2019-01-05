<template>
  <div id="app">
    <geo-header class="mb-3"/>
    <b-container fluid id="main-content">
      <router-view/>
    </b-container>
  </div>
</template>

<script>
import Header from '@/components/Header.vue'
export default {
  name: 'App',
  components: {
    geoHeader: Header
  }
}
</script>
<style lang="scss">
$primary: #002366;
$secondary: #002366;
$info: #01579b;
$dark: #001838;
$white: #fff;
$navbar-dark-color:                 rgba($white, .8);
$navbar-dark-hover-color:           rgba($white, .9);
@import '~proboot/scss/proboot';

body {
  background-color: #f0f0f0;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, "Open Sans", "Helvetica Neue", sans-serif;
  color: #333;
}

#main-content {
  font-size: 14px!important;
}

.earthworks-brand {
  font-weight: 600;
  font-size: 1.8rem;
}
.menu-heading a {
  color: rgb(61, 61, 61)!important;
}
.menu-item a {
  color: rgb(61, 61, 61);
}
.active-menu {
  color:#1f548a!important;
  background-color: #e2e6ff;
  font-weight: 600;
}
.table-heading {
  color: #333!important;
  text-transform: none!important;
}
</style>
