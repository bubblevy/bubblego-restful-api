package test

import (
	"bubblevy/restful-api/app"
	"bubblevy/restful-api/controller"
	"bubblevy/restful-api/helper"
	"bubblevy/restful-api/middleware"
	"bubblevy/restful-api/model/domain"
	"bubblevy/restful-api/repository"
	"bubblevy/restful-api/service"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func testDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_golang_restful_api_test")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, validate)
	productController := controller.NewProductController(productService)
	router := app.NewRouter(productController)

	return middleware.NewAuthMiddleware(router)
}

func truncateProduct(db *sql.DB) {
	db.Exec("TRUNCATE products")
}

func TestCreateProductSuccess(t *testing.T) {
	db := testDB()
	truncateProduct(db)
	router := setupRouter(db)
	requestBody := strings.NewReader(`{"product_name" : "Cokelat", "price" : 9500}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 201, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 201, int(responseBody["code"].(float64)))
	assert.Equal(t, false, responseBody["error"])
	assert.Equal(t, "Cokelat", responseBody["data"].(map[string]interface{})["product_name"])
	assert.Equal(t, 9500, int(responseBody["data"].(map[string]interface{})["price"].(float64)))
}

func TestCreateProductFailed(t *testing.T) {
	db := testDB()
	truncateProduct(db)
	router := setupRouter(db)
	requestBody := strings.NewReader(`{"product_name" : "", "price" : 9500}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, true, responseBody["error"])
}

func TestUpdateProductSuccess(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		ProductName: "Cokelat",
		Price:       9500,
	})
	tx.Commit()

	router := setupRouter(db)
	requestBody := strings.NewReader(`{"product_name" : "Cokelat", "price" : 9500}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/products/"+strconv.Itoa(product.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, false, responseBody["error"])
	assert.Equal(t, product.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Cokelat", responseBody["data"].(map[string]interface{})["product_name"])
	assert.Equal(t, 9500, int(responseBody["data"].(map[string]interface{})["price"].(float64)))
}

func TestUpdateProductFailed(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		ProductName: "Cokelat",
		Price:       9500,
	})
	tx.Commit()

	router := setupRouter(db)
	requestBody := strings.NewReader(`{"product_name" : "", "price" : 9500}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/products/"+strconv.Itoa(product.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, true, responseBody["error"])
}

func TestGetProductSuccess(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		ProductName: "Cokelat",
		Price:       9500,
	})
	tx.Commit()

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products/"+strconv.Itoa(product.Id), nil)
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, false, responseBody["error"])
	assert.Equal(t, product.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, product.ProductName, responseBody["data"].(map[string]interface{})["product_name"])
	assert.Equal(t, product.Price, int(responseBody["data"].(map[string]interface{})["price"].(float64)))
}

func TestGetProductFailed(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products/9999", nil)
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, true, responseBody["error"])
}

func TestDeleteProductSuccess(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		ProductName: "Cokelat",
		Price:       9500,
	})
	tx.Commit()

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/products/"+strconv.Itoa(product.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, false, responseBody["error"])
}

func TestDeleteProductFailed(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/products/999", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, true, responseBody["error"])
}

func TestGetAllProductSuccess(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product1 := productRepository.Save(context.Background(), tx, domain.Product{
		ProductName: "Cokelat",
		Price:       9500,
	})
	product2 := productRepository.Save(context.Background(), tx, domain.Product{
		ProductName: "Kentang",
		Price:       5000,
	})
	tx.Commit()

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products", nil)
	request.Header.Add("API-Key", "BUBBLEKEY")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, false, responseBody["error"])

	var products = responseBody["data"].([]interface{})
	productsResponse1 := products[0].(map[string]interface{})
	productsResponse2 := products[1].(map[string]interface{})

	assert.Equal(t, product1.Id, int(productsResponse1["id"].(float64)))
	assert.Equal(t, product1.ProductName, productsResponse1["product_name"])
	assert.Equal(t, product1.Price, int(productsResponse1["price"].(float64)))

	assert.Equal(t, product2.Id, int(productsResponse2["id"].(float64)))
	assert.Equal(t, product2.ProductName, productsResponse2["product_name"])
	assert.Equal(t, product2.Price, int(productsResponse2["price"].(float64)))
}

func TestUnauthorized(t *testing.T) {
	db := testDB()
	truncateProduct(db)

	router := setupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products", nil)
	request.Header.Add("API-Key", "KEYSALAH")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, true, responseBody["error"])
}
