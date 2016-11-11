package main

import "github.com/gin-gonic/gin"
import elastic "gopkg.in/olivere/elastic.v3"
import models "go-server/models"

import (
	"reflect"
)

func main() {

	// connect elasticsearch
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)
	if err != nil {
		panic("Connection to elasticsearch failed")
	}
	
	r := gin.Default()
	r.GET("/products", func(c *gin.Context) {

		// search
		searchResult, err := client.Search().
			Index("odoo").// search in index "odoo"
			Type("product").// search in type "product"
			From(0).
			Size(10000).// take documents 0-9
			Pretty(true).// pretty print request and response JSON
			Do()          // execute

		products := []models.Product{}
		if err != nil {
			// Handle error
			panic(err)
		} else {
			var ttyp models.Product
			for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
				if t, ok := item.(models.Product); ok {
					// fmt.Printf("Tweet by %s: %f\n", t.Name, t.Price)
					products = append(products, t)
				}
			}
		}
		if err != nil {
			return
		}
		c.JSON(200, products)
	})
	r.Run() // listen and server on 0.0.0.0:8080
}