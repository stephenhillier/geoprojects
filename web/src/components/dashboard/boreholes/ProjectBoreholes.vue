<template>
  <b-card>

      <b-row class="mb-3">
        <b-col cols="12" xl="6">
                <h4>Boreholes</h4>
                <h6 class="text-muted">{{project.name}}</h6>
        </b-col>
        <b-col>
          <multi-marker-map :locations="boreholes"></multi-marker-map>
        </b-col>
      </b-row>
    <ag-grid-vue style="height: 500px;"
            :enableSorting="true"
            :enableFilter="true"
            rowHeight="32"
            class="ag-theme-balham mb-3"
            :columnDefs="columnDefs"
            :rowData="boreholes"/>
    <b-btn variant="secondary" size="sm" :to="{ name: 'new-borehole' }">New borehole</b-btn>

  </b-card>
</template>

<script>
import { AgGridVue } from 'ag-grid-vue'

import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'

import BoreholeLink from '@/components/gridcells/BoreholeLink.vue'
import Coords from '@/components/gridcells/Coords.vue'

export default {
  name: 'Boreholes',
  props: ['project'],
  components: {
    MultiMarkerMap,
    AgGridVue
  },
  data () {
    return {
      boreholes: [],
      numberOfRecords: 0,
      columnDefs: [
        { headerName: 'Name', field: 'name', filter: 'agTextColumnFilter', cellRendererFramework: BoreholeLink },
        { headerName: 'Started Drilling', field: 'start_date', filter: 'agDateColumnFilter' },
        { headerName: 'Finished Drilling', field: 'end_date', filter: 'agDateColumnFilter' },
        { headerName: 'Field Engineer', field: 'field_eng', filter: 'agTextColumnFilter' },
        { headerName: 'Location', field: 'location', cellRendererFramework: Coords }
      ]
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
