<template>
  <div class="login-card">
    <font-awesome-icon :icon="['fas', 'spinner']" spin pulse class="text-muted" size="lg"></font-awesome-icon>
  </div>
</template>

<script>
export default {
  name: 'Auth0Login',
  data () {
    return {}
  },
  created () {
    const next = this.$route.query.redirect
    this.$auth.login(next)
  }
}
</script>

<style>
.login-card {
  max-width: 20rem;
  margin-left: auto;
  margin-right: auto;
  margin-top: 8rem;
}
</style>
