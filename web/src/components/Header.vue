<template>
  <b-navbar toggleable="sm" variant="primary" type="dark">
    <b-navbar-brand href="/" class="earthworks-brand">Earthworks</b-navbar-brand>
    <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
    <b-collapse is-nav id="nav_collapse">
      <b-navbar-nav>
        <b-nav-item :to="{ name: 'projects'}">Projects</b-nav-item>
      </b-navbar-nav>
      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown right v-if="username">
          <!-- Using button-content slot -->
          <template slot="button-content">
            {{ username }}
          </template>
          <b-dropdown-item href="#" @click="logout">Logout</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>

  </b-navbar>
</template>

<script>
export default {
  name: 'Header',
  data () {
    return {
      username: null,
      authenticated: false
    }
  },
  watch: {
    $auth () {
      this.username = this.$auth.name
    }
  },
  methods: {
    logout () {
      this.$auth.logout()
    }
  },
  created () {
    this.$auth.authNotifier.on('authChange', authState => {
      this.authenticated = authState.authenticated
      this.username = authState.authenticated ? this.$auth.name : null
    })

    this.$auth.renewSession()
  }
}
</script>

<style>
.earthworks-header-nav.a {
  color: #cfd8dc!important;
}
</style>
