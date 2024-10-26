package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Order struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

func main() {
	// Чтение флагов командной строки
	candyType := flag.String("k", "", "Type of candy (two-letter abbreviation)")
	count := flag.Int("c", 0, "Count of candy to buy")
	money := flag.Int("m", 0, "Amount of money")
	flag.Parse()

	// Проверка, что параметры указаны
	if *candyType == "" || *count <= 0 || *money <= 0 {
		log.Fatal("All parameters (-k, -c, -m) must be provided with valid values")
	}

	// Загрузка клиентских сертификатов
	cert, err := tls.LoadX509KeyPair("localhost/cert.pem", "localhost/key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	// Загрузка CA сертификата
	caCert, err := os.ReadFile("minica.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройка конфигурации TLS для клиента
	//InsecureSkipVerify true даже если он самоподписанный или не соответствует имени хоста
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Формирование URL с параметрами запроса
	url := fmt.Sprintf("https://localhost:3333/buy_candy?k=%s&c=%d&m=%d", *candyType, *count, *money)

	// Создание структуры заказа
	order := Order{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *count,
	}

	// Преобразование структуры в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	// Отправка POST-запроса с данными в формате JSON
	response, err := client.Post(url, "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Вывод ответа от сервера
	fmt.Println(string(body))
}

// func main() {
// 	candyType := flag.String("k", "", "Candy type")
// 	count := flag.Int("c", 0, "Count of candy")
// 	money := flag.Int("m", 0, "Amount of money")
// 	flag.Parse()
// 	if *candyType == "" || *count <= 0 || *money <= 0 {
// 		fmt.Println("Invalid input. Please provide valid candy type, count, and money.")
// 		return
// 	}
// 	// Настройка TLS
// 	cert, err := tls.LoadX509KeyPair("localhost/cert.pem", "localhost/key.pem")
// 	if err != nil {
// 		log.Fatalf("Failed to load client certificate: %v", err)
// 	}

// 	caCert, err := os.ReadFile("minica.pem")
// 	if err != nil {
// 		log.Fatalf("Failed to read CA certificate: %v", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	if !caCertPool.AppendCertsFromPEM(caCert) {
// 		log.Fatal("Failed to append CA certificate")
// 	}

// 	tlsConfig := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 		RootCAs:      caCertPool,
// 	}

// 	transport := &http.Transport{TLSClientConfig: tlsConfig}
// 	client := &http.Client{Transport: transport}
// 	order := Order{
// 		Money:      *money,
// 		CandyType:  *candyType,
// 		CandyCount: *count,
// 	}
// 	orderJSON, err := json.Marshal(order)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal order to JSON: %v", err)
// 	}
// 	// Отправка запроса
// 	req, err := http.NewRequest("POST", "https://localhost:3333/buy_candy", bytes.NewBuffer(orderJSON))
// 	if err != nil {
// 		log.Fatalf("Failed to create request: %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Fatalf("Request failed with status: %s", resp.Status)
// 	}

// 	// Обработка ответа
// 	var responseBody struct {
// 		Change int `json:"change"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
// 		log.Fatalf("Failed to decode response: %v", err)
// 	}

// 	log.Printf("Thank you! Your change is %d\n", responseBody.Change)

// }
