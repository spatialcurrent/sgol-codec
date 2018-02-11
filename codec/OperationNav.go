package codec

type OperationNav struct {
	*AbstractOperation
	Seeds                 []string `json:"seeds" bson:"seeds" yaml:"seeds" hcl:"seeds"`
	Source                  string `json:"source" bson:"source" yaml:"source" hcl:"source"`
	Destination             string `json:"destination" bson:"destination" yaml:"destination" hcl:"destination"`
	Entities              []string `json:"entities" bson:"entities" yaml:"entities" hcl:"entities"`
	Edges                 []string `json:"edges" bson:"edges" yaml:"edges" hcl:"edges"`
	Direction               string `json:"direction" bson:"direction" yaml:"direction" hcl:"direction"`
	EdgeIdentifierToExtract string `json:"edgeIdentifierToExtract" bson:"edgeIdentifierToExtract" yaml:"edgeIdentifierToExtract" hcl:"edgeIdentifierToExtract"`
	SeedMatching            string `json:"seedMatching" bson:"seedMatching" yaml:"seedMatching" hcl:"seedMatching"`
	Depth                   int `json:"depth" bson:"depth" yaml:"depth" hcl:"depth"`
	FilterFunctions       map[string][]FilterFunction `json:"filter_functions" bson:"filter_functions" yaml:"filter_functions" hcl:"filter_functions"`
	UpdateKey             string `json:"update_key" bson:"update_key" yaml:"update_key" hcl:"update_key"`
	UpdateFilterFunctions map[string][]FilterFunction `json:"update_filter_functions" bson:"update_filter_functions" yaml:"update_filter_functions" hcl:"update_filter_functions"`
}


func (op OperationNav) Sgol() (string, error) {
  return "NAV", nil
}
