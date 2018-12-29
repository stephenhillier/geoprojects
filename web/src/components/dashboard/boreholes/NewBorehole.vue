<template>
  <b-card>

    <b-form @submit.prevent="handleFormSubmit" @reset.prevent="resetForm">
      <b-row>
        <b-col cols="12" xl="6">
                <h4>New Borehole</h4>
                <h6 class="text-muted">{{project.name}}</h6>
          <b-row>
            <b-col cols="12" md="4" xl="3">
              <form-input label="Name" id="newBoreholeName" hint="e.g. BH18-1" v-model="form.name" required></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12" md="6">
              <form-input label="Date drilling started" type="date" id="newBoreholeDateStarted" v-model="form.start_date" required></form-input>
            </b-col>
            <b-col cols="12" md="6">
              <form-input label="Date drilling finished" type="date" id="newBoreholeDateStarted" v-model="form.end_date"></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12" md="6">
              <form-input label="Field technician/engineer" id="newBoreholeFieldEng" v-model.number="form.field_eng" required></form-input>
            </b-col>
          </b-row>
          <fieldset class="my-3">
            <legend class="h5">Location</legend>
            <b-row>
              <b-col cols="12" md="6">
                <form-input label="Latitude" id="newBoreholeLatitude" v-model.number="form.location[1]"></form-input>
              </b-col>
              <b-col cols="12" md="6">
                <form-input label="Longitude" id="newBoreholeLongitude" v-model.number="form.location[0]"></form-input>
              </b-col>
            </b-row>
          </fieldset>
        </b-col>
        <b-col cols="12" xl="6">
          <ew-map :longitude="form.location[0]" :latitude="form.location[1]" @update-coordinates="updateCoords" :add-mode="true"></ew-map>
        </b-col>

      </b-row>

      <div class="mb-3 mt-5">
        <b-btn type="submit" variant="dark">Submit</b-btn>
        <b-btn type="reset" class="mx-3" variant="light">Reset</b-btn>
      </div>
    </b-form>
  </b-card>
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
