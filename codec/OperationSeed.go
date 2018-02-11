package codec

type OperationSeed struct {
	*AbstractOperation
}

func (op OperationSeed) Sgol() (string, error) {
	return "SEED", nil
}
