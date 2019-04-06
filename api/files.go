package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// FileRequest is a request that contains a file to be uploaded
type FileRequest struct {
	File      []byte
	Filename  string
	CreatedBy string
	Category  string
	Project   int
}

// File represents a file and some metadata about the file
type File struct {
	ID         int64     `json:"id"`
	Project    int64     `json:"project"`
	Filename   string    `json:"filename"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	Superseded bool      `json:"superseded"`
	Archived   NullDate  `json:"archived,omitempty"`
}

// FileObject is a struct that contains a byte slice,
// for reading into and out of a database/file store.
type FileObject struct {
	File     []byte
	Filename string
}

// FileFilter allows searching for files based on different criteria (project, file category etc)
type FileFilter struct {
	ID       int64  `json:"id" schema:"id"`
	Project  int    `json:"project" schema:"project"`
	Category string `json:"category" schema:"category"`
	Archived bool   `json:"archived" schema:"archived"`
}

func (s *server) NewFile(w http.ResponseWriter, r *http.Request) {

	// get project context
	ctx := r.Context()
	project, ok := ctx.Value(projectCtx).(Project)
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

	newFile := FileRequest{
		File:     contents,
		Filename: name,
		Project:  project.ID,
		Category: "other",
	}

	created, err := s.datastore.NewFile(newFile)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, created)
}

func (s *server) ListFiles(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	var filter FileFilter

	err = decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	filter.Project = project.ID

	files, err := s.datastore.ListFiles(filter)
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

func (s *server) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileID"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	file, err := s.datastore.GetFile(fileID, project.ID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", file.Filename))
	render.Data(w, r, file.File)
}

func (s *server) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileID"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = s.datastore.DeleteFile(fileID, project.ID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	}

	render.Status(r, http.StatusNoContent)
	return
}

func (s *server) RestoreFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileID"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = s.datastore.RestoreFile(fileID, project.ID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	}

	render.Status(r, http.StatusOK)
	return
}
