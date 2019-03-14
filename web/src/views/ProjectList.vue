<template>
<div>
  <b-row>
    <b-col cols="12" lg="5" xl="3" style="height: 90vh" class="d-flex flex-column">
      <b-row class="mb-3">
        <b-col>
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
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <div class="form-group horizontal">
            <label class="form-label sr-only">Search for projects</label>
            <div class="input-icon mb-3">
              <span class="input-icon-addon">
                <font-awesome-icon :icon="['fas', 'search']" class="text-dark"></font-awesome-icon>
              </span>
              <input type="text" class="form-control" v-model="searchParamsInput.project_name" @input="handleSearchInput" placeholder="Search by project name">
            </div>
          </div>
        </b-col>
      </b-row>
      <b-row class="flex-grow-1">
        <b-col style="height: 100%">
          <!-- <ag-grid-vue style="height: 100%;"
              :enableSorting="true"
              :enableFilter="true"
              rowHeight="32"
              class="ag-theme-balham"
              :columnDefs="columnDefs"
              :rowData="projects"/> -->

          <b-table striped hover :items="projects" :fields="fields">
            <template slot="project" slot-scope="data">
              <router-link :to="`/projects/${data.item.id}`">{{ data.item.number ? `${data.item.number} - ` :'' }}{{ data.item.name }}</router-link>
            </template>
          </b-table>
          <b-pagination v-if="projects.length > perPage" v-model="currentPage" :total-rows="projects.length" :per-page="perPage" size="md" />
          <b-modal centered title="Creating a project" ref="tutorialProjectModal" cancel-title="Don't show again" @cancel="handleCancelProjectTutorial">
            <div class="d-block text-center">
              <h5>Welcome!</h5>
              <p>This is the project list, where you can search for and select a project.</p>
              <p>We'll show you how to create projects. To get started, click on the <router-link :to="{name: 'new-project'}">New project</router-link> option on the left sidebar.</p>
            </div>
          </b-modal>
        </b-col>
      </b-row>
    </b-col>
    <b-col style="height: 90vh;">
      <multi-marker-map :locations="locations"></multi-marker-map>
    </b-col>
  </b-row>
</div>

</template>

<script>
import debounce from 'lodash.debounce'
import querystring from 'querystring'
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
      fields: [
        {
          key: 'project',
          thClass: 'table-heading'
        },
        {
          key: 'location',
          thClass: 'table-heading'
        },
        {
          key: 'borehole_count',
          thClass: 'table-heading'
        }
      ],
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
        { headerName: 'Location', field: 'location', filter: 'agTextColumnFilter' }
        // { headerName: 'Boreholes', field: 'borehole_count' }
      ]
    }
  },
  methods: {
    fetchProjects () {
      this.$http.get('projects?' + querystring.stringify(this.searchParams)).then((response) => {
        this.locations = []
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
      this.fetchProjects()
    },
    clearSearchFilter (key) {
      this.searchParams[key] = null
      this.searchParamsInput[key] = ''
      this.$root.$emit('bv::refresh::table', 'projectSearchTable')
    },
    handleSearchInput () {
      this.debouncedSearch()
    },
    handleCancelProjectTutorial () {
      localStorage.setItem('earthworks-tutorial-projects', JSON.stringify(true))
    }
  },
  created () {
    this.fetchProjects()
    this.debouncedSearch = debounce(() => {
      this.onSearchHandler()
    }, 300)

    if (!JSON.parse(localStorage.getItem('earthworks-tutorial-projects'))) {
      setTimeout(() => {
        this.$refs.tutorialProjectModal.show()
      }, 1000)
    }
  }
}
</script>

<style>
.table-heading {
  color: #333!important;
  text-transform: none!important;
}
</style>
