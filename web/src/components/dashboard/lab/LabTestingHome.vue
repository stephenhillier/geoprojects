<template>
  <div>
    <b-row>
      <b-col>
        <h1>Lab Testing</h1>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <lab-test-grid :labTestRowData="labTestingRowData" :sampleOptions="sampleRowData" @labtest-update="fetchLabTests" />
      </b-col>
    </b-row>
  </div>
</template>

<script>
import LabTestGrid from '@/components/dashboard/boreholes/grids/LabTestGrid.vue'

export default {
  name: 'LabTesting',
  props: ['project'],
  components: {
    LabTestGrid
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
