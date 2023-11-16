package controller

import (
	"encoding/json"
	"errors"
	"flukis/invokiss/app/model"
	"flukis/invokiss/database/querier"
	"flukis/invokiss/lib/httpresponse"
	"net/http"
	"strings"

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
	r.Put("/{id}", p.Change)
	r.Post("/", p.Create)
	r.Patch("/{id}/inventory", p.AssignQuantity)

	return r
}

type assignQtyBodyRequest struct {
	Quantity int `json:"qty"`
}

func (p assignQtyBodyRequest) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Quantity, validation.Required),
	)
}

func (p *ProductController) AssignQuantity(w http.ResponseWriter, req *http.Request) {
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

	var data assignQtyBodyRequest
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
	err = p.writeProduct.AssignQuantity(ctx, id, data.Quantity)
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

	httpresponse.WriteData(w, http.StatusCreated, data.Quantity, nil)
}

type createProductBodyRequest struct {
	Sku         string      `json:"sku"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       *[]byte     `json:"image"`
	Amount      float64     `json:"amount"`
	Quantity    int         `json:"quantity"`
	Categories  []ulid.ULID `json:"categories"`
}

func (p createProductBodyRequest) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Sku, validation.Required),
		validation.Field(&p.Description, validation.Required),
		validation.Field(&p.Amount, validation.Required),
		validation.Field(&p.Quantity, validation.Required),
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
		data.Quantity,
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

type changeProductBodyRequest struct {
	Sku         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       *[]byte `json:"image"`
	Amount      float64 `json:"amount"`
}

func (p changeProductBodyRequest) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Sku, validation.Required),
		validation.Field(&p.Description, validation.Required),
		validation.Field(&p.Amount, validation.Required),
	)
}

func (p *ProductController) Change(w http.ResponseWriter, req *http.Request) {
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

	var data changeProductBodyRequest
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
	newProduct := model.Product{
		ID:          id,
		Sku:         data.Sku,
		Name:        data.Name,
		Description: data.Description,
		Image:       data.Image,
		Amount:      data.Amount,
	}
	err = p.writeProduct.Edit(ctx, newProduct)
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

	httpresponse.WriteData(w, http.StatusCreated, newProduct.ID, nil)
}

func (p *ProductController) GetAll(w http.ResponseWriter, req *http.Request) {
	filterInString := req.URL.Query().Get("category")
	IdsInArray := strings.Split(filterInString, ",")

	var IdsUlid = make([]ulid.ULID, len(IdsInArray))
	for i := range IdsInArray {
		idInUlid, err := ulid.Parse(IdsInArray[i])
		if err != nil {
			continue
		}
		IdsUlid[i] = idInUlid
	}

	ctx := req.Context()
	data, err := p.readProduct.FetchByCategoryID(ctx, IdsUlid)
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
