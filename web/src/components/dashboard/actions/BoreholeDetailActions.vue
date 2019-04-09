<template>
  <div class="box">
    <div>
      <a :href="`${fileHost}/logs/${$route.params.id}/boreholes/${$route.params.bh}/log.pdf`" target="_blank">
        <font-awesome-icon :icon="['far', 'file-alt']" class="text-muted"></font-awesome-icon>
        Publish to PDF
      </a>
    </div>
    <div class="mt-1">
      <a href="#" class="is-warning" @click.prevent="handleDelete">
        <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
        Delete borehole
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BoreholeDetailActions',
  data () {
    return {
      deleteError: false,
      fileHost: process.env.VUE_APP_FILE_URL || 'http://localhost:8081'
    }
  },
  methods: {
    deleteBorehole () {
      this.$http.delete(`boreholes/${this.$route.params.bh}`).then(() => {
        this.$router.push({ name: 'project-boreholes', params: { id: this.$route.params.id } })
        this.$noty.success('Borehole deleted')
      }).catch((e) => {
        this.$noty.error(`An error occured while deleting borehole (${e.response.status})`)
      })
    },
    handleGetBoreholePdf () {
      this.$file.get(`/logs/boreholes/${this.$route.params.id}.pdf`)
        .then(response => {
          console.log(response)

          let blob = new Blob([response.data], { type: 'application/pdf' })
          let url = window.URL.createObjectURL(blob)

          window.open(url)
        }).catch((e) => {
          console.error(e)
        })
    },
    handleDelete () {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this borehole?',
        onConfirm: this.deleteBorehole
      })
    }
  }
}
</script>

<style>

</style>
