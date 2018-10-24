import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'

import ProjectList from './views/ProjectList.vue'
import NewProject from './components/dashboard/projects/NewProject.vue'

import ProjectDashboard from './views/ProjectDashboard.vue'
import ProjectDetails from './components/dashboard/projects/ProjectDetails.vue'
import ProjectDetailActions from './components/dashboard/actions/ProjectDetailActions.vue'

import ProjectBoreholes from './components/dashboard/boreholes/ProjectBoreholes.vue'
import NewBorehole from './components/dashboard/boreholes/NewBorehole.vue'

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
      path: '/projects',
      name: 'projects',
      component: ProjectList
    },
    {
      path: '/',
      name: 'home',
      component: Home
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
