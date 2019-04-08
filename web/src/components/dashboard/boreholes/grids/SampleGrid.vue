<template>
  <div>
    <h2 class="subtitle">
      Soil Samples
    </h2>
    <button ref="newSampleBtn" @click="handleNewSampleModal" class="button">Add soil sample</button>

      <b-table
          :data="samples"
          paginated
          :per-page="perPage"
          :current-page.sync="currentPage"
      >
        <template slot-scope="props">
            <b-table-column field="start" label="From" class="is-narrow">
                {{ props.row.start }}
            </b-table-column>
            <b-table-column field="end" label="To" class="is-narrow">
                {{ props.row.end }}
            </b-table-column>
            <b-table-column field="name" label="Name">
                {{ props.row.name }}
            </b-table-column>
            <b-table-column field="borehole" label="Sampled from">
                {{ props.row.borehole_name }}
            </b-table-column>
            <b-table-column field="actions" label="Actions" class="is-narrow">
              <button class="button is-small" @click="handleEdit(props.row)"><font-awesome-icon :icon="['far', 'edit']"></font-awesome-icon></button>
              <button class="button is-small ml" @click="handleDelete(props.row.id)"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></button>
            </b-table-column>
        </template>
      </b-table>

    <!-- New sample modal -->
    <b-modal :active.sync="newSampleModal" id="newSampleModal" title="Add a new sample" @close="handleCloseNewSampleModal">
      <div class="modal-card" style="width: auto">
        <form @submit.stop.prevent="handleSubmit">
          <fieldset :disabled="loading">
            <header class="modal-card-head">
              <p class="modal-card-title">Add sample</p>
            </header>
            <section class="modal-card-body">
                <div class="columns">
                  <div class="column">
                    <b-field label="From">
                      <b-input
                        required
                        ref="sampleStartInput"
                        id="sampleStartInput"
                        message="Depth (m)"
                        v-model="form.start"
                      ></b-input>
                    </b-field>
                  </div>
                  <div class="column">
                    <b-field label="To">
                      <b-input
                        required
                        id="sampleEndInput"
                        message="Depth (m)"
                        v-model="form.end"
                    ></b-input>
                    </b-field>
                  </div>
                  <div class="column" v-if="!currentBorehole">
                    <b-field label="Sampled from">
                      <b-select id="sampleLocationInput" required v-model="sampledFrom" message="Sampling location (e.g. borehole)">
                          <option v-for="(option, i) in boreholes" :key="`boreholeOption${i}`" :value="option.id">{{ option.name }}</option>
                      </b-select>
                    </b-field>
                  </div>
                  <div class="column">
                    <b-field label="Sample name">
                      <b-input
                        id="sampleNameInput"
                        message="Sample name, e.g. SA-1"
                        required
                        v-model="form.name"
                    ></b-input>
                    </b-field>
                  </div>
                </div>

            </section>
            <footer class="modal-card-foot">
                <button class="button" type="button" @click="newSampleModal = false">Close</button>
                <button class="button is-primary">Add sample</button>
            </footer>
          </fieldset>
        </form>
      </div>
    </b-modal>

    <b-modal :active.sync="editSampleModal" id="editSampleModal" ref="editSampleModal" @close="handleResetEdit">
      <div class="modal-card">
        <form @submit.stop.prevent="submitEdit">
          <fieldset :disabled="loading">
            <header class="modal-card-head">
              <p class="modal-card-title">Edit sample</p>
            </header>
            <section class="modal-card-body">
            <div class="columns">
              <div class="column">
                <b-field label="From">
                  <b-input
                    required
                    id="sampleEditStartInput"
                    ref="sampleStartEditInput"
                    message="Depth (m)"
                    v-model="editForm.start"
                ></b-input>
                </b-field>
              </div>
              <div class="column">
                <b-field label="To">
                  <b-input
                    required
                    id="sampleEditEndInput"
                    message="Depth (m)"
                    v-model="editForm.end"
                ></b-input>
                </b-field>
              </div>
              <div class="column" v-if="!currentBorehole">
                <b-field label="Sampled from">
                  <b-select id="sampleLocationInput" required v-model="editForm.borehole" message="Sampling location (e.g. borehole)">
                      <option v-for="(option, i) in boreholes" :key="`boreholeOption${i}`" :value="option.id">{{ option.name }}</option>
                  </b-select>
                </b-field>
              </div>
              <div class="column">
                <b-field label="Name">
                  <b-input
                        id="sampleNameInput"
                        message="Sample name, e.g. SA-1"
                        required
                        v-model="editForm.name"
                ></b-input>
                </b-field>
              </div>
            </div>
            </section>
            <footer class="modal-card-foot">
                <button class="button" type="button" @click="editSampleModal = false">Cancel</button>
                <button class="button is-primary">Save</button>
            </footer>
          </fieldset>
        </form>
      </div>
    </b-modal>
  </div>
</template>

<script>

export default {
  name: 'SampleGrid',

  props: {
    samples: {
      type: Array,
      default: () => ([])
    },
    borehole: null
  },
  data () {
    return {
      newSampleModal: false,
      editSampleModal: false,
      sampleIsBusy: false,
      addNewSample: false,
      fields: [],
      form: {
        start: '',
        end: '',
        name: ''
      },
      success: false,
      loading: false,
      editForm: {},
      perPage: 10,
      currentPage: 1,
      boreholes: []
    }
  },
  computed: {
    currentBorehole () {
      return this.$route.params.bh || null
    }
  },
  methods: {
    deleteSample (id) {
      this.$http.delete(`boreholes/${this.$route.params.bh}/samples/${id}`).then((response) => {
        this.$emit('sample-update')
        this.$noty.success('Soil sample deleted.')
      }).catch((e) => {
        console.error(e)
        this.$noty.error('An error occurred while deleting soil sample.')
      })
    },
    handleSubmit () {
      const data = Object.assign({}, this.form)
      data.borehole = this.$route.params.bh

      this.loading = true
      this.$http.post(`boreholes/${this.currentBorehole || this.sampledFrom}/samples`, data).then((response) => {
        this.$noty.success('Soil sample added.')
        this.loading = false
        this.resetForm()
        this.$emit('sample-update')
        this.handleCloseNewSampleModal()
      }).catch((e) => {
        console.log(e)
        this.loading = false
        this.$noty.error('An error occurred while adding soil sample.')
      })
    },
    submitEdit () {
      const data = Object.assign({}, this.toStrings(this.editForm))
      const sampleId = data.id
      delete data.id

      this.loading = true
      this.$http.put(`boreholes/${this.currentBorehole || data.borehole}/samples/${sampleId}`, data).then((response) => {
        this.$noty.success('Soil sample updated.')
        this.loading = false
        this.$emit('sample-update')
        this.editSampleModal = false
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while updating soil sample.')
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
      this.editSampleModal = true
      this.$nextTick(() => {
        if (this.$refs.sampleStartEditInput) {
          this.$refs.sampleStartEditInput.focus()
        }
      })
    },
    handleDelete (id) {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this soil sample record?',
        onConfirm: () => this.deleteSample(id)
      })
    },
    handleNewSampleModal () {
      this.newSampleModal = true
      this.$nextTick(() => {
        if (this.$refs.sampleStartInput) {
          this.$refs.sampleStartInput.focus()
        }
      })
    },
    handleCloseNewSampleModal () {
      this.newSampleModal = false
      this.resetForm()
      this.$refs.newSampleBtn.focus()
    },
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
    // if a borehole is not defined on the route, then fetch a list
    // of available boreholes for displaying samples across a whole project.
    // If a borehole is defined on the route, we only display samples for that
    // borehole and loading all boreholes isn't necessary.
    if (!this.$route.params.bh) {
      this.fetchBoreholes()
    }
  }
}
</script>

<style>

</style>
