package rdfjs

import (
	"regexp"
	"strings"
)

const (
	wso      = "[ \\t]*"
	iri      = "(?:<([^:]+:[^>]*)>)"
	bnode    = "(_:[a-zA-Z0-9]+)"
	variable = "(\\?[a-zA-Z0-9]+)"
	plain    = "\"([^\"\\\\]*(?:\\\\.[^\"\\\\]*)*)\""
	datatype = "(?:\\^\\^" + iri + ")"
	language = "(?:@([a-z]+(?:-[a-zA-Z0-9]+)*))"
	literal  = "(?:" + plain + "(?:" + datatype + "|" + language + ")?)"
	ws       = "[ \\t]+"
	node     = "(?:" + iri + "|" + bnode + "|" + literal + "|" + variable + ")"
	graph    = "(?:\\.|(?:" + node + wso + "\\.))"
)

var regexLiteral = regexp.MustCompile("^" + plain + "(?:" + datatype + "|" + language + ")?$")
var regexQuad = regexp.MustCompile("^" + wso + node + ws + node + ws + node + ws + graph + wso + "$")

var regexWSO = regexp.MustCompile(wso)

var regexEOLN = regexp.MustCompile("(?:\\r\\n)|(?:\\n)|(?:\\r)")

var regexEmpty = regexp.MustCompile("^" + wso + "$")

func escape(str string) string {
	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	str = strings.Replace(str, "\n", "\\n", -1)
	str = strings.Replace(str, "\r", "\\r", -1)
	str = strings.Replace(str, "\t", "\\t", -1)
	return str
}

func unescape(str string) string {
	str = strings.Replace(str, "\\\\", "\\", -1)
	str = strings.Replace(str, "\\\"", "\"", -1)
	str = strings.Replace(str, "\\n", "\n", -1)
	str = strings.Replace(str, "\\r", "\r", -1)
	str = strings.Replace(str, "\\t", "\t", -1)
	return str
}

// ParseQuad parses a Quad out of a string
func ParseQuad(s string) *Quad {
	match := regexQuad.FindStringSubmatch(s)
	if match == nil || len(match) != 25 {
		return nil
	}

	q := &Quad{
		parseQuadTerm(match[1:7]),
		parseQuadTerm(match[7:13]),
		parseQuadTerm(match[13:19]),
		parseQuadTerm(match[19:25]),
	}

	if q[3] == nil {
		q[3] = Default
	}
	return q
}

func parseQuadTerm(match []string) Term {
	if match[0] != "" {
		return NewNamedNode(match[0])
	} else if match[1] != "" {
		return NewBlankNode(match[1])
	} else if match[2] != "" {
		if match[3] != "" && match[3] != XSDString.value {
			return NewLiteral(match[2], "", NewNamedNode(match[3]))
		} else if match[4] != "" {
			return NewLiteral(match[2], match[4], RDFLangString)
		} else {
			return NewLiteral(match[2], "", nil)
		}
	} else if match[5] != "" {
		return NewVariable(match[5])
	}
	return nil
}

// ParseTerm parses a Term value out of an N-Quads string
func ParseTerm(t string) (Term, error) {
	if t == "" {
		return Default, nil
	} else if l := len(t); t[0] == '<' && t[l-1] == '>' {
		return NewNamedNode(t[1 : l-1]), nil
	} else if t[0:2] == "_:" {
		return NewBlankNode(t), nil
	} else if match := regexLiteral.FindStringSubmatch(t); match != nil {
		if match[2] != "" && match[2] != XSDString.value {
			return NewLiteral(unescape(match[1]), "", NewNamedNode(match[2])), nil
		} else if match[3] != "" {
			return NewLiteral(unescape(match[1]), match[3], RDFLangString), nil
		}
		return NewLiteral(unescape(match[1]), "", nil), nil
	} else if t[0] == '?' {
		return NewVariable(t), nil
	} else {
		return nil, ErrParseTerm
	}
}
