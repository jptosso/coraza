// Copyright 2022 Juan Pablo Tosso and the OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package coraza

import (
	"io/fs"

	"github.com/corazawaf/coraza/v3/internal/corazawaf"
	"github.com/corazawaf/coraza/v3/loggers"
	"github.com/corazawaf/coraza/v3/types"
)

// WAFConfig controls the behavior of the WAF.
//
// Note: WAFConfig is immutable. Each WithXXX function returns a new instance including the corresponding change.
type WAFConfig interface {
	// WithRules adds rules to the WAF.
	WithRules(rules ...*corazawaf.Rule) WAFConfig

	// WithDirectives parses the directives from the given string and adds them to the WAF.
	WithDirectives(directives string) WAFConfig

	// WithDirectivesFromFile parses the directives from the given file and adds them to the WAF.
	WithDirectivesFromFile(path string) WAFConfig

	// WithAuditLog configures audit logging.
	WithAuditLog(config AuditLogConfig) WAFConfig

	// WithContentInjection enables content injection.
	WithContentInjection() WAFConfig

	// WithRequestBodyAccess enables access to the request body.
	WithRequestBodyAccess() WAFConfig

	// WithRequestBodyLimit sets the maximum number of bytes that can be read from the request body. Bytes beyond that set
	// in WithInMemoryLimit will be buffered to disk.
	WithRequestBodyBytesLimit(limit int64) WAFConfig

	// WithRequestBodyInMemoryLimit sets the maximum number of bytes that can be read from the request body and buffered in memory.
	WithRequestBodyInMemoryBytesLimit(limit int64) WAFConfig

	// WithResponseBodyAccess enables access to the response body.
	WithResponseBodyAccess() WAFConfig

	// WithResponseBodyLimit sets the maximum number of bytes that can be read from the response body and buffered in memory.
	WithResponseBodyBytesLimit(limit int64) WAFConfig

	// WithResponseBodyMimeTypes sets the mime types of responses that will be processed.
	WithResponseBodyMimeTypes(mimeTypes []string) WAFConfig

	// WithDebugLogger configures a debug logger.
	WithDebugLogger(logger loggers.DebugLogger) WAFConfig

	// WithErrorLogger configures an error logger.
	WithErrorLogger(logger corazawaf.ErrorLogCallback) WAFConfig

	// WithRootFS configures the root file system.
	WithRootFS(fs fs.FS) WAFConfig
}

// NewWAFConfig creates a new WAFConfig with the default settings.
func NewWAFConfig() WAFConfig {
	return &wafConfig{
		requestBodyLimit:         -1,
		requestBodyInMemoryLimit: -1,
		responseBodyLimit:        -1,
	}
}

// AuditLogConfig controls audit logging.
type AuditLogConfig interface {
	// LogRelevantOnly enables audit logging only for relevant events.
	LogRelevantOnly() AuditLogConfig

	// WithParts configures the parts of the request/response to be logged.
	WithParts(parts types.AuditLogParts) AuditLogConfig

	// WithLogger configures the loggers.LogWriter to write logs to.
	WithLogger(logger loggers.LogWriter) AuditLogConfig
}

// NewAuditLogConfig returns a new AuditLogConfig with the default settings.
func NewAuditLogConfig() AuditLogConfig {
	return &auditLogConfig{}
}

type wafRule struct {
	rule *corazawaf.Rule
	str  string
	file string
}

type wafConfig struct {
	rules                    []wafRule
	auditLog                 *auditLogConfig
	contentInjection         bool
	requestBodyAccess        bool
	requestBodyLimit         int64
	requestBodyInMemoryLimit int64
	responseBodyAccess       bool
	responseBodyLimit        int64
	responseBodyMimeTypes    []string
	debugLogger              loggers.DebugLogger
	errorLogger              corazawaf.ErrorLogCallback
	fsRoot                   fs.FS
}

func (c *wafConfig) WithRules(rules ...*corazawaf.Rule) WAFConfig {
	if len(rules) == 0 {
		return c
	}

	ret := c.clone()
	for _, r := range rules {
		ret.rules = append(ret.rules, wafRule{rule: r})
	}
	return ret
}

func (c *wafConfig) WithDirectivesFromFile(path string) WAFConfig {
	ret := c.clone()
	ret.rules = append(ret.rules, wafRule{file: path})
	return ret
}

func (c *wafConfig) WithDirectives(directives string) WAFConfig {
	ret := c.clone()
	ret.rules = append(ret.rules, wafRule{str: directives})
	return ret
}

func (c *wafConfig) WithAuditLog(config AuditLogConfig) WAFConfig {
	ret := c.clone()
	ret.auditLog = config.(*auditLogConfig)
	return ret
}

func (c *wafConfig) WithContentInjection() WAFConfig {
	ret := c.clone()
	ret.contentInjection = true
	return ret
}

func (c *wafConfig) WithRequestBodyAccess() WAFConfig {
	ret := c.clone()
	ret.requestBodyAccess = true
	return ret
}

func (c *wafConfig) WithResponseBodyAccess() WAFConfig {
	ret := c.clone()
	ret.responseBodyAccess = true
	return ret
}

func (c *wafConfig) WithDebugLogger(logger loggers.DebugLogger) WAFConfig {
	ret := c.clone()
	ret.debugLogger = logger
	return ret
}

func (c *wafConfig) WithErrorLogger(logger corazawaf.ErrorLogCallback) WAFConfig {
	ret := c.clone()
	ret.errorLogger = logger
	return ret
}

func (c *wafConfig) WithRootFS(fs fs.FS) WAFConfig {
	ret := c.clone()
	ret.fsRoot = fs
	return ret
}

func (c *wafConfig) clone() *wafConfig {
	ret := *c // copy
	rules := make([]wafRule, len(c.rules))
	copy(rules, c.rules)
	ret.rules = rules
	return &ret
}

func (c *wafConfig) WithRequestBodyBytesLimit(limit int64) WAFConfig {
	ret := c.clone()
	ret.requestBodyLimit = limit
	return ret
}

func (c *wafConfig) WithRequestBodyInMemoryBytesLimit(limit int64) WAFConfig {
	ret := c.clone()
	ret.requestBodyInMemoryLimit = limit
	return ret
}

func (c *wafConfig) WithResponseBodyBytesLimit(limit int64) WAFConfig {
	ret := c.clone()
	ret.responseBodyLimit = limit
	return ret
}

func (c *wafConfig) WithResponseBodyMimeTypes(mimeTypes []string) WAFConfig {
	ret := c.clone()
	ret.responseBodyMimeTypes = mimeTypes
	return ret
}

type auditLogConfig struct {
	relevantOnly bool
	parts        types.AuditLogParts
	logger       loggers.LogWriter
}

func (c *auditLogConfig) LogRelevantOnly() AuditLogConfig {
	ret := c.clone()
	c.relevantOnly = true
	return ret
}

func (c *auditLogConfig) WithParts(parts types.AuditLogParts) AuditLogConfig {
	ret := c.clone()
	ret.parts = parts
	return ret
}

func (c *auditLogConfig) WithLogger(logger loggers.LogWriter) AuditLogConfig {
	ret := c.clone()
	ret.logger = logger
	return ret
}

func (c *auditLogConfig) clone() *auditLogConfig {
	ret := *c // copy
	return &ret
}
