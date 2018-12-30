<template>
  <b-card title="Actions">
    <b-row class="mt-2">
      <b-col>
        <b-btn variant="link" size="sm" :href="`${fileHost}/logs/boreholes/${$route.params.id}.pdf`" target="_blank">
          <font-awesome-icon :icon="['far', 'file-alt']" class="text-muted"></font-awesome-icon>
          Publish to PDF
        </b-btn>
      </b-col>
    </b-row>
    <b-row class="mt-5">
      <b-col>
        <b-btn variant="outline-danger" size="sm" v-b-modal.deleteBoreholeModal>
          <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
          Delete borehole
        </b-btn>
      </b-col>
    </b-row>
    <b-modal id="deleteBoreholeModal" centered @ok="deleteBorehole" title="Confirm delete">
      Are you sure you want to delete this borehole?
    </b-modal>
  </b-card>
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
    }
  }
}
</script>

<style>

</style>
