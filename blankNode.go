package rdfjs

import "encoding/json"

// BlankNodeType is the TermType blank nodes
const BlankNodeType = "BlankNode"

var blankNodeTermType = termType{BlankNodeType}

// A BlankNode is labelled blank node
type BlankNode struct{ value string }

// NewBlankNode creates a new blank node
func NewBlankNode(value string) *BlankNode { return &BlankNode{value} }

// TermType of a blank node is "BlankNode"
func (node *BlankNode) TermType() string { return BlankNodeType }

// Value of a blank node is its label, which will always begin with "_:"
func (node *BlankNode) Value() string { return node.value }

// Equal checks for functional equivalence of terms
func (node *BlankNode) Equal(term Term) bool {
	return term.TermType() == BlankNodeType && term.Value() == node.value
}

// MarshalJSON marshals the blank node into a byte slice
func (node *BlankNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(value{blankNodeTermType, node.value})
}

// UnmarshalJSON umarshals a byte slice into the blank node
func (node *BlankNode) UnmarshalJSON(data []byte) error {
	v := &value{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	} else if v.TermType != BlankNodeType {
		return ErrTermType
	}
	node.value = v.Value
	return nil
}
