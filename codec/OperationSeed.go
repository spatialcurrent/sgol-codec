package codec

type OperationSeed struct {
	*AbstractOperation
}

func NewOperationSeed() OperationSeed {
	return OperationSeed{
		&AbstractOperation{
			Type: "SEED",
		},
	}
}

func (op OperationSeed) Sgol() (string, error) {
	return "SEED", nil
}
