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
$primary: #01579b;
$secondary: #263238;
$dark: #002f6c;
$white: #fff;
$navbar-dark-color:                 rgba($white, .8);
$navbar-dark-hover-color:           rgba($white, .9);

@import "~bootstrap/scss/bootstrap.scss";

#app {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, "Open Sans", "Helvetica Neue", sans-serif;
  color: #333;
}

#main-content {
  font-size: 14px;
}

.earthworks-brand {
  font-weight: 600;
  font-size: 1.8rem;
}

</style>
