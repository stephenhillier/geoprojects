<template>
  <section class="section">
    <h2 class="subtitle">Files</h2>

    <b-table
        :data="filteredFiles"
        paginated
        :per-page="perPage"
        :current-page.sync="currentPage"
    >
      <template slot-scope="props">
          <b-table-column field="filename" label="Filename">
            <a href="#">
              <span :class="props.row.superseded ? 'file-superseded' : ''">{{ props.row.filename }} {{ props.row.superseded ? '[superseded]' : '' }}</span>
            </a>
          </b-table-column>
          <b-table-column field="created_at" label="Uploaded">
            {{ props.row.created_at | moment("dddd, MMMM Do YYYY, h:mm:ss a") }}
          </b-table-column>
          <b-table-column field="actions" label="Actions" class="is-narrow">
            <button class="button is-small"><font-awesome-icon :icon="['fas', 'download']"></font-awesome-icon></button>
            <button class="button is-small ml"><font-awesome-icon :icon="['far', 'edit']"></font-awesome-icon></button>
            <button class="button is-small ml" @click="handleDelete(props.row.id)"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></button>
          </b-table-column>
      </template>
    </b-table>

    <form @submit.prevent="handleFileUpload">
      <fieldset :disabled="loading" :class="loading ? 'upload-progress' : ''">
        <b-field>
            <b-upload v-model="dropFiles"
                multiple
                drag-drop>
                <section class="section">
                    <div class="content has-text-centered">
                        <p>
                            <b-icon
                                icon="upload"
                                size="is-large">
                            </b-icon>
                        </p>
                        <p>Drop your files here or click to upload</p>
                    </div>
                </section>
            </b-upload>
        </b-field>

        <div class="tags">
          <transition-group name="list" tag="span">
              <span v-for="(file, index) in dropFiles"
                  :key="index"
                  class="tag is-primary list-item" >
                  {{file.name}}
                  <button class="delete is-small"
                      type="button"
                      @click="deleteDropFile(index)">
                  </button>
              </span>
            </transition-group>
        </div>
        <b-field><!-- Label left empty for spacing -->
            <p class="control">
                <button type="submit" class="button is-primary" :disabled="!dropFiles.length">
                  Upload files
                </button>
                <font-awesome-icon v-if="loading" :icon="['fas', 'spinner']" spin pulse class="loading-spinner" size="lg"></font-awesome-icon>
            </p>
        </b-field>
      </fieldset>
    </form>
  </section>
</template>

<script>
export default {
  name: 'ProjectFiles',
  props: ['project', 'files'],
  data () {
    return {
      dropFiles: [],
      loading: false,
      showSuperseded: false,
      perPage: 10,
      currentPage: 1
    }
  },
  computed: {
    filteredFiles () {
      const files = this.files || []
      if (!this.showSuperseded) {
        return files.filter((f) => {
          return f.superseded === false
        })
      }
      return files
    }
  },
  methods: {
    handleFileUpload () {
      this.loading = true

      const formData = new FormData()
      formData.append('file', this.dropFiles[0])
      const config = {
        headers: {
          'content-type': 'multipart/form-data'
        }
      }

      this.$http.post(`projects/${this.$route.params.id}/files`, formData, config).then((response) => {
        this.dropFiles = []
        this.$noty.success('Files successfully uploaded.')
        this.$emit('updated', true)
      }).catch((e) => {
        this.$noty.error('Error uploading files. Please try again later.')
      }).finally(() => {
        this.loading = false
      })
    },
    // delete a file staged for upload
    deleteDropFile (i) {
      this.dropFiles.splice(i, 1)
    },
    // delete a file on the server
    handleDelete (id) {
      console.log(id)
    }
  }
}
</script>

<style>
.list-item {
  transition: all 1s;
  display: inline-block;
  margin-right: 10px;
}
.list-enter, .list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
.list-leave-active {
  position: absolute;
}
.upload-progress {
  opacity: 0.6;
}
.loading-spinner {
  margin-left: 1rem;
}
.file-superseded {
  opacity: 0.6;
}
</style>
