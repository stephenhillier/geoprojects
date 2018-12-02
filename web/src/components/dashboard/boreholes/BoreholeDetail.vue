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
            <b-tab title="Stratigraphy">
              <h5>Soil Stratigraphy</h5>

              <ag-grid-vue style="height: 500px;"
                      :enableSorting="true"
                      :enableFilter="true"
                      rowHeight="32"
                      class="ag-theme-balham mb-3"
                      :columnDefs="strataColumnDefs"
                      :rowData="strata"
                      :enableColResize="true"
                      :cellValueChanged="onCellValueChanged"/>

              <b-btn size="sm" :variant="addNewStrata ? 'secondary' : 'primary'" @click="addNewStrata = !addNewStrata">{{ addNewStrata ? 'Cancel' : 'Add strata'}}</b-btn>
              <new-strata v-if="addNewStrata" :borehole="borehole.id" @strata-update="fetchStrata" @strata-dismiss="addNewStrata = false"></new-strata>
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

              <ag-grid-vue style="height: 500px;"
                 class="ag-theme-balham"
                 :columnDefs="labTestingColumnDefs"
                 :rowData="labTestingRowData">
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
import StrataDelete from '@/components/gridcells/StrataDelete.vue'
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
      strata: [],
      strataIsBusy: false,
      samplesIsBusy: false,
      addNewStrata: false,
      addNewSample: false,
      labTestingColumnDefs: [
        { headerName: 'Sample', field: 'sample' },
        { headerName: 'Test', field: 'test' }
      ],
      labTestingRowData: [
        { sample: 'SA-1', test: 'Moisture content' },
        { sample: 'SA-1', test: 'Grain size analysis' },
        { sample: 'SA-1', test: 'Hydrometer' }
      ],
      strataColumnDefs: [
        { headerName: 'From', field: 'start', filter: 'agNumberColumnFilter', width: 90, editable: true },
        { headerName: 'To', field: 'end', filter: 'agNumberColumnFilter', width: 90, editable: true },
        { headerName: 'Description', field: 'description', filter: 'agTextColumnFilter', width: 400, editable: true },
        { headerName: 'Soil tags', field: 'soils', filter: 'agTextColumnFilter' },
        { headerName: 'Moisture', field: 'moisture', filter: 'agTextColumnFilter', width: 140 },
        { headerName: 'Consistency', field: 'consistency', filter: 'agTextColumnFilter', width: 140 },
        { headerName: 'Actions', width: 100, cellRendererFramework: StrataDelete }

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
    onCellValueChanged (evt) {
      // this event fires even if the value didn't actually change.
      // if this is the case, stop here.
      if (evt.oldValue === evt.newValue) {
        return
      }

      const strataData = {
        borehole: String(evt.data.borehole),
        start: String(evt.data.start),
        end: String(evt.data.end),
        description: evt.data.description
      }

      this.$http.put(`strata/${evt.data.id}`, strataData).then((response) => {
        this.fetchStrata()
      }).catch((e) => {
        console.error(e)
      })
    },
    fetchBorehole () {
      this.$http.get(`boreholes/${this.$route.params.bh}`).then((response) => {
        this.borehole = response.data
      }).catch((e) => {
        // in future, set an error message
      })
    },
    fetchStrata () {
      this.$http.get(`boreholes/${this.$route.params.bh}/strata`).then((response) => {
        this.strata = response.data
      }).catch((e) => {
        console.error(e)
      })
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
    this.fetchStrata()
  }
}
</script>

<style>

</style>
