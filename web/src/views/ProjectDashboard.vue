<template>
  <div id="dashboard">
    <div class="columns is-fullheight" v-if="!!project && project != {}">
      <div class="column is-narrow is-fullheight">
        <project-menu :project="project" :project-id="project.id"></project-menu>

      </div>
      <div class="column is-fullheight">
        <div class="box is-fullheight">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li v-for="(breadcrumb, i) in breadcrumbs" :key="`breadcrumb${i}`" :class="`${i === breadcrumbs.length - 1 ? 'is-active':''}`"><router-link :to="breadcrumb.to">{{breadcrumb.text}}</router-link></li>
            </ul>
          </nav>
          <router-view :project="project" @update-project="fetchProjectData"></router-view>
        </div>
      </div>
      <div class="column is-narrow">
        <router-view :project="project" name="actions"></router-view>
      </div>
    </div>
  </div>
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
      this.$http.get(`projects/${this.$route.params.id}`).then((response) => {
        this.project = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving project details.')
      })
    }
  },
  created () {
    this.fetchProjectData()
  }
}
</script>

<style>
#dashboard {
  height: 100%;
}

</style>
