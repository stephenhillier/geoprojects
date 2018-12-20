<template>
  <b-card>
    <b-row>
      <b-col>
        <h4 class="card-title">Lab Testing: {{ sample.name }} {{ sample.name ? '-': ''}} Moisture content</h4>
        <h6 class="text-muted">{{project.name}}</h6>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <b-form @submit.prevent="handleSubmit">
          <b-row>
            <b-col cols="12" lg="4">
              <form-input
                id="moistureContentTare"
                label="Tare mass (g)"
                required
                v-model="moistureData.tare_mass"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12" lg="4">
              <form-input
                id="moistureContentSampleMass"
                label="Sample mass (g)"
                required
                hint="Sample mass (including tare) in grams"
                v-model="moistureData.sample_plus_tare"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12" lg="4">
              <form-input
                id="moistureContentDryMass"
                label="Dry mass (g)"
                required
                hint="Dried sample mass (including tare) in grams"
                v-model="moistureData.dry_plus_tare"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col>
              <b-btn type="submit" variant="primary">{{ testExists ? 'Update' : 'Submit'}}</b-btn>
            </b-col>
          </b-row>
        </b-form>
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
        name: ''
      },
      moistureData: {
        tare_mass: '',
        sample_plus_tare: '',
        dry_plus_tare: ''
      },
      testExists: false,
      loading: false
    }
  },
  methods: {
    fetchTest () {
      this.loading = true
      this.$http.get(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}`).then((response) => {
        this.sample = response.data
      }).catch((e) => {
      })
      this.$http.get(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}/moisture`).then((response) => {
        this.moistureData = response.data
        this.testExists = true
        this.loading = false
      }).catch((e) => {
        this.loading = false
      })
    },
    handleSubmit () {
      if (this.testExists) {
        this.loading = true
        this.$http.put(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}/moisture`, this.toStrings(this.moistureData)).then((response) => {
          this.moistureData = response.data
          this.testExists = true
          this.loading = false
        }).catch((e) => {
          this.loading = false
        })
      } else {
        this.loading = true
        this.$http.post(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}/moisture`, this.toStrings(this.moistureData)).then((response) => {
          this.moistureData = response.data
          this.testExists = true
          this.loading = false
        }).catch((e) => {
          this.loading = false
        })
      }
    },
    toStrings (o) {
      Object.keys(o).forEach((k) => {
        o[k] = '' + o[k]
      })
      return o
    }
  },
  created () {
    this.fetchTest()
  }
}
</script>

<style>

</style>
