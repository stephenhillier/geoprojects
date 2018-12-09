<template>
  <b-card title="Actions">
    <b-row class="mt-2">
      <b-col>
        <b-btn variant="link" size="sm" :href="`/logs/boreholes/${this.$route.params.id}`">
          <font-awesome-icon :icon="['far', 'file-alt']" class="text-muted"></font-awesome-icon>
          View PDF
        </b-btn>
      </b-col>
    </b-row>
    <b-row class="mt-2">
      <b-col>
        <b-btn variant="link" size="sm">
          <font-awesome-icon :icon="['fas', 'link']" class="text-muted"></font-awesome-icon>
          Share link
        </b-btn>
      </b-col>
    </b-row>
    <b-row class="mt-2">
      <b-col>
        <b-btn variant="link" size="sm">
          <font-awesome-icon :icon="['fas', 'print']" class="text-muted"></font-awesome-icon>
          Print borehole log
        </b-btn>
      </b-col>
    </b-row>
    <b-row class="mt-5">
      <b-col>
        <b-btn variant="outline-danger" size="sm" @click="deleteBorehole">
          <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
          Delete borehole
        </b-btn>
      </b-col>
    </b-row>
    <b-row class="mt-2">
      <b-col>
        <b-alert variant="danger" dismissible v-if="deleteError" @dismissed="deleteError=false">
          Unable to delete borehole.
        </b-alert>
      </b-col>
    </b-row>

  </b-card>
</template>

<script>
export default {
  name: 'BoreholeDetailActions',
  data () {
    return {
      deleteError: false
    }
  },
  methods: {
    deleteBorehole () {
      this.$http.delete(`boreholes/${this.$route.params.bh}`).then(() => {
        this.$router.push({ name: 'project-boreholes', params: { id: this.$route.params.id } })
      }).catch(() => {
        this.deleteError = true
      })
    },
    handleGetBoreholePdf () {

    }
  }
}
</script>

<style>

</style>
