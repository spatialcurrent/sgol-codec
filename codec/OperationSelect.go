package codec

import (
	//"fmt"
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-graph/graph"
)

type OperationSelect struct {
	*AbstractOperation
	*AbstractOperationOutput
	*AbstractOperationView
	*AbstractOperationUpdate
	Seeds  SeedCollection  `json:"seeds" bson:"seeds" yaml:"seeds" hcl:"seeds"`
	Groups GroupCollection `json:"groups" bson:"groups" yaml:"groups" hcl:"groups"`
}

func (op *OperationSelect) SetSeeds(seeds SeedCollection) {
	op.Seeds = seeds
}

func (op *OperationSelect) SetGroups(groups GroupCollection) {
	op.Groups = groups
}

func (op OperationSelect) HasGroup(group string) bool {
	return op.Groups.HasGroup(group)
}

func (op OperationSelect) HasSeeds() bool {
	return len(op.Seeds.Vertices) > 0
}

func (op OperationSelect) GetSeeds() []string {
	return op.Seeds.Vertices
}

func (op OperationSelect) GetGroupNames() []string {
	return op.Groups.GetGroups()
}

func (op OperationSelect) Validate(schema graph.Schema) error {
	schema_entities := schema.GetEntityGroupNames()
	schema_edges := schema.GetEdgeGroupNames()
	for _, g := range op.GetGroupNames() {
		if ! (StringSliceContains(schema_entities, g) || StringSliceContains(schema_edges, g)) {
			return errors.New("Group "+g+" was not found in schema entities or edges.")
		}
	}
  return nil
}

func (op OperationSelect) Sgol() (string, error) {

	if op.HasSeeds() {
		return "SELECT " + op.Seeds.Sgol(), nil
	}

	s := "SELECT " + op.Groups.Sgol()

	if op.HasFilters() {
		s += " " + op.View.Sgol()
	}

	if op.HasUpdate() {
		s += " " + op.Update.Sgol()
	}

	return s, nil
}

func NewOperationSelect() *OperationSelect {
	return &OperationSelect{
		AbstractOperation: &AbstractOperation{
			Type: "SELECT",
		},
		AbstractOperationOutput: &AbstractOperationOutput{
			OutputType: "elements",
		},
		AbstractOperationView:   &AbstractOperationView{},
		AbstractOperationUpdate: &AbstractOperationUpdate{},
	}
}
