<template>
  <div>
    <div class="columns">
      <div class="column">
        <h1 class="title">Lab Testing</h1>
      </div>
    </div>
    <div class="columns">
      <div class="column">
        <sample-grid :samples="sampleRowData" @sample-update="fetchSamples"/>
      </div>
    </div>
  </div>
</template>

<script>
import SampleGrid from '@/components/dashboard/boreholes/grids/SampleGrid.vue'

export default {
  name: 'SamplesHome',
  props: ['project'],
  components: {
    SampleGrid
  },
  data () {
    return {
      sampleRowData: []
    }
  },
  methods: {
    fetchSamples () {
      this.$http.get(`projects/${this.$route.params.id}/samples`).then((response) => {
        this.sampleRowData = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving samples.')
      })
    }
  },
  created () {
    this.fetchSamples()
  }
}
</script>

<style>

</style>
