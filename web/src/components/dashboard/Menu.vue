<template>
  <b-card>
    <nav>
      <div>
        <ul class="nav flex-column">
          <b-nav-item class="menu-item" exact :to="{ name: 'projects' }"><font-awesome-icon class="text-muted mr-3" :icon="['fas', 'th-list']"></font-awesome-icon>Back to project list</b-nav-item>
          <b-nav-item active-class="active-menu" class="menu-item mt-3" exact :to="{ name: 'project-dashboard', params: { id: this.$route.params.id } }">
              Project: {{ project.name }}
          </b-nav-item>
          <b-nav-item active-class="active-menu" class="menu-item" :to="{ name: 'project-boreholes', params: { id: this.$route.params.id } }">
            <font-awesome-icon :icon="['far', 'dot-circle']" class="text-muted mr-3"></font-awesome-icon>Boreholes<b-badge variant="primary" class="float-right mt-1" pill>{{ project.borehole_count }}</b-badge>
          </b-nav-item>
          <b-nav-item active-class="active-menu" class="menu-item" :to="{ name: 'lab-home', params: { id: this.$route.params.id } }">
            <font-awesome-icon :icon="['far', 'clipboard']" class="text-muted mr-3"></font-awesome-icon>Lab tests<b-badge variant="primary" pill></b-badge>
          </b-nav-item>
        </ul>
      </div>
    </nav>
  </b-card>
</template>

<script>
export default {
  name: 'DashboardMenu',
  props: ['project']
}
</script>

<style>
  .menu-heading a {
    color: rgb(61, 61, 61)!important;
  }
  .menu-item a {
    color: rgb(61, 61, 61);
  }
</style>
