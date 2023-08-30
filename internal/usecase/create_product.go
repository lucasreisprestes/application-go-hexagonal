package usecase

import "github.com/lucasreisprestes/application-go-hexagonal/internal/entity"

type CreateProductInputDto struct {
	Name  string
	Price float64
}

type CreateProductOutputDto struct {
	ID    string
	Name  string
	Price float64
}

type CreateProductUserCase struct {
	ProductRepository entity.ProductRepository
}

func NewCreateProductsUseCase(productRepository entity.ProductRepository) *CreateProductUserCase {
	return &CreateProductUserCase{ProductRepository: productRepository}
}

func (u *CreateProductUserCase) Execute(input CreateProductInputDto) (*CreateProductOutputDto, error) {
	product := entity.NewProduct(input.Name, input.Price)
	err := u.ProductRepository.Create(product)
	if err != nil {
		return nil, err
	}

	return &CreateProductOutputDto{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}
