<template>
  <b-card title="Project Details" :sub-title="project.name">
    <div>Location: {{ project.location }}</div>
    <div>PM: {{ project.pm }}</div>
  </b-card>
</template>

<script>
export default {
  name: 'ProjectDetails',
  props: ['project'],
  data () {
    return {}
  }
}
</script>

<style>

</style>
