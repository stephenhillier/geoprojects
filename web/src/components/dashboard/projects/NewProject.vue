<template>
  <div>
    <b-row class="mb-3">
      <b-col>
        <b-breadcrumb :items="breadcrumbs"></b-breadcrumb>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <b-card title="Add New Project">
          <b-form @submit.prevent="submit">
            <b-row>
              <b-col cols="12" md="6" lg="4">
                <form-input
                  id="projectName"
                  label="Project Name"
                  v-model="form.name"
                ></form-input>
              </b-col>
            </b-row>
            <b-row>
              <b-col cols="12" md="6" lg="4">
                <form-input
                  id="projectLocation"
                  label="Location"
                  v-model="form.location"
                ></form-input>
              </b-col>
            </b-row>
            <b-row>
              <b-col cols="12" md="6" lg="4">
                <form-input
                  id="projectPM"
                  label="Project manager"
                  v-model="form.pm"
                ></form-input>
              </b-col>
            </b-row>
            <b-row>
              <b-col>
                <b-btn type="submit">Submit</b-btn>
              </b-col>
            </b-row>
          </b-form>
        </b-card>
      </b-col>
    </b-row>
  </div>

</template>

<script>
import FormInput from '@/components/common/FormInput.vue'
export default {
  name: 'NewProject',
  components: {
    FormInput
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
        pm: ''
      }
    }
  },
  methods: {
    submit () {
      this.$http.post('/api/v1/projects/', this.form)
        .then(() => {
          this.$router.push({ name: 'projects' })
        }).catch((e) => {
          console.log(e)
        })
    }
  }

}
</script>

<style>

</style>
