<template>
  <b-card title="Projects">
    <b-table :items="projects" :fields="fields">
      <template slot="name" slot-scope="data">
        <router-link :to="{ name: 'project-details', params: { id: data.item.id }}">{{ data.value }}</router-link>
      </template>
    </b-table>
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
