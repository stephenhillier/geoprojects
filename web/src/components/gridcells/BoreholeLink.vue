<template>
  <router-link :to="`/projects/${$route.params.id}/boreholes/${params.data.id}`">{{ params.value }}</router-link>
</template>

<script>
import Vue from 'vue'
export default Vue.extend({

})
</script>

<style>

</style>
