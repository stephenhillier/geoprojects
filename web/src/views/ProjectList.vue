<template>
<b-row>
  <b-col cols="12" md="4" lg="3" xl="2">
    <b-card title="Actions">
      <div><router-link :to="{name: 'new-project'}">New project</router-link></div>
    </b-card>
  </b-col>
  <b-col>
    <b-card title="Projects">
      <div v-if="!loading">

      </div>
      <b-table :items="projects" :fields="fields">
        <template slot="name" slot-scope="data">
          <router-link :to="{ name: 'project-details', params: { id: data.item.id }}">{{ data.value }}</router-link>
        </template>
      </b-table>
      <b-btn variant="success" size="sm" :to="{ name: 'new-project' }">New project</b-btn>
    </b-card>
  </b-col>
</b-row>

</template>

<script>
export default {
  name: 'ProjectList',
  data () {
    return {
      projects: [],
      loading: false,
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
