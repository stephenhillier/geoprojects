<template>
  <div>
    <b-row class="mt-4">
      <b-col>
        <h4 class="card-title">Lab Testing: Grain size analysis</h4>
        <h5 class="font-weight-bold">{{ sample.borehole_name }} {{ sample.sample_name }} <b-btn class="float-right">Help</b-btn></h5>
        <h6 class="text-muted">{{project.name}}</h6>
      </b-col>
    </b-row>
    <b-row class="mt-5">
      <b-col cols="12" md="5">
        <b-form @submit.prevent="handleSubmit">
          <b-row class="mt-3">
            <b-col>
              <h5>Sample details</h5>
            </b-col>
          </b-row>
          <b-row>
            <b-col>
              <form-input id="sampleTare" label="Tare mass" v-model="sample.tare_mass"></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col>
              <form-input id="sampleDryMass" label="Dry sample (plus tare)" v-model="sample.dry_plus_tare"></form-input>
            </b-col>
          </b-row>
          <b-row>
            <b-col>
              <form-input id="sampleWashedMass" label="Washed sample (plus tare)" v-model="sample.washed_plus_tare"></form-input>
            </b-col>
          </b-row>
          <b-row class="mt-3">
            <b-col>
              <h5>Sieve masses</h5>
            </b-col>
          </b-row>
          <b-row>
            <b-col>
              <table class="table">
                <thead>
                  <th class="text-dark table-heading col-5">Size</th>
                  <th class="text-dark table-heading col-5">Mass retained</th>
                  <th class="col-2"><b-btn size="sm" variant="outline-primary" @click="handleAddSieve">Add <span class="d-none d-sm-inline"> sieve</span></b-btn></th>
                </thead>
                <tbody>
                  <template v-for="(sieve, i) in sample.sieves">
                    <tr :key="`sieve${i}`">
                      <td>
                        <form-input :id="`size${i}`" groupClass="p-0 m-0" v-model="sample.sieves[i].size" :disabled="sample.sieves[i].pan"></form-input>
                      </td>
                      <td>
                        <form-input :id="`mass${i}`" groupClass="p-0 m-0" v-model="sample.sieves[i].mass_retained"></form-input>
                      </td>
                      <td class="align-middle text-center">
                        <b-btn size="sm" @click="handleRemoveSieve(i)" v-if="!sample.sieves[i].pan"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></b-btn>
                      </td>
                    </tr>
                  </template>
                </tbody>
              </table>
            </b-col>
          </b-row>
          <b-row>
            <b-col class="text-right">
              <b-btn type="submit" variant="primary">Update</b-btn>
            </b-col>
          </b-row>
        </b-form>
      </b-col>
      <b-col cols="12" md="7" class="pl-md-5">
        <!-- <apexchart width="100%" type="line" :options="sieveOptions" :series="sieveSeries"></apexchart> -->
        <gsa-chart :chart-data="chartData" :options="chartOptions" v-if="chartData.datasets && chartData.datasets.length && chartData.datasets[0].data.length"></gsa-chart>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import { SieveTest } from 'geotech-utils/sieve'
import GSAChart from '@/components/charts/GSAChart.vue'

export default {
  name: 'GrainSizeAnalysis',
  props: ['project'],
  components: {
    'gsa-chart': GSAChart
  },
  data () {
    return {
      sample: {
        tare_mass: '',
        sample_plus_tare: '',
        washed_plus_tare: '',
        dry_plus_tare: '',
        sieves: [
          {
            pan: true,
            size: 'Pan',
            mass_retained: '0'
          }
        ]
      },
      loading: false,
      defaultSpecSizes: [0.08, 0.1, 0.16, 0.2, 0.32, 0.640, 1, 2, 5, 10, 12, 16, 20, 50, 75, 100]
    }
  },
  computed: {
    sieveResult () {
      const wetMass = Number(this.sample.sample_plus_tare || 0)
      const dryMass = Number(this.sample.dry_plus_tare || 0)
      const washedMass = Number(this.sample.washed_plus_tare || 0)
      const tareMass = Number(this.sample.tare_mass || 0)

      if (!dryMass || !washedMass) {
        return []
      }

      const sample = {
        wetMass: wetMass - tareMass,
        washedMass: washedMass - tareMass,
        dryMass: dryMass - tareMass
      }

      const test = new SieveTest({ sample })

      this.sample.sieves.forEach((s) => {
        const size = Number(s.size)
        const retained = Number(s.mass_retained)

        if (s.size === 'Pan') {
          test.addSieve('Pan')
          test.sieve('Pan').retained(retained)
        } else if (size && !Number.isNaN(size) && !Number.isNaN(retained)) {
          test.addSieve(size)
          test.sieve(size).retained(retained)
        }
      })

      return test.passing()
    },
    chartOptions () {
      return {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          yAxes: [{
            type: 'linear',
            ticks: {
              min: 0,
              max: 100,
              callback: function (value, index, values) {
                return value + '%'
              }
            }
          }],
          xAxes: [{
            type: 'logarithmic',
            ticks: {
              min: 0,
              max: 100,
              callback: (value, index, values) => {
                if (!this.defaultSpecSizes.includes(value)) {
                  return ''
                }
                return value + ' mm'
              }
            }
          }]
        }
      }
    },
    chartData () {
      return {
        datasets: [
          {
            label: 'Grain size distribution',
            lineTension: 0,
            fill: false,
            borderColor: '#1f548a',
            data: this.sieveSeries.map((i) => {
              return { x: i.x, y: i.y }
            })
          }
        ]
      }
    },
    sieveSeries () {
      const passing = this.sieveResult.filter((i) => {
        return (i.size !== 'Pan' && !Number.isNaN(Number(i.size)))
      }).map((i) => {
        return { x: Number(i.size), y: Number(i.percentPassing) }
      }).reverse()

      return passing
    }
  },
  methods: {
    fetchTest () {
      this.loading = true
      this.$http.get(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}`).then((response) => {
        this.sample = this.toStrings(response.data)
        this.initializeSieveArray()
        this.loading = false
      }).catch((e) => {
        this.loading = false
        this.$noty.error('An error occurred while retrieving lab test data.')
      })
    },
    // removes a sieve as position i
    handleRemoveSieve (i) {
      this.sample.sieves.splice(i, 1)
    },
    handleAddSieve () {
      this.initializeSieveArray()

      this.sample.sieves.unshift(
        {
          pan: false,
          size: '',
          mass_passing: ''
        }
      )
    },
    handleSubmit () {
      const sample = JSON.parse(JSON.stringify(this.sample))

      if (sample.sieves && sample.sieves.length) {
        sample.sieves.forEach((sieve, i) => {
          if (sieve.size === 'Pan') {
            sieve.size = '0'
            sieve.pan = true
          }

          if (sieve.size === '' || sieve.mass_retained === '') {
            sample.sieves.splice(i, 1)
          }
        })
      }

      this.loading = true
      this.$http.put(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}`, sample).then((repsonse) => {
        this.$noty.success('Updated grain size test')
        this.fetchTest()
      }).catch((e) => {
        this.loading = false
        this.$noty.error('Error updating grain size test')
      })
    },
    toStrings (o) {
      Object.keys(o).forEach((k) => {
        if (typeof o[k] === 'number') {
          o[k] = '' + o[k]
        }
      })
      return o
    },
    nullToStrings (o) {
      Object.keys(o).forEach((k) => {
        if (o[k] === null) o[k] = ''
      })
      return o
    },
    //  start sieve array with default values, or add Pan if not exists (it's required for the test)
    initializeSieveArray () {
      if (!this.sample.sieves) {
        this.sample.sieves = [
          {
            pan: true,
            size: 'Pan',
            mass_passing: '0'
          }
        ]
      }

      const pan = this.sample.sieves.findIndex((i) => {
        return i.pan === true
      })

      if (!~pan) {
        this.sample.sieves.push({
          pan: true,
          size: 'Pan',
          mass_passing: '0'
        })
      } else {
        this.sample.sieves[pan].size = 'Pan'
      }

      // convert numbers to strings for use in html forms
      this.sample.sieves.forEach((o) => {
        this.toStrings(o)
      })
    }
  },
  created () {
    this.fetchTest()
  }
}
</script>

<style>
.table-heading {
  text-transform: none;
}
</style>
