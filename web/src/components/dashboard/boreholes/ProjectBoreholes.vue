<template>
  <b-card title="Boreholes" :sub-title="project.name">
    <b-table :items="boreholes" show-empty></b-table>
    <b-btn variant="info" size="sm" :to="{ name: 'new-borehole' }">New borehole</b-btn>

  </b-card>
</template>

<script>
export default {
  name: 'Boreholes',
  props: ['project'],
  data () {
    return {
      boreholes: []
    }
  },
  created () {
    this.$http.get(`api/v1/field/boreholes?project=${this.$route.params.id}`).then((response) => {
      this.boreholes = response.data
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style>

</style>
