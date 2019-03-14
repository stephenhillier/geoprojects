<template>
  <div>
    <h5>
      Tests
      <b-btn v-b-modal.newLabTestModal size="sm" variant="info" class="ml-5">New test</b-btn>
      <!-- <b-btn v-b-modal.deleteLabTestModal size="sm" variant="dark" class="ml-2" :disabled="!selectedRow">Delete test</b-btn> -->
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

    <!-- <ag-grid-vue style="height: 400px;"
        class="ag-theme-balham"
        :columnDefs="labColumnDefs"
        :rowData="labTestRowData"
        :gridReady="onLabGridReady"
        rowSelection="single"
        :selectionChanged="onSelectionChanged"
        >
    </ag-grid-vue> -->

    <b-table striped hover :items="labTestRowData" :fields="fields">
      <template slot="test" slot-scope="data">
        <router-link :to="{ name: getTestRoute(data.item.test_type), params: { id: $route.params.id, test: data.item.id }}">
          {{ data.item.borehole_name }} {{ data.item.sample_name }} ({{formatTestName(data.item.test_type)}}) {{ data.item.name }}
        </router-link>
      </template>
      <template slot="test_type" slot-scope="data">{{ formatTestName(data.value)}}</template>
      <template slot="actions" slot-scope="data">
        <router-link :to="{ name: getTestRoute(data.item.test_type), params: { id: $route.params.id, test: data.item.id }}">Details</router-link>
      </template>
    </b-table>

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
        { headerName: 'Sample', field: 'sample_name', width: 120 },
        { headerName: 'Borehole', field: 'borehole_name', width: 120 },
        { headerName: 'Test', field: 'test_type' },
        { headerName: 'Started', field: 'result_inputted', width: 120 },
        { headerName: 'Completed', field: 'result_completed', width: 120 },
        { headerName: 'Checked', field: 'result_checked', width: 120 }

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
          id: 'moisture_content',
          description: 'Moisture content'
        },
        {
          id: 'grain_size_analysis',
          description: 'Grain size (sieve)'
        }
      ],
      fields: [
        {
          key: 'test',
          sortable: false
        },
        {
          key: 'sample_name',
          label: 'Sample',
          sortable: true
        },
        {
          key: 'borehole_name',
          sortable: true
        },
        {
          key: 'test_type',
          sortable: true
        },
        {
          key: 'test_completed',
          label: 'Completed',
          sortable: true
        },
        {
          key: 'test_checked',
          label: 'Checked',
          sortable: true
        },
        {
          key: 'actions',
          sortable: false
        }
      ]
    }
  },
  computed: {
    selectedTestRoute () {
      // mapping of implemented test_types to frontend routes
      const testMap = {
        home: 'lab-home',
        moisture_content: 'lab-moisture',
        grain_size_analysis: 'lab-grainsize'
      }

      if (!this.selectedRow) {
        return 'lab-home'
      }

      const testObj = this.labTestRowData.find((test) => {
        return this.selectedRow === test.id
      }) || {}

      // return the route corresponding to the test type (defaulting to lab test home)
      return testMap[testObj['test_type'] || 'home']
    }
  },
  methods: {
    getTestRoute (type) {
      const testMap = {
        home: 'lab-home',
        moisture_content: 'lab-moisture',
        grain_size_analysis: 'lab-grainsize'
      }

      return testMap[type] || 'home'
    },
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
        this.$noty.success('Lab test created.')
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while creating lab test.')
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
        this.$noty.success('Lab test updated.')
      }).catch((e) => {
        this.$noty.error('An error occurred while updating lab test.')
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
        this.$noty.success('Lab test deleted.')
        this.loading = false
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while deleting lab test.')
      })
    },
    formatTestName (codeName) {
      const tests = {
        grain_size_analysis: 'Grain size / sieve',
        moisture_content: 'Moisture content'
      }
      return tests[codeName] || codeName
    }
  }
}
</script>

<style>

</style>
