<template>
  <b-card title="Boreholes" :sub-title="project.name">
    <b-table
      id="boreholeSearchTable"
      ref="boreholeSearchTable"
      :busy.sync="isBusy"
      :items="boreholeSearch"
      :fields="fields"
      :per-page="perPage"
      :current-page="currentPage"
      >
      <template slot="project" slot-scope="data">
        <router-link :to="{ name: 'project-dashboard', params: { id: data.item.id }}">{{data.item.id}} - {{ data.item.name }}</router-link>
      </template>
    </b-table>
    <b-btn variant="info" size="sm" :to="{ name: 'new-borehole' }">New borehole</b-btn>

  </b-card>
</template>

<script>
import querystring from 'querystring'
export default {
  name: 'Boreholes',
  props: ['project'],
  data () {
    return {
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      fields: ['name', 'start_date', 'end_date', 'field_eng']
    }
  },
  methods: {
    boreholeSearch (ctx = { perPage: this.perPage, currentPage: this.currentPage }) {
      /**
      * projectSearch() is a table items provider function
      * https://bootstrap-vue.js.org/docs/components/table/
      *
      * a refresh can be triggered by this.$root.$emit('bv::refresh::table', 'projectSearchTable')
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
