<template>
  <b-card>
    <b-row>
      <b-col>
        <h4 class="card-title">Lab Testing</h4>
        <h6 class="text-muted">{{project.name}}</h6>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <lab-test-grid :labTestRowData="labTestingRowData" :sampleOptions="sampleRowData" @labtest-update="fetchLabTests" />
      </b-col>
    </b-row>
  </b-card>
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
        console.error(e)
      })
    },
    fetchSamples () {
      this.$http.get(`projects/${this.$route.params.id}/samples`).then((response) => {
        this.sampleRowData = response.data
      }).catch((e) => {
        console.error(e)
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
