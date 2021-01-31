// Package log defines common structured-logging fields and functions.
package log

// Logging categories.
const (
	CategoryHTTPAccess  = "access"
	CategoryApplication = "application"
)

// Better zero values.
const (
	FieldValueUnknown = "unknown"
)

// Common structured-logging fields:
//
//  Field name     Description
//  ==========     ===========
//  msg            Left empty (reserved, automatically set by the Logger)
//  time           Time request was processed (reserved, auto set by the Logger)
//  level          Logging level that is always set to INFO (reserved, auto set by the Logger)
//
//  category       Category of log message (e.g., access, application)
//  service        Name of the service responsible for the log entry
//  localaddr      Network address of a local service
//  localport      Network port of a local service
//  duration       Duration in milliseconds
//  reqlen         Request body length in bytes (e.g., http.Request.Headers[Content-Length])
//  method         Request method (e.g., http.Request.Method)
//  uri            Request URI (e.g., http.Request.RequestURI)
//  proto          Request protocol (e.g., http.Request.Proto)
//  remoteaddr     Network address of a remote address (e.g., http.Request.RemoteAddr)
//  remoteport     Network port of a remote service
//  status         Response status code (loggedResponseWriter.code)
//  authuser       Authenticated username
//  authgroups     Authenticated groups
//  authroles      Authenticated roles
//
//
// Notes:
// 1) authuser,authgroups and authroles fields require handler.Authn.
const (
	FieldCategory           = "category"
	FieldApplication        = "application"
	FieldService            = "service"
	FieldVersion            = "version"
	FieldLocalAddress       = "localaddr"
	FieldLocalPort          = "localport"
	FieldDuration           = "duration"
	FieldRemoteAddress      = "remoteaddr"
	FieldRemotePort         = "remoteport"
	FieldRequestBodyLen     = "reqlen"
	FieldRequestMethod      = "method"
	FieldRequestURI         = "uri"
	FieldRequestProtocol    = "proto"
	FieldResponseStatusCode = "status"
	FieldAuthnUsername      = "authuser"
	FieldAuthnGroups        = "authgroups"
	FieldAuthnRoles         = "authnroles"
	FieldUserAgent          = "useragent"
	FieldFunc               = "func"
	FieldEntityType         = "entitytype"
	FieldTraceID            = "traceid"
	FieldSpanID             = "spanid"
	FieldParentSpanID       = "parentspanid"
	FieldRequestUUID        = "requestuuid"
)

// Deprecated fields.
const (
	FieldHost = "host"
)
