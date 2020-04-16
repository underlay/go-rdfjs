package rdfjs

import (
	"encoding/json"
	"fmt"
)

// NamedNodeType is the TermType for IRIs
const NamedNodeType = "NamedNode"

var namedNodeTermType = termType{NamedNodeType}

// A NamedNode is an IRI term
type NamedNode struct{ value string }

// NewNamedNode creates a new IRI
func NewNamedNode(value string) *NamedNode { return &NamedNode{value} }

func (node *NamedNode) String() string { return fmt.Sprintf("<%s>", node.value) }

// TermType of an IRI is "NamedNode"
func (node *NamedNode) TermType() string { return NamedNodeType }

// Value of an IRI is its string value
func (node *NamedNode) Value() string { return node.value }

// Equal checks for functional equivalence of terms
func (node *NamedNode) Equal(term Term) bool {
	return term.TermType() == NamedNodeType && term.Value() == node.Value()
}

// MarshalJSON marshals the named node into a byte slice
func (node *NamedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(value{namedNodeTermType, node.value})
}

// UnmarshalJSON umarshals a byte slice into the named node
func (node *NamedNode) UnmarshalJSON(data []byte) error {
	v := &value{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	} else if v.TermType != NamedNodeType {
		return ErrTermType
	}
	node.value = v.Value
	return nil
}
