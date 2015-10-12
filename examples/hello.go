// go run hello.go
package main

import (
	"fmt"

	"github.com/ng-vu/graphql-go"
	"github.com/ng-vu/graphql-go/ql"
)

func main() {
	schema, err := graphql.NewSchema(
		ql.Object{
			Name: "RootQueryType",
			Fields: ql.FieldMap{
				"hello": {
					Type: ql.String,
					Resolve: func(v struct {
						Name string `graphql:"name"`
					}) string {
						return "world"
					},
				},
			},
		})

	query := `{ hello }`

	req, err := graphql.NewRequest(schema, query)
	if err != nil {
		fmt.Println("Req Error:", err)
		return
	}
	fmt.Println("Print:")
	fmt.Println(req.Print())

	var v struct {
		Hello string `graphql:"hello"`
	}
	err = req.Run(&v)
	if err != nil {
		fmt.Println("Exec Error:", err)
	}
	fmt.Printf("OK %v\n", v)
}
