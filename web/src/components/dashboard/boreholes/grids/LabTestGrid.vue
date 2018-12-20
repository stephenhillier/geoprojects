<template>
  <div>
    <h5>
      Lab Testing
      <b-btn v-b-modal.newLabTestModal size="sm" variant="secondary" class="ml-5">New test</b-btn>
      <b-btn :to="{ name: 'lab-moisture', params: { id: $route.params.id, test: selectedRow }}" size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Test details</b-btn>
      <b-btn v-b-modal.deleteLabTestModal size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Delete test</b-btn>
    </h5>

    <!-- New lab test form -->
    <b-modal id="newLabTestModal" title="Add new lab test" centered @ok="handleSubmit">
      <b-form @submit.prevent="handleSubmit">
        <b-row>
          <b-col cols="12">
            <form-input
              id="labTestSample"
              label="Sample"
              required
              select
              value-field="id"
              text-field="name"
              placeholder="Select sample"
              :options="sampleOptions"
              v-model.number="form.sample"
            ></form-input>
          </b-col>
          <b-col cols="12">
            <form-input
              id="labTestType"
              label="Type of test"
              required
              select
              value-field="id"
              text-field="description"
              placeholder="Select test to perform"
              :options="testOptions"
              v-model.number="form.test_type"
            ></form-input>
          </b-col>
        </b-row>
      </b-form>
    </b-modal>

    <!-- Edit lab test form -->
    <b-modal id="editLabTestModal" title="Edit lab test" centered @ok="handleEdit" @cancel="handleResetEdit">
      <b-form @submit.prevent="handleEdit">
        <b-row>
          <b-col cols="12">
            <form-input
              id="labTestEditSample"
              label="Sample"
              disabled
              required
              select
              value-field="id"
              text-field="name"
              placeholder="Select sample"
              :options="sampleOptions"
              v-model.number="editForm.sample"
            ></form-input>
          </b-col>
          <b-col cols="12">
            <form-input
              id="labTestEditType"
              label="Type of test"
              required
              select
              value-field="id"
              text-field="description"
              placeholder="Select test to perform"
              :options="testOptions"
              v-model.number="editForm.test_type"
            ></form-input>
          </b-col>
        </b-row>
      </b-form>
    </b-modal>

    <!-- Delete lab test confirmation -->
    <b-modal id="deleteLabTestModal" centered @ok="handleDelete" title="Confirm delete">
      Are you sure you want to delete this lab test?
    </b-modal>

    <ag-grid-vue style="height: 500px;"
        class="ag-theme-balham"
        :columnDefs="labColumnDefs"
        :rowData="labTestRowData"
        :gridReady="onLabGridReady"
        rowSelection="single"
        :selectionChanged="onSelectionChanged"
        >
    </ag-grid-vue>
  </div>

</template>

<script>
import { AgGridVue } from 'ag-grid-vue'

export default {
  name: 'LabTestGrid',
  components: {
    AgGridVue
  },
  props: {
    labTestRowData: {
      type: Array,
      default: () => ([])
    },
    sampleOptions: {
      type: Array,
      default: () => ([])
    }
  },
  data () {
    return {
      labGridApi: null,
      labColumnApi: null,
      labColumnDefs: [
        { headerName: 'Sample', field: 'sample_name' },
        { headerName: 'Borehole', field: 'borehole_name' },
        { headerName: 'Test', field: 'test_type' }
      ],
      form: {
        name: '',
        start_date: '',
        end_date: '',
        test_type: '',
        sample: '',
        performed_by: ''
      },
      editForm: {
        name: '',
        start_date: '',
        end_date: '',
        test_type: '',
        sample: '',
        performed_by: ''
      },
      success: false,
      loading: false,
      selectedRow: null,
      testOptions: [
        {
          id: 'moisture',
          description: 'Moisture content'
        },
        {
          id: 'grainsize',
          description: 'Grain size (sieve)'
        }
      ]
    }
  },
  methods: {
    onLabGridReady (params) {
      this.labGridApi = params.api
      this.labColumnApi = params.columnApi
    },
    handleSubmit () {
      const data = Object.assign({}, this.form)

      this.loading = true
      this.$http.post(`projects/${this.$route.params.id}/lab/tests`, data).then((response) => {
        this.success = true
        this.loading = false
        this.$emit('labtest-update')
      }).catch((e) => {
        this.loading = false
      })
    },
    handleEdit () {
      const data = Object.assign({}, this.editForm)
      const testId = data.id
      delete data.id

      this.loading = true
      this.$http.put(`projects/${this.$route.params.id}/lab/tests/${testId}`, data).then((response) => {
        this.success = true
        this.loading = false
        this.$emit('labtest-update')
      }).catch((e) => {
        this.loading = false
      })
    },
    handleResetEdit () {
      const selection = this.labGridApi.getSelectedNodes()
      const rowData = selection.map((item) => (item.data))
      if (rowData && rowData.length) {
        this.editForm = Object.assign({}, rowData[0])
      }
    },
    onSelectionChanged () {
      const selection = this.labGridApi.getSelectedNodes()
      const rowData = selection.map((item) => (item.data))
      if (rowData && rowData.length) {
        this.selectedRow = rowData[0].id
        this.editForm = Object.assign({}, rowData[0])
      } else {
        this.selectedRow = null
      }
    },
    handleDelete () {
      this.loading = true
      this.$http.delete(`projects/${this.$route.params.id}/lab/tests/${this.selectedRow}`).then((response) => {
        this.$emit('labtest-update')
        this.loading = false
      }).catch((e) => {
        this.loading = false
      })
    }
  }
}
</script>

<style>

</style>
