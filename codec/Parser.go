package codec

import (
	"regexp"
	"strconv"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-graph/graph"
)

type Parser struct {
	Clauses  []string `json:"clauses" hcl:"clauses"`
	Entities []string `json:"entities" hcl:"entities"`
	Edges    []string `json:"edges" hcl:"edges"`
}

func (p *Parser) Rejoin(block []string) string {
	text := ""

	for _, token := range block {
		if strings.Contains(token, " ") && !(strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"")) {
			text += "\"" + token + "\""
		} else {
			text += token + " "
		}
	}

	return text
}

func (p *Parser) ParseEntities(text string) []string {
	if len(text) > 0 {
		if strings.HasPrefix(text, "$") {
			if len(text) > 1 {
				if strings.Contains(text, ",") {
					parts := strings.Split(text, ",")
					result := make([]string, len(parts))
					for i, x := range parts {
						result[i] = x[1:len(x)]
					}
					return result
				} else {
					return []string{text[1:len(text)]}
				}
			} else {
				return p.Entities
			}
		} else if text == "_" || "text" == "-" {
			return []string{}
		} else {
			return []string{text}
		}
	} else {
		return []string{}
	}
}

func (p *Parser) ParseEdges(text string) []string {
	if len(text) > 0 {
		if strings.HasPrefix(text, "$") {
			if len(text) > 1 {
				if strings.Contains(text, ",") {
					parts := strings.Split(text, ",")
					result := make([]string, len(parts))
					for i, x := range parts {
						result[i] = x[1:len(x)]
					}
					return result
				} else {
					return []string{text[1:len(text)]}
				}
			} else {
				return p.Edges
			}
		} else {
			return []string{}
		}
	} else {
		return []string{}
	}
}

func (p *Parser) ParseUpdate(block []string) (string, map[string][]FilterFunction, error) {
	updateKey := block[0]
	filterFunctionsByEntity := map[string][]FilterFunction{}
	if len(block) > 3 {
		ff, err := p.ParseFilterFunctions(p.Rejoin(block[4:]))
		if err != nil {
			return updateKey, filterFunctionsByEntity, err
		}
		if len(ff) > 0 {
			filterFunctionsByEntity[block[2][1:]] = ff
		}
	}

	/*if len(block) > 6 {
		filterFunctions, err := p.ParseFilterFunctions(p.Rejoin(block[7:len(block)]))
		if err != nil {
			return op, err
		}
		if len(filterFunctions) > 0 {
			op.UpdateFilterFunctions = map[string][]FilterFunction{}
			op.UpdateFilterFunctions[block[5][1:len(block[5])]] = filterFunctions
		}
	}*/

	return updateKey, filterFunctionsByEntity, nil
}

func (p *Parser) ParseFilter(block []string, entities []string) (map[string][]FilterFunction, string, map[string][]FilterFunction, error) {

  filterFunctionsByEntity := map[string][]FilterFunction{}
  updateKey := ""
	updateFilterFunctionsByEntity := map[string][]FilterFunction{}

	if StringSliceContains(block, "UPDATE") {
		blockUpdateIndex := StringSliceIndex(block, "UPDATE")
		ff, err := p.ParseFilterFunctions(p.Rejoin(block[2:blockUpdateIndex]))
		if err != nil {
			return filterFunctionsByEntity, updateKey, updateFilterFunctionsByEntity, err
		}
		if len(ff) > 0 {
			filterFunctionsByEntity[block[0][1:len(block[0])]] = ff
		}
    updateKey, updateFilterFunctionsByEntity, err = p.ParseUpdate(block[blockUpdateIndex+1:])
	} else {
		if len(block) > 2 && block[1] == "WITH" {
			ff, err := p.ParseFilterFunctions(p.Rejoin(block[2:len(block)]))
			if err != nil {
				return filterFunctionsByEntity, updateKey, updateFilterFunctionsByEntity, err
			}
			if len(ff) > 0 {
				filterFunctionsByEntity[block[0][1:len(block[0])]] = ff
			}
		} else {
			ff, err := p.ParseFilterFunctions(p.Rejoin(block))
			if err != nil {
				return filterFunctionsByEntity, updateKey, updateFilterFunctionsByEntity, err
			}
			if len(ff) > 0 {
				for _, entity := range entities {
					filterFunctionsByEntity[entity] = ff
				}
			}
		}
	}
	return filterFunctionsByEntity, updateKey, updateFilterFunctionsByEntity, nil
}

func (p *Parser) ParseQueryFunctions(text string) ([]QueryFunction, error) {

	functions := make([]QueryFunction, 0)
	if len(text) > 0 {

		re, err := regexp.Compile("(\\s*)(?P<name>([a-zA-Z]+))(\\s*)\\((\\s*)(?P<args>(.)*?)(\\s*)\\)(\\s*)")
		if err != nil {
			return functions, err
		}

		matches := re.FindAllStringSubmatch(strings.TrimSpace(text), -1)
		for _, m := range matches {
			g := map[string]string{}
			for i, name := range re.SubexpNames() {
				if i != 0 {
					g[name] = m[i]
				}
			}
			fn := QueryFunction{Name: strings.TrimSpace(g["name"])}
			if args_text, ok := g["args"]; ok {
				if len(args_text) > 0 {
					args := []string{}

					re2, err := regexp.Compile("(\\s*)(?P<value>((\"([^\"]+?)\")|([^,\\s]+)))(\\s*)")
					if err != nil {
						return functions, err
					}

					matches2 := re2.FindAllStringSubmatch(args_text, -1)
					for _, m2 := range matches2 {
						g2 := map[string]string{}
						for i, name := range re2.SubexpNames() {
							if i != 0 {
								g2[name] = m2[i]
							}
						}
						if value, ok := g2["value"]; ok {
							value = strings.TrimSpace(value)
							if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
								args = append(args, value[1:len(value)-1])
							} else {
								args = append(args, value)
							}
						}
					}
					fn.Args = args
				}
			}
			functions = append(functions, fn)
		}

	}

	return functions, nil
}

func (p *Parser) ParseFilterFunctions(text string) ([]FilterFunction, error) {

	filterFunctions := []FilterFunction{}
	queryFunctions, err := p.ParseQueryFunctions(text)
	if err != nil {
		return filterFunctions, err
	}
	if len(queryFunctions) > 0 {
		for _, qf := range queryFunctions {
			qf_name_lc := strings.ToLower(qf.Name)
			if qf_name_lc == "collectioncontains" {
				ff := FilterFunctionCollectionContains{
					AbstractFilterFunction: &AbstractFilterFunction{Name: qf.Name},
					PropertyName: qf.Args[0],
					PropertyValue: qf.Args[1],
				}
				filterFunctions = append(filterFunctions, ff)
			}
		}
		return filterFunctions, nil
	} else {
		return []FilterFunction{}, nil
	}
}

func (p *Parser) ParseTokens(text string) ([]string, error) {
	tokens := []string{}

	re, err := regexp.Compile("(?P<token>((\"([^\"]+)\")|(\\S+)))")
	if err != nil {
		return tokens, err
	}

	matches := re.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		g := map[string]string{}
		for i, name := range re.SubexpNames() {
			if i != 0 {
				g[name] = m[i]
			}
		}
		tokens = append(tokens, strings.TrimSpace(g["token"]))
	}

	return tokens, nil

}

func (p *Parser) ParseBlocks(tokens []string) [][]string {
	blocks := [][]string{}
	block := []string{}

	for _, token := range tokens {
		token_uc := strings.ToUpper(token)
		if StringSliceContains(p.Clauses, token_uc) {
			if len(block) > 0 {
				blocks = append(blocks, block)
			}
			block = []string{token_uc}
		} else {
			block = append(block, token)
		}
	}

	if len(block) > 0 {
		blocks = append(blocks, block)
	}

	return blocks
}

func (p *Parser) ParseSelect(block []string) (OperationSelect, error) {
	op := OperationSelect{
		AbstractOperation: &AbstractOperation{Type: "SELECT"},
	}

	if block[1] == "$" {
		op.Edges = p.Edges
	}

	if strings.HasPrefix(block[1], "$") {
		op.Entities = p.ParseEntities(block[1])
	} else {
		op.Entities = p.Entities
		op.Seeds = strings.Split(block[1], ",")
	}

	if len(block) > 2 {
		block_2_uc := strings.ToUpper(block[2])
		if block_2_uc == "UPDATE" {
			updateKey, updateFilterFunctions, err := p.ParseUpdate(block[3:])
			if err != nil {
				return op, err
			}
			op.UpdateKey = updateKey
			op.UpdateFilterFunctions = updateFilterFunctions
		} else if block_2_uc == "FILTER" {
			filterFunctions, updateKey, updateFilterFunctions, err := p.ParseFilter(block[3:], p.Entities)
			if err != nil {
				return op, err
			}
			op.UpdateKey = updateKey
			op.FilterFunctions = filterFunctions
			op.UpdateFilterFunctions = updateFilterFunctions
		}
	}

	return op, nil
}

func (p *Parser) ParseNav(block []string) (OperationNav, error) {
	op := OperationNav{
		AbstractOperation: &AbstractOperation{Type: "NAV"},
		Source: block[1],
		Destination: block[3],
		SeedMatching: "RELATED",
	}

	if strings.HasPrefix(block[1], "$") {
		op.Entities = p.ParseEntities(block[1])
		op.Direction = "INCOMING"
		op.EdgeIdentifierToExtract = "SOURCE"
		if block[3] != "INPUT" {
			op.Seeds = []string{block[3]}
		}
	} else {
		op.Entities = p.ParseEntities(block[3])
		op.Direction = "OUTGOING"
		op.EdgeIdentifierToExtract = "DESTINATION"
		if strings.HasPrefix(block[3], "$") && block[1] != "INPUT" {
			op.Seeds = []string{block[1]}
		}
	}

	if strings.Contains(block[2], "*") {
		block_2_parts := strings.Split(block[2], "*")
		op.Edges = p.ParseEdges(block_2_parts[0])
		depth, err := strconv.Atoi(block_2_parts[1])
		if err != nil {
			return op, err
		}
		op.Depth = depth
	} else {
		op.Edges = p.ParseEdges(block[2])
		op.Depth = 1
	}

	if len(block) > 4 {
		block_4_uc := strings.ToUpper(block[4])
		if block_4_uc == "UPDATE" {
			updateKey, updateFilterFunctions, err := p.ParseUpdate(block[5:])
			if err != nil {
				return op, err
			}
			op.UpdateKey = updateKey
			op.UpdateFilterFunctions = updateFilterFunctions
		} else if block_4_uc == "FILTER" {
			filterFunctions, updateKey, updateFilterFunctions, err := p.ParseFilter(block[5:], p.Entities)
			if err != nil {
				return op, err
			}
			op.UpdateKey = updateKey
			op.FilterFunctions = filterFunctions
			op.UpdateFilterFunctions = updateFilterFunctions
		}
	}

	return op, nil
}

func (p *Parser) ParseHas(block []string) (OperationHas, error) {
	op := OperationHas{
		AbstractOperation: &AbstractOperation{Type: "HAS"},
		Source: block[1],
		Destination: block[3],
		Entities: p.Entities,
		Edges: p.ParseEdges(block[2]),
		SeedMatching: "RELATED",
		Depth: 1,
	}

	if strings.HasPrefix(block[1], "$") {
		op.Direction = "INCOMING"
	} else {
		op.Direction = "OUTGOING"
	}

	if block[1] == "INPUT" {
		op.EdgeIdentifierToExtract = "SOURCE"
		op.SeedsB = []string{block[3]}
	} else if block[3] == "INPUT" {
		op.EdgeIdentifierToExtract = "DESTINATION"
		op.SeedsB = []string{block[1]}
	} else {
		return op, errors.New("HAS clause required INPUT.")
	}

	if len(block) > 4 {
		block_4_uc := strings.ToUpper(block[4])
		if block_4_uc == "UPDATE" {
			updateKey, updateFilterFunctions, err := p.ParseUpdate(block[5:])
			if err != nil {
				return op, err
			}
			op.UpdateKey = updateKey
			op.UpdateFilterFunctions = updateFilterFunctions
		} else if block_4_uc == "FILTER" {
			filterFunctions, updateKey, updateFilterFunctions, err := p.ParseFilter(block[5:], p.Entities)
			if err != nil {
				return op, err
			}
			op.UpdateKey = updateKey
			op.FilterFunctions = filterFunctions
			op.UpdateFilterFunctions = updateFilterFunctions
		}
	}

	return op, nil
}

func (p *Parser) ParseOperations(blocks [][]string) ([]graph.Operation, string, error) {
	operations := make([]graph.Operation, 0)
	output_type := ""
	for _, block := range blocks {
		switch block[0] {
		case "INIT":
			operations = append(operations, OperationInit{
				&AbstractOperationKey{
					AbstractOperation: &AbstractOperation{Type: "INIT"},
					Key: block[1],
				},
			})
		case "ADD":
			operations = append(operations, OperationAdd{
				&AbstractOperationKey{
					AbstractOperation: &AbstractOperation{Type: "ADD"},
					Key: block[1],
				},
			})
		case "DISCARD":
			operations = append(operations, OperationDiscard{
				&AbstractOperationKey{
					AbstractOperation: &AbstractOperation{Type: "DISCARD"},
					Key: block[1],
				},
			})
		case "RELATE":
			operations = append(operations, OperationRelate{
				AbstractOperation: &AbstractOperation{Type: "RELATE"},
				Keys: []string{block[1]},
			})
		case "LIMIT":
			limit_int, err := strconv.Atoi(block[1])
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, OperationLimit{
				AbstractOperation: &AbstractOperation{Type: "LIMIT"},
				Limit: limit_int,
			})
		case "SELECT":
			op, err := p.ParseSelect(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, op)
		case "NAV":
			op, err := p.ParseNav(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, op)
		case "HAS":
			op, err := p.ParseHas(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, op)
		case "FETCH":
			operations = append(operations, OperationFetch{
				&AbstractOperationKey{
					AbstractOperation: &AbstractOperation{Type: "DISCARD"},
					Key: block[1],
				},
			})
		case "SEED":
			operations = append(operations, OperationSeed{&AbstractOperation{Type: "SEED"}})
		case "RUN":
			operations = append(operations, OperationRun{
				AbstractOperation: &AbstractOperation{Type: "RUN"},
				Operations: []string{},
			})
		case "OUTPUT":
			output_type = block[1]
		}
	}
	return operations, output_type, nil
}

func (p *Parser) ParseQuery(q string) (graph.OperationChain, error) {

	chain := OperationChain{}

	tokens, err := p.ParseTokens(q)
	if err != nil {
		return OperationChain{}, err
	}

	blocks := p.ParseBlocks(tokens)

	operations, outputType, err := p.ParseOperations(blocks)
	if err != nil {
		return chain, err
	}

	return NewOperationChain("chain", operations, outputType), nil
}
