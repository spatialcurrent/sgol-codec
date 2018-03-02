package codec

import (
	"strconv"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-graph/graph"
)

type OperationHas struct {
	*AbstractOperation
	*AbstractOperationOutput
	*AbstractOperationView
	*AbstractOperationUpdate
	Seeds                   SeedCollection  `json:"seeds" bson:"seeds" yaml:"seeds" hcl:"seeds"`
	SourceGroups                  GroupCollection `json:"source" bson:"source" yaml:"source" hcl:"source"`
	Edges                   GroupCollection `json:"edges" bson:"edges" yaml:"edges" hcl:"edges"`
	DestinationGroups             GroupCollection `json:"destination" bson:"destination" yaml:"destination" hcl:"destination"`
	Direction               string          `json:"direction" bson:"direction" yaml:"direction" hcl:"direction"`
	EdgeIdentifierToExtract string          `json:"edgeIdentifierToExtract" bson:"edgeIdentifierToExtract" yaml:"edgeIdentifierToExtract" hcl:"edgeIdentifierToExtract"`
	SeedMatching            string          `json:"seedMatching" bson:"seedMatching" yaml:"seedMatching" hcl:"seedMatching"`
	Depth                   int             `json:"depth" bson:"depth" yaml:"depth" hcl:"depth"`
}

func (op *OperationHas) SetSeeds(seeds SeedCollection) {
	op.Seeds = seeds
}

func (op *OperationHas) SetSource(groups GroupCollection) {
	op.SourceGroups = groups
}

func (op *OperationHas) SetEdges(groups GroupCollection) {
	op.Edges = groups
}

func (op *OperationHas) SetDestination(groups GroupCollection) {
	op.DestinationGroups = groups
}

func (op OperationHas) Validate(schema graph.Schema) error {
	schema_entities := schema.GetEntityGroupNames()
	schema_edges := schema.GetEdgeGroupNames()

	for _, g := range op.GetSourceGroups() {
		if ! StringSliceContains(schema_entities, g) {
			return errors.New("Source group "+g+" was not found in schema entities.")
		}
	}

	for _, g := range op.GetEdgeGroups() {
		if ! StringSliceContains(schema_edges, g) {
			return errors.New("Edge group "+g+" was not found in schema edges.")
		}
	}

	for _, g := range op.GetDestinationGroups() {
		if ! StringSliceContains(schema_entities, g) {
			return errors.New("Destination group "+g+" was not found in schema entities.")
		}
	}

  return nil
}


func (op OperationHas) HasSeeds() bool {
	return len(op.Seeds.Vertices) > 0
}

func (op OperationHas) GetSourceGroups() []string {
	return op.SourceGroups.GetGroups()
}

func (op OperationHas) GetEdgeGroups() []string {
	return op.Edges.GetGroups()
}

func (op OperationHas) GetDestinationGroups() []string {
	return op.DestinationGroups.GetGroups()
}

func (op OperationHas) HasGroup(group string) bool {
	return op.SourceGroups.HasGroup(group) || op.Edges.HasGroup(group) || op.DestinationGroups.HasGroup(group)
}

func (op OperationHas) HasEntityGroup(group string) bool {
	return op.SourceGroups.HasGroup(group) || op.DestinationGroups.HasGroup(group)
}

func (op OperationHas) HasEdgeGroup(group string) bool {
	return op.Edges.HasGroup(group)
}

func (op OperationHas) Sgol() (string, error) {

	s := "HAS"

	if op.HasSeeds() && op.Direction == "OUTGOING" {
		s += " " + op.Seeds.Sgol()
	} else {
		s += " " + op.SourceGroups.Sgol()
	}

	s += " " + op.Edges.Sgol()

	if op.Depth > 1 {
		s += "*" + strconv.Itoa(op.Depth)
	}

	if op.HasSeeds() && op.Direction == "INCOMING" {
		s += " " + op.Seeds.Sgol()
	} else {
		s += " " + op.DestinationGroups.Sgol()
	}

	if op.HasFilters() {
		s += " " + op.View.Sgol()
	}

	if op.HasUpdate() {
		s += " " + op.Update.Sgol()
	}

	return s, nil
}


func NewOperationHas() *OperationHas {
	return &OperationHas{
		AbstractOperation: &AbstractOperation{
			Type: "HAS",
		},
		AbstractOperationOutput: &AbstractOperationOutput{
			OutputType: "elements",
		},
		AbstractOperationView:   &AbstractOperationView{},
		AbstractOperationUpdate: &AbstractOperationUpdate{},
		SeedMatching:           "RELATED",
	}
}
