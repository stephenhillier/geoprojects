<template>
  <b-table
    :data="boreholes"
    :columns="fields"
    paginated
    :per-page="perPage"
    :current-page.sync="currentPage"
  >
    <template slot-scope="props">
      <b-table-column field="name" label="Borehole">
          <router-link :to="`/projects/${$route.params.id}/boreholes/${props.row.id}`">{{ props.row.name }}</router-link>
      </b-table-column>

      <b-table-column field="start_date" label="Start date">
          {{ props.row.start_date | moment('YYYY-MM-DD') }}
      </b-table-column>
      <b-table-column field="end_date" label="End date">
          {{ props.row.end_date | moment('YYYY-MM-DD') }}
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
  name: 'BoreholeTable',
  props: ['project', 'boreholes'],
  data () {
    return {
      currentPage: 1,
      perPage: 10,
      isBusy: false,
      numberOfRecords: 0,
      fields: [
        {
          field: 'name',
          label: 'Name'
        },
        {
          field: 'start_date',
          label: 'Start date'
        },
        {
          field: 'end_date',
          label: 'End date'
        },
        {
          field: 'field_eng',
          label: 'Drilled by'
        },
        {
          field: 'location',
          label: 'Location'
        }
      ]
    }
  }
}
</script>

<style>

</style>
