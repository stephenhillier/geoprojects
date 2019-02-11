<template>
            <b-card title="Actions">
              <div>
                <b-btn variant="outline-primary" size="sm" :to="{ name: 'new-borehole' }"><font-awesome-icon :icon="['far', 'plus-square']"></font-awesome-icon> New borehole</b-btn>
              </div>
              <div class="mt-3">
                <b-btn v-b-modal.deleteProjectModal variant="outline-danger" size="sm">
                  <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
                  Delete project
                </b-btn>
              </div>

                <!-- Delete project confirmation -->
              <b-modal id="deleteProjectModal" centered @ok="deleteProject" title="Confirm delete">
                Are you sure you want to delete this project?
              </b-modal>
            </b-card>
</template>

<script>
export default {
  name: 'ProjectDetailActions',
  data () {
    return {
      deleteError: false
    }
  },
  methods: {
    deleteProject () {
      this.$http.delete(`projects/${this.$route.params.id}`).then(() => {
        this.$router.push({ name: 'projects' })
        this.$noty.success('Project deleted')
      }).catch((e) => {
        this.$noty.error(`An error occurred while deleting project (${e.response.status})`)
      })
    }
  }
}
</script>

<style>

</style>
