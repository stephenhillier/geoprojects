<template>
  <div>
    <h5>
      Soil Samples
      <b-btn size="sm" class="ml-5" :variant="addNewSample ? 'primary' : 'secondary'" @click="addNewSample = !addNewSample">{{ addNewSample ? 'Cancel' : 'Add sample'}}</b-btn>
      <b-btn size="sm" variant="dark" class="ml-2" disabled>Edit sample</b-btn>
      <b-btn size="sm" variant="dark" class="ml-2" disabled>Delete sample</b-btn>
      <b-btn size="sm" variant="dark" class="ml-2" disabled>New lab test</b-btn>
    </h5>
    <new-sample v-if="addNewSample" :borehole="borehole.id" @sample-update="$emit('sample-update')" @sample-dismiss="addNewSample = false"></new-sample>
      <ag-grid-vue style="height: 500px;"
          class="ag-theme-balham mb-3"
          rowSelection="single"
          :columnDefs="sampleColumnDefs"
          :rowData="sampleRowData"
          :gridReady="onSampleGridReady"

          >
      </ag-grid-vue>
  </div>

</template>

<script>
import NewSample from '@/components/dashboard/boreholes/NewSample.vue'
import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'SampleGrid',
  components: {
    AgGridVue,
    NewSample
  },
  props: {
    sampleRowData: {
      type: Array,
      default: () => ([])
    }
  },
  data () {
    return {
      addNewSample: false,
      sampleGridApi: null,
      sampleColumnApi: null,
      sampleColumnDefs: [
        { headerName: 'From (m)', field: 'start', width: 110 },
        { headerName: 'To (m)', field: 'end', width: 110 },
        { headerName: 'Name', field: 'name', width: 150 }
      ]
    }
  },
  methods: {
    onSampleGridReady (params) {
      this.sampleGridApi = params.api
      this.sampleColumnApi = params.columnApi
    }
  }
}
</script>

<style>

</style>
