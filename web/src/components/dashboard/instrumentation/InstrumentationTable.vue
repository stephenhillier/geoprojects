<template>
  <b-table
    :data="instruments"
    :columns="fields"
    paginated
    :per-page="perPage"
    :current-page.sync="currentPage"
  >
    <template slot-scope="props">
      <b-table-column field="name" label="Instrument">
          <router-link :to="`/projects/${$route.params.id}/instrumentation/${props.row.id}`">{{ props.row.name }}</router-link>
      </b-table-column>

      <b-table-column field="device_id" label="Device ID">
          {{ props.row.device_id }}
      </b-table-column>
      <b-table-column field="install_date" label="Install date">
          {{ props.row.install_date | moment('YYYY-MM-DD') }}
      </b-table-column>
      <b-table-column field="field_eng" label="Field technician/engineer">
          {{ props.row.field_eng }}
      </b-table-column>
      <b-table-column field="location" label="Location">
                  {{ props.row.location[0] ? props.row.location[0].toFixed(6) : null }},
                  {{ props.row.location[1] ? props.row.location[1].toFixed(6) : null }}
      </b-table-column>
  </template>

</b-table>
</template>

<script>
export default {
  name: 'InstrumentationTable',
  props: ['project', 'instruments'],
  data () {
    return {
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      fields: []
    }
  }
}
</script>

<style>

</style>
