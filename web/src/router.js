import Vue from 'vue'
import Router from 'vue-router'

import Login from './views/Login.vue'
import LoginCallback from '@/components/common/AuthCallback.vue'

import ProjectList from './views/ProjectList.vue'
import NewProject from './components/dashboard/projects/NewProject.vue'

import ProjectDashboard from './views/ProjectDashboard.vue'
import ProjectDetails from './components/dashboard/projects/ProjectDetails.vue'
import ProjectDetailActions from './components/dashboard/actions/ProjectDetailActions.vue'

import ProjectBoreholes from './components/dashboard/boreholes/ProjectBoreholes.vue'
import NewBorehole from './components/dashboard/boreholes/NewBorehole.vue'
import BoreholeDetail from './components/dashboard/boreholes/BoreholeDetail.vue'
import BoreholeDetailActions from './components/dashboard/actions/BoreholeDetailActions.vue'
import SamplesHome from './components/dashboard/lab/SamplesHome.vue'

import InstrumentationHome from './components/dashboard/instrumentation/InstrumentationHome.vue'
import NewInstrument from './components/dashboard/instrumentation/NewInstrument.vue'
import InstrumentationDetails from './components/dashboard/instrumentation/InstrumentationDetails.vue'

import LabTestingHome from './components/dashboard/lab/LabTestingHome.vue'
import MoistureContent from './components/dashboard/lab/MoistureContent.vue'
import GrainSizeAnalysis from './components/dashboard/lab/GrainSizeAnalysis.vue'

Vue.use(Router)

const guard = (to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // this route requires auth, check if logged in
    // if not, redirect to login page.
    if (!router.app.$auth.isAuthenticated()) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next() // make sure to always call next()!
  }
}

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/projects/new',
      name: 'new-project',
      component: NewProject,
      meta: { requiresAuth: true }
    },
    {
      path: '/projects/:id',
      component: ProjectDashboard,
      meta: { requiresAuth: true },
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
          path: 'lab/grainsize/:test',
          name: 'lab-grainsize',
          component: GrainSizeAnalysis,
          meta: {
            breadcrumbs: [
              {
                text: 'Lab testing',
                to: { name: 'lab-home' }
              },
              {
                text: 'Grain size analysis',
                to: { name: 'lab-grainsize' }
              }
            ]
          }
        },
        {
          path: 'lab/moisture/:test',
          name: 'lab-moisture',
          component: MoistureContent,
          meta: {
            breadcrumbs: [
              {
                text: 'Lab testing',
                to: { name: 'lab-home' }
              },
              {
                text: 'Moisture content',
                to: { name: 'lab-moisture' }
              }
            ]
          }
        },
        {
          path: 'lab',
          name: 'lab-home',
          component: LabTestingHome,
          meta: {
            breadcrumbs: [
              {
                text: 'Lab testing',
                to: { name: 'lab-home' }
              }
            ]
          }
        },
        {
          path: 'samples',
          name: 'samples-home',
          component: SamplesHome,
          meta: {
            breadcrumbs: [
              {
                text: 'Samples',
                to: { name: 'samples-home' }
              }
            ]
          }
        },
        {
          path: 'instrumentation/:instr',
          name: 'instrumentation-details',
          component: InstrumentationDetails,
          meta: {
            breadcrumbs: [
              {
                text: 'Instrumentation',
                to: { name: 'instrumentation-home' }
              },
              {
                text: 'Instrument summary',
                to: { name: 'instrumentation-details' }
              }
            ]
          }
        },
        {
          path: 'instrumentation/new',
          name: 'instrumentation-new',
          component: NewInstrument,
          meta: {
            breadcrumbs: [
              {
                text: 'Instrumentation',
                to: { name: 'instrumentation-home' }
              },
              {
                text: 'New instrument',
                to: { name: 'instrumentation-new' }
              }
            ]
          }
        },
        {
          path: 'instrumentation',
          name: 'instrumentation-home',
          component: InstrumentationHome,
          meta: {
            breadcrumbs: [
              {
                text: 'Instrumentation',
                to: { name: 'instrumentation-home' }
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
      path: '/callback',
      name: 'login-callback',
      component: LoginCallback
    },
    {
      path: '/',
      name: 'projects',
      meta: { requiresAuth: true },
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

router.beforeEach(guard)

export default router
