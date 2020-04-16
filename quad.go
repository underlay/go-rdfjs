package rdfjs

import "encoding/json"

// This internal struct is used for serialization and deserialization.
// Since we don't know what types each of the terms are, we need to
// first deserialize them into json.RawMessages and then dispatch that
// to UnmarshalTerm([]byte): Term
type quad struct {
	Subject   json.RawMessage `json:"subject"`
	Predicate json.RawMessage `json:"predicate"`
	Object    json.RawMessage `json:"object"`
	Graph     json.RawMessage `json:"graph"`
}

// A Quad is a 4-tuple of terms
type Quad [4]Term

// NewQuad creates a new quad. If graph is nil, Default will be used.
func NewQuad(subject, predicate, object, graph Term) *Quad {
	if graph == nil {
		graph = Default
	}
	return &Quad{subject, predicate, object, graph}
}

// Subject returns the first term
func (q *Quad) Subject() Term { return q[0] }

// Predicate returns the second term
func (q *Quad) Predicate() Term { return q[1] }

// Object returns the third term
func (q *Quad) Object() Term { return q[2] }

// Graph returns the fourth term
func (q *Quad) Graph() Term { return q[3] }

// UnmarshalJSON unmarshals a byte slice into a quad
func (q *Quad) UnmarshalJSON(data []byte) error {
	Q := &quad{}
	err := json.Unmarshal(data, Q)
	if err != nil {
		return err
	}

	q[0], err = UnmarshalTerm(Q.Subject)
	if err != nil {
		return err
	}

	q[1], err = UnmarshalTerm(Q.Predicate)
	if err != nil {
		return err
	}

	q[2], err = UnmarshalTerm(Q.Object)
	if err != nil {
		return err
	}

	q[3], err = UnmarshalTerm(Q.Graph)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON marshals a quad into a byte slice
func (q *Quad) MarshalJSON() ([]byte, error) {
	s, err := q[0].MarshalJSON()
	if err != nil {
		return nil, err
	}
	p, err := q[1].MarshalJSON()
	if err != nil {
		return nil, err
	}
	o, err := q[2].MarshalJSON()
	if err != nil {
		return nil, err
	}
	g, err := q[3].MarshalJSON()
	if err != nil {
		return nil, err
	}
	result := &quad{s, p, o, g}
	return json.Marshal(result)
}
