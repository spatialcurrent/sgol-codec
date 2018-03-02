package codec

type AbstractOperationInput struct {
	InputType string `json:"input_type" bson:"input_type" yaml:"input_type" hcl:"input_type"`
}

func (op AbstractOperationInput) GetInputType() string {
	return op.InputType
}
