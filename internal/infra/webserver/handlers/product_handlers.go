package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/leobelini-studies/go_expert_api/internal/dto"
	"github.com/leobelini-studies/go_expert_api/internal/entity"
	entityPkg "github.com/leobelini-studies/go_expert_api/pkg/entity"
	"github.com/leobelini-studies/go_expert_api/internal/infra/database"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct Create Product godoc
// @Summary     Create product
// @Description Create products
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       resquest body dto.CreateProductInput true "product request"
// @Success     201
// @Failure     500 {object} dto.ErrorOutput
// @Router      /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetProduct Get Product godoc
// @Summary     Get product
// @Description Get product
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       id path string true "product ID" Format(uuid)
// @Success     200 {array} entity.Product
// @Failure     404 {object} dto.ErrorOutput
// @Failure     500 {object} dto.ErrorOutput
// @Router      /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "Invalid ID"}
		json.NewEncoder(w).Encode(error)
		return
	}

	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct Update Product godoc
// @Summary     Update product
// @Description Update product
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       id path string true "product ID" Format(uuid)
// @Param       resquest body entity.Product true "product update"
// @Success     200
// @Failure     404 {object} dto.ErrorOutput
// @Failure     500 {object} dto.ErrorOutput
// @Router      /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "Invalid ID"}
		json.NewEncoder(w).Encode(error)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct Delete Product godoc
// @Summary     Delete product
// @Description Delete product
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       id path string true "product ID" Format(uuid)
// @Success     200
// @Failure     404 {object} dto.ErrorOutput
// @Failure     500 {object} dto.ErrorOutput
// @Router      /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: "Invalid ID"}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetProducts List all products godoc
// @Summary     List products
// @Description Get all product
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       page query string false "page number"
// @Param       limit query string false "limit"
// @Success     200 {array} entity.Product
// @Failure     404 {object} dto.ErrorOutput
// @Failure     500 {object} dto.ErrorOutput
// @Router      /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
