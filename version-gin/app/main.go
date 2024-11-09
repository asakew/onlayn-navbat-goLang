package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type QueueItem struct {
	Number int
	Time   time.Time
}

var queue []QueueItem

type TemplateData struct {
	Number int
	Queue  []QueueItem
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := gin.Default()

	// Asosiy sahifa uchun marshrutni o'rnatish
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", TemplateData{
			Queue: queue,
		})
	})

	// Navbatga yozish uchun so'rovni qayta ishlash
	router.POST("/join", func(c *gin.Context) {
		number := rand.Intn(1000) // Генерация случайного номера
		queue = append(queue, QueueItem{Number: number, Time: time.Now()})

		c.HTML(http.StatusOK, "index.html", TemplateData{
			Number: number,
			Queue:  queue,
		})
	})

	// Shablonni o'rnatish
	router.LoadHTMLGlob("web/templates/*")

	// Serverni ishga tushirish
	fmt.Fprint(os.Stderr, "Server is starting: http://localhost:8080\n")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
