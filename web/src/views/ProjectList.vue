<template>
<div>
  <b-row class="mb-3">
    <b-col>
      <b-breadcrumb :items="breadcrumbs"></b-breadcrumb>
    </b-col>
  </b-row>
  <b-row>
    <b-col cols="12" md="4" lg="3" xl="2">
      <b-card title="Actions">
        <div><router-link :to="{name: 'new-project'}">New project</router-link></div>
      </b-card>
    </b-col>
    <b-col>
      <b-card title="Projects">
        <div v-if="loading" class="my-5 text-center">
          <font-awesome-icon icon="spinner"></font-awesome-icon>
        </div>
        <div v-else>
          <b-table :items="projects" :fields="fields">
            <template slot="name" slot-scope="data">
              <router-link :to="{ name: 'project-dashboard', params: { id: data.item.id }}">{{ data.value }}</router-link>
            </template>
          </b-table>
          <b-btn variant="success" size="sm" :to="{ name: 'new-project' }">New project</b-btn>
        </div>

      </b-card>
    </b-col>
  </b-row>
</div>

</template>

<script>
export default {
  name: 'ProjectList',
  data () {
    return {
      projects: [],
      loading: false,
      fields: [ 'id', 'name', 'location', 'pm' ],
      breadcrumbs: [
        {
          text: 'Projects',
          to: { name: 'projects' }
        }
      ]
    }
  },
  created () {
    this.loading = true
    this.$http.get('api/v1/projects/').then((r) => {
      this.projects = r.data
    }).catch((e) => {
      console.log(e)
    }).finally(() => {
      this.loading = false
    })
  }
}
</script>

<style>

</style>
