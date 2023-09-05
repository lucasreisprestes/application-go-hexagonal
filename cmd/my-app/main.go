package my_app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi"
	"github.com/lucasreisprestes/application-go-hexagonal/internal/infra/akafka"
	"github.com/lucasreisprestes/application-go-hexagonal/internal/infra/repository"
	"github.com/lucasreisprestes/application-go-hexagonal/internal/infra/web"
	"github.com/lucasreisprestes/application-go-hexagonal/internal/usecase"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306/products)")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := repository.NewProductRepositoryMysql(db)
	createProductUserCase := usecase.NewCreateProductUseCase(repository)
	listProductUserCase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUserCase, listProductUserCase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductsHandler)

	go http.ListenAndServe(":8000", r)

	//create channel
	msgChan := make(chan *kafka.Message)
	//go-routine consume
	go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)

		if err != nil {
			fmt.Printf("Error in message s%", err)
		}
		_, err = createProductUserCase.Execute(dto)
	}
}
