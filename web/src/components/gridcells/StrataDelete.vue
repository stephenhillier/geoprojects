<template>
<div>
  <b-badge :id="`strataDelete${params.data.id}`" href="#" variant="light"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></b-badge>
    <b-popover :target="`strataDelete${params.data.id}`"
              placement="top"
              :show.sync="popoverShow"
              title="Delete strata?"
              triggers="click blur">
              <b-btn @click="onDeleteRow(params.data.id)" size="sm" variant="outline-danger">Ok</b-btn>
              <b-btn @click="onClose" size="sm" variant="primary" class="ml-2">Cancel</b-btn>

  </b-popover>
</div>

</template>

<script>
import Vue from 'vue'
export default Vue.extend({
  data () {
    return {
      popoverShow: false
    }
  },
  methods: {
    onDeleteRow (id) {
      this.params.context.componentParent.onDeleteStrataRow(id)
      this.onClose()
    },
    onClose () {
      this.popoverShow = false
    }
  }
})
</script>

<style>

</style>
