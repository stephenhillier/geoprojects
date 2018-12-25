<template>
  <b-card>
    <b-row>
      <b-col>
        <h4 class="card-title">Lab Testing: Moisture content</h4>
        <h5 class="font-weight-bold">{{ sample.borehole_name }} {{ sample.sample_name }}</h5>
        <h6 class="text-muted">{{project.name}}</h6>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <b-form @submit.prevent="handleSubmit">
          <b-row>
            <b-col cols="12">
              <form-input
                id="moistureContentTare"
                label="Tare mass (g)"
                required
                v-model="sample.tare_mass"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              <form-input
                id="moistureContentSampleMass"
                label="Sample mass (g)"
                required
                hint="Sample mass (including tare) in grams"
                v-model="sample.sample_plus_tare"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              <form-input
                id="moistureContentDryMass"
                label="Dry mass (g)"
                required
                hint="Dried sample mass (including tare) in grams"
                v-model="sample.dry_plus_tare"
              ></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col>
              <b-btn type="submit" variant="primary">Submit</b-btn>
            </b-col>
          </b-row>
        </b-form>
      </b-col>
      <b-col>
        <div class="h6">Results</div>
        <div>
          Moisture content:
          <span v-if="calculatedMoisture">
            {{ calculatedMoisture }}
          </span>
        </div>
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
        tare_mass: '',
        sample_plus_tare: '',
        dry_plus_tare: ''
      },
      loading: false
    }
  },
  computed: {
    calculatedMoisture () {
      let sample = this.sample_plus_tare
      let dry = this.dry_plus_tare
      let tare = this.tare_mass

      if (!sample || !dry || tare === '') return null

      sample = Number(sample)
      dry = Number(dry)
      tare = Number(tare)

      return (dry - tare) > 0
        ? (sample - dry) / (dry - tare)
        : null
    }
  },
  methods: {
    fetchTest () {
      this.loading = true
      this.$http.get(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}/moisture`).then((response) => {
        this.sample = this.toStrings(this.nullToStrings(response.data))
        this.loading = false
      }).catch((e) => {
        this.loading = false
      })
    },
    handleSubmit () {
      this.loading = true
      this.$http.put(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}/moisture`, this.toStrings(this.sample)).then((response) => {
        this.sample = this.toStrings(this.nullToStrings(response.data))
        this.loading = false
      }).catch((e) => {
        this.loading = false
      })
    },
    toStrings (o) {
      Object.keys(o).forEach((k) => {
        o[k] = '' + o[k]
      })
      return o
    },
    nullToStrings (o) {
      Object.keys(o).forEach((k) => {
        if (o[k] === null) o[k] = ''
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
