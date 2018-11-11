<template>
<div>
  <b-card no-body class="mb-3">
    <b-list-group flush>
      <b-list-group-item class="menu-heading font-weight-bold">{{ project.name }}</b-list-group-item>
      <b-list-group-item exact :to="{ name: 'project-dashboard', params: { id: this.$route.params.id } }">Project Details</b-list-group-item>
    </b-list-group>
  </b-card>
  <b-card no-body class="mb-3">
    <b-list-group flush>
      <b-list-group-item class="menu-heading">Data:</b-list-group-item>
      <b-list-group-item class="d-flex justify-content-between align-items-center" :to="{ name: 'project-boreholes', params: { id: this.$route.params.id } }">Boreholes<b-badge variant="light" pill>{{ project.borehole_count }}</b-badge></b-list-group-item>
      <b-list-group-item>Instrumentation</b-list-group-item>
      <b-list-group-item>Soil Samples</b-list-group-item>
    </b-list-group>
  </b-card>

</div>

</template>

<script>
export default {
  name: 'DashboardMenu',
  props: ['project']
}
</script>

<style>

</style>
