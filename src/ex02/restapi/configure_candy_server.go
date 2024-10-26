// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"cow/restapi/operations"
)

//go:generate swagger generate server --target ../../ex01 --name CandyServer --spec ../swagger.yaml --principal interface{}

func configureFlags(api *operations.CandyServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CandyServerAPI) http.Handler {
	// Настройка обработки ошибок
	api.ServeError = errors.ServeError
	// Подключение Swagger UI
	api.UseSwaggerUI()
	// Если нужно использовать redoc, раскомментируйте следующую строку

	// Настройка JSON Consumer и Producer
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Реализация BuyCandyHandler
	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(func(params operations.BuyCandyParams) middleware.Responder {
		order := params.Order

		// Проверка на отрицательное количество конфет
		if *order.CandyCount < 0 {
			return operations.NewBuyCandyBadRequest().WithPayload(&operations.BuyCandyBadRequestBody{
				Error: "Candy count cannot be negative",
			})
		}

		// Определение типов конфет и их стоимости
		validCandyTypes := map[string]int{
			"CE": 5,  // Цена конфеты "CE"
			"AA": 15, // Цена конфеты "AA"
			"NT": 8,  // Цена конфеты "NT"
			"DE": 10, // Цена конфеты "DE"
			"YR": 20, // Цена конфеты "YR"
		}

		// Разыменование указателя на CandyType
		candyPrice, valid := validCandyTypes[*order.CandyType]
		if !valid {
			return operations.NewBuyCandyBadRequest().WithPayload(&operations.BuyCandyBadRequestBody{
				Error: "Invalid candy type",
			})
		}

		// Разыменование указателя на Money и расчет стоимости
		totalCost := candyPrice * int(*order.CandyCount)
		if *order.Money < int64(totalCost) { // Сравниваем *int64 с int64
			return operations.NewBuyCandyPaymentRequired().WithPayload(&operations.BuyCandyPaymentRequiredBody{
				Error: "Insufficient money for this purchase",
			})
		}

		// Возврат успешного ответа
		change := *order.Money - int64(totalCost) // Выполняем разницу между *int64 и int64
		return operations.NewBuyCandyCreated().WithPayload(&operations.BuyCandyCreatedBody{
			Thanks: "Thank you!",
			Change: change,
		})
	})

	// Возвращаем обработчик с глобальными middleware
	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
