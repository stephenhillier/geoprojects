<template>
<div>
  <b-row>
    <b-col>
      <b-breadcrumb :items="breadcrumbs"></b-breadcrumb>
    </b-col>
  </b-row>
  <b-row>
    <b-col cols="12" md="3" lg="2" xl="2">
      <b-card no-body class="mb-3">
        <b-list-group flush>
          <b-list-group-item exact :to="{name: 'projects'}">Project List</b-list-group-item>
          <b-list-group-item exact :to="{name: 'new-project'}">New Project</b-list-group-item>
        </b-list-group>
      </b-card>
      <b-card class="my-3">
        <div class="card-title">
          <span class="h4">Search</span>
        </div>
        <b-form @submit.prevent="onSearchHandler">
          <div>
            <form-input id="projectSearchNumber" label="Project number" v-model="searchParamsInput.project_number"></form-input>
          </div>
          <div>
            <form-input id="projectSearchName" label="Project name" v-model="searchParamsInput.project_name"></form-input>
          </div>
          <b-btn type="submit">Search</b-btn>
        </b-form>
      </b-card>
    </b-col>
    <b-col cols="12" md="6" lg="8" xl="8">
      <b-card class="mb-3">
        <div class="card-title">
          <h1 class="h1">Projects</h1>
          <b-row class="my-3">
            <b-col>
              <multi-marker-map :locations="locations"></multi-marker-map>
            </b-col>
          </b-row>

          <div class="float-right">Filters:

            <span
              v-for="(value, key) in searchParams"
              :key="`searchFilterChip${key}`"
            >
              <b-badge
                  variant="info"
                  v-if="value"
                  class="ml-2 pr-1"
                  pill
                  :id="`searchFilterChip${key}`"
              >
                {{ key | readable }}: {{ value }}
                <a href="#" @click="clearSearchFilter(key)" class="text-white"><font-awesome-icon :icon="['far', 'times-circle']" size="lg" class="m-0 p-0 ml-2"></font-awesome-icon></a>
              </b-badge>
            </span>
          </div>
        </div>
        <div>
          <b-table
            id="projectSearchTable"
            ref="projectSearchTable"
            :busy.sync="isBusy"
            responsive
            :items="projectSearch"
            :fields="fields"
            :per-page="perPage"
            :current-page="currentPage"
            show-empty
            >
            <template slot="project" slot-scope="data">
              <router-link :to="{ name: 'project-dashboard', params: { id: data.item.id }}">{{data.item.id}} - {{ data.item.name }}</router-link>
            </template>
          </b-table>

          <div>
            <b-pagination :disabled="isBusy" size="md" :total-rows="numberOfRecords" v-model="currentPage" :per-page="perPage"></b-pagination>
          </div>
        </div>
      </b-card>
    </b-col>
    <b-col cols="12" md="3" lg="2" xl="2">
      <b-card title="Actions" class="mb-3">
        <b-row class="mt-2">
          <b-col>
            <b-btn variant="link" size="sm" :to="{name: 'new-project'}">
              <font-awesome-icon :icon="['far', 'plus-square']" class="text-muted"></font-awesome-icon>
              New project
            </b-btn>
          </b-col>
        </b-row>
      </b-card>
    </b-col>
  </b-row>
</div>

</template>

<script>
import querystring from 'querystring'
import FormInput from '@/components/common/FormInput.vue'
import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'

export default {
  name: 'ProjectList',
  components: {
    FormInput,
    MultiMarkerMap
  },
  data () {
    return {
      projects: [],
      locations: [],
      loading: false,
      fields: [ 'project', 'location', 'borehole_count' ],
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      searchParamsInput: {
        project_number: null,
        project_name: null
      },
      searchParams: {},
      breadcrumbs: [
        {
          text: 'Projects',
          to: { name: 'projects' }
        }
      ]
    }
  },
  methods: {
    projectSearch (ctx = { perPage: this.perPage, currentPage: this.currentPage }) {
      /**
      * projectSearch() is a table items provider function
      * https://bootstrap-vue.js.org/docs/components/table/
      *
      * a refresh can be triggered by this.$root.$emit('bv::refresh::table', 'projectSearchTable')
      */

      const params = {
        limit: ctx.perPage,
        offset: ctx.perPage * (ctx.currentPage - 1)
      }

      this.locations = []

      // add other search parameters into the params object.
      // these will be urlencoded and the API will filter on these values.
      Object.assign(params, this.searchParams)

      return this.$http.get('projects' + '?' + querystring.stringify(params)).then((response) => {
        this.numberOfRecords = response.data.count
        response.data.results.forEach((item) => {
          if (item.centroid[0] !== 0 && item.centroid[1] !== 0) {
            this.locations.push({ name: item.name, location: item.centroid })
          }
        })
        return response.data.results || []
      }).catch((e) => {
        return []
      })
    },
    onSearchHandler () {
      this.searchParams = Object.assign({}, this.searchParamsInput)
      this.$root.$emit('bv::refresh::table', 'projectSearchTable')
    },
    clearSearchFilter (key) {
      this.searchParams[key] = null
      this.searchParamsInput[key] = ''
      this.$root.$emit('bv::refresh::table', 'projectSearchTable')
    }
  }
}
</script>

<style>

</style>
