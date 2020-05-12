package rdf

import (
	"bufio"
	"io"
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

var regexNode = regexp.MustCompile("^" + node + "$")
var regexQuad = regexp.MustCompile("^" + wso + node + ws + node + ws + node + ws + graph + wso + "\n?$")

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
		parseNode(match[1:7]),
		parseNode(match[7:13]),
		parseNode(match[13:19]),
		parseNode(match[19:25]),
	}

	if q[3] == nil {
		q[3] = Default
	}
	return q
}

func parseNode(match []string) Term {
	if match[0] != "" {
		return NewNamedNode(match[0])
	} else if match[1] != "" {
		return NewBlankNode(match[1][2:])
	} else if match[2] != "" {
		value := unescape(match[2])
		if match[3] != "" && match[3] != XSDString.value {
			return NewLiteral(value, "", NewNamedNode(match[3]))
		} else if match[4] != "" {
			return NewLiteral(value, match[4], RDFLangString)
		} else {
			return NewLiteral(value, "", nil)
		}
	} else if match[5] != "" {
		return NewVariable(match[5][1:])
	}
	return nil
}

// ParseTerm parses a Term value out of an N-Quads string
func ParseTerm(t string) (Term, error) {
	if t == "" {
		return Default, nil
	}

	match := regexNode.FindStringSubmatch(t)
	if match == nil || len(match) != 7 {
		return nil, ErrParseTerm
	}

	return parseNode(match[1:]), nil
}

// ReadQuads parses an io.Reader of serialize n-quads into a slice of *Quads
func ReadQuads(input io.Reader) ([]*Quad, error) {
	quads := []*Quad{}
	reader := bufio.NewReader(input)
	line, err := reader.ReadString('\n')
	for ; err == nil; line, err = reader.ReadString('\n') {
		if line != "" {
			quad := ParseQuad(line)
			if quad != nil {
				quads = append(quads, quad)
			}
		}
	}

	if err != io.EOF {
		return nil, err
	}

	return quads, nil
}
