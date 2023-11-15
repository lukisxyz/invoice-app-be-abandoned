package controller

import (
	"encoding/json"
	"errors"
	"flukis/invokiss/app/model"
	"flukis/invokiss/database/querier"
	"flukis/invokiss/lib/httpresponse"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oklog/ulid/v2"
)

type CategoryController struct {
	writeCategory querier.CategoryWriteModel
	readCategory  querier.CategoryReadModel
}

func NewCategoryController(
	writeCategory querier.CategoryWriteModel,
	readCategory querier.CategoryReadModel,
) *CategoryController {
	return &CategoryController{writeCategory, readCategory}
}

func (p *CategoryController) Routes() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", p.GetAll)
	r.Get("/{id}", p.GetOneByID)
	r.Post("/", p.Create)

	return r
}

type createCategoryBodyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p createCategoryBodyRequest) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Description, validation.Required),
	)
}

func (p *CategoryController) Create(w http.ResponseWriter, req *http.Request) {
	var data createCategoryBodyRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if err := data.Validate(); err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	ctx := req.Context()
	newCategory := model.NewCategory(
		data.Name,
		data.Description,
	)
	err := p.writeCategory.Save(ctx, newCategory)
	if err != nil {
		if errors.Is(err, model.ErrCategoryNotFound) {
			httpresponse.WriteError(
				w,
				http.StatusBadRequest,
				err,
			)
			return
		}
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	httpresponse.WriteData(w, http.StatusCreated, newCategory.ID, nil)
}

func (p *CategoryController) GetAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	data, err := p.readCategory.Fetch(ctx)
	if err != nil {
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}

	meta.Total = data.Count

	httpresponse.WriteData(w, http.StatusOK, data.Data, meta)
}

func (p *CategoryController) GetOneByID(w http.ResponseWriter, req *http.Request) {
	var idStr = chi.URLParam(req, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	ctx := req.Context()
	data, err := p.readCategory.GetOneByID(ctx, id)
	if err != nil {
		httpresponse.WriteError(
			w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	httpresponse.WriteData(w, http.StatusOK, data, nil)
}
