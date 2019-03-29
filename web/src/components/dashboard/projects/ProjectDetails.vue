<template>
  <div class="is-fullheight">
      <div class="columns is-desktop">
        <div class="column is-half-desktop is-one-third-widescreen is-one-quarter-fullhd">
          <div class="columns">
            <span class="column title">{{project.name}}</span>
            <div class="column is-narrow">
              <b-dropdown aria-role="list">
                  <button class="button is-primary is-action-button" slot="trigger">
                      <span><font-awesome-icon :icon="['fas', 'cog']"></font-awesome-icon> Actions</span>
                      <b-icon icon="menu-down"></b-icon>
                  </button>

                  <b-dropdown-item aria-role="listitem">
                    <a @click="handleEditProject" href="#"><font-awesome-icon :icon="['far', 'edit']" class="text-muted"></font-awesome-icon> Edit</a>
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

        </div>
        <div class="column is-480-map">
          <multi-marker-map :locations="boreholes" class="is-fullheight"></multi-marker-map>
        </div>
      </div>
      <div class="section">
        <h2 class="subtitle">Boreholes</h2>
        <b-table
          :data="boreholes"
          :columns="fields"
          paginated
          :per-page="perPage"
          :current-page.sync="currentPage"
          >
          <template slot-scope="props">
              <b-table-column field="name" label="Borehole">
                 <router-link :to="`/projects/${$route.params.id}/boreholes/${props.row.id}`">{{ props.row.name }}</router-link>
              </b-table-column>

              <b-table-column field="start_date" label="Start date">
                  {{ props.row.start_date }}
              </b-table-column>
              <b-table-column field="end_date" label="End date">
                  {{ props.row.end_date }}
              </b-table-column>
              <b-table-column field="field_eng" label="Field technician/engineer">
                  {{ props.row.field_eng }}
              </b-table-column>
              <b-table-column field="location" label="Location">
                  {{ props.row.location[0] }},
                  {{ props.row.location[1] }}
              </b-table-column>
          </template>

        </b-table>
      </div>

      <project-files :files="projectFiles" :project="project" @updated="fetchFiles"></project-files>

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
                        <button class="button is-primary" type="submit">Save</button>
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
import ProjectFiles from './ProjectFiles.vue'

export default {
  name: 'ProjectDetails',
  props: ['project'],
  components: {
    MultiMarkerMap,
    AgGridVue,
    ProjectFiles
  },
  data () {
    return {
      isEditModalActive: false,
      boreholes: [],
      projectFiles: [],
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      fields: [
        {
          field: 'name',
          label: 'Name'
        },
        {
          field: 'start_date',
          label: 'Start date'
        },
        {
          field: 'end_date',
          label: 'End date'
        },
        {
          field: 'field_eng',
          label: 'Drilled by'
        },
        {
          field: 'location',
          label: 'Location'
        }
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
    },
    fetchFiles () {
      this.$http.get(`projects/${this.$route.params.id}/files`).then((r) => {
        this.projectFiles = r.data || []
      }).catch((e) => {
        this.$noty.error('Error retrieving project files. Please try again later.')
      })
    },
    handleEditProject () {
      this.editForm = Object.assign({}, this.project)
      this.isEditModalActive = true
    }
  },
  created () {
    this.fetchBoreholes()
    this.fetchFiles()
    // if (!JSON.parse(localStorage.getItem('earthworks-tutorial-project-summary'))) {
    //   setTimeout(() => {
    //     this.$refs.tutorialProjectSummaryModal.show()
    //   }, 1000)
    // }
  }
}
</script>

<style>
.is-action-button {
  margin-left: 0.5rem;
  margin-right: 0.5rem;
}
</style>
