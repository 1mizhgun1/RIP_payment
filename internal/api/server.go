package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const password = "bf438279f8hnc28497g8fhvbv3682739fhbvc2332f9ch2438bgnv0v348ng0fgg"

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/pay", func(c *gin.Context) {
		var data PayData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orderID := data.OrderID

		// Запуск горутины для отправки статуса
		go sendStatus(orderID, password, fmt.Sprintf("http://localhost:8080/orders/%d/status/", orderID))

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	router.Run(":5000")

	log.Println("Server down")
}

func genRandomStatus(password string) Result {
	time.Sleep(8 * time.Second)
	status := "A"
	if rand.Intn(100) < 20 {
		status = "W"
	}
	return Result{status, password}
}

// Функция для отправки статуса в отдельной горутине
func sendStatus(orderID int, password string, url string) {
	// Выполнение расчётов с randomStatus
	result := genRandomStatus(password)

	// Отправка PUT-запроса к основному серверу
	_, err := performPUTRequest(url, result)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return
	}

	fmt.Println("Status sent successfully for orderID:", orderID)
}

type Result struct {
	Status   string `json:"status"`
	Password string `json:"password"`
}

type PayData struct {
	OrderID int `json:"order_id"`
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Создание PUT-запроса
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
