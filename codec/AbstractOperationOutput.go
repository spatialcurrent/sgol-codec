package codec

type AbstractOperationOutput struct {
	OutputType string `json:"output_type" bson:"output_type" yaml:"output_type" hcl:"output_type"`
}

func (op AbstractOperationOutput) GetOutputType() string {
	return op.OutputType
}
