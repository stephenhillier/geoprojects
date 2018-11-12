<template>
  <b-card>

      <b-row class="mb-3">
        <b-col cols="12" xl="6">
                <h4>Boreholes</h4>
                <h6 class="text-muted">{{project.name}}</h6>
        </b-col>
        <b-col>
          <multi-marker-map :locations="locations"></multi-marker-map>
        </b-col>
      </b-row>
    <b-table
      id="boreholeSearchTable"
      ref="boreholeSearchTable"
      :busy.sync="isBusy"
      responsive
      :items="boreholeSearch"
      :fields="fields"
      :per-page="perPage"
      :current-page="currentPage"
      >
      <template slot="name" slot-scope="data">
        <router-link :to="{ name: 'borehole-detail', params: { bh: data.item.id }}">{{ data.item.name }}</router-link>
      </template>
      <template slot="location" slot-scope="data">
        {{ data.value[0].toFixed(6) }}, {{ data.value[1].toFixed(6) }}
      </template>
    </b-table>
    <b-btn variant="info" size="sm" :to="{ name: 'new-borehole' }">New borehole</b-btn>

  </b-card>
</template>

<script>
import querystring from 'querystring'
import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'
export default {
  name: 'Boreholes',
  props: ['project'],
  components: {
    MultiMarkerMap
  },
  data () {
    return {
      locations: [],
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      fields: ['name', 'start_date', 'end_date', 'field_eng', 'location']
    }
  },
  methods: {
    boreholeSearch (ctx = { perPage: this.perPage, currentPage: this.currentPage }) {
      /**
      * table items provider function
      * https://bootstrap-vue.js.org/docs/components/table/
      *
      * a refresh can be triggered by this.$root.$emit('bv::refresh::table', 'boreholeSearchTable')
      */

      const params = {
        project: this.$route.params.id,
        limit: ctx.perPage,
        offset: ctx.perPage * (ctx.currentPage - 1)
      }

      // add other search parameters into the params object.
      // these will be urlencoded and the API will filter on these values.
      Object.assign(params, this.searchParams)

      return this.$http.get('boreholes' + '?' + querystring.stringify(params)).then((response) => {
        this.numberOfRecords = response.data.count
        this.locations = response.data.results || []
        return response.data.results || []
      }).catch((e) => {
        return []
      })
    }
  },
  created () {
    this.$http.get(`boreholes?project=${this.$route.params.id}`).then((response) => {
      this.boreholes = response.data
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style>

</style>
