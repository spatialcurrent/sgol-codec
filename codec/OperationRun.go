package codec

type OperationRun struct {
	*AbstractOperation
	Operations []string `json:"operations" bson:"operations" yaml:"operations" hcl:"operations"`
}

func NewOperationRun(operations []string) OperationRun {
	return OperationRun{
		AbstractOperation: &AbstractOperation{Type: "RUN"},
		Operations:        operations,
	}
}

func (op OperationRun) Sgol() (string, error) {
	return "RUN", nil
}
