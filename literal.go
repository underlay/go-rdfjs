package rdfjs

import (
	"encoding/json"
	"fmt"
)

// LiteralType is the TermType literals
const LiteralType = "Literal"

var literalTermType = termType{LiteralType}

// A Literal is literal term
type Literal struct {
	value    string
	language string
	datatype *NamedNode
}

// NewLiteral creates a new literal
func NewLiteral(value string, language string, datatype *NamedNode) *Literal {
	return &Literal{value, language, datatype}
}

func (node *Literal) String() string {
	if node.datatype == nil || node.datatype.value == XSDString.value {
		return fmt.Sprintf("\"%s\"", escape(node.value))
	} else if node.language != "" && node.datatype.value == RDFLangString.value {
		return fmt.Sprintf("\"%s\"@%s", escape(node.value), node.language)
	}
	return fmt.Sprintf("\"%s\"^^<%s>", escape(node.value), node.datatype.value)
}

// TermType of a literal is "Literal"
func (node *Literal) TermType() string { return LiteralType }

// Value of a literal is its string value, without the datatype and langauge
func (node *Literal) Value() string { return node.value }

// Language of a literal is the literal's language tag
func (node *Literal) Language() string { return node.language }

// Datatype of a literal returns the literal's datatype
func (node *Literal) Datatype() *NamedNode {
	if node.datatype != nil {
		return node.datatype
	}
	return XSDString
}

// Equal checks for functional equivalence of terms
func (node *Literal) Equal(term Term) bool {
	if term.TermType() != LiteralType || term.Value() != node.value {
		return false
	}

	switch term := term.(type) {
	case TermLiteral:
		return term.Language() == node.language && term.Datatype().Equal(node.datatype)
	default:
		return false
	}
}

// MarshalJSON marshals the literal into a byte slice
func (node *Literal) MarshalJSON() ([]byte, error) {
	result := &term{value{literalTermType, node.value}, node.language, nil}
	if node.datatype != nil {
		result.Datatype = &value{namedNodeTermType, node.datatype.value}
	}
	return json.Marshal(result)
}

// UnmarshalJSON umarshals a byte slice into the literal
func (node *Literal) UnmarshalJSON(data []byte) error {
	t := &term{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return err
	} else if t.TermType != LiteralType {
		return ErrTermType
	}
	node.value = t.Value
	if t.Datatype != nil {
		if t.Datatype.TermType != NamedNodeType {
			return ErrTermType
		}
		node.datatype = &NamedNode{t.Datatype.Value}
		if t.Datatype.Value == RDFLangString.value {
			node.language = t.Language
		}
	}
	return nil
}
