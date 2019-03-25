<template>
 <div></div>
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
    },
    handleDelete () {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this project?',
        onConfirm: this.deleteProject()
      })
    }
  }
}
</script>

<style>

</style>
