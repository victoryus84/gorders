package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/service"
	"github.com/victoryus84/gorders/internal/utils" 	
	"github.com/victoryus84/gorders/internal/middleware"
)

type ProductHandler struct {
	// [CORECȚIE 2]: Fără steluță (*) la interfață!
	svc service.ProductService 
}

// [CORECȚIE 2]: Fără steluță (*) la interfață în argument!
func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (hdl *ProductHandler) CreateProduct(ctx *gin.Context) {
	requests, err := utils.ParseBody[dto.ProductDTO](ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	// [ATENȚIE]: Asigură-te că în svc_product.go (în interfață) funcția se numește exact așa!
	result := hdl.svc.ProcessProductImport(requests)

	ctx.JSON(http.StatusCreated, result)
}

func (hdl *ProductHandler) CreateProductGroup(ctx *gin.Context) {
	requests, err := utils.ParseBody[dto.ProductGroupDTO](ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	// [ATENȚIE]: Asigură-te că funcția există în interfața ProductService!
	result := hdl.svc.ProcessProductGroupImport(requests)

	ctx.JSON(http.StatusCreated, result)
}

// GetProducts - Obține primii 1000 de produse
func (hdl *ProductHandler) GetProducts(ctx *gin.Context) {
	// [ATENȚIE]: Asigură-te că funcția există în interfața ProductService!
	products, err := hdl.svc.GetFirst1000Products()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func RegisterProductRoutes(r *gin.Engine, h *ProductHandler) {
	products := r.Group("/api/v1/products")
	products.Use(middleware.AuthJWT())
	{
		products.GET("", h.GetProducts)
		products.POST("", h.CreateProduct)
		products.POST("/groups", h.CreateProductGroup)
	}
}