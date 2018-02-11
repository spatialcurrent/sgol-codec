package codec

type OperationRun struct {
	*AbstractOperation
	Operations []string `json:"operations" bson:"operations" yaml:"operations" hcl:"operations"`
}

func (op OperationRun) Sgol() (string, error) {
	return "RUN", nil
}
