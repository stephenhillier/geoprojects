<template>
  <div>
    <div class="columns">
      <div class="column">
        <h1 class="title">{{ instrument.name }}</h1>
        <h2 class="subtitle">Summary</h2>
        <div v-if="instrument.location && instrument.location.length">Location: {{ latitude.toFixed(6) }}, {{ longitude.toFixed(6) }}</div>
        <div>
          Date installed: {{ instrument.start_date | moment('YYYY-MM-DD') }}
        </div>
        <div>Installed by: {{ instrument.field_eng }}</div>
      </div>
      <div class="column">
        <div class="is-480-map">
          <single-marker-map :latitude="latitude" :longitude="longitude"></single-marker-map>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'

export default {
  name: 'InstrumentDetails',
  components: {
    SingleMarkerMap
  },
  data () {
    return {
      instrument: {
        location: []
      }
    }
  },
  computed: {
    latitude () {
      return this.instrument.location[0] || '49'
    },
    longitude () {
      return this.instrument.location[1] || '-123'
    }
  },
  methods: {
    fetchInstrument () {
      this.$http.get(`projects/${this.$route.params.id}/instrumentation/${this.$route.params.instr}`).then((response) => {
        this.instrument = response.data
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving instrument summary.')
      })
    }
  },
  created () {
    this.fetchInstrument()
  }
}
</script>

<style>

</style>
