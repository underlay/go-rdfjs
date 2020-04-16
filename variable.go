package rdfjs

import "encoding/json"

// VariableType is the TermType for variables
const VariableType = "Variable"

var variableTermType = termType{VariableType}

// A Variable is labelled variable
type Variable struct{ value string }

// NewVariable creates a new variable
func NewVariable(value string) *Variable { return &Variable{value} }

// TermType of a variable is "Variable"
func (node *Variable) TermType() string { return VariableType }

// Value of a variable is its label, which will always begin with "_:"
func (node *Variable) Value() string { return node.value }

// Equal checks for functional equivalence of terms
func (node *Variable) Equal(term Term) bool {
	return term.TermType() == VariableType && term.Value() == node.value
}

// MarshalJSON marshals the variable into a byte slice
func (node *Variable) MarshalJSON() ([]byte, error) {
	return json.Marshal(value{variableTermType, node.value})
}

// UnmarshalJSON umarshals a byte slice into the variable
func (node *Variable) UnmarshalJSON(data []byte) error {
	v := &value{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	} else if v.TermType != VariableType {
		return ErrTermType
	}
	node.value = v.Value
	return nil
}
