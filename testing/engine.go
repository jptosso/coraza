// Copyright 2021 Juan Pablo Tosso
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	b64 "encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	engine "github.com/jptosso/coraza-waf/v2"
	"github.com/jptosso/coraza-waf/v2/types/variables"
)

type test struct {
	waf         *engine.Waf
	transaction *engine.Transaction
	magic       bool
	name        string
	body        string

	// public variables
	RequestAddress   string
	RequestPort      int
	RequestUri       string
	RequestMethod    string
	RequestProtocol  string
	RequestHeaders   map[string]string
	ResponseHeaders  map[string]string
	ResponseCode     int
	ResponseProtocol string
	ServerAddress    string
	ServerPort       int
	ExpectedOutput   expectedOutput
}

func (t *test) SetWaf(waf *engine.Waf) {
	t.waf = waf
}

func (t *test) DisableMagic() {
	t.magic = false
}

func (t *test) SetEncodedRequest(request string) error {
	sDec, err := b64.StdEncoding.DecodeString(request)
	if err != nil {
		return err
	}
	return t.SetRawRequest(sDec)
}

func (t *test) SetRawRequest(request []byte) error {
	return nil
}

func (t *test) SetRequestBody(body interface{}) error {
	if body == nil {
		return nil
	}
	data := ""
	v := reflect.ValueOf(body)
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			data += fmt.Sprintf("%s\r\n", v.Index(i))
		}
		data += "\r\n"
	case reflect.String:
		data = body.(string)
	}
	lbody := len(data)
	if lbody == 0 {
		return nil
	}
	t.body = data
	if t.magic {
		t.RequestHeaders["content-length"] = strconv.Itoa(lbody)
	}
	if _, err := t.transaction.RequestBodyBuffer.Write([]byte(data)); err != nil {
		return err
	}
	return nil
}

func (t *test) RunPhases() error {
	t.transaction.ProcessConnection(t.RequestAddress, t.RequestPort, t.ServerAddress, t.ServerPort)
	t.transaction.ProcessUri(t.RequestUri, t.RequestMethod, t.RequestProtocol)
	for k, v := range t.RequestHeaders {
		t.transaction.AddRequestHeader(k, v)
	}
	t.transaction.ProcessRequestHeaders()
	if _, err := t.transaction.ProcessRequestBody(); err != nil {
		return err
	}
	for k, v := range t.ResponseHeaders {
		t.transaction.AddResponseHeader(k, v)
	}
	t.transaction.ProcessResponseHeaders(t.ResponseCode, t.ResponseProtocol)
	if _, err := t.transaction.ProcessResponseBody(); err != nil {
		return err
	}
	t.transaction.ProcessLogging()
	return nil
}

func (t *test) OutputErrors() []string {
	var errors []string
	if lc := t.ExpectedOutput.LogContains; lc != "" {
		if !t.LogContains(lc) {
			errors = append(errors, fmt.Sprintf("Expected log to contain '%s'", lc))
		}
	}
	if lc := t.ExpectedOutput.NoLogContains; lc != "" {
		if t.LogContains(lc) {
			errors = append(errors, fmt.Sprintf("Expected log to not contain '%s'", lc))
		}
	}
	if rc := t.ExpectedOutput.Status; rc != 0 {
		// do nothing
	}
	if tr := t.ExpectedOutput.TriggeredRules; tr != nil {
		for _, rule := range tr {
			if !t.LogContains(fmt.Sprintf("id \"%d\"", rule)) {
				errors = append(errors, fmt.Sprintf("Expected rule '%d' to be triggered", rule))
			}
		}
	}
	if tr := t.ExpectedOutput.NonTriggeredRules; tr != nil {
		for _, rule := range tr {
			if t.LogContains(fmt.Sprintf("id \"%d\"", rule)) {
				errors = append(errors, fmt.Sprintf("Expected rule '%d' to not be triggered", rule))
			}
		}
	}

	return errors
}

func (t *test) LogContains(log string) bool {
	for _, mr := range t.transaction.MatchedRules {
		if strings.Contains(mr.ErrorLog(t.ResponseCode), log) {
			return true
		}
	}
	return false
}

func (t *test) Transaction() *engine.Transaction {
	return t.transaction
}

func (test *test) String() string {
	tx := test.transaction
	res := "======DEBUG======\n"
	for v := byte(1); v < 100; v++ {
		vr := variables.RuleVariable(v)
		if vr.Name() == "UNKNOWN" {
			break
		}
		res += fmt.Sprintf("%s:\n", vr.Name())
		data := tx.GetCollection(vr).Data()
		for k, d := range data {
			if k != "" {
				res += fmt.Sprintf("-->%s: %s\n", k, strings.Join(d, ","))
			} else {
				res += fmt.Sprintf("-->%s\n", strings.Join(d, ","))
			}
		}
	}
	return res
}

func (test *test) Request() string {
	str := fmt.Sprintf("%s %s %s\r\n", test.RequestMethod, test.RequestUri, test.RequestProtocol)
	for k, v := range test.RequestHeaders {
		str += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	str += "\r\n"
	if test.body != "" {
		str += test.body
	}
	return str
}

func newTest(name string, waf *engine.Waf) *test {
	t := &test{
		name:            name,
		waf:             waf,
		transaction:     waf.NewTransaction(),
		RequestHeaders:  map[string]string{},
		ResponseHeaders: map[string]string{},
		RequestMethod:   "GET",
		RequestProtocol: "HTTP/1.1",
		RequestUri:      "/",
		RequestAddress:  "127.0.0.1",
		RequestPort:     80,
		magic:           true,
	}
	return t
}
