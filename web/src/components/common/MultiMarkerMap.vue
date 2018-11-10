<template>
  <div id="map" class="map"></div>
</template>

<script>
import L from 'leaflet'

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
      markers: []
    }
  },
  computed: {
    filteredLocations () {
      return this.locations.filter((point) => {
        return ((!!point.latitude || point.latitude === 0) && (!!point.longitude || point.longitude === 0))
      })
    },
    centroid () {
      if (this.filteredLocations.length === 0) {
        return L.latLng(49, -123)
      }

      let latSum = 0
      let lngSum = 0
      for (let i = 0; i < this.filteredLocations.length; i++) {
        latSum = latSum + Number(this.filteredLocations[i].latitude)
        lngSum = lngSum + Number(this.filteredLocations[i].longitude)
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
      this.map = L.map('map').setView([this.lat, this.long], 5)

      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png?').addTo(this.map)
    },
    createMarkers (latlng) {
      this.filteredLocations.forEach((item) => {
        this.markers.push({

        })
      })
    }
  },
  watch: {
    centroid () {
      this.markers = []
      this.createMarkers()
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
.map {
  height: 400px;
}
</style>
