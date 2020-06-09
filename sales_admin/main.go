package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	r := mux.NewRouter()

	svc, err := NewSalesService()
	if err != nil {
		log.Fatal(err)
	}

	err = svc.CreateSalesTable()
	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(20000000)
		if err != nil {
			respondWithError(w, 404, "file exceeds max 20mb")
			return
		}

		_, header, err := r.FormFile("file")
		if err != nil {
			respondWithError(w, 404, "could not find file with key 'file'")
			return
		}

		fileName := header.Filename
		fmt.Println(fileName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		//Debug: true,
	})

	handler := c.Handler(r)

	fmt.Println("server running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
