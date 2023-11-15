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

type ProductController struct {
	writeProduct querier.ProductWriteModel
	readProduct  querier.ProductReadModel
}

func NewProductController(
	writeProduct querier.ProductWriteModel,
	readProduct querier.ProductReadModel,
) *ProductController {
	return &ProductController{writeProduct, readProduct}
}

func (p *ProductController) Routes() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", p.GetAll)
	r.Get("/{id}", p.GetOneByID)
	r.Post("/", p.Create)

	return r
}

type createProductBodyRequest struct {
	Sku         string      `json:"sku"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       *[]byte     `json:"image"`
	Amount      float64     `json:"amount"`
	Categories  []ulid.ULID `json:"categories"`
}

func (p createProductBodyRequest) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Sku, validation.Required),
		validation.Field(&p.Description, validation.Required),
		validation.Field(&p.Amount, validation.Required),
	)
}

func (p *ProductController) Create(w http.ResponseWriter, req *http.Request) {
	var data createProductBodyRequest
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
	newProduct := model.NewProduct(
		data.Sku,
		data.Name,
		data.Description,
		data.Image,
		data.Amount,
	)
	err := p.writeProduct.Save(ctx, newProduct)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) || errors.Is(err, model.ErrProductSKUDuplicated) {
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

	if len(data.Categories) > 0 {
		err := p.writeProduct.AssignCategories(ctx, newProduct.ID, data.Categories)
		if err != nil {
			httpresponse.WriteError(
				w,
				http.StatusInternalServerError,
				err,
			)
			return
		}
	}

	httpresponse.WriteData(w, http.StatusCreated, newProduct.ID, nil)
}

func (p *ProductController) GetAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	data, err := p.readProduct.Fetch(ctx)
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

func (p *ProductController) GetOneByID(w http.ResponseWriter, req *http.Request) {
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
	data, err := p.readProduct.GetOneByID(ctx, id)
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
