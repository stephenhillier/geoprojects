<template>
  <nav class="menu earthworks-project-menu">
    <p class="menu-label">
      Project: {{ project.name }}
    </p>
    <ul class="menu-list">
      <li><router-link :to="{ name: 'project-dashboard', params: { id: this.$route.params.id } }">Summary</router-link></li>
      <li><router-link :to="{ name: 'projects' }">Back to project list</router-link></li>
    </ul>
    <p class="menu-label">
      Field data
    </p>
    <ul class="menu-list">
      <li><router-link :to="{ name: 'project-boreholes', params: { id: this.$route.params.id } }">
        Boreholes</router-link></li>
    </ul>
    <p class="menu-label">
      Laboratory
    </p>
    <ul class="menu-list">
      <li><router-link :to="{ name: 'samples-home', params: { id: this.$route.params.id } }">
        Samples</router-link></li>
      <li><router-link :to="{ name: 'lab-home', params: { id: this.$route.params.id } }">
        Lab tests</router-link></li>
    </ul>
  </nav>
</template>

<script>
export default {
  name: 'DashboardMenu',
  props: ['project']
}
</script>

<style>
  .menu-heading a {
    color: rgb(61, 61, 61)!important;
  }
  .menu-item a {
    color: rgb(61, 61, 61);
  }
  .earthworks-project-menu {
    margin-top: 8px;
    margin-left: 15px;
    width: 300px;
  }
  ul.menu-list > li > .router-link-exact-active {
    border-left: 2px solid #41438d;
  }
</style>
