package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/petr-baranov/imageservice/internal/handler"
	"github.com/petr-baranov/imageservice/internal/services"
)

const config = `[
{
	"Name":"small",
	"Factor":0.1 
},
{
	"Name":"medium",
	"Factor":0.5 
},
{
	"Name":"large",
	"Factor":1.5 
}
]`
const imageStoreDir = "images"

func main() {
	var scaleConfig []services.ScaleConfig
	if err := json.Unmarshal([]byte(config), &scaleConfig); err != nil {
		log.Fatal("failed to read config")
	}
	handler := handler.NewHandler(services.NewImageService(scaleConfig), services.NewImageStore(imageStoreDir))
	http.HandleFunc("/images", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			handler.HandlePost(w, req)
			return
		case http.MethodGet:
			handler.HandleGet(w, req)
			return
		default:
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}
	})
	fmt.Println("Starting IMAGE server on PORT 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
