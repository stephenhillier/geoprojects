import Vue from 'vue'
import Router from 'vue-router'

import Login from './views/Login.vue'

import ProjectList from './views/ProjectList.vue'
import NewProject from './components/dashboard/projects/NewProject.vue'

import ProjectDashboard from './views/ProjectDashboard.vue'
import ProjectDetails from './components/dashboard/projects/ProjectDetails.vue'
import ProjectDetailActions from './components/dashboard/actions/ProjectDetailActions.vue'

import ProjectBoreholes from './components/dashboard/boreholes/ProjectBoreholes.vue'
import NewBorehole from './components/dashboard/boreholes/NewBorehole.vue'
import BoreholeDetail from './components/dashboard/boreholes/BoreholeDetail.vue'
import BoreholeDetailActions from './components/dashboard/actions/BoreholeDetailActions.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/projects/new',
      name: 'new-project',
      component: NewProject
    },
    {
      path: '/projects/:id',
      component: ProjectDashboard,
      children: [
        {
          path: 'boreholes/new',
          name: 'new-borehole',
          component: NewBorehole,
          meta: {
            breadcrumbs: [
              {
                text: 'Boreholes',
                to: { name: 'project-boreholes' }
              },
              {
                text: 'New borehole',
                to: { name: 'new-project' }
              }
            ]
          }
        },
        {
          path: 'boreholes/:bh',
          name: 'borehole-detail',
          components: {
            default: BoreholeDetail,
            actions: BoreholeDetailActions
          },
          meta: {
            breadcrumbs: [
              {
                text: 'Boreholes',
                to: { name: 'project-boreholes' }
              },
              {
                text: 'Borehole Details',
                to: { name: 'project-boreholes' }
              }
            ]
          }
        },
        {
          path: 'boreholes',
          name: 'project-boreholes',
          component: ProjectBoreholes,
          meta: {
            breadcrumbs: [
              {
                text: 'Boreholes',
                to: { name: 'project-boreholes' }
              }
            ]
          }
        },
        {
          path: '/',
          name: 'project-dashboard',
          components: {
            default: ProjectDetails,
            actions: ProjectDetailActions
          },
          meta: {
            breadcrumbs: []
          }
        }
      ]
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/',
      name: 'projects',
      component: ProjectList
    }

    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (about.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import(/* webpackChunkName: "about" */ './views/About.vue')
    // }
  ]
})
