<template>
  <div>
    <h1 class="title">Instrumentation</h1>
    <div style="height: 400px">
      <multi-marker-map :locations="instruments"></multi-marker-map>
    </div>
    <instrumentation-table class="mt-1" :instruments="instruments"></instrumentation-table>
    <router-link tag="button" class="button is-primary" :to="{ name: 'instrumentation-new' }">New instrument</router-link>
  </div>
</template>

<script>
import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'
import InstrumentationTable from './InstrumentationTable.vue'

export default {
  name: 'InstrumentationHome',
  props: ['project'],
  components: {
    MultiMarkerMap,
    InstrumentationTable
  },
  data () {
    return {
      instruments: [],
      numberOfRecords: 0
    }
  },
  methods: {
    fetchInstruments () {
      this.$http.get(`projects/${this.$route.params.id}/instrumentation`).then((response) => {
        this.numberOfRecords = response.data.length
        this.instruments = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving instrumentation.')
      })
    }
  },
  created () {
    this.fetchInstruments()
  }
}
</script>

<style>

</style>
