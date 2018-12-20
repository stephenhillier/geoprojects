<template>
  <b-card>
    <b-row>
      <b-col>
        <h4 class="card-title">Lab Testing: {{ sample }} - Moisture content</h4>
        <h6 class="text-muted">{{project.name}}</h6>
      </b-col>
    </b-row>
    <b-row>
      <b-col>

      </b-col>
    </b-row>
  </b-card>
</template>

<script>
export default {
  name: 'MoistureContent',
  props: ['project'],
  data () {
    return {
      sample: {
        name: 'SA-1'
      },
      loading: false
    }
  },
  methods: {
    fetchTest () {
      this.loading = true
      this.$http.get(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}`).then((response) => {
        this.sample = response.data
        this.loading = false
      }).catch((e) => {
        this.loading = false
      })
    }
  }
}
</script>

<style>

</style>
