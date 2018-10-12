<template>
  <div id="app">
    <geo-header class="mb-3"/>
    <b-container fluid>
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
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #333;
}

.earthworks-brand {
  font-weight: 600;
  font-size: 1.8rem;
}

</style>
