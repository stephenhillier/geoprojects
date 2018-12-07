<template>
  <div class="login-card mx-auto mt-5">
    <font-awesome-icon :icon="['fas', 'spinner']" class="text-muted"></font-awesome-icon>
  </div>
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
