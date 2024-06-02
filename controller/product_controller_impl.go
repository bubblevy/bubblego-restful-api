package controller

import (
	"bubblevy/restful-api/helper"
	"bubblevy/restful-api/model/web"
	"bubblevy/restful-api/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type productControllerImpl struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &productControllerImpl{
		ProductService: productService,
	}
}

func (controller *productControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productCreateRequest := web.ProductCreateRequest{}
	helper.ReadFromRequestBody(request, &productCreateRequest)

	productResponse := controller.ProductService.Create(request.Context(), productCreateRequest)
	webResponse := web.WebResponse{
		Code:    http.StatusCreated,
		Error:   false,
		Message: "Create product successfully",
		Data:    productResponse,
	}

	writer.WriteHeader(http.StatusCreated)

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *productControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productUpdateRequest := web.ProductUpdateRequest{}
	helper.ReadFromRequestBody(request, &productUpdateRequest)

	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	productUpdateRequest.Id = id

	productResponse := controller.ProductService.Update(request.Context(), productUpdateRequest)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Update product successfully",
		Data:    productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *productControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	controller.ProductService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Delete product successfully",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *productControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	productResponse := controller.ProductService.FindById(request.Context(), id)
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully retrieved a single product",
		Data:    productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *productControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productResponse := controller.ProductService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully retrieved all products",
		Data:    productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
