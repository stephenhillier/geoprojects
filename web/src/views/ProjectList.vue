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
    </b-col>
    <b-col cols="12" md="6" lg="8" xl="8">
      <b-card class="mb-3" no-body>
        <b-row class="no-gutters">
          <b-col>
            <multi-marker-map :locations="locations"></multi-marker-map>
          </b-col>
        </b-row>
      </b-card>
      <b-card no-body class="mb-3">
        <ag-grid-vue style="height: 500px;"
            :enableSorting="true"
            :enableFilter="true"
            rowHeight="32"
            class="ag-theme-balham"
            :columnDefs="columnDefs"
            :rowData="projects"/>
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
import { AgGridVue } from 'ag-grid-vue'

import FormInput from '@/components/common/FormInput.vue'
import MultiMarkerMap from '@/components/common/MultiMarkerMap.vue'
import NameWithLink from '@/components/gridcells/NameWithLink.vue'

export default {
  name: 'ProjectList',
  components: {
    FormInput,
    MultiMarkerMap,
    AgGridVue
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
      ],
      columnDefs: [
        { headerName: 'Project', field: 'name', filter: 'agTextColumnFilter', cellRendererFramework: NameWithLink, colId: 'params' },
        { headerName: 'Location', field: 'location', filter: 'agTextColumnFilter' },
        { headerName: 'Boreholes', field: 'borehole_count' }
      ]
    }
  },
  methods: {
    fetchProjects () {
      this.$http.get('projects').then((response) => {
        this.numberOfRecords = response.data.length
        response.data.forEach((item) => {
          this.locations.push({ name: item.name, location: (item.centroid[0] === 0 && item.centroid[1] === 0) ? item.default_coords : item.centroid })
        })
        this.projects = response.data || []
      }).catch((e) => {
        this.$noty.error('An error occurred while retrieving projects.')
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
  },
  created () {
    this.fetchProjects()
  }
}
</script>

<style>

</style>
