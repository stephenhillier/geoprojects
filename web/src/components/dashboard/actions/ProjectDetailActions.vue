<template>
            <b-card title="Actions">
              <b-btn variant="outline-danger" size="sm" @click="deleteProject">
                <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
                Delete project
              </b-btn>
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
