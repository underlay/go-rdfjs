# go-rdfjs

> Go interfaces for RDF terms and quads

This is a zero-dependency module implementing the RDF data model. It's a faithful and idiomatic adaptation of the [RDFJS JSON-based data model](http://rdf.js.org/data-model-spec/) and it comes with JSON marshalers and unmarshalers for easy interop with JavaScript libraries like [n3.js](https://github.com/rdfjs/N3.js), [jsonld.js](https://github.com/digitalbazaar/jsonld.js), and [graphy.js](https://github.com/blake-regalia/graphy.js).

## Terms

The five term types are the structs `NamedNode`, `BlankNode`, `Literal`, `DefaultGraph`, and `Variable`. They all satisfy the `Term` interface:

```golang
type Term interface {
  TermType() string
  Value() string
  Equal(term Term) bool
  MarshalJSON() ([]byte, error)
  UnmarshalJSON(data []byte) error
}
```

Named nodes, blank nodes, and variables all have single string value and can be created with the constructors:

- `NewNamedNode(value: string) *NamedNode`
- `NewBlankNode(value: string) *BlankNode`
- `NewVariable(value: string) *Variable`

It's up to the user to validate that named node values are valid IRIs, that blank node values begin with `_:`, and that variable values begin with `?`.

Default graphs have no value - the `Value() string` method always returns the empty string, and the type `DefaultGraph` is just `struct{}`. A default graph can be created with the constructor:

- `NewDefaultGraph() *DefaultGraph`

Literals have a string value, a string language, and a named node datatype, which may be `nil` (indicating the default datatype of `xsd:string`). A literal can be created with the constructor:

- `NewLiteral(value: string, language: string, datatype: *DefaultGraph) *Literal`

If the given datatype does not have a value of `rdf:langString`, then the resulting `*Literal` will have no langauge, even if one is passed. You can use the exported `var RDFLangString *NamedNode` value to avoid repeatedly constructing an `rdf:langString` term.

### Marshal and Unmarshal generic terms

Each of the term structs implements `MarshalJSON` and `UnmarshalJSON`; however it is often necessary to marshal and unmarshal a term without knowing its type in advance:

- `MarshalTerm(t Term) ([]byte, error)`
- `UnmarshalTerm(data []byte) (Term, error)`

## Quads

```golang
type Quad = [4]Term
```

Quads are represented interally as 4-tuples of `Term` interfaces. This was chosen instead of a struct type to support advanced uses like arithmetic or permutations of term positions. Quad terms can be accessed by name with the `.Subject(): Term`, `.Predicate(): Term`, `.Object(): Term`, and `.Graph(): Term` methods.

Quads implement also `MarshalJSON` and `UnmarshalJSON`.

A new quad can be created with the constructor:

- `NewQuad(subject, predicate, object, graph Term) *Quad`

The user is responsible for checking that the terms of a quad are valid for their positions (literals as subjects, etc). If `graph` is `nil`, the exported default graph `var Default *DefaultGraph` will be used.
