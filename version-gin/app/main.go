package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
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
			Queue: sortQueueByTime(queue),
		})
	})

	// Navbatga yozish uchun so'rovni qayta ishlash
	router.POST("/join", func(c *gin.Context) {
		number := rand.Intn(1000) // Tasodiy raqam yaratish
		queue = append(queue, QueueItem{Number: number, Time: time.Now()})

		c.HTML(http.StatusOK, "index.html", TemplateData{
			Number: number,
			Queue:  sortQueueByTime(queue),
		})
	})

	// Shablonni o'rnatish
	router.LoadHTMLGlob("web/templates/*")

	// Serverni ishga tushirish
	_, err := fmt.Fprint(os.Stderr, "Server is starting: http://localhost:8080\n")
	if err != nil {
		return
	}
	err = router.Run(":8080")
	if err != nil {
		return
	}
}

// Yangi vaqtlarni qo'shish voqtida navbatlarni yuqoridan saralash
func sortQueueByTime(queue []QueueItem) []QueueItem {
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].Time.After(queue[j].Time)
	})
	return queue
}
