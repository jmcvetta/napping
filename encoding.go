// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"encoding/json"
	"encoding/xml"
)

// these enum constants are used to declaratively assign session encoding
// default is JSON
const (
	JSON = iota
	XML
)

type EncodingType int

type EncodingMarshaller struct {
	Encoding EncodingType
}

type EncodingInterface interface {
	marshal(em *EncodingMarshaller, v interface{}) ([]byte, error)
	unmarshal(em *EncodingMarshaller, data []byte, v interface{}) error
}

// Marshal parses the JSON/XML-encoded data and returns its
// raw bytes representation.
func (em *EncodingMarshaller) Marshal(v interface{}) ([]byte, error) {
	switch em.Encoding {
	case JSON:
		return json.Marshal(v)
	case XML:
		return xml.Marshal(v)
	}
	panic("invalid encoding")
}

// Unmarshal parses the JSON/XML-encoded data and stores
// the result in the value pointed to by v.
func (em *EncodingMarshaller) Unmarshal(data []byte, v interface{}) error {
	switch em.Encoding {
	case JSON:
		return json.Unmarshal(data, v)
	case XML:
		return xml.Unmarshal(data, v)
	}
	panic("invalid encoding")
}
