package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/zioufang/mltrackapi/pkg/api/apiutil"
	"github.com/zioufang/mltrackapi/pkg/api/model"
)

// CreateModelRun creates the entity in the database
func (s *Server) CreateModelRun(w http.ResponseWriter, r *http.Request) {
	m := model.ModelRun{}
	err := apiutil.ReadReqBody(w, r, s.db, &m)
	if err != nil {
		apiutil.RespError(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = m.Create(s.db)
	if err != nil {
		apiutil.RespError(w, http.StatusInternalServerError, err)
		return
	}
	apiutil.RespSuccess(w, m)
}

// GetAllModelRuns gets all the models from the database
func (s *Server) GetAllModelRuns(w http.ResponseWriter, r *http.Request) {
	m := model.ModelRun{}
	models, err := m.GetAll(s.db)
	if err != nil {
		apiutil.RespError(w, http.StatusInternalServerError, err)
		return
	}
	log.Println(models)
	apiutil.RespSuccess(w, models)
}

// GetModelRunByID gets one model given an ID from the database
func (s *Server) GetModelRunByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apiutil.RespError(w, http.StatusBadRequest, err)
		return
	}
	m := model.ModelRun{}
	modelGet, err := m.GetByID(s.db, id)
	if err != nil {
		apiutil.RespError(w, http.StatusBadRequest, err)
		return
	}
	apiutil.RespSuccess(w, modelGet)
}

// DeleteModelRunByID deletes one model given an ID from the database
func (s *Server) DeleteModelRunByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apiutil.RespError(w, http.StatusBadRequest, err)
		return
	}
	m := model.ModelRun{}
	err = m.DeleteByID(s.db, id)
	if err != nil {
		apiutil.RespError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	w.WriteHeader(http.StatusNoContent)
	apiutil.RespSuccess(w, "")
}