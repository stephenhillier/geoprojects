<template>
<div class="columns content-wrapper is-fullheight">
  <div class="column is-narrow project-list is-hidden-mobile">
    <div class="box project-list-menu-panel">
      <nav class="menu">
        <p class="menu-label">
          Projects
        </p>
        <ul class="menu-list">
          <li class="project-list-menu-item"><router-link :to="{name: 'projects'}">
            <font-awesome-icon :icon="['fas', 'th-list']"></font-awesome-icon> Project List
          </router-link></li>
          <li class="project-list-menu-item"><router-link :to="{name: 'new-project'}">
            <font-awesome-icon :icon="['far', 'plus-square']"></font-awesome-icon> New Project
          </router-link></li>
        </ul>
      </nav>
      <b-field class="project-list-search">
        <b-input placeholder="Search..."
            type="search"
            @input="handleSearchInput"
            v-model="searchParamsInput.search"
            icon="search">
        </b-input>
      </b-field>
      <b-table
          :data="projects"
          :columns="fields"
          paginated
          :per-page="perPage"
          :current-page.sync="currentPage"
      >
        <template slot-scope="props">
            <b-table-column field="project" label="Project">
                <router-link :to="`/projects/${props.row.id}`">{{ props.row.number ? `${props.row.number} -` : '' }} {{ props.row.name }}</router-link>
            </b-table-column>
            <b-table-column field="location" label="Location">
                {{ props.row.location }}
            </b-table-column>
            <b-table-column field="borehole_count" label="Boreholes">
                {{ props.row.borehole_count }}
            </b-table-column>
        </template>
      </b-table>
      <div class="project-helper">
        <b-message title="Navigating projects" type="is-info" v-if="projectTutorial">
              <h4 class="help-header">Welcome!</h4>
              <p class="help-text">This is the project list, where you can search for and select a project.</p>
              <p class="help-text">We'll show you how to create projects. To get started, click on the <router-link :to="{name: 'new-project'}">New project</router-link> option on the left sidebar.</p>
              <div style="display: flex; flex-direction: row;"><a href="#" @click="handleCancelProjectTutorial" style="margin-left: auto;">Don't show again</a></div>
        </b-message>
      </div>
    </div>

  </div>
  <div class="column is-fullheight project-map">
    <multi-marker-map :locations="locations"></multi-marker-map>
  </div>
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
      projectTutorial: false,
      projects: [],
      locations: [],
      loading: false,
      fields: [
        {
          field: 'project',
          label: 'Project'
        },
        {
          field: 'location',
          label: 'Location'
        },
        {
          field: 'borehole_count',
          label: 'Boreholes'
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
      this.projectTutorial = false

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
        this.projectTutorial = true
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
.project-list {
  height: 100%;
}
.project-map {
  height: 100%;
}
.project-list-menu-item {
  margin-top: 12px;
}
.project-list-search {
  margin-top: 25px;
}
.project-list-menu-panel {
  width: 600px;
  height: 100%!important;
  display: flex;
  flex-direction: column;
  margin-left: 20px;
}
.project-helper {
  margin-top: auto;
}
.help-text {
  margin-top: 10px;
}
.help-header {
  font-weight: 600;
}
.content-wrapper {
  height: 100%;
}
</style>
