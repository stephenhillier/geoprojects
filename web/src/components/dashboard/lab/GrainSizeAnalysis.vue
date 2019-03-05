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
      <b-col cols="12" md="4">
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
      </b-col>
      <b-col cols="12" md="6" offset-md="1">
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
                  <form-input :id="`mass${i}`" groupClass="p-0 m-0" v-model="sample.sieves[i].massg"></form-input>
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
  </div>
</template>

<script>
export default {
  name: 'GrainSizeAnalysis',
  props: ['project'],
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
            mass: ''
          }
        ]
      },
      loading: false
    }
  },
  methods: {
    fetchTest () {
      this.loading = true
      this.$http.get(`projects/${this.$route.params.id}/lab/tests/${this.$route.params.test}`).then((response) => {
        this.sample = this.toStrings(this.nullToStrings(response.data))
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

      this.sample.sieves.splice(
        this.sample.sieves.length - 2,
        0,
        {
          pan: false,
          size: '',
          mass_passing: ''
        }
      )
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
    },
    initializeSieveArray () {
      if (!this.sample.sieves || !~this.sample.sieves.findIndex((i) => {
        return i.pan === true
      })) {
        this.sample.sieves = [
          {
            pan: true,
            size: 'Pan',
            mass_passing: ''
          }
        ]
      }
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
