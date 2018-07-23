package codec

type AbstractOperationSeeds struct {
	Seeds                   SeedCollection  `json:"seeds" bson:"seeds" yaml:"seeds" hcl:"seeds"`
}

func (op *AbstractOperationSeeds) SetSeeds(seeds SeedCollection) {
	op.Seeds = seeds
}

func (op AbstractOperationSeeds) HasSeeds() bool {
	return len(op.Seeds.Vertices) > 0
}

func (op AbstractOperationSeeds) GetSeeds() []string {
	return op.Seeds.Vertices
}
