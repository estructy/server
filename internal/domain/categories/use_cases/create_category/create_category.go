// Package createcategory provides functionality to create a new category.
package createcategory

type CreateCategoryUseCase struct{}

func NewCreateCategoryUseCase() *CreateCategoryUseCase {
	return &CreateCategoryUseCase{}
}

func (u *CreateCategoryUseCase) Execute(name string) error {
	return nil
}
