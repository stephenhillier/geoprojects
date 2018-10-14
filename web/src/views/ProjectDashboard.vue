<template>
  <b-container fluid id="dashboard" class="mb-3">
    <b-row>
      <b-col>
        <b-breadcrumb :items="breadcrumbs"></b-breadcrumb>
      </b-col>
    </b-row>
    <b-row v-if="!!project && project != {}">
      <b-col cols="12" lg="3" xl="2">
        <project-menu :project="project" :project-id="project.id"></project-menu>
      </b-col>
      <b-col>
        <router-view :project="project" @update-project="fetchProjectData"></router-view>
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
      project: {}
    }
  },
  computed: {
    breadcrumbs () {
      const crumbs = [
        {
          text: 'Projects',
          to: { name: 'projects' }
        },
        {
          text: this.project ? this.project.name : '',
          to: { name: 'project-dashboard', params: { id: this.$route.params.id } }
        }
      ]

      if (this.$route.meta && this.$route.meta.breadcrumbs) {
        this.$route.meta.breadcrumbs.forEach((crumb) => {
          if (crumb.to) {
            crumb.to.params = { id: this.$route.params.id }
          }
          crumbs.push(crumb)
        })
      }

      return crumbs
    }
  },
  methods: {
    fetchProjectData () {
      this.$http.get(`api/v1/projects/${this.$route.params.id}`).then((response) => {
        this.project = response.data
      }).catch((e) => {
        console.log(e)
      })
    }
  },
  created () {
    this.fetchProjectData()
  }
}
</script>

<style>

</style>
