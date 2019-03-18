<template>
  <div>
    <h1 class="title">Boreholes</h1>
    <div style="height: 400px">
      <multi-marker-map :locations="boreholes"></multi-marker-map>
    </div>
    <borehole-grid class="mt-1" :project="$route.params.id" :boreholes="boreholes"></borehole-grid>
    <router-link tag="button" class="button" :to="{ name: 'new-borehole' }">New borehole</router-link>

  </div>
</template>

<script>
import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'

import BoreholeGrid from '@/components/dashboard/boreholes/grids/BoreholeGrid.vue'

export default {
  name: 'Boreholes',
  props: ['project'],
  components: {
    MultiMarkerMap,
    BoreholeGrid
  },
  data () {
    return {
      boreholes: [],
      numberOfRecords: 0
    }
  },
  methods: {
    fetchBoreholes () {
      this.$http.get(`boreholes?project=${this.$route.params.id}`).then((response) => {
        this.numberOfRecords = response.data.results.length
        this.boreholes = response.data.results
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving boreholes.')
      })
    }
  },
  created () {
    this.fetchBoreholes()
  }
}
</script>

<style>

</style>
