package codec

type OperationDiscard struct {
	*AbstractOperation
}

func NewOperationDiscard() OperationDiscard {
	return OperationDiscard{
		AbstractOperation: &AbstractOperation{Type: "DISCARD"},
	}
}

func (op OperationDiscard) Sgol() (string, error) {
	return "DISCARD", nil
}
