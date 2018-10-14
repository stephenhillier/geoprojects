<template>
  <b-card title="New Borehole" :sub-title="project.name">
    <b-form class="mt-4" @submit.prevent="handleFormSubmit" @reset.prevent="resetForm">
      <b-row>
        <b-col cols="12" md="4" xl="3">
          <form-input label="Name" id="newBoreholeName" hint="e.g. BH18-1" v-model="form.name"></form-input>
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="12" md="4" lg="3" xl="2">
          <form-input label="Date drilling started" type="date" id="newBoreholeDateStarted" v-model="form.start_date"></form-input>
        </b-col>
        <b-col cols="12" md="4" lg="3" xl="2">
          <form-input label="Date drilling finished" type="date" id="newBoreholeDateStarted" v-model="form.end_date"></form-input>
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="12" md="4" xl="3">
          <form-input label="Field technician/engineer" id="newBoreholeFieldEng" v-model.number="form.field_eng"></form-input>
        </b-col>
      </b-row>
      <fieldset class="my-3">
        <legend>Location</legend>
        <b-row>
          <b-col cols="12" md="4" lg="3" xl="2">
            <form-input label="Longitude" id="newBoreholeLongitude" v-model.number="form.location[0]"></form-input>
          </b-col>
          <b-col cols="12" md="4" lg="3" xl="2">
            <form-input label="Latitude" id="newBoreholeLatitude" v-model.number="form.location[1]"></form-input>
          </b-col>
        </b-row>
      </fieldset>

      <fieldset class="my-3">
        <legend class="mb-0">Soil strata</legend>
        <small class="mt-0">Note: you can add/modify soil strata later.</small>

      </fieldset>

      <div class="mb-3 mt-5">
        <b-btn type="submit" variant="dark">Submit</b-btn>
        <b-btn type="reset" class="mx-3" variant="light">Reset</b-btn>
      </div>
    </b-form>
  </b-card>
</template>

<script>
export default {
  name: 'NewBorehole',
  props: ['project'],
  data () {
    return {
      form: {},
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
      this.formSubmitSuccess = false
      this.formSubmitError = false
      this.$http.post('api/v1/field/boreholes', this.form).then((resp) => {
        this.formSubmitSuccess = true
        this.$emit('update-project')
        this.$router.push({ name: 'project-boreholes', params: { id: this.$route.params.id } })
      }).error((e) => {
        this.formSubmitError = true
      })
    }
  },
  created () {
    this.resetForm()
  }
}
</script>

<style>

</style>
