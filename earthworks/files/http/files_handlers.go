package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"github.com/stephenhillier/geoprojects/earthworks"
)

var decoder = schema.NewDecoder()

// NewFile adds a new file to the file service repo
func (svc *FileSvc) NewFile(w http.ResponseWriter, r *http.Request) {

	// get project context
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	var Buf bytes.Buffer

	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	name := header.Filename

	io.Copy(&Buf, file)
	contents := Buf.Bytes()

	newFile := earthworks.FileRequest{
		File:     contents,
		Filename: name,
		Project:  project.ID,
		Category: "other",
	}

	created, err := svc.Repo.NewFile(newFile)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, created)
}

// ListFiles returns a list of all files in the file service repo
func (svc *FileSvc) ListFiles(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	var filter earthworks.FileFilter

	err = decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	filter.Project = project.ID

	files, err := svc.Repo.ListFiles(filter)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fileVersions := make(map[string]int)

	for i, f := range files {
		if fileVersions[f.Filename] >= 1 {
			f.Superseded = true
			files[i] = f
		}
		fileVersions[f.Filename]++
	}

	render.JSON(w, r, files)
}

// GetFile retrieves a file from the file service repo
func (svc *FileSvc) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileID"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	file, err := svc.Repo.GetFile(fileID, project.ID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", file.Filename))
	render.Data(w, r, file.File)
}

// DeleteFile expires a file from the file service repo
func (svc *FileSvc) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileID"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = svc.Repo.DeleteFile(fileID, project.ID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	}

	w.WriteHeader(http.StatusNoContent)
}

// RestoreFile unexpires an expired file from the file service repo
func (svc *FileSvc) RestoreFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileID"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = svc.Repo.RestoreFile(fileID, project.ID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	}

	w.WriteHeader(http.StatusOK)
}
