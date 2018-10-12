import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'

import ProjectDashboard from './views/ProjectDashboard.vue'
import ProjectList from './views/ProjectList.vue'
import NewProject from './components/dashboard/projects/NewProject.vue'
import ProjectDetails from './components/dashboard/projects/ProjectDetails.vue'

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
      name: 'project-dashboard',
      component: ProjectDashboard,
      children: [
        {
          path: '/',
          name: 'project-details',
          component: ProjectDetails
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
