// Copyright 2022 Juan Pablo Tosso and the OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package corazawaf

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	waf := NewWAF()
	waf.RequestBodyAccess = true
	waf.ResponseBodyAccess = true
	waf.RequestBodyLimit = 1044

	tx := waf.NewTransactionWithOptions(Options{ID: "test"})
	if !tx.RequestBodyAccess {
		t.Error("Request body access not enabled")
	}
	if !tx.ResponseBodyAccess {
		t.Error("Response body access not enabled")
	}
	if tx.RequestBodyLimit != 1044 {
		t.Error("Request body limit not set")
	}
	if tx.id != "test" {
		t.Error("ID not set")
	}
	tx = waf.NewTransactionWithOptions(Options{ID: ""})
	if tx.id == "" {
		t.Error("ID not set")
	}
	tx = waf.NewTransaction()
	if tx.id == "" {
		t.Error("ID not set")
	}
}

func TestSetDebugLogPath(t *testing.T) {
	tests := map[string]struct {
		path  string
		match func(io.Writer, bool) bool
	}{
		"empty path": {path: "", match: func(w io.Writer, b bool) bool { return w == io.Discard && b == false }},
		"stdout":     {path: "/dev/stdout", match: func(w io.Writer, b bool) bool { return w == os.Stdout && b == false }},
		"stderr":     {path: "/dev/stderr", match: func(w io.Writer, b bool) bool { return w == os.Stderr && b == false }},
		"file":       {path: filepath.Join(t.TempDir(), "debuglogtest.log"), match: func(w io.Writer, b bool) bool { fl, ok := w.(io.Closer); fl.Close(); return ok && b == true }},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w, toClose, err := resolveLogPath(test.path)
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if !test.match(w, toClose) {
				t.Errorf("unexpected return from check")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := map[string]struct {
		customizer func(*WAF)
		expectErr  bool
	}{
		"default": {
			expectErr:  false,
			customizer: func(w *WAF) {},
		},
		"request body limit less than zero": {
			expectErr:  true,
			customizer: func(w *WAF) { w.RequestBodyLimit = -1 },
		},
		"request body limit greater than 1gb": {
			expectErr:  true,
			customizer: func(w *WAF) { w.RequestBodyLimit = _1gb + 1 },
		},
		"request body in memory limit less than zero": {
			expectErr:  true,
			customizer: func(w *WAF) { w.SetRequestBodyInMemoryLimit(-1) },
		},
		"request body limit less than request body in memory limit": {
			expectErr: true,
			customizer: func(w *WAF) {
				w.RequestBodyLimit = 10
				w.SetRequestBodyInMemoryLimit(11)
			}},
		"response body limit less than zero": {
			expectErr:  true,
			customizer: func(w *WAF) { w.ResponseBodyLimit = -1 },
		},
		"response body limit greater than 1gb": {
			expectErr:  true,
			customizer: func(w *WAF) { w.ResponseBodyLimit = _1gb + 1 },
		},
		"argument limit greater than 0": {
			expectErr:  false,
			customizer: func(w *WAF) { w.ArgumentLimit = 1000 },
		},
		"argument limit less than 0": {
			expectErr:  true,
			customizer: func(w *WAF) { w.ArgumentLimit = -1 },
		},
	}

	for name, tCase := range testCases {
		t.Run(name, func(t *testing.T) {
			waf := NewWAF()
			tCase.customizer(waf)
			err := waf.Validate()
			if tCase.expectErr {
				if err == nil {
					t.Fatalf("expected error: %s", err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}
		})
	}
}
