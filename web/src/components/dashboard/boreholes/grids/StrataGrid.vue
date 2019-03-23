<template>
  <div>
    <h2 class="subtitle">
      Soil Stratigraphy
    </h2>
    <button ref="newStrataBtn" @click="handleNewStrataModal" class="button" size="sm" variant="info">Add soil layer</button>

      <b-table
          :data="strataRowData"
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
            <b-table-column field="description" label="Description">
                {{ props.row.description }}
            </b-table-column>
            <b-table-column field="actions" label="Actions" class="is-narrow">
              <button class="button is-small" @click="handleEdit(props.row)"><font-awesome-icon :icon="['far', 'edit']"></font-awesome-icon></button>
              <button class="button is-small ml" @click="handleDelete(props.row.id)"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></button>
            </b-table-column>
        </template>
      </b-table>

    <!-- New strata modal -->
    <b-modal :active.sync="newStrataModal" id="newStrataModal" title="Add a new strata" @close="handleCloseNewStrataModal">
      <div class="modal-card" style="width: auto">
        <form @submit.stop.prevent="handleSubmit">
          <header class="modal-card-head">
            <p class="modal-card-title">Add strata</p>
          </header>
          <section class="modal-card-body">
              <div class="columns">
                <div class="column">
                  <b-field label="From">
                    <b-input
                      required
                      ref="strataStartInput"
                      id="strataStartInput"
                      message="Depth (m)"
                      v-model="form.start"
                    ></b-input>
                  </b-field>
                </div>
                <div class="column">
                  <b-field label="To">
                    <b-input
                      required
                      id="strataEndInput"
                      message="Depth (m)"
                      v-model="form.end"
                  ></b-input>
                  </b-field>
                </div>
                <div class="column">
                  <b-field label="Description">
                    <b-input
                      required
                      id="strataDescInput"
                      v-model="form.description"
                  ></b-input>
                  </b-field>
                </div>
              </div>

          </section>
          <footer class="modal-card-foot">
              <button class="button" type="button" @click="newStrataModal = false">Close</button>
              <button class="button is-primary">Add strata</button>
          </footer>
        </form>
      </div>
    </b-modal>

    <b-modal :active.sync="editStrataModal" id="editStrataModal" ref="editStrataModal" @close="handleResetEdit">
      <div class="modal-card">
        <form @submit.stop.prevent="handleEdit">
          <header class="modal-card-head">
            <p class="modal-card-title">Edit strata</p>
          </header>
          <section class="modal-card-body">
          <div class="columns">
            <div class="column">
              <b-field label="From">
                <b-input
                  required
                  id="strataEditStartInput"
                  message="Depth (m)"
                  v-model="editForm.start"
              ></b-input>
              </b-field>
            </div>
            <div class="column">
              <b-field label="To">
                <b-input
                  required
                  id="strataEditEndInput"
                  message="Depth (m)"
                  v-model="editForm.end"
              ></b-input>
              </b-field>
            </div>
            <div class="column">
              <b-field label="Description">
                <b-input
                  required
                  id="strataEditDescInput"
                  v-model="editForm.description"
              ></b-input>
              </b-field>
            </div>
          </div>
          </section>
          <footer class="modal-card-foot">
              <button class="button" type="button" @click="editStrataModal = false">Cancel</button>
              <button class="button is-primary">Save</button>
          </footer>
        </form>
      </div>
    </b-modal>
  </div>

</template>

<script>

export default {
  name: 'StrataGrid',

  props: {
    strataRowData: {
      type: Array,
      default: () => ([])
    },
    borehole: null
  },
  data () {
    return {
      newStrataModal: false,
      editStrataModal: false,
      strataIsBusy: false,
      addNewStrata: false,
      fields: [],
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
      editForm: {},
      perPage: 10,
      currentPage: 1
    }
  },
  methods: {
    onStrataGridReady (params) {
      this.strataGridApi = params.api
      this.strataColumnApi = params.columnApi
    },
    deleteStrata (id) {
      this.$http.delete(`strata/${id}`).then((response) => {
        this.$emit('strata-update')
        this.$noty.success('Soil layer deleted.')
      }).catch((e) => {
        console.error(e)
        this.$noty.error('An error occurred while deleting soil layer.')
      })
    },
    handleSubmit () {
      const data = Object.assign({}, this.form)
      data.borehole = this.$route.params.bh

      this.loading = true
      this.$http.post(`strata`, data).then((response) => {
        this.$noty.success('Soil layer added.')
        this.loading = false
        this.resetForm()
        this.$emit('strata-update')
        this.handleCloseNewStrataModal()
      }).catch((e) => {
        console.log(e)
        this.loading = false
        this.$noty.error('An error occurred while adding soil layer.')
      })
    },
    submitEdit () {
      const data = Object.assign({}, this.toStrings(this.editForm))
      const strataId = data.id
      delete data.id

      this.loading = true
      this.$http.put(`strata/${strataId}`, data).then((response) => {
        this.$noty.success('Soil layer updated.')
        this.loading = false
        this.$emit('strata-update')
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while updating soil layer.')
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

    },
    toStrings (o) {
      Object.keys(o).forEach((k) => {
        o[k] = '' + o[k]
      })
      return o
    },
    handleEdit (row) {
      this.editForm = Object.assign({}, row)
      this.editStrataModal = true
    },
    handleDelete (id) {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this soil strata record?',
        onConfirm: () => this.deleteStrata(id)
      })
    },
    handleNewStrataModal () {
      this.newStrataModal = true
      this.$nextTick(() => {
        this.$refs.strataStartInput.focus()
      })
    },
    handleCloseNewStrataModal () {
      this.newStrataModal = false
      this.resetForm()
      this.$refs.newStrataBtn.focus()
    }
  }
}
</script>

<style>

</style>
