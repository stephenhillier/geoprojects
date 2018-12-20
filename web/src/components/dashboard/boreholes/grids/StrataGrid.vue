<template>
  <div>
    <h5>
      Soil Stratigraphy
      <b-btn v-b-modal.newStrataModal class="ml-5" size="sm" variant="secondary">Add soil layer</b-btn>
      <b-btn v-b-modal.editStrataModal size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Edit layer</b-btn>
      <b-btn v-b-modal.deleteStrataModal size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Delete layer</b-btn>
    </h5>
    <new-strata v-if="addNewStrata" :borehole="borehole.id" @strata-update="fetchStrata" @strata-dismiss="addNewStrata = false"></new-strata>
    <ag-grid-vue style="height: 500px;"
            :enableSorting="true"
            :enableFilter="true"
            rowHeight="32"
            class="ag-theme-balham mb-3"
            :columnDefs="strataColumnDefs"
            :rowData="strataRowData"
            :enableColResize="true"
            :gridReady="onStrataGridReady"
            :gridOptions="gridOptions"
            :selectionChanged="onSelectionChanged"
            rowSelection="single"
            />

    <!-- New strata modal -->
    <b-modal centered id="newStrataModal" title="Add a new strata" @ok="handleSubmit" @cancel="resetForm">
      <b-container fluid>
        <b-form @submit.stop.prevent="handleSubmit">
          <b-row>
            <b-col cols="12" lg="12" xl="6">
              <form-input
                id="strataStartInput"
                label="From"
                required
                v-model="form.start"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="2" xl="6">
              <form-input
                id="strataEndInput"
                label="To"
                required
                v-model="form.end"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="12" xl="12">
              <form-input
                id="strataDescInput"
                label="Description"
                required
                v-model="form.description"
              ></form-input>
            </b-col>
          </b-row>
        </b-form>
      </b-container>
    </b-modal>

    <b-modal centered id="editStrataModal" ref="editStrataModal" title="Edit strata" @ok="handleEdit" @cancel="handleResetEdit">
      <b-container fluid>
        <b-form @submit.stop.prevent="">
          <b-row>
            <b-col cols="12" lg="12" xl="6">
              <form-input
                id="strataStartEditInput"
                label="From"
                required
                v-model="editForm.start"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="2" xl="6">
              <form-input
                id="strataEndEditInput"
                label="To"
                required
                v-model="editForm.end"
                hint="Depth (m)"
              ></form-input>
            </b-col>
            <b-col cols="12" lg="12" xl="12">
              <form-input
                id="strataDescriptionEditInput"
                label="Visual Description"
                required
                v-model="editForm.description"
              ></form-input>
            </b-col>
          </b-row>
        </b-form>
      </b-container>
    </b-modal>

    <!-- Delete strata confirmation -->
    <b-modal id="deleteStrataModal" centered @ok="handleDelete" title="Confirm delete">
      Are you sure you want to delete this soil strata?
    </b-modal>

  </div>

</template>

<script>
import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'StrataGrid',
  components: {
    AgGridVue
  },
  props: {
    strataRowData: {
      type: Array,
      default: () => ([])
    },
    borehole: null
  },
  data () {
    return {
      strataIsBusy: false,
      addNewStrata: false,
      strataColumnDefs: [
        { headerName: 'From (m)', field: 'start', filter: 'agNumberColumnFilter', width: 110 },
        { headerName: 'To (m)', field: 'end', filter: 'agNumberColumnFilter', width: 110 },
        { headerName: 'Description', field: 'description', filter: 'agTextColumnFilter', width: 400 }
      ],
      strataGridApi: null,
      strataColumnApi: null,
      form: {
        start: '',
        end: '',
        description: ''
      },
      success: false,
      loading: false,
      selectedRow: null,
      gridOptions: {},
      editForm: {}
    }
  },
  methods: {
    onStrataGridReady (params) {
      this.strataGridApi = params.api
      this.strataColumnApi = params.columnApi
    },
    handleDelete () {
      this.$http.delete(`strata/${this.selectedRow}`).then((response) => {
        this.$emit('strata-update')
      }).catch((e) => {
        console.error(e)
      })
    },
    handleSubmit () {
      const data = Object.assign({}, this.form)
      data.borehole = this.$route.params.bh

      this.loading = true
      this.$http.post(`strata`, data).then((response) => {
        this.success = true
        this.loading = false
        this.resetForm()
        this.$emit('strata-update')
      }).catch((e) => {
        console.log(e)
        this.loading = false
      })
    },
    handleEdit () {
      const data = Object.assign({}, this.toStrings(this.editForm))
      const strataId = data.id
      delete data.id

      this.loading = true
      this.$http.put(`strata/${strataId}`, data).then((response) => {
        this.success = true
        this.loading = false
        this.$emit('strata-update')
      }).catch((e) => {
        this.loading = false
      })
    },
    resetForm () {
      this.form = {
        start: '',
        end: '',
        description: ''
      }
    },
    onSelectionChanged () {
      const selection = this.strataGridApi.getSelectedNodes()
      const rowData = selection.map((item) => (item.data))
      if (rowData && rowData.length) {
        this.selectedRow = rowData[0].id
        this.editForm = Object.assign({}, rowData[0])
      } else {
        this.selectedRow = null
      }
    },
    handleResetEdit () {
      const selection = this.strataGridApi.getSelectedNodes()
      const rowData = selection.map((item) => (item.data))
      if (rowData && rowData.length) {
        this.editForm = Object.assign({}, rowData[0])
      }
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
