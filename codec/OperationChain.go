package codec

import (
  "fmt"
  "strconv"
  "encoding/json"
  "crypto/md5"
)

import (
  "github.com/pkg/errors"
)

import (
  "github.com/spatialcurrent/go-graph/graph"
)

type OperationChain struct {
  Name string `json:"name" yaml:"name" hcl:"name"`
  Operations []graph.Operation `json:"operations" yaml:"operations" hcl:"operations"`
  OutputType string `json:"output_type" "yaml":"output_type" hcl:"output_type"`
  Limit int `json:"limit" bson:"limit" yaml:"limit" hcl:"limit"`
}

func (oc OperationChain) GetName() string {
  return oc.Name
}

func (oc OperationChain) GetOperations() []graph.Operation {
  return oc.Operations
}

func (oc OperationChain) GetOutputType() string {
  return oc.OutputType
}

func (oc OperationChain) GetLimit() int {
  return oc.Limit
}

func (oc OperationChain) GetLastOperation() (graph.Operation, error) {
  if len(oc.Operations) > 0 {
    return oc.Operations[len(oc.Operations) -1], nil
  } else {
    return OperationPlaceholder{}, errors.New("Error: The operation chain does not have any operations.")
  }
}

func (oc OperationChain) Hash() (string, error) {
  data, err := json.Marshal(oc)
  if err != nil {
    return "", err
  }
  hash := fmt.Sprintf("%x", md5.Sum(data))
  return hash, nil
}

func (chain OperationChain) Sgol() (string, error) {
  q := ""

  for i, op := range chain.Operations {
    if i > 0 {
      q += " "
    }
    x, err := op.Sgol()
    if err != nil {
      return "", err
    }
    q += x
  }

  if chain.Limit > 0 {
    q += " LIMIT "+strconv.Itoa(chain.Limit)
  }

  if len(chain.OutputType) > 0 {
    q += " OUTPUT "+chain.OutputType
  } else {
    q += " OUTPUT elements"
  }

  return q, nil
}


func NewOperationChain(name string, operations []graph.Operation, outputType string) graph.OperationChain {

  chain := OperationChain{
    Name: name,
    OutputType: outputType,
  }

  if len(operations) > 0 {
    chain.Operations = operations
    if op, err := chain.GetLastOperation(); err == nil {
      if op.GetTypeName() == "LIMIT" {
        chain.Limit = op.(OperationLimit).Limit
      }
    }
  }

  return chain

}
