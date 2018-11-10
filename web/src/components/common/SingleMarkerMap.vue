<template>
  <div id="map" class="map"></div>
</template>

<script>
import L from 'leaflet'

export default {
  name: 'Map',
  props: {
    latitude: null,
    longitude: null,
    addMode: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      map: null,
      marker: null
    }
  },
  computed: {
    lat () {
      if (this.latitude) {
        const lat = Number(this.latitude)
        if (!Number.isNaN(lat)) {
          return lat
        }
      }
      return 49
    },
    long () {
      if (this.longitude) {
        const long = Number(this.longitude)
        if (!Number.isNaN(long)) {
          return long
        }
      }
      return -123
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

      // we only need double click to add marker functionality when not in multiple mode
      if (this.addMode) {
        this.map.on('dblclick', (e) => {
          this.updateMarker(e.latlng)
          this.$emit('update-coordinates', e.latlng)
        })
      }
    },
    handleDrag (ev) {
      const markerLatLng = this.marker.getLatLng()
      this.$emit('update-coordinates', markerLatLng)
    },
    createMarker (latlng) {
      if (latlng || (this.latitude && this.longitude)) {
        latlng = latlng || L.latLng(this.lat, this.long)
        this.marker = L.marker(latlng, { draggable: this.addMode, autoPan: true })

        if (this.addMode) {
          this.marker.on('dragend', (ev) => {
            this.handleDrag(ev)
          })
        }

        this.marker.addTo(this.map)
        this.map.setView(latlng, 13)
      }
    },
    updateMarker (latlng) {
      if (this.marker) {
        this.marker.setLatLng(latlng)
        this.map.panTo(latlng)
      } else {
        this.createMarker(latlng)
      }
    }
  },
  watch: {
    latitude () {
      this.updateMarker(L.latLng(this.lat, this.long))
    },
    longitude () {
      this.updateMarker(L.latLng(this.lat, this.long))
    }
  },
  created () {
    this.$nextTick(function () {
      this.initLeaflet()
      this.initMap()
      this.createMarker()
    })
  }

}
</script>

<style lang="scss">
.map {
  height: 400px;
}
</style>
