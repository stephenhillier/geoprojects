<template>
  <b-card>
    <div><router-link :to='{name: "projects"}'>Select project</router-link></div>
    <div>Field programs</div>
    <div>Lab programs</div>
    <div><strong>Data:</strong></div>
    <div>Boreholes</div>
    <div>Instrumentation</div>
    <div>Soil samples</div>
  </b-card>
</template>

<script>
export default {
  name: 'DashboardMenu'
}
</script>

<style>

</style>
