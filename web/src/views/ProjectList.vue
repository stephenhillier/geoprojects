<template>
<div>
  <b-row>
    <b-col cols="12" md="3" lg="2" xl="2">
      <b-card>
        <nav>
          <div>
            <ul class="nav flex-column">
              <b-nav-item active-class="active-menu" class="menu-item mt-3" exact :to="{name: 'projects'}">
                <font-awesome-icon :icon="['fas', 'th-list']" class="text-muted mr-3"></font-awesome-icon>Project List
              </b-nav-item>
              <b-nav-item active-class="active-menu" class="menu-item" exact :to="{name: 'new-project'}">
                <font-awesome-icon :icon="['far', 'plus-square']" class="text-muted mr-3"></font-awesome-icon>New Project
              </b-nav-item>
            </ul>
          </div>
        </nav>
      </b-card>
    </b-col>
    <b-col>
      <b-row>
        <b-col>
          <b-card class="mb-3">
            <div class="card-status bg-blue"></div>
            <b-row>
              <b-col>
                <h1>Projects</h1>
              </b-col>
              <b-col>
                <div class="form-group horizontal">
                  <label class="form-label sr-only">Search for projects</label>
                  <div class="input-icon mb-3">
                    <span class="input-icon-addon">
                      <font-awesome-icon :icon="['fas', 'search']" class="text-dark"></font-awesome-icon>
                    </span>
                    <input type="text" class="form-control" placeholder="Search by name or project number">
                  </div>
                </div>
              </b-col>
            </b-row>
            <b-row class="no-gutters">
              <b-col>
                <multi-marker-map :locations="locations"></multi-marker-map>
              </b-col>
            </b-row>
            <b-row>
              <b-col>
                <ag-grid-vue style="height: 400px;"
                    :enableSorting="true"
                    :enableFilter="true"
                    rowHeight="32"
                    class="ag-theme-balham"
                    :columnDefs="columnDefs"
                    :rowData="projects"/>
              </b-col>
            </b-row>
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
