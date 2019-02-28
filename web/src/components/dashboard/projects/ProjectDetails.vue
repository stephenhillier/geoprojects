<template>
  <div>
      <b-row class="mb-3">
        <b-col cols="12" xl="6">
          <h1>{{project.name}}</h1>
          <table class="table">
            <thead>
              <th>Overview</th>
              <th><b-btn v-b-modal.editProjectModal size="sm" class="float-right" variant="outline-info"><font-awesome-icon :icon="['far', 'edit']" class="text-muted"></font-awesome-icon> Edit</b-btn></th>
            </thead>
              <tbody class="font-weight-bold">
                  <tr>
                    <th class="table-heading">Project number</th>
                    <td>{{ project.number }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading">Client</th>
                    <td>{{ project.client }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading">Location</th>
                    <td>{{ project.location }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading">Project manager</th>
                    <td>{{ project.pm }}</td>
                  </tr>
                  <tr>
                    <th class="table-heading">Boreholes</th>
                    <td>{{ project.borehole_count }}</td>
                  </tr>
              </tbody>
          </table>
        </b-col>
        <b-col>
          <multi-marker-map :locations="boreholes"></multi-marker-map>
        </b-col>
      </b-row>
      <ag-grid-vue style="height: 400px;"
              :enableSorting="true"
              :enableFilter="true"
              rowHeight="32"
              class="ag-theme-balham mb-3"
              :columnDefs="columnDefs"
              :rowData="boreholes"/>
      <b-btn variant="info" size="sm" :to="{ name: 'new-borehole' }">New borehole</b-btn>

      <!-- Edit project form -->
      <b-modal id="editProjectModal" title="Edit project information" centered @ok="handleEdit" @cancel="handleResetEdit" @shown="handleShowEditForm">
        <b-form @submit.prevent="handleEdit">
          <b-row>
            <b-col cols="12">
              <form-input
                id="projectEditName"
                label="Name"
                required
                v-model="editForm.name"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              <form-input
                id="projectEditNumber"
                label="Number"
                required
                v-model="editForm.number"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              <form-input
                id="projectEditClient"
                label="Client"
                required
                v-model="editForm.client"
              ></form-input>
            </b-col>
            <b-col cols="12">
              <form-input
                id="projectEditLocation"
                label="Location"
                required
                v-model="editForm.location"
              ></form-input>
            </b-col>
            <b-col cols="12">
              <form-input
                id="projectEditPM"
                label="Project manager"
                required
                v-model="editForm.pm"
              ></form-input>
            </b-col>
          </b-row>
        </b-form>
      </b-modal>
      <b-modal centered title="Creating boreholes" ref="tutorialProjectSummaryModal" cancel-title="Don't show again" @cancel="handleCancelProjectSummaryTutorial">
        <div class="d-block text-center">
          <p>This is a summary of the project {{ project.name }}.</p>
          <p>You'll be able to create and view field data from this screen.</p>
          <p>Go ahead and click the "New Borehole" button to give it a try (don't worry, you can remove it afterward).</p>
        </div>
      </b-modal>
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
    handleShowEditForm () {
      this.editForm = Object.assign({}, this.project)
    },
    handleCancelProjectSummaryTutorial () {
      localStorage.setItem('earthworks-tutorial-project-summary', JSON.stringify(true))
    },
    handleEdit () {
      console.log(this.editForm)
    },
    handleResetEdit () {
      this.editForm = {}
    }
  },
  created () {
    this.fetchBoreholes()

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
