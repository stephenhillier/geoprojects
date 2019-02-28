<template>
  <div>
    <b-row>
      <b-col cols="12" md="3" lg="2" xl="2">
        <b-card>
          <nav>
            <div>
              <ul class="nav flex-column">
                <b-nav-item active-class="active-menu" class="menu-item mt-3" exact :to="{name: 'projects'}">
                  <font-awesome-icon :icon="['fas', 'th-list']" class="text-muted mr-3"></font-awesome-icon>Project List
                </b-nav-item>
                <b-nav-item active-class="active-menu" class="menu-item" exact :to="{name: 'new-project'}">
                  <font-awesome-icon :icon="['far', 'plus-square']" class="text-muted mr-3"></font-awesome-icon>New Project
                </b-nav-item>
              </ul>
            </div>
          </nav>
        </b-card>
      </b-col>
      <b-col>
        <b-row>
          <b-col>
          </b-col>
        </b-row>
        <b-row>
          <b-col>
            <b-card>
              <div class="card-status bg-success"></div>

              <b-breadcrumb class="bg-light m-0 p-0 mb-3" :items="breadcrumbs"></b-breadcrumb>
              <b-row>
                <b-col>
                  <h1>Start new project</h1>
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="12" xl="6">
                  <b-form @submit.prevent="submit">
                    <b-row>
                      <b-col cols="12">
                        <form-input
                          id="projectNumber"
                          label="Project number"
                          v-model="form.number"
                        ></form-input>
                      </b-col>
                    </b-row>
                    <b-row>
                      <b-col cols="12">
                        <form-input
                          required
                          id="projectName"
                          label="Project Name *"
                          v-model="form.name"
                        ></form-input>
                      </b-col>
                    </b-row>
                    <b-row>
                      <b-col cols="12">
                        <form-input
                          id="projectClient"
                          label="Client"
                          v-model="form.client"
                        ></form-input>
                      </b-col>
                    </b-row>
                    <b-row>
                      <b-col cols="12">
                        <form-input
                          id="projectManager"
                          label="Project manager"
                          v-model="form.pm"
                        ></form-input>
                      </b-col>
                    </b-row>
                    <b-row>
                      <b-col cols="12">
                        <form-input
                          id="projectLocation"
                          label="Location Description"
                          hint="Enter a city or other geographic description such as a highway"
                          v-model="form.location"
                        ></form-input>
                      </b-col>
                    </b-row>
                    <fieldset class="my-3">
                      <legend class="h5">Default Location</legend>
                      <p>This is the default location for this project.  The project's map marker will be placed here when there are no data locations (such as boreholes) to display.</p>
                      <p class="small"><span class="font-weight-bold">Hint:</span> double click the map to place a marker.</p>
                      <b-row>
                        <b-col cols="12" md="6">
                          <form-input label="Latitude" id="newProjectLatitude" v-model.number="form.default_coords[1]"></form-input>
                        </b-col>
                        <b-col cols="12" md="6">
                          <form-input label="Longitude" id="newProjectLongitude" v-model.number="form.default_coords[0]"></form-input>
                        </b-col>
                      </b-row>
                    </fieldset>
                    <b-row>
                      <b-col>
                        <b-btn type="submit" variant="info">Create project</b-btn>
                      </b-col>
                    </b-row>
                  </b-form>
                </b-col>
                <b-col cols="12" xl="6">
                  <ew-map :longitude="form.default_coords[0]" :latitude="form.default_coords[1]" @update-coordinates="updateCoords" :add-mode="true"></ew-map>
                </b-col>
              </b-row>
            </b-card>
          </b-col>
        </b-row>
        <b-modal centered title="Creating a project" ref="tutorialProjectCreateModal" cancel-title="Don't show again" @cancel="handleCancelProjectCreateTutorial">
          <div class="d-block text-center">
            <p>This is where you fill out background information about your new project. Remember: you can always come back and edit the details later.</p>
            <p>Double click on the map to place a marker. This is where your project will appear on the map. Give it a try!</p>
          </div>
        </b-modal>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import FormInput from '@/components/common/FormInput.vue'
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'

export default {
  name: 'NewProject',
  components: {
    FormInput,
    'ew-map': SingleMarkerMap
  },
  data () {
    return {
      breadcrumbs: [
        {
          text: 'Projects',
          to: { name: 'projects' }
        },
        {
          text: 'New project',
          to: { name: 'new-project' }
        }
      ],
      form: {
        number: '',
        name: '',
        location: '',
        client: '',
        pm: '',
        default_coords: ['', '']
      }
    }
  },
  methods: {
    submit () {
      this.$http.post('projects/', this.form)
        .then(() => {
          this.$router.push({ name: 'projects' })
          this.$noty.success('Project created.')
        }).catch((e) => {
          this.$noty.error('An error occurred while creating project.')
        })
    },
    updateCoords (val) {
      const { lat, lng } = val
      this.form.default_coords = [lng, lat]
    },
    handleCancelProjectCreateTutorial () {
      localStorage.setItem('earthworks-tutorial-project-creation', JSON.stringify(true))
    }
  },
  created () {
    localStorage.setItem('earthworks-tutorial-projects', JSON.stringify(true))
    if (!JSON.parse(localStorage.getItem('earthworks-tutorial-project-creation'))) {
      setTimeout(() => {
        this.$refs.tutorialProjectCreateModal.show()
      }, 1000)
    }
  }

}
</script>

<style>

</style>
