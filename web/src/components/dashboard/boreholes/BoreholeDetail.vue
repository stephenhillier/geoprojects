<template>
  <div>
    <div class="columns">
      <div class="column">
        <h1 class="title">{{ borehole.name }}</h1>
        <h2 class="subtitle">Summary</h2>
        <div v-if="borehole.location && borehole.location.length">Location: {{ latitude.toFixed(6) }}, {{ longitude.toFixed(6) }}</div>
        <div>
          Started: {{ borehole.start_date | moment('YYYY-MM-DD') }}
        </div>
        <div>Completed: {{ borehole.end_date | moment('YYYY-MM-DD') }}</div>
        <div>Logged by: {{ borehole.field_eng }}</div>
        <div>Logged soil strata: {{ borehole.strata_count }}</div>
        <div>Samples: 0</div>
        <div>Lab tests: 0</div>
      </div>
      <div class="column">
        <div class="is-480-map">
          <single-marker-map :latitude="latitude" :longitude="longitude"></single-marker-map>
        </div>
      </div>
    </div>

    <strata-grid :strataRowData="strataRowData" @strata-update="fetchStrata"/>

    <sample-grid :samples="sampleRowData" @sample-update="fetchSamples();fetchLabTests()"/>

    <test-grid :testRowData="labTestingRowData" :sampleOptions="sampleRowData" @test-update="fetchLabTests" />

  </div>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'
import StrataGrid from '@/components/dashboard/boreholes/grids/StrataGrid.vue'
import SampleGrid from '@/components/dashboard/boreholes/grids/SampleGrid.vue'
import TestGrid from '@/components/dashboard/boreholes/grids/TestGrid.vue'

import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'BoreholeDetails',
  components: {
    SingleMarkerMap,
    StrataGrid,
    TestGrid,
    AgGridVue,
    SampleGrid
  },
  props: {
    project: {
      type: Object,
      default: () => ({})
    }
  },
  data () {
    return {
      borehole: {
        location: []
      },
      strataRowData: [],

      labTestingRowData: [],
      sampleRowData: []
    }
  },
  computed: {
    latitude () {
      return this.borehole.location[0] || '49'
    },
    longitude () {
      return this.borehole.location[1] || '-123'
    }
  },
  methods: {
    fetchBorehole () {
      this.$http.get(`boreholes/${this.$route.params.bh}`).then((response) => {
        this.borehole = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving borehole data.')
      })
    },
    fetchStrata () {
      this.$http.get(`boreholes/${this.$route.params.bh}/strata`).then((response) => {
        this.strataRowData = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving soil data.')
      })
    },
    fetchSamples () {
      this.$http.get(`boreholes/${this.$route.params.bh}/samples`).then((response) => {
        this.sampleRowData = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving sample data.')
      })
    },
    fetchLabTests () {
      this.$http.get(`projects/${this.$route.params.id}/lab/tests?borehole=${this.$route.params.bh}`).then((response) => {
        this.labTestingRowData = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving laboratory data.')
      })
    }
  },
  created () {
    this.fetchBorehole()
    this.fetchStrata()
    this.fetchSamples()
    this.fetchLabTests()
  },
  beforeMount () {
    this.gridOptions = {
      context: {
        componentParent: this
      }
    }
  }

}
</script>

<style>

</style>
