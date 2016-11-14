package main

import (
	"go-server/models"
	"log"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

var products []models.Product

// function to get all products from elasticsearch
func GetProducts(client *elastic.Client, result []models.Product) {
	// search
	searchResult, err := client.Search().
		Index("odoo").   // search in index "odoo"
		Type("product"). // search in type "product"
		From(0).
		Size(10000).  // take documents 0-9
		Pretty(true). // pretty print request and response JSON
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
	copy(result, products)
}

//CreateSchema
func CreateSchema() graphql.Schema {
	var userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Product",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"product": &graphql.Field{
					Type: userType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						idQuery, isOK := p.Args["id"].(string)
						if isOK {
							return idQuery, nil
						}
						return nil, nil
					},
				},
			},
		},
	)
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	return schema
}

//Query
func Query(schema graphql.Schema, query string) *graphql.Result {
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	return r
}

func main() {
	router := gin.Default()
	// connect elasticsearch
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)
	if err != nil {
		panic("Connection to elasticsearch failed")
	}
	GetProducts(client, products)
	schema := CreateSchema()
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/query", func(c *gin.Context) {
		query := c.Query("query") // shortcut for c.Request.URL.Query().Get("lastname")
		result := Query(schema, query)
		c.JSON(200, result)
	})
	router.Run(":8080")
}

/*func main() {

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


	})
	r.Run() // listen and server on 0.0.0.0:8080
}*/
