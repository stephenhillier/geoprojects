<template>
  <div class="box is-mapbox is-fullheight">
    <div id="map" class="project-map"></div>
  </div>
</template>

<script>
import L from 'leaflet'
import 'leaflet.markercluster/dist/MarkerCluster.css'
import 'leaflet.markercluster/dist/MarkerCluster.Default.css'
import 'leaflet.markercluster/dist/leaflet.markercluster.js'

export default {
  name: 'Map',
  props: {
    locations: {
      type: Array,
      default: () => []
    }
  },
  data () {
    return {
      map: null,
      markers: [],
      cluster: null
    }
  },
  computed: {
    filteredLocations () {
      return this.locations.filter((point) => {
        return ((!!point.location[0] || point.location[0] === 0) && (!!point.location[1] || point.location[1] === 0))
      })
    },
    centroid () {
      if (this.filteredLocations.length === 0) {
        return L.latLng(49, -123)
      }

      let latSum = 0
      let lngSum = 0
      for (let i = 0; i < this.filteredLocations.length; i++) {
        latSum = latSum + Number(this.filteredLocations[i].location[0])
        lngSum = lngSum + Number(this.filteredLocations[i].location[1])
      }

      const latAvg = latSum / this.filteredLocations.length
      const lngAvg = lngSum / this.filteredLocations.length

      return L.latLng(latAvg, lngAvg)
    }
  },
  methods: {
    /**
     * Methods borrow from github.com/bcgov/gwells (the Government of the Province of British Columbia)
     * GWELLS source code is made available under the Apache 2.0 License
     */
    initLeaflet () {
      delete L.Icon.Default.prototype._getIconUrl
      L.Icon.Default.mergeOptions({
        iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
        iconUrl: require('leaflet/dist/images/marker-icon.png'),
        shadowUrl: require('leaflet/dist/images/marker-shadow.png')
      })
    },
    initMap () {
      this.map = L.map('map').setView([this.centroid.lat, this.centroid.lng], 3)
      const osmAttrib = 'Map data © <a href="https://openstreetmap.org">OpenStreetMap</a> contributors'
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png?', { attribution: osmAttrib }).addTo(this.map)
    },
    createMarkers (latlng) {
      // clear markers

      if (!this.cluster) {
        this.cluster = L.markerClusterGroup()
        this.map.addLayer(this.cluster)
      }

      this.cluster.clearLayers()
      const markers = this.filteredLocations.map((item) => {
        return L.marker(L.latLng(item.location[0], item.location[1])).bindPopup(item.name)
      })
      this.cluster.addLayers(markers)
    }
  },
  watch: {
    centroid () {
      if (this.map) {
        this.createMarkers()
        this.map.panTo(this.centroid)
      }
    }
  },
  created () {
    this.$nextTick(function () {
      this.initLeaflet()
      this.initMap()
      this.createMarkers()
    })
  }

}
</script>

<style lang="scss">
.map-card {
  height: 100%;
}
.project-map {
  height: 100%;
  z-index: 1;
}
</style>
