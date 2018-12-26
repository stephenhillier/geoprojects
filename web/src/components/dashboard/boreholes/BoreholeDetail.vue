<template>
  <b-card>

    <b-row>
      <b-col cols="12" lg="6" xl="6">
        <h2>{{ borehole.name }}</h2>
        <h6 class="text-muted">{{project.name}}</h6>
        <h5>Summary</h5>
        <div v-if="borehole.location && borehole.location.length">Location: {{ latitude.toFixed(6) }}, {{ longitude.toFixed(6) }}</div>
        <div>
          Started: {{ borehole.start_date }}
        </div>
        <div>Completed: {{ borehole.end_date }}</div>
        <div>Logged by: {{ borehole.field_eng }}</div>
        <div>Logged soil strata: {{ borehole.strata_count }}</div>
        <div>Samples: 0</div>
        <div>Lab tests: 0</div>
      </b-col>
      <b-col cols="12" lg="6" xl="6">
        <b-card no-body>
          <single-marker-map :latitude="latitude" :longitude="longitude"></single-marker-map>
        </b-card>
      </b-col>
    </b-row>
    <b-row class="mt-3 no-gutters">
      <b-col class="p-0 m-0">
        <b-card no-body>
          <b-tabs pills card>
            <b-tab title="Stratigraphy" class="p-2 p-lg-3">
              <strata-grid :strataRowData="strataRowData" @strata-update="fetchStrata"/>
            </b-tab>
            <b-tab title="Samples" class="p-2 p-lg-3">
              <sample-grid :sampleRowData="sampleRowData" @sample-update="fetchSamples();fetchLabTests()"/>
            </b-tab>
            <b-tab title="Lab testing" class="p-2 p-lg-3">
              <lab-test-grid :labTestRowData="labTestingRowData" :sampleOptions="sampleRowData" @labtest-update="fetchLabTests" />
            </b-tab>
          </b-tabs>
        </b-card>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'
import StrataGrid from '@/components/dashboard/boreholes/grids/StrataGrid.vue'
import SampleGrid from '@/components/dashboard/boreholes/grids/SampleGrid.vue'
import LabTestGrid from '@/components/dashboard/boreholes/grids/LabTestGrid.vue'

import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'BoreholeDetails',
  components: {
    SingleMarkerMap,
    StrataGrid,
    LabTestGrid,
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
