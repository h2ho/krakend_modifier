package querystring

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/martian"
	"github.com/google/martian/parse"
)

func init() {
	parse.Register("Body.Retrive", bodyModifierFromJSON)
}

// BodyModifier contains the private and public Marvel API key
type BodyModifier struct {
	source       string
	target, keys []string
}

// BodyModifierJSON to Unmarshal the JSON configuration
type BodyModifierJSON struct {
	Source string               `json:"source"`
	Target []string             `json:"target"`
	Keys   []string             `json:"keys"`
	Scope  []parse.ModifierType `json:"scope"`
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *BodyModifier) ModifyRequest(req *http.Request) error {

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	data := url.Values{}
	query := req.URL.Query()
	if m.source == "header" {
		for i := 0; i < len(m.target); i++ {
			data.Set(m.keys[i], req.Header.Get(m.target[i]))
		}
		req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))
	}

	if m.source == "query" {
		for i := 0; i < len(m.target); i++ {
			data.Set(m.keys[i], query.Get(m.target[i]))
		}
		req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))
	}

	// req.Header.Set("Content-Type", "plain/text")
	// req.Body = ioutil.NopCloser(strings.NewReader(strings.Join(m.target, " ")))

	return nil
}

// BodyNewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func BodyNewModifier(source string, target []string, keys []string) martian.RequestModifier {
	return &BodyModifier{
		source: source,
		target: target,
		keys:   keys,
	}
}

// marvelModifierFromJSON takes a JSON message as a byte slice and returns
// a querystring.modifier and an error.
//
// Example JSON:
// {
//  "public": "apikey",
//  "private": "apikey",
//  "scope": ["request", "response"]
// }
func bodyModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &BodyModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}
	return parse.NewResult(BodyNewModifier(msg.Source, msg.Target, msg.Keys), msg.Scope)
}
