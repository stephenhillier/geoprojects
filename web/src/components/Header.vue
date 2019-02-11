<template>
  <b-navbar toggleable="sm" variant="primary" type="dark" class="border-bottom">
    <b-navbar-brand href="/" class="earthworks-brand">Earthworks</b-navbar-brand>
    <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
    <b-collapse is-nav id="nav_collapse">

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-dropdown text="Create" right class="mr-5 mt-1 d-none d-md-block" variant="success">
          <b-dropdown-item :to="{ name: 'new-project' }">Project</b-dropdown-item>
        </b-dropdown>

        <b-nav-item-dropdown right v-if="username">
          <!-- Using button-content slot -->
          <template slot="button-content">
            <span v-if="picture" class="avatar avatar-blue mr-2" :style="{ 'background-image': 'url(' + picture+ ')' }"></span> <span class="username-text">{{ username }}</span>
          </template>
          <b-dropdown-item href="#">Settings</b-dropdown-item>
          <b-dropdown-item href="#">Need help?</b-dropdown-item>
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
      picture: null,
      authenticated: false
    }
  },
  watch: {
    $auth () {
      this.username = this.$auth.name
      this.picture = this.$auth.picture
    }
  },
  methods: {
    logout () {
      this.$auth.logout()
    }
  },
  created () {
    this.username = this.$auth.authenticated ? this.$auth.name : null
    this.picture = this.$auth.authenticated ? this.$auth.picture : null

    this.$auth.authNotifier.on('authChange', authState => {
      this.authenticated = authState.authenticated
      this.username = authState.authenticated ? this.$auth.name : null
      this.picture = authState.authenticated ? this.$auth.picture : null
    })
  }
}
</script>

<style>
/* .earthworks-brand {
  color: #34495e!important;
}
.username-text {
  color: #01579b!important;
} */
</style>
