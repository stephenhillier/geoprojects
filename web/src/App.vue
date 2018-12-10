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
$primary: #34495e;
$secondary: #01579b;
$dark: #002f6c;
$white: #fff;
$navbar-dark-color:                 rgba($white, .8);
$navbar-dark-hover-color:           rgba($white, .9);
@import '~bootswatch/dist/cosmo/_variables.scss';
@import "~bootstrap/scss/bootstrap.scss";
@import '~bootswatch/dist/cosmo/_bootswatch.scss';

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

</style>
