package rdf

import (
	"encoding/json"
	"errors"
)

// ErrTermType indicates an unexpected or mismatching term types
var ErrTermType = errors.New("Mismatching term types")

// ErrParseTerm indicates that a string could not parse into a term
var ErrParseTerm = errors.New("Error parsing term")

// XSDString is the default datatype for literals
var XSDString = &NamedNode{"http://www.w3.org/2001/XMLSchema#string"}

// RDFLangString is the datatype for language-tagged literals
var RDFLangString = &NamedNode{"http://www.w3.org/1999/02/22-rdf-syntax-ns#langString"}

// Term is the interface that all terms satisfy
type Term interface {
	String() string
	TermType() string
	Value() string
	Equal(term Term) bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// TermLiteral is an interface for literal terms.
// This is only used in the Equal() method of the Literal struct.
// You shouldn't have to worry about this unless you're going to
// use external implementations of Term, in which case you need
// to make sure that literal terms implement this interface.
type TermLiteral interface {
	Term
	Language() string
	Datatype() Term
}

type termType struct {
	TermType string `json:"termType"`
}

type value struct {
	termType
	Value string `json:"value"`
}

type term struct {
	value
	Language string `json:"language"`
	Datatype *value `json:"datatype,omitempty"`
}

// UnmarshalTerm unmarshals a byte slice into a Term
func UnmarshalTerm(data []byte) (Term, error) {
	t := &term{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}

	switch t.TermType {
	case NamedNodeType:
		return &NamedNode{t.Value}, nil
	case BlankNodeType:
		return &BlankNode{t.Value}, nil
	case LiteralType:
		literal := &Literal{value: t.Value}
		if t.Datatype != nil &&
			t.Datatype.TermType == NamedNodeType &&
			t.Datatype.Value != XSDString.value {
			literal.datatype = &NamedNode{t.Datatype.Value}
			if t.Datatype.Value == RDFLangString.value {
				literal.language = t.Language
			}
		}
		return literal, nil
	case DefaultGraphType:
		return &DefaultGraph{}, nil
	case VariableType:
		return &Variable{t.Value}, nil
	}
	return nil, nil
}

// MarshalTerm marshals a Term into a byte slice
func MarshalTerm(t Term) ([]byte, error) {
	var result interface{}
	tt := termType{t.TermType()}
	switch tt.TermType {
	case NamedNodeType:
		result = &value{tt, t.Value()}
	case BlankNodeType:
		result = &value{tt, t.Value()}
	case LiteralType:
		switch t := t.(type) {
		case TermLiteral:
			d := t.Datatype()
			datatype := &value{termType{d.TermType()}, d.Value()}
			result = &term{value{tt, t.Value()}, t.Language(), datatype}
		default:
			return nil, ErrTermType
		}
	case DefaultGraphType:
		result = tt
	case VariableType:
		result = &value{tt, t.Value()}
	}
	return json.Marshal(result)
}
