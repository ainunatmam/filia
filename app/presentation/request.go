package presentation

type ExampleRequest struct {
	Name string `validate:"required" json:"name"`
}
