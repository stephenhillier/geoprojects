<template>
  <div class="is-fullheight">
    <div class="columns">
      <div class="column">
        <h4 class="card-title">Lab Testing: Moisture content</h4>
        <h5 class="font-weight-bold">{{ sample.borehole_name }} {{ sample.sample_name }}</h5>
        <h6 class="text-muted">{{project.name}}</h6>
      </div>
    </div>
    <div class="columns">
      <div class="column">
        <form @submit.prevent="handleSubmit">
          <b-field label="Tare mass (g)">
            <b-input
              id="moistureContentTare"
              type="text"
              required
              v-model="sample.tare_mass"
            ></b-input>
          </b-field>
          <b-field label="Sample mass (g)">
            <b-input
              id="moistureContentSampleMass"
              type="text"
              message="Sample mass (including tare) in grams"
              required
              v-model="sample.sample_plus_tare"
            ></b-input>
          </b-field>
          <b-field label="Dry mass (g)">
            <b-input
              id="moistureContentDryMass"
              type="text"
              message="Dried sample mass (including tare) in grams"
              required
              v-model="sample.dry_plus_tare"
            ></b-input>
          </b-field>

          <button class="button is-primary" type="submit">Save</button>

        </form>
      </div>
      <div class="column">
        <h2 class="subtitle">Results</h2>
        <div>
          Moisture content:
          <span v-if="calculatedMoisture" class="font-weight-bold">
            {{ calculatedMoisture.toFixed(2) }} %
          </span>
        </div>
      </div>
    </div>
  </div>
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
      let sample = Number(this.sample.sample_plus_tare || 0)
      let dry = Number(this.sample.dry_plus_tare || 0)
      let tare = Number(this.sample.tare_mass || 0)

      if (!sample || !dry || Number.isNaN(sample) || Number.isNaN(dry) || Number.isNaN(tare)) { return null }

      sample = Number(sample)
      dry = Number(dry)
      tare = Number(tare)

      return (dry - tare) > 0
        ? (sample - dry) / (dry - tare) * 100
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
        this.$noty.error('An error occurred while retrieving lab test data.')
      })
    },
    handleSubmit () {
      this.loading = true
      this.$http.put(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}/moisture`, this.toStrings(this.sample)).then((response) => {
        this.sample = this.toStrings(this.nullToStrings(response.data))
        this.$noty.success('Moisture content test updated.')
        this.loading = false
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while updating moisture content test.')
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
