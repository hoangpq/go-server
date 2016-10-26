package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	elastic "gopkg.in/olivere/elastic.v3"
)

// Product is a structure used for serizlizing/deserializing data in ES.
type Product struct {
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Categories []string `json:"categories"`
}

func main() {

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://localhost:9200"),
	)
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// search

		searchResult, err := client.Search().
			Index("odoo").   // search in index "odoo"
			Type("product"). // search in type "product"
			From(0).
			Size(10000).  // take documents 0-9
			Pretty(true). // pretty print request and response JSON
			Do()          // execute

		products := []Product{}
		if err != nil {
			// Handle error
			panic(err)
		} else {
			var ttyp Product
			for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
				if t, ok := item.(Product); ok {
					// fmt.Printf("Tweet by %s: %f\n", t.Name, t.Price)
					products = append(products, t)
				}
			}
		}

		json, err := json.Marshal(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})

	http.ListenAndServe(":8080", r)
}
