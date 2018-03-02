package codec

type OperationDiscard struct {
	*AbstractOperationKey
}

func NewOperationDiscard(key string) OperationDiscard {
	return OperationDiscard{
		&AbstractOperationKey{
			AbstractOperation: &AbstractOperation{Type: "DISCARD"},
			Key:               key,
		},
	}
}

func (op OperationDiscard) Sgol() (string, error) {
	return "DISCARD " + op.Key, nil
}
