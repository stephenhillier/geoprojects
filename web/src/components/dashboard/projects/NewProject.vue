<template>
  <div>
    <b-row>
      <b-col>
        <b-breadcrumb :items="breadcrumbs"></b-breadcrumb>
      </b-col>
    </b-row>
    <b-row>
      <b-col cols="12" md="3" lg="2" xl="2">
        <b-card no-body class="mb-3">
          <b-list-group flush>
            <b-list-group-item exact :to="{name: 'projects'}">Project List</b-list-group-item>
            <b-list-group-item exact :to="{name: 'new-project'}">New Project</b-list-group-item>
          </b-list-group>
        </b-card>
      </b-col>
      <b-col cols="12" md="7" lg="8">
        <b-card title="Add New Project">
          <b-row>
            <b-col cols="12" xl="6">
              <b-form @submit.prevent="submit">
                <b-row>
                  <b-col cols="12">
                    <form-input
                      id="projectName"
                      label="Project Name"
                      v-model="form.name"
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
                  <p><span class="font-weight-bold">Hint:</span> double click the map to place a marker.</p>
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
                    <b-btn type="submit">Submit</b-btn>
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
        name: '',
        location: '',
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
    }
  }

}
</script>

<style>

</style>
