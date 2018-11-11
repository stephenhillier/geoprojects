<template>
            <b-card title="Actions">
              <b-btn variant="outline-danger" size="sm" @click="deleteProject">
                <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
                Delete project
              </b-btn>
              <b-alert variant="danger" dismissible v-if="deleteError" @dismissed="deleteError=false">
                Unable to delete project.
              </b-alert>
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
      }).catch(() => {
        this.deleteError = true
      })
    }
  }
}
</script>

<style>

</style>
