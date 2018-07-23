package codec

type OperationFetch struct {
	*AbstractOperationKey
}

func NewOperationFetch(key string) OperationFetch {
	return OperationFetch{
		&AbstractOperationKey{
			AbstractOperation: &AbstractOperation{Type: "FETCH"},
			Key:               key,
		},
	}
}

func (op OperationFetch) Sgol() (string, error) {
	return "FETCH " + op.Key, nil
}
