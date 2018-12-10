<template>
  <b-card class="my-3" title="Add new soil sample">
    <b-form @submit.prevent="handleSubmit">
      <b-row>
        <b-col cols="12" lg="2" xl="1">
          <form-input
            id="sampleStartInput"
            label="From"
            required
            v-model="form.start"
            hint="Depth (m)"
          ></form-input>
        </b-col>
        <b-col cols="12" lg="2" xl="1">
          <form-input
            id="sampleEndInput"
            label="To"
            required
            v-model="form.end"
            hint="Depth (m)"
          ></form-input>
        </b-col>
        <b-col cols="12" lg="6" xl="8">
          <form-input
            id="sampleNameInput"
            label="Name"
            hint="Sample name, e.g. SA-1"
            required
            v-model="form.name"
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
  name: 'NewSample',
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
      const data = Object.assign({}, this.form)

      this.loading = true
      this.$http.post(`boreholes/${this.$route.params.bh}/samples`, data).then((response) => {
        this.success = true
        this.loading = false
        this.$emit('sample-update')
        this.$emit('sample-dismiss')
      }).catch((e) => {
        console.log(e)
        this.loading = false
      })
    }
  }
}
</script>

<style>

</style>
