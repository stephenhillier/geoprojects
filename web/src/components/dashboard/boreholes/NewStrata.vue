<template>
  <b-card class="my-3" title="Add new soil strata/layer">
    <b-form @submit.prevent="handleSubmit">
      <b-row>
        <b-col cols="12" lg="2" xl="1">
          <form-input
            id="strataStartInput"
            label="From"
            required
            v-model="form.start"
            hint="Depth (m)"
          ></form-input>
        </b-col>
        <b-col cols="12" lg="2" xl="1">
          <form-input
            id="strataEndInput"
            label="To"
            required
            v-model="form.end"
            hint="Depth (m)"
          ></form-input>
        </b-col>
        <b-col cols="12" lg="6" xl="8">
          <form-input
            id="strataDescriptionInput"
            label="Description"
            required
            v-model="form.description"
            hint="Soil visual description"
          ></form-input>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <b-btn size="sm" variant="primary" type="submit" :disabled="loading">Done</b-btn>
        </b-col>
      </b-row>
    </b-form>
  </b-card>
</template>

<script>
export default {
  name: 'NewStrata',
  props: ['borehole'],
  data () {
    return {
      form: {
        start: '',
        end: '',
        description: ''
      },
      success: false,
      loading: false
    }
  },
  methods: {
    handleSubmit () {
      const borehole = {
        borehole: String(this.$route.params.bh)
      }
      const data = Object.assign(borehole, this.form)

      this.loading = true
      this.$http.post('strata/', data).then((response) => {
        this.loading = false
        this.$noty.error('Soil layer added.')
        this.$emit('strata-update')
        this.$emit('strata-dismiss')
      }).catch((e) => {
        this.$noty.error('An error occurred while adding soil layer.')
      })
    }
  }
}
</script>

<style>

</style>
