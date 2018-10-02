import Vue from 'vue'
import Router from 'vue-router'
import Dashboard from './views/Dashboard.vue'
import ProjectList from './components/dashboard/projects/ProjectList.vue'
import NewProject from './components/dashboard/projects/NewProject.vue'
import ProjectDetails from './components/dashboard/projects/ProjectDetails.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: Dashboard,
      children: [
        {
          path: 'projects/new',
          name: 'newProject',
          component: NewProject
        },
        {
          path: 'projects/:id',
          name: 'project-details',
          component: ProjectDetails
        },
        {
          path: 'projects',
          name: 'projects',
          component: ProjectList
        }
      ]
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
