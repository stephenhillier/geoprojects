<template>
  <div id="app">
    <geo-header class="geoprojects-header"/>
    <div id="main-content">
      <router-view/>
    </div>
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
$primary: #1f548a;
$secondary: #002366;
$dark: #001838;
$white: #fff;
$navbar-dark-color:                 rgba($white, .8);
$navbar-dark-hover-color:           rgba($white, .9);

// Import Bulma's core
@import "~bulma/sass/utilities/_all";

// Set your colors
$primary-invert: findColorInvert($primary);
$twitter: #4099FF;
$twitter-invert: findColorInvert($twitter);

// Setup $colors to use as bulma classes (e.g. 'is-twitter')
$colors: (
    "white": ($white, $black),
    "black": ($black, $white),
    "light": ($light, $light-invert),
    "dark": ($dark, $dark-invert),
    "primary": ($primary, $primary-invert),
    "info": ($info, $info-invert),
    "success": ($success, $success-invert),
    "warning": ($warning, $warning-invert),
    "danger": ($danger, $danger-invert),
    "twitter": ($twitter, $twitter-invert)
);

// Links
$link: $primary;
$link-invert: $primary-invert;
$link-focus-border: $primary;

$input-border-color: $grey-light;
$input-icon-color: $grey-light;

// Import Bulma and Buefy styles
@import "~bulma";
@import "~buefy/src/scss/buefy";

body {
  background-color: #f0f0f0;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, "Open Sans", "Helvetica Neue", sans-serif;
  color: #333;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

#main-content {
  font-size: 14px!important;
  flex: 1;
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
  font-weight: 600;
}
.table-heading {
  color: #333!important;
  font-weight: 600!important;
  text-transform: none!important;
}
.table-heading-horizontal {
  color: #333!important;
  font-weight: 600!important;
  text-transform: none!important;
  margin-right: 15rem;
}
.geoprojects-header {
  margin-bottom: 15px;
  flex: none;
}
.is-fullheight {
  height: 100%;
}
.is-map-container {
  height: 400px;
}
.is-action-button {
  margin-bottom: 10px;
}
.is-button-right {
  margin-left: 10px;
}
.is-full-map {
  height: 95%;
}
</style>
