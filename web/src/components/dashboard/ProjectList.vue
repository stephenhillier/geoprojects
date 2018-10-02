<template>
  <b-card title="Projects">
    <b-table :items="projects" :fields="fields"></b-table>
  </b-card>
</template>

<script>
export default {
  name: 'ProjectList',
  data () {
    return {
      projects: [],
      fields: [ 'id', 'name', 'location', 'pm' ]
    }
  },
  created () {
    this.$http.get('api/v1/projects/').then((r) => {
      this.projects = r.data
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style>

</style>
