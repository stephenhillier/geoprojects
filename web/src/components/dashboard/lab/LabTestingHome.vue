<template>
  <div>
    <div class="columns">
      <div class="column">
        <h1 class="title">Lab Testing</h1>
      </div>
    </div>
    <div class="columns">
      <div class="column">
        <test-grid :testRowData="labTestingRowData" :sampleOptions="sampleRowData" @test-update="fetchLabTests" />
      </div>
    </div>
  </div>
</template>

<script>
import TestGrid from '@/components/dashboard/boreholes/grids/TestGrid.vue'

export default {
  name: 'LabTesting',
  props: ['project'],
  components: {
    TestGrid
  },
  data () {
    return {
      labTestingRowData: [],
      sampleRowData: []
    }
  },
  methods: {
    fetchLabTests () {
      this.$http.get(`projects/${this.$route.params.id}/lab/tests`).then((response) => {
        this.labTestingRowData = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving lab tests.')
      })
    },
    fetchSamples () {
      this.$http.get(`projects/${this.$route.params.id}/samples`).then((response) => {
        this.sampleRowData = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving samples.')
      })
    }
  },
  created () {
    this.fetchLabTests()
    this.fetchSamples()
  }
}
</script>

<style>

</style>
