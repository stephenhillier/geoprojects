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
    <div class="container mt-3">
      <div class="card">
        <ThermChart
          v-if="chartLoaded"
          :chart-data="chartData"
          :options="options"
        />
      </div>
    </div>
  </div>
</template>

<script>
import SingleMarkerMap from '@/components/common/SingleMarkerMap.vue'
import ThermChart from '@/components/charts/ThermChart.js'

export default {
  name: 'InstrumentDetails',
  components: {
    SingleMarkerMap,
    ThermChart
  },
  data () {
    return {
      instrument: {
        location: []
      },
      instrData: [],
      instrFields: [
        {
          field: 'timestamp',
          label: 'Time'
        },
        {
          field: 'value',
          label: 'Value'
        }
      ],
      chartData: null,
      chartLoaded: false,
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          xAxes: [{
            type: 'time',
            time: {
              unit: 'day',
              round: 'hour',
              displayFormats: {
                day: 'MMM D'
              }
            }
          }]
        }
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
    },
    fetchData () {
      this.$http.get(`projects/${this.$route.params.id}/instrumentation/${this.$route.params.instr}/data`).then((response) => {
        this.instrData = response.data
        this.chartData = this.buildChartData(this.instrData)
        this.chartLoaded = true
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving instrument data.')
      })
    },
    buildChartData (dataset) {
      const labels = []
      const set = []
      for (let i = 0; i < dataset.length; i++) {
        console.log(dataset[i].timestamp, new Date(dataset[i].timestamp))

        labels.push(new Date(dataset[i].timestamp))
        set.push(dataset[i].value)
      }
      return {
        labels: labels,
        datasets: [
          {
            backgroundColor: '#fc6c71',
            borderColor: '#FC2525',
            borderWidth: 2,
            label: 'Resistance (Ohms)',
            fill: false,
            pointBackgroundColor: '#FC2525',
            pointBorderColor: '#FC2525',
            data: set
          }
        ]
      }
    }
  },
  created () {
    this.fetchInstrument()
    this.fetchData()
  }
}
</script>

<style>

</style>
