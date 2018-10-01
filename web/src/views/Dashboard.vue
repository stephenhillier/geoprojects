<template>
  <b-container fluid id="dashboard" class="mb-3">
    <b-row>
      <b-col cols="12" lg="3" xl="2" class="mt-3">
        <geo-menu></geo-menu>
      </b-col>
      <b-col class="mt-3">
        <router-view></router-view>
      </b-col>
      <b-col cols="12" lg="3" xl="2" class="d-none d-xl-block mt-3">
        <div>
          <b-card title="Actions">
            <div><router-link :to="{name: 'newProject'}">New project</router-link></div>
          </b-card>
        </div>
        <div>
          <b-card title="Latest activity" class="mt-3">

          </b-card>
        </div>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import ProjectList from '@/components/dashboard/ProjectList.vue'
import Menu from '@/components/dashboard/Menu.vue'
export default {
  name: 'Dashboard',
  components: {
    ProjectList,
    geoMenu: Menu
  }
}
</script>

<style>

</style>
