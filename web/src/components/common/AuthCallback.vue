<template>
  <div class="callback">Logging in...</div>
</template>

<script>

export default {
  name: 'callback',
  created () {
    this.$auth.handleAuthentication()
  }
}
</script>
