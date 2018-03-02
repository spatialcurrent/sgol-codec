package codec

type OperationPlaceholder struct {
	*AbstractOperation
}

func NewOperationPlaceholder(key string) OperationPlaceholder {
	return OperationPlaceholder{
		AbstractOperation: &AbstractOperation{Type: "DISCARD"},
	}
}

func (op OperationPlaceholder) Sgol() (string, error) {
	return "PLACEHOLDER", nil
}
