package models

// Product is a structure used for serilizing/deserializing data in ES.
type Product struct {
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Categories []string `json:"categories"`
}

