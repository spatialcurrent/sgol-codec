package codec

import (
	"strings"
)

type OperationSelect struct {
	*AbstractOperation
	Seeds                 []string `json:"seeds" bson:"seeds" yaml:"seeds" hcl:"seeds"`
	Entities              []string `json:"entities" bson:"entities" yaml:"entities" hcl:"entities"`
	Edges                 []string `json:"edges" bson:"edges" yaml:"edges" hcl:"edges"`
	FilterFunctions       map[string][]FilterFunction `json:"filter_functions" bson:"filter_functions" yaml:"filter_functions" hcl:"filter_functions"`
	UpdateKey             string `json:"update_key" bson:"update_key" yaml:"update_key" hcl:"update_key"`
	UpdateFilterFunctions map[string][]FilterFunction `json:"update_filter_functions" bson:"update_filter_functions" yaml:"update_filter_functions" hcl:"update_filter_functions"`
}

func (op OperationSelect) Sgol() (string, error) {

	if len(op.Seeds) > 0 {
		return "SELECT "+strings.Join(op.Seeds, ","), nil
	}

  s := "SELECT "

  if len(op.Entities) > 0 {
		for i, x := range op.Entities {
			if i > 0 {
				s += ","
			}
			s += "$"+x
		}
	} else if len(op.Edges) > 0 {
		for i, x := range op.Edges {
			if i > 0 {
				s += ","
			}
			s += "$"+x
		}
	}

  if len(op.FilterFunctions) > 0 {
		//for k, v := range op.FilterFunctions {
			// TBD
		//}
	}

	if len(op.UpdateKey) > 0 {
		s += "UPDATE "+op.UpdateKey+" FILTER"
		//for k, v := range op.UpdateFilterFunctions {
			// TBD
		//}
	}

	return s, nil
}
