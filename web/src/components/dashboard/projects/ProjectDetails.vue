<template>
  <div class="is-fullheight">
      <div class="columns is-desktop is-fullheight">
        <div class="column is-half">
          <div class="columns">
            <span class="column title">{{project.name}}</span>
            <div class="column is-narrow">
              <b-dropdown aria-role="list">
                  <button class="button is-primary is-button-right" slot="trigger">
                      <span><font-awesome-icon :icon="['far', 'edit']"></font-awesome-icon> Actions</span>
                      <b-icon icon="menu-down"></b-icon>
                  </button>

                  <b-dropdown-item aria-role="listitem">
                    <a @click="isEditModalActive = true" href="#"><font-awesome-icon :icon="['far', 'edit']" class="text-muted"></font-awesome-icon> Edit</a>
                  </b-dropdown-item>
                  <b-dropdown-item aria-role="listitem">
                    <router-link :to="{ name: 'new-borehole' }">
                      <font-awesome-icon :icon="['far', 'plus-square']"></font-awesome-icon>
                       New borehole
                    </router-link>
                  </b-dropdown-item>
                  <b-dropdown-item aria-role="listitem">
                    <a href="#" class="has-text-danger" @click="handleDelete">
                      <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
                      Delete project
                    </a>
                  </b-dropdown-item>
              </b-dropdown>
            </div>

          </div>
          <table class="table">
            <thead>
              <th><h2 class="subtitle">Overview</h2></th>
              <th></th>
            </thead>
              <tbody class="font-weight-bold">
                  <tr>
                    <th class="table-heading-horizontal">Project number</th>
                    <td>{{ project.number }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading-horizontal">Client</th>
                    <td>{{ project.client }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading-horizontal">Location</th>
                    <td>{{ project.location }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading-horizontal">Project manager</th>
                    <td>{{ project.pm }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading-horizontal">Boreholes</th>
                    <td>{{ project.borehole_count }}</td>
                  </tr>
              </tbody>
          </table>
          <div>
            <router-link tag="button" class="button is-link is-action-button" :to="{ name: 'new-borehole' }"><font-awesome-icon :icon="['far', 'plus-square']"></font-awesome-icon> New borehole</router-link>
          </div>
          <div class="mt-3">
            <button class="button is-danger is-action-button" @click="handleDelete">
              <font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon>
              Delete project
            </button>
          </div>
        </div>
        <div class="column is-full-map">
          <multi-marker-map :locations="boreholes" class="is-fullheight"></multi-marker-map>
        </div>
      </div>
      <!-- <ag-grid-vue style="height: 400px;"
              :enableSorting="true"
              :enableFilter="true"
              rowHeight="32"
              class="ag-theme-balham mb-3"
              :columnDefs="columnDefs"
              :rowData="boreholes"/> -->
      <router-link tag="button" class="button is-primary" id="button" :to="{ name: 'new-borehole' }">New borehole</router-link>
      <b-table
        :data="boreholes"
        :fields="fields">
      </b-table>

      <b-modal :active.sync="isEditModalActive" @close="handleResetEdit">
        <form action="" @submit.prevent="handleEdit">
                <div class="modal-card" style="width: auto">
                    <header class="modal-card-head">
                        <p class="modal-card-title">Edit project</p>
                    </header>
                    <section class="modal-card-body">
                        <b-field label="Name">
                            <b-input
                                id="projectEditName"
                                type="text"
                                v-model="editForm.name"
                                placeholder="Project name"
                                required>
                            </b-input>
                        </b-field>

                        <b-field label="Number">
                            <b-input
                                id="projectEditNumber"
                                type="text"
                                v-model="editForm.number"
                                placeholder="Project number"
                                required>
                            </b-input>
                        </b-field>
                        <b-field label="Client">
                          <b-input
                            id="projectEditClient"
                            type="text"
                            placeholder="Client"
                            required
                            v-model="editForm.client"
                          ></b-input>
                        </b-field>
                        <b-field label="Location">
                          <b-input
                            id="projectEditLocation"
                            type="text"
                            placeholder="Location"
                            required
                            v-model="editForm.location"
                          ></b-input>
                        </b-field>
                        <b-field label="Project manager">
                          <b-input
                            id="projectEditPM"
                            type="text"
                            placeholder="Project manager"
                            required
                            v-model="editForm.pm"
                          ></b-input>
                        </b-field>
                    </section>
                    <footer class="modal-card-foot">
                        <button class="button" type="button" @click="isEditModalActive = false">Close</button>
                        <button class="button is-primary" type="submit">Login</button>
                    </footer>
                </div>
            </form>
      </b-modal>

      <!-- <b-modal centered title="Creating boreholes" ref="tutorialProjectSummaryModal" cancel-title="Don't show again" @cancel="handleCancelProjectSummaryTutorial">
        <div class="d-block text-center">
          <p>This is a summary of the project {{ project.name }}.</p>
          <p>You'll be able to create and view field data from this screen.</p>
          <p>Go ahead and click the "New Borehole" button to give it a try (don't worry, you can remove it afterward).</p>
        </div>
      </b-modal> -->
  </div>
</template>

<script>
import { AgGridVue } from 'ag-grid-vue'
import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'
import BoreholeLink from '@/components/gridcells/BoreholeLink.vue'
import Coords from '@/components/gridcells/Coords.vue'

export default {
  name: 'ProjectDetails',
  props: ['project'],
  components: {
    MultiMarkerMap,
    AgGridVue
  },
  data () {
    return {
      isEditModalActive: false,
      boreholes: [],
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      fields: ['name', 'start_date', 'end_date', 'field_eng', 'location'],
      columnDefs: [
        { headerName: 'Name', field: 'name', filter: 'agTextColumnFilter', cellRendererFramework: BoreholeLink },
        { headerName: 'Started Drilling', field: 'start_date', filter: 'agDateColumnFilter' },
        { headerName: 'Finished Drilling', field: 'end_date', filter: 'agDateColumnFilter' },
        { headerName: 'Field Engineer', field: 'field_eng', filter: 'agTextColumnFilter' },
        { headerName: 'Location', field: 'location', cellRendererFramework: Coords }

      ],
      editForm: {}
    }
  },
  methods: {
    fetchBoreholes () {
      this.$http.get(`boreholes?project=${this.$route.params.id}`).then((response) => {
        this.boreholes = response.data.results
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving boreholes.')
      })
    },
    handleCancelProjectSummaryTutorial () {
      localStorage.setItem('earthworks-tutorial-project-summary', JSON.stringify(true))
    },
    handleEdit () {
      this.isEditModalActive = false
      console.log(this.editForm)
    },
    handleResetEdit () {
      this.editForm = Object.assign({}, this.project)
    },
    deleteProject () {
      this.$http.delete(`projects/${this.$route.params.id}`).then(() => {
        this.$router.push({ name: 'projects' })
        this.$noty.success('Project deleted')
      }).catch((e) => {
        this.$noty.error(`An error occurred while deleting project (${e.response.status})`)
      })
    },
    handleDelete () {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this project?',
        onConfirm: this.deleteProject
      })
    }
  },
  created () {
    this.fetchBoreholes()
    this.editForm = Object.assign({}, this.project)
    if (!JSON.parse(localStorage.getItem('earthworks-tutorial-project-summary'))) {
      setTimeout(() => {
        this.$refs.tutorialProjectSummaryModal.show()
      }, 1000)
    }
  }
}
</script>

<style>

</style>
