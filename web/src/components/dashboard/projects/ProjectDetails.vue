<template>
  <b-card :title="project.name">
    <div>Location: {{ project.location }}</div>
    <div>PM: {{ project.pm }}</div>
  </b-card>
</template>

<script>
export default {
  name: 'ProjectDetails',
  data () {
    return {
      project: {}
    }
  },
  created () {
    this.$http.get(`api/v1/projects/${this.$route.params.id}`).then((response) => {
      this.project = response.data
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style>

</style>
