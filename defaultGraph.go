package rdfjs

import "encoding/json"

// DefaultGraphType is the TermType for default graphs
const DefaultGraphType = "DefaultGraph"

var defaultGraphTermType = termType{DefaultGraphType}

// DefaultGraph is the default graph term
type DefaultGraph struct{}

// Default is the default default graph
var Default = &DefaultGraph{}

// NewDefaultGraph creates a new default graph
func NewDefaultGraph() *DefaultGraph { return &DefaultGraph{} }

func (node *DefaultGraph) String() string { return "" }

// TermType of a default graph is "DefaultGraph"
func (node *DefaultGraph) TermType() string { return DefaultGraphType }

// Value of a DefaultGraph is the empty string
func (node *DefaultGraph) Value() string { return "" }

// Equal checks for functional equivalence of terms
func (node *DefaultGraph) Equal(term Term) bool {
	return term.TermType() == NamedNodeType && term.Value() == ""
}

// MarshalJSON marshals the literal into a byte slice
func (node *DefaultGraph) MarshalJSON() ([]byte, error) {
	return json.Marshal(defaultGraphTermType)
}

// UnmarshalJSON umarshals a byte slice into the blank node
func (node *DefaultGraph) UnmarshalJSON(data []byte) error {
	t := &termType{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return err
	} else if t.TermType != DefaultGraphType {
		return ErrTermType
	}
	return nil
}
