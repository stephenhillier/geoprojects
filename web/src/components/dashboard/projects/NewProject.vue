<template>
  <div>
    <div class="columns">
      <div class="column is-narrow">
        <nav class="menu earthworks-project-menu">
          <p class="menu-label">
            Projects
          </p>
          <ul class="menu-list">
            <li class="project-list-menu-item"><router-link :to="{name: 'projects'}">
              <font-awesome-icon :icon="['fas', 'th-list']"></font-awesome-icon> Project List
            </router-link></li>
            <li class="project-list-menu-item"><router-link :to="{name: 'new-project'}">
              <font-awesome-icon :icon="['far', 'plus-square']"></font-awesome-icon> New Project
            </router-link></li>
          </ul>
        </nav>
      </div>
      <div class="column">
        <div class="box">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li v-for="(breadcrumb, i) in breadcrumbs" :key="`breadcrumb${i}`" :class="`${i === breadcrumbs.length - 1 ? 'is-active':''}`"><router-link :to="breadcrumb.to">{{breadcrumb.text}}</router-link></li>
            </ul>
          </nav>
              <h1 class="title">Start new project</h1>
          <div class="columns">
            <div class="column">
              <form @submit.prevent="submit">
                <b-field label="Project number" class="input-width-medium">
                  <b-input
                    id="projectNumber"
                    v-model="form.number"
                  >

                  </b-input>
                </b-field>
                <b-field label="Project name *">
                  <b-input
                      required
                      id="projectName"
                      v-model="form.name"
                  >

                  </b-input>
                </b-field>
                <b-field label="Client">
                  <b-input
                      id="projectClient"
                      v-model="form.client"
                  ></b-input>
                </b-field>
                <b-field label="Project manager">
                  <b-input
                      id="projectManager"
                      v-model="form.pm"
                  ></b-input>
                </b-field>
                <b-field label="Location">
                  <b-input
                      id="projectLocation"
                      label="Location Description"
                      message="Enter a city or other geographic description such as a highway"
                      v-model="form.location"
                  ></b-input>
                </b-field>

                <fieldset class="form-section">
                  <legend class="subtitle">Default Location</legend>
                  <p>This is the default location for this project.  The project's map marker will be placed here when there are no data locations (such as boreholes) to display.</p>
                  <p class="small"><span class="font-weight-bold">Hint:</span> double click the map to place a marker.</p>
                  <div class="columns mt-1">
                    <div class="column">
                      <b-field label="Latitude">
                        <b-input
                            id="newProjectLatitude"
                            v-model.number="form.default_coords[1]"
                        ></b-input>
                      </b-field>
                    </div>
                    <div class="column">
                      <b-field label="Longitude">
                        <b-input
                            id="newProjectLongitude"
                            v-model.number="form.default_coords[0]"
                        ></b-input>
                      </b-field>
                    </div>
                  </div>
                </fieldset>
                <button class="button mt-1 is-primary" type="submit" variant="info">Create project</button>
              </form>
            </div>
            <div class="column" style="max-height: 800px;">
              <ew-map :longitude="form.default_coords[0]" :latitude="form.default_coords[1]" @update-coordinates="updateCoords" :add-mode="true"></ew-map>
            </div>
          </div>
        </div>
        <b-modal centered title="Creating a project" ref="tutorialProjectCreateModal" cancel-title="Don't show again" @cancel="handleCancelProjectCreateTutorial">
          <div class="d-block text-center">
            <p>This is where you fill out background information about your new project. Remember: you can always come back and edit the details later.</p>
            <p>Double click on the map to place a marker. This is where your project will appear on the map. Give it a try!</p>
          </div>
        </b-modal>
      </div>
    </div>
  </div>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'

export default {
  name: 'NewProject',
  components: {
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
