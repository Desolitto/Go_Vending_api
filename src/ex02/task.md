
```
swagger generate server -f swagger.yaml
go mod tidy
```
```
minica -domains localhost
```
```
in restapi/configure_candy_server

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
```
```
go run cmd/candy-server-server/main.go --tls-certificate localhost/cert.pem --tls-key localhost/key.pem --tls-port=3333
```
```

go build -o candy-client main_client.go && ./candy-client -k AA -c 2 -m 50
./candy-client -k AA -c 2 -m 50
```