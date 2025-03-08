package types

type Student struct {
	ID    int    `validate:"required"`
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required"`
}
