<template>
  <div>
    <h2 class="subtitle">
      Lab Tests
    </h2>
    <button ref="newTestBtn" @click="handleNewTestModal" class="button">Add lab test</button>

      <b-table
          :data="testRowData"
          paginated
          :per-page="perPage"
          :current-page.sync="currentPage"
      >
        <template slot-scope="props">
            <b-table-column field="name" label="Name" class="is-narrow">
              <router-link :to="{ name: getTestRoute(props.row.test_type), params: { id: $route.params.id, test: props.row.id }}">
                {{ props.row.borehole_name }} {{ props.row.sample_name }} ({{formatTestName(props.row.test_type)}}) {{ props.row.name }}
              </router-link>
            </b-table-column>
            <b-table-column field="test_type" label="Test type" class="is-narrow">
                {{ formatTestName(props.row.test_type)}}
            </b-table-column>
            <b-table-column field="start_date" label="Started" class="is-narrow">
                {{ props.row.start_date }}
            </b-table-column>
            <b-table-column field="end_date" label="Completed" class="is-narrow">
                {{ props.row.end_date }}
            </b-table-column>
            <b-table-column field="checked_date" label="Checked" class="is-narrow">
                {{ props.row.checked_date }} {{ props.row.checked_by }}
            </b-table-column>
            <b-table-column field="summary" label="Summary">
                {{ props.row.description }}
            </b-table-column>
            <b-table-column field="actions" label="Actions" class="is-narrow">
              <button class="button is-small"><font-awesome-icon :icon="['far', 'edit']"></font-awesome-icon></button>
              <button class="button is-small ml" @click="handleDelete(props.row.id)"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></button>
            </b-table-column>
        </template>
      </b-table>

    <!-- New test modal -->
    <b-modal :active.sync="newTestModal" id="newTestModal" title="Add a new test" @close="handleCloseNewTestModal">
      <div class="modal-card" style="width: auto">
        <form @submit.stop.prevent="handleSubmit">
          <header class="modal-card-head">
            <p class="modal-card-title">Add test</p>
          </header>
          <section class="modal-card-body">
            <b-field grouped>
              <b-field label="Sample">
                  <b-select v-model="form.sample" message="Choose a sample to test">
                      <option v-for="(option, i) in sampleOptions" :key="`soilNameOption${i}`" :value="option.id">{{ option.name }}</option>
                  </b-select>
              </b-field>
              <b-field label="Test type">
                  <b-select v-model="form.test_type" message="Select type of test">
                      <option v-for="(option, i) in testOptions" :key="`testTypeOption${i}`" :value="option.id">{{ option.description }}</option>
                  </b-select>
              </b-field>
          </b-field>
          </section>
          <footer class="modal-card-foot">
              <button class="button" type="button" @click="newTestModal = false">Close</button>
              <button class="button is-primary">Add test</button>
          </footer>
        </form>
      </div>
    </b-modal>
<!--
    <b-modal :active.sync="editTestModal" id="editTestModal" ref="editTestModal" @close="handleResetEdit">
      <div class="modal-card">
        <form @submit.stop.prevent="submitEdit">
          <header class="modal-card-head">
            <p class="modal-card-title">Edit test</p>
          </header>
          <section class="modal-card-body">
          <div class="columns">
            <div class="column">
              <b-field label="From">
                <b-input
                  required
                  id="testEditStartInput"
                  ref="testStartEditInput"
                  message="Depth (m)"
                  v-model="editForm.start"
              ></b-input>
              </b-field>
            </div>
            <div class="column">
              <b-field label="To">
                <b-input
                  required
                  id="testEditEndInput"
                  message="Depth (m)"
                  v-model="editForm.end"
              ></b-input>
              </b-field>
            </div>
            <div class="column">
              <b-field label="Description">
                <b-input
                  required
                  id="testEditDescInput"
                  v-model="editForm.description"
              ></b-input>
              </b-field>
            </div>
          </div>
          </section>
          <footer class="modal-card-foot">
              <button class="button" type="button" @click="editTestModal = false">Cancel</button>
              <button class="button is-primary">Save</button>
          </footer>
        </form>
      </div>
    </b-modal> -->
  </div>
</template>

<script>

export default {
  name: 'TestGrid',

  props: {
    testRowData: {
      type: Array,
      default: () => ([])
    },
    sampleOptions: {
      type: Array,
      default: () => ([])
    },
    borehole: null
  },
  data () {
    return {
      newTestModal: false,
      editTestModal: false,
      testIsBusy: false,
      addNewTest: false,
      fields: [],
      form: {
        start: '',
        end: '',
        description: ''
      },
      success: false,
      loading: false,
      editForm: {},
      perPage: 10,
      currentPage: 1,
      testOptions: [
        {
          id: 'moisture_content',
          description: 'Moisture content'
        },
        {
          id: 'grain_size_analysis',
          description: 'Grain size (sieve)'
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
    deleteTest (id) {
      this.$http.delete(`projects/${this.$route.params.id}/lab/tests/${id}`).then((response) => {
        this.$emit('test-update')
        this.$noty.success('Lab test deleted.')
      }).catch((e) => {
        console.error(e)
        this.$noty.error('An error occurred while deleting lab test.')
      })
    },
    handleSubmit () {
      const data = Object.assign({}, this.form)
      data.borehole = this.$route.params.bh

      this.loading = true
      this.$http.post(`projects/${this.$route.params.id}/lab/tests`, data).then((response) => {
        this.$noty.success('Lab test added.')
        this.loading = false
        this.resetForm()
        this.$emit('test-update')
        this.handleCloseNewTestModal()
      }).catch((e) => {
        console.log(e)
        this.loading = false
        this.$noty.error('An error occurred while adding lab test.')
      })
    },
    submitEdit () {
      const data = Object.assign({}, this.toStrings(this.editForm))
      const testId = data.id
      delete data.id

      this.loading = true
      this.$http.put(`projects/${this.$route.params.id}/lab/tests/${testId}`, data).then((response) => {
        this.$noty.success('Lab test updated.')
        this.loading = false
        this.$emit('test-update')
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while updating lab test.')
      })
    },
    resetForm () {
      this.form = {
        start: '',
        end: '',
        description: ''
      }
    },
    handleResetEdit () {

    },
    toStrings (o) {
      Object.keys(o).forEach((k) => {
        o[k] = '' + o[k]
      })
      return o
    },
    handleEdit (row) {
      this.editForm = Object.assign({}, row)
      this.editTestModal = true
      this.$nextTick(() => {
        if (this.$refs.testStartEditInput) {
          this.$refs.testStartEditInput.focus()
        }
      })
    },
    handleDelete (id) {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this test record?',
        onConfirm: () => this.deleteTest(id)
      })
    },
    handleNewTestModal () {
      this.newTestModal = true
      this.$nextTick(() => {
        if (this.$refs.testStartInput) {
          this.$refs.testStartInput.focus()
        }
      })
    },
    handleCloseNewTestModal () {
      this.newTestModal = false
      this.resetForm()
      this.$refs.newTestBtn.focus()
    },
    formatTestName (codeName) {
      const tests = {
        grain_size_analysis: 'Sieve analysis',
        moisture_content: 'Moisture content'
      }
      return tests[codeName] || codeName
    }
  }
}
</script>

<style>

</style>
