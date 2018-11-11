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
      this.map = L.map('map').setView([this.centroid.lat, this.centroid.lng], 7)
      const osmAttrib = 'Map data © <a href="https://openstreetmap.org">OpenStreetMap</a> contributors'
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png?', { attribution: osmAttrib }).addTo(this.map)
    },
    createMarkers (latlng) {
      this.filteredLocations.forEach((item) => {
        const loc = L.latLng(item.location[0], item.location[1])
        const marker = L.marker(loc)
        this.markers.push(marker)
        marker.addTo(this.map)
        marker.bindPopup(item.name)
      })
    }
  },
  watch: {
    centroid () {
      this.markers = []
      this.createMarkers()
      this.map.panTo(this.centroid)
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
