<template>
  <div>
    <h5>
      Soil Samples
      <b-btn v-b-modal.newSampleModal size="sm" class="ml-5" variant="secondary">Add sample</b-btn>
      <b-btn v-b-modal.editSampleModal size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Edit sample</b-btn>
      <b-btn v-b-modal.deleteSampleModal size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Delete sample</b-btn>
    </h5>

    <!-- New sample modal -->
    <b-modal id="newSampleModal" title="Add a new sample" @ok="handleSubmit" @cancel="resetForm" @keydown.native.enter="handleSubmit">
      <b-container fluid>
        <b-form @submit.stop.prevent="handleSubmit">
          <b-row>
            <b-col cols="12" lg="12" xl="6">
              <form-input
                id="sampleStartInput"
                label="From"
                required
                v-model="form.start"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="2" xl="6">
              <form-input
                id="sampleEndInput"
                label="To"
                required
                v-model="form.end"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="12" xl="12">
              <form-input
                id="sampleNameInput"
                label="Name"
                hint="Sample name, e.g. SA-1"
                required
                v-model="form.name"
              ></form-input>
            </b-col>
          </b-row>
        </b-form>
      </b-container>
    </b-modal>

    <b-modal id="editSampleModal" ref="editSampleModal" title="Edit sample" @ok="handleEdit" @cancel="handleResetEdit" @keydown.native.enter="handleEdit;$refs.editSampleModal.hide()">
      <b-container fluid>
        <b-form @submit.stop.prevent="">
          <b-row>
            <b-col cols="12" lg="12" xl="6">
              <form-input
                id="sampleStartEditInput"
                label="From"
                required
                v-model="editForm.start"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="2" xl="6">
              <form-input
                id="sampleEndEditInput"
                label="To"
                required
                v-model="editForm.end"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="12" xl="12">
              <form-input
                id="sampleNameEditInput"
                label="Name"
                hint="Sample name, e.g. SA-1"
                required
                v-model="editForm.name"
              ></form-input>
            </b-col>
          </b-row>
        </b-form>
      </b-container>
    </b-modal>

    <!-- Delete sample confirmation -->
    <b-modal id="deleteSampleModal" centered @ok="handleDelete" title="Confirm delete">
      Are you sure you want to delete this sample?
    </b-modal>

    <ag-grid-vue style="height: 500px;"
          class="ag-theme-balham mb-3"
          rowSelection="single"
          :columnDefs="sampleColumnDefs"
          :rowData="sampleRowData"
          :gridReady="onSampleGridReady"
          :selectionChanged="onSelectionChanged"
          >
      </ag-grid-vue>
  </div>

</template>

<script>
import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'SampleGrid',
  components: {
    AgGridVue
  },
  props: {
    sampleRowData: {
      type: Array,
      default: () => ([])
    },
    borehole: null
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
      ],
      form: {
        start: '',
        end: '',
        name: ''
      },
      success: false,
      loading: false,
      selectedRow: null,
      editForm: {}
    }
  },
  methods: {
    onSampleGridReady (params) {
      this.sampleGridApi = params.api
      this.sampleColumnApi = params.columnApi
    },
    handleSubmit () {
      const data = Object.assign({}, this.form)

      this.loading = true
      this.$http.post(`boreholes/${this.$route.params.bh}/samples`, data).then((response) => {
        this.loading = false
        this.resetForm()
        this.$emit('sample-update')
        this.$emit('sample-dismiss')
        this.$noty.success('Sample created.')
      }).catch((e) => {
        console.log(e)
        this.loading = false
      })
    },
    handleEdit () {
      const data = Object.assign({}, this.toStrings(this.editForm))
      const sampleId = data.id
      delete data.id

      this.loading = true
      this.$http.put(`boreholes/${this.$route.params.bh}/samples/${sampleId}`, data).then((response) => {
        this.loading = false
        this.$noty.success('Sample updated.')
        this.$emit('sample-update')
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while updating sample.')
      })
    },
    resetForm () {
      this.form = {
        start: '',
        end: '',
        name: ''
      }
    },
    onSelectionChanged () {
      const selection = this.sampleGridApi.getSelectedNodes()
      const rowData = selection.map((item) => (item.data))
      if (rowData && rowData.length) {
        this.selectedRow = rowData[0].id
        this.editForm = Object.assign({}, rowData[0])
      } else {
        this.selectedRow = null
      }
    },
    handleResetEdit () {
      const selection = this.sampleGridApi.getSelectedNodes()
      const rowData = selection.map((item) => (item.data))
      if (rowData && rowData.length) {
        this.editForm = Object.assign({}, rowData[0])
      }
    },
    handleDelete () {
      this.loading = true
      this.$http.delete(`boreholes/${this.$route.params.bh}/samples/${this.selectedRow}`).then((response) => {
        this.$emit('sample-update')
        this.loading = false
        this.$noty.success('Sample deleted.')
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while deleting sample.')
      })
    },
    toStrings (o) {
      Object.keys(o).forEach((k) => {
        o[k] = '' + o[k]
      })
      return o
    }
  }
}
</script>

<style>

</style>
