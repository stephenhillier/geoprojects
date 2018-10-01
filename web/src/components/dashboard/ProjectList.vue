<template>
  <b-card title="Projects">
    <b-table :items="projects" :fields="fields"></b-table>
  </b-card>
</template>

<script>
import axios from 'axios'
export default {
  name: 'ProjectList',
  data () {
    return {
      projects: [],
      fields: [ 'id', 'name', 'location', 'pm' ]
    }
  },
  created () {
    axios.get('http://localhost:8000/api/v1/projects/').then((r) => {
      this.projects = r.data
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style>

</style>
