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
        location: ''
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
    }
  }

}
</script>

<style>

</style>
