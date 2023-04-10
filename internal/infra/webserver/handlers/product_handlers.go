package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luiscovelo/goexpert-api-rest/internal/dto"
	"github.com/luiscovelo/goexpert-api-rest/internal/entity"
	"github.com/luiscovelo/goexpert-api-rest/internal/infra/database"
	entityPkg "github.com/luiscovelo/goexpert-api-rest/pkg/entity"
)

type ProductHandler struct {
	Response
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// CreateProduct godoc
//	@Summary		Create Product
//	@Description	Create Product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateProductInput	true	"product request"
//	@Success		201		{object}	entity.Product
//	@Failure		500		{object}	handlers.Response
//	@Router			/products [post]
//	@Security		ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, req *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}
	if err := h.ProductDB.Create(p); err != nil {
		h.JSON(w, http.StatusInternalServerError, err)
		return
	}

	newProduct, err := h.ProductDB.FindByID(p.ID.String())
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	h.JSON(w, http.StatusCreated, newProduct)
}

// GetProduct godoc
//	@Summary		Get a product
//	@Description	Get a product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"product ID"	Format(uuid)
//	@Success		200	{object}	entity.Product
//	@Failure		404	{object}	handlers.Response
//	@Failure		500	{object}	handlers.Response
//	@Router			/products/{id} [get]
//	@Security		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		h.JSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		h.JSON(w, http.StatusNotFound, err)
		return
	}

	h.JSON(w, http.StatusOK, product)
}

// ListProducts godoc
//	@Summary		List Products
//	@Description	get all products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			page	query		string	false	"page number"
//	@Param			limit	query		string	false	"limit offset"
//	@Success		200		{array}		entity.Product
//	@Failure		500		{object}	handlers.Response
//	@Router			/products [get]
//	@Security		ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, req *http.Request) {
	page := req.URL.Query().Get("page")
	limit := req.URL.Query().Get("limit")
	sort := req.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	h.JSON(w, http.StatusOK, products)
}

// UpdateProduct godoc
//	@Summary		Update Product
//	@Description	Update Product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateProductInput	true	"product request"
//	@Param			id		path		string					true	"Product ID"	Format(uuid)
//	@Success		200		{object}	entity.Product
//	@Failure		404		{object}	handlers.Response
//	@Failure		500		{object}	handlers.Response
//	@Router			/products/{id} [put]
//	@Security		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		h.JSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}

	var uuid entityPkg.ID
	uuid, err := entityPkg.ParseID(id)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	if _, err := h.ProductDB.FindByID(id); err != nil {
		h.JSON(w, http.StatusNotFound, err)
		return
	}

	var product entity.Product
	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	product.ID = uuid

	if err := h.ProductDB.Update(&product); err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	h.JSON(w, http.StatusOK, product)
}

// DeleteProduct godoc
//	@Summary		Delete Product
//	@Description	Delete Product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Product ID"	Format(uuid)
//	@Success		200
//	@Failure		404	{object}	handlers.Response
//	@Failure		500	{object}	handlers.Response
//	@Router			/products/{id} [delete]
//	@Security		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		h.JSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}

	if _, err := h.ProductDB.FindByID(id); err != nil {
		h.JSON(w, http.StatusNotFound, err)
		return
	}

	if err := h.ProductDB.Delete(id); err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	h.JSON(w, http.StatusOK, nil)
}
