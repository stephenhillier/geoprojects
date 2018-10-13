<template>
  <b-container fluid id="dashboard" class="mb-3">
    <b-row class="mb-3">
      <b-col>
        <b-breadcrumb :items="breadcrumbs"></b-breadcrumb>
      </b-col>
    </b-row>
    <b-row v-if="!!project">
      <b-col cols="12" lg="3" xl="2">
        <project-menu :project-name="project.name"></project-menu>
      </b-col>
      <b-col>
        <router-view></router-view>
      </b-col>
      <b-col cols="12" lg="3" xl="2" class="d-none d-xl-block">
        <div>
          <b-card title="Project stats" class="mb-3">
            <dl>
              <dt>Boreholes:</dt>
              <dd>{{ project.borehole_count }}</dd>
            </dl>
          </b-card>
        </div>
        <div>
          <b-card title="Latest activity" class="mb-3">

          </b-card>
        </div>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>

import Menu from '@/components/dashboard/Menu.vue'
export default {
  name: 'Dashboard',
  components: {
    projectMenu: Menu
  },
  data () {
    return {
      project: null
    }
  },
  computed: {
    breadcrumbs () {
      return [
        {
          text: 'Projects',
          to: { name: 'projects' }
        },
        {
          text: this.project ? this.project.name : '',
          to: { name: 'project-dashboard', params: { id: this.$route.params.id } }
        }
      ]
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
