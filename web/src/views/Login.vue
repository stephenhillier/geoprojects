<template>
  <b-card title="Login to Earthworks" class="login-card mx-auto mt-5">
  </b-card>
</template>

<script>
export default {
  name: 'Auth0Login',
  data () {
    return {}
  },
  created () {
    this.$auth.login()
  }
}
</script>

<style>
.login-card {
  max-width: 20rem;
}
</style>
