# go-rdfjs

> Go implementation of the RDF data model

This is a zero-dependency module implementing the RDF data model. It's a faithful and idiomatic adaptation of the [RDFJS JSON-based interface](http://rdf.js.org/data-model-spec/) and it comes with RDFJS-compatible JSON marshalers and unmarshalers for easy interop with JavaScript libraries like [n3.js](https://github.com/rdfjs/N3.js), [jsonld.js](https://github.com/digitalbazaar/jsonld.js), and [graphy.js](https://github.com/blake-regalia/graphy.js).

## Terms

The five term types are the structs `NamedNode`, `BlankNode`, `Literal`, `DefaultGraph`, and `Variable`. They all satisfy the `Term` interface:

```golang
type Term interface {
  String() string
  TermType() string
  Value() string
  Equal(term Term) bool
  MarshalJSON() ([]byte, error)
  UnmarshalJSON(data []byte) error
}
```

Named nodes, blank nodes, and variables all have single string value and can be created with the constructors:

```golang
NewNamedNode(value string) *NamedNode
NewBlankNode(value string) *BlankNode
NewVariable(value string) *Variable
```

It's up to the user to validate that named node values are valid IRIs. [Per the RDFJS spec](http://rdf.js.org/data-model-spec/#blanknode-interface), blank node values should **not** begin with `_:`, and variable values should **not** begin with `?`.

Default graphs have no value - the `Value() string` method always returns the empty string, and the type `DefaultGraph` is just `struct{}`. There is a "default default graph" value `var Default *DefaultGraph` that is recommended for most purposes, although new default graphs can also be created with the constructor:

```golang
NewDefaultGraph() *DefaultGraph
```

Literals have a string value, a string language, and a named node datatype, which may be `nil` (interpreted as the default datatype of `xsd:string`). A literal can be created with the constructor:

```golang
NewLiteral(value, language string, datatype *NamedNode) *Literal
```

If the given datatype does not have a value of `rdf:langString`, then the resulting `*Literal` will have no langauge, even if one is passed. You can use the exported `var RDFLangString *NamedNode` value to avoid repeatedly constructing an `rdf:langString` term.

The term structs do not have internal term type fields - the `TermType() string` method is a constant function on each struct _type_. This is done to save memory.

### Unmarshal generic terms

Each of the term structs implements `MarshalJSON` and `UnmarshalJSON`; however it is often necessary to unmarshal a term without knowing its type in advance:

```golang
UnmarshalTerm(data []byte) (Term, error)
```

### Serialize and parse strings

Terms also implement a `.String() string` method that return their N-Quads term representation (e.g. `"example"^^<http://example.com>`). The `ParseTerm(s: string): Term` function parses terms back from this format.

## Quads

```golang
type Quad [4]Term
```

Quads are represented interally as 4-tuples of `Term` interfaces. This was chosen instead of a struct type to support advanced uses like arithmetic or permutations of term positions. Quad terms can be accessed by name with the `.Subject(): Term`, `.Predicate(): Term`, `.Object(): Term`, and `.Graph(): Term` methods.

Quads also implement `MarshalJSON` and `UnmarshalJSON`, which serialize to and from the RDFJS object representation of quads (`{"subject": { }, ...}`). Internally, `Quad.UnmarshalJSON` calls the generic `UnmarshalTerm` for each of its components.

A new quad can be created with the constructor:

```golang
NewQuad(subject, predicate, object, graph Term) *Quad
```

The user is responsible for checking that the terms of a quad are valid for their positions (no literals as subjects, etc). If `graph` is `nil`, the "default default graph" `var Default *DefaultGraph` will be used.

### Serialize and parse strings

Quads also have a `.String() string` method that returns the N-Quads representation of the quad, **including a trailing period, but not including a newline**. `ParseQuad(s: string): *Quad` parses a quad back from this format.

You can also parse a slice of quads out of an `io.Reader` using `ReadQuads(input io.Reader) ([]*Quad, error)`.
