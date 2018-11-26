<template>
  <b-card>
    <h2>{{ borehole.name }}</h2>
    <h6 class="text-muted">{{project.name}}</h6>
    <b-row class="mt-3 no-gutters">
      <b-col class="p-0 m-0">
        <b-card no-body>
          <b-tabs pills card>
            <b-tab title="Summary" active>
              <b-row>
                <b-col cols="12" lg="6" xl="6">
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

            </b-tab>
            <b-tab title="Stratigraphy">
              <h5>Soil Stratigraphy</h5>
              <b-table
                id="strataTable"
                ref="strataTable"
                responsive
                :busy.sync="strataIsBusy"
                :items="fetchStrata"
                show-empty
                :fields="['start', 'end', 'description', 'soils', 'moisture', 'consistency']"
                >
              </b-table>
              <b-btn size="sm" :variant="addNewStrata ? 'secondary' : 'primary'" @click="addNewStrata = !addNewStrata">{{ addNewStrata ? 'Cancel' : 'Add strata'}}</b-btn>
              <new-strata v-if="addNewStrata" :borehole="borehole.id" @strata-update="refreshStrata" @strata-dismiss="addNewStrata = false"></new-strata>
            </b-tab>
            <b-tab title="Samples">
              <h5>Soil Samples</h5>
              <b-table
                id="sampleTable"
                ref="sampleTable"
                responsive
                :busy.sync="samplesIsBusy"
                :items="fetchSamples"
                show-empty
                :fields="['start', 'end', 'sample_name', 'sample_type', 'tests_ordered', 'tests_completed']"
                >
              </b-table>
              <b-btn size="sm" :variant="addNewSample ? 'secondary' : 'primary'" @click="addNewSample = !addNewSample">{{ addNewSample ? 'Cancel' : 'Add sample'}}</b-btn>
            </b-tab>
            <b-tab title="Lab testing">
              <h5>Lab Testing</h5>
              <!-- <b-table
                id="labTable"
                ref="labTable"
                responsive
                show-empty
                :fields="['from', 'to', 'sample_name', 'test']"
                >
              </b-table> -->

              <ag-grid-vue style="height: 200px;"
                 class="ag-theme-balham"
                 :columnDefs="columnDefs"
                 :rowData="rowData">
              </ag-grid-vue>
            </b-tab>
          </b-tabs>
        </b-card>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'
import NewStrata from '@/components/dashboard/boreholes/NewStrata.vue'
import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'BoreholeDetails',
  components: {
    SingleMarkerMap,
    NewStrata,
    AgGridVue
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
      strataIsBusy: false,
      samplesIsBusy: false,
      addNewStrata: false,
      addNewSample: false,
      columnDefs: [
        { headerName: 'Sample', field: 'sample' },
        { headerName: 'Test', field: 'test' }
      ],
      rowData: [
        { sample: 'SA-1', test: 'Moisture content' },
        { sample: 'SA-1', test: 'Grain size analysis' },
        { sample: 'SA-1', test: 'Hydrometer' }
      ]
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
        // in future, set an error message
      })
    },
    fetchStrata (ctx = { perPage: this.perPage, currentPage: this.currentPage }) {
      /**
      * table items provider function
      * https://bootstrap-vue.js.org/docs/components/table/
      *
      * a refresh can be triggered by this.$root.$emit('bv::refresh::table', 'strataTable')
      */

      return this.$http.get(`boreholes/${this.$route.params.bh}/strata`).then((response) => {
        return response.data || []
      }).catch((e) => {
        return []
      })
    },
    refreshStrata () {
      this.$root.$emit('bv::refresh::table', 'strataTable')
      this.addNewStrata = false
    },
    fetchSamples (ctx = { perPage: this.perPage, currentPage: this.currentPage }) {
      /**
      * table items provider function
      * https://bootstrap-vue.js.org/docs/components/table/
      *
      * a refresh can be triggered by this.$root.$emit('bv::refresh::table', 'samplesTable')
      */

      return this.$http.get(`boreholes/${this.$route.params.bh}/samples`).then((response) => {
        return response.data || []
      }).catch((e) => {
        return []
      })
    }
  },
  created () {
    this.fetchBorehole()
  }
}
</script>

<style>

</style>
