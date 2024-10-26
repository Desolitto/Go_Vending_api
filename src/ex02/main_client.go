package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char* ask_cow(char phrase[]) {
    int phrase_len = strlen(phrase);
    char *buf = (char *)malloc(sizeof(char) * (160 + (phrase_len + 2) * 3));
    strcpy(buf, " ");

    for (unsigned int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "_");
    }

    strcat(buf, "\n< ");
    strcat(buf, phrase);
    strcat(buf, " ");
    strcat(buf, ">\n ");

    for (unsigned int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "-");
    }
    strcat(buf, "\n");
    strcat(buf, "        \\   ^__^\n");
    strcat(buf, "         \\  (oo)\\_______\n");
    strcat(buf, "            (__)\\       )\\/\\\n");
    strcat(buf, "                ||----w |\n");
    strcat(buf, "                ||     ||\n");
    return buf;
}
*/
import "C"
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
	"unsafe"
)

type Order struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

// Структура для ответа с изменением и благодарностью
type Response struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
	Error  string `json:"error,omitempty"` // Добавлено для обработки ошибок
}

// Функция для генерации фразы коровы
func generateCowPhrase(msg string) string {
	cowPhrase := C.CString(msg)
	defer C.free(unsafe.Pointer(cowPhrase))
	cow := C.ask_cow(cowPhrase)
	defer C.free(unsafe.Pointer(cow))
	return C.GoString(cow)
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

	// Создание структуры заказа
	order := Order{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *count,
	}

	// Преобразование структуры заказа в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order to JSON: %v", err)
	}

	// Отправка POST-запроса с данными в формате JSON
	url := "https://localhost:3333/buy_candy"
	response, err := client.Post(url, "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		log.Fatalf("Failed to send POST request: %v", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Обработка JSON-ответа сервера
	var serverResponse Response
	if err := json.Unmarshal(body, &serverResponse); err != nil {
		log.Fatalf("Failed to unmarshal response from JSON: %v", err)
	}

	// Проверка на наличие ошибки
	if serverResponse.Error != "" {
		// Если ошибка, формируем сообщение для коровы с текстом ошибки
		cowMessage := serverResponse.Error
		serverResponse.Thanks = generateCowPhrase(cowMessage)
		fmt.Println("Error:")
		fmt.Println(serverResponse.Thanks) // Выводим текст коровы с ошибкой
	} else {
		// В противном случае, благодарим за успешный заказ
		cowMessage := "Thank you!"
		serverResponse.Thanks = generateCowPhrase(cowMessage)
		fmt.Printf("Change: %d\nThanks:\n%s\n", serverResponse.Change, serverResponse.Thanks)
	}

}
