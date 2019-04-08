<template>
  <section class="section">
    <h2 class="subtitle">Files</h2>

    <b-table
        :data="filteredFiles"
        paginated
        :per-page="perPage"
        :current-page.sync="currentPage"
    >
      <template slot="empty">
        No files for this project.
      </template>
      <template slot-scope="props">
          <b-table-column field="filename" label="Filename">
            <a href="#" @click.prevent="handleDownload(props.row.id, props.row.filename)">
              <span :class="props.row.superseded ? 'file-superseded' : ''">{{ props.row.filename }} {{ props.row.superseded ? '[superseded]' : '' }}</span>
            </a>
          </b-table-column>
          <b-table-column field="created_at" label="Uploaded">
            {{ props.row.created_at | moment("dddd, MMMM Do YYYY, h:mm:ss a") }}
          </b-table-column>
          <b-table-column field="actions" label="Actions" class="is-narrow">
            <button class="button is-small" @click="handleDownload(props.row.id, props.row.filename)"><font-awesome-icon :icon="['fas', 'download']"></font-awesome-icon></button>
            <button class="button is-small ml" @click="handleDelete(props.row.id)"><font-awesome-icon :icon="['far', 'trash-alt']"></font-awesome-icon></button>
          </b-table-column>
      </template>
    </b-table>
    <button type="button" class="button is-primary" @click="isUploadModalActive=true">Upload files</button>
    <b-modal :active.sync="isUploadModalActive">
      <form @submit.prevent="handleFileUpload">
        <div class="modal-card" style="width: 20rem;">
          <section class="modal-card-body">

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
          </section>
        </div>
      </form>
    </b-modal>

  </section>
</template>

<script>
import { saveAs } from 'file-saver'
export default {
  name: 'ProjectFiles',
  props: ['project', 'files'],
  data () {
    return {
      dropFiles: [],
      loading: false,
      showSuperseded: false,
      perPage: 10,
      currentPage: 1,
      isUploadModalActive: false
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

      const config = {
        headers: {
          'content-type': 'multipart/form-data'
        }
      }

      const staged = [...this.dropFiles]
      const uploads = []

      staged.forEach((file, i) => {
        const formData = new FormData()
        formData.append('file', file)
        const req = this.$http.post(`projects/${this.$route.params.id}/files`, formData, config)
        uploads.push(req)
      })

      Promise.all(uploads).then((response) => {
        this.dropFiles = []
        this.$noty.success('Files successfully uploaded.')
      }).catch((e) => {
        this.$noty.error('Error uploading files. Please try again later.')
      }).finally(() => {
        this.loading = false
        this.isUploadModalActive = false
        this.$emit('updated', true)
      })
    },
    // delete a file staged for upload
    deleteDropFile (i) {
      this.dropFiles.splice(i, 1)
    },
    // delete a file on the server
    deleteFile (id) {
      this.$http.delete(`projects/${this.$route.params.id}/files/${id}`).then((response) => {
        this.$emit('updated', true)
        this.$noty.success('File deleted')
      }).catch((e) => {
        this.$noty.error('Error deleting file. Please try again later.')
      })
    },
    handleDelete (id) {
      this.$dialog.confirm({
        message: 'Are you sure you want to delete this file?',
        onConfirm: () => this.deleteFile(id)
      })
    },
    handleDownload (id, filename) {
      this.$http.get(`projects/${this.$route.params.id}/files/${id}`, { responseType: 'blob' }).then((response) => {
        const file = new Blob([response.data])
        saveAs(file, filename)
      }).catch((e) => {
        console.error(e)
      })
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
