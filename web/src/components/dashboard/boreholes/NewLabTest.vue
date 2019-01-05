<template>
    <b-form @submit.prevent="handleSubmit">
      <b-row>
        <b-col cols="12" lg="2" xl="1">
          <form-input
            id="labTestName"
            label="Test name"
            required
            v-model="form.name"
            hint="Depth (m)"
          ></form-input>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <b-btn size="sm" variant="info" type="submit" :disabled="loading">Done</b-btn>
        </b-col>
      </b-row>
    </b-form>
</template>

<script>
export default {
  name: 'NewLabTest',
  props: ['sampleId'],
  data () {
    return {
      form: {
        name: '',
        start_date: '',
        end_date: '',
        description: '',
        sample: this.sampleId,
        performed_by: ''
      },
      success: false,
      loading: false
    }
  },
  methods: {
    handleSubmit () {
      const data = Object.assign({}, this.form)

      this.loading = true
      this.$http.post(`projects/${this.$route.params.id}/lab/tests`, data).then((response) => {
        this.$noty.success('Lab test added.')

        this.loading = false
        this.$emit('labtest-update')
        this.$emit('labtest-dismiss')
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while adding lab test.')
      })
    }
  }
}
</script>

<style>

</style>
