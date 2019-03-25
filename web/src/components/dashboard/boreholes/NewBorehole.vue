<template>
  <div>
    <form @submit.prevent="handleFormSubmit" @reset.prevent="resetForm">
      <div class="columns">
        <div class="column">
          <h1 class="title">New Borehole</h1>
          <h2 class="subtitle">{{project.name}}</h2>
          <b-field label="Name">
            <b-input
                id="newBoreholeName"
                type="text"
                message="e.g. BH19-1"
                placeholder="e.g. BH19-1"
                v-model="form.name"
            ></b-input>
          </b-field>
          <div class="columns">
            <div class="column">
              <b-field label="Date drilling started">
                <b-datepicker
                    id="newBoreholeDateStarted"
                    placeholder="Click to select..."
                    icon="calendar"
                    v-model="form.start_date"
                    required
                    >
                </b-datepicker>
              </b-field>
            </div>
            <div class="column">
              <b-field label="Date drilling finished">
                <b-datepicker
                    id="newBoreholeDateFinished"
                    placeholder="Click to select..."
                    icon="calendar"
                    v-model="form.end_date"
                    >
                </b-datepicker>
              </b-field>
            </div>
          </div>
          <b-field label="Field technician/engineer">
            <b-input
                id="newBoreholeFieldEng"
                type="text"
                v-model.number="form.field_eng"
            ></b-input>
          </b-field>
          <fieldset class="my-3">
            <legend class="h5">Location</legend>
            <p class="small"><span class="font-weight-bold">Hint:</span> double click the map the place a marker.</p>

            <div class="columns mt-1">
              <div class="column">
                <b-field label="Latitude">
                  <b-input
                      id="newBoreholeLatitude"
                      v-model.number="form.location[1]"
                  ></b-input>
                </b-field>
              </div>
              <div class="column">
                <b-field label="Longitude">
                  <b-input
                      id="newBoreholeLongitude"
                      v-model.number="form.location[0]"
                  ></b-input>
                </b-field>
              </div>
            </div>
          </fieldset>
        </div>
        <div class="column">
          <ew-map :longitude="form.location[0]" :latitude="form.location[1]" @update-coordinates="updateCoords" :add-mode="true"></ew-map>
        </div>

      </div>

      <div class="mb-3 mt-5">
        <button class="button is-primary" type="submit">Create borehole</button>
        <button class="button mx-1" type="reset">Reset</button>
      </div>
    </form>
  </div>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'
export default {
  name: 'NewBorehole',
  components: {
    'ew-map': SingleMarkerMap
  },
  props: ['project'],
  data () {
    return {
      form: {},
      strata: [],
      formSubmitSuccess: false,
      formSubmitError: false
    }
  },
  methods: {
    resetForm () {
      this.form = {
        project: Number(this.$route.params.id),
        name: '',
        start_date: null,
        end_date: null,
        field_eng: null,
        location: ['', '']
      }
    },
    handleFormSubmit () {
      this.formSubmitError = false
      this.$http.post('boreholes', this.form).then((resp) => {
        this.$emit('update-project')
        this.$noty.success('Borehole added.')
        this.$router.push({ name: 'project-boreholes', params: { id: this.$route.params.id } })
      }).catch((e) => {
        this.formSubmitError = true
        this.$noty.error('An error occurred while adding borehole.')
      })
    },
    addStrataRow () {
      this.strata.push({})
    },
    removeStrataRow (index) {
      this.strata.splice(index, 1)
    },
    updateCoords (val) {
      const { lat, lng } = val
      this.form.location = [lng, lat]
    }
  },
  created () {
    this.resetForm()
  }
}
</script>

<style>

</style>
