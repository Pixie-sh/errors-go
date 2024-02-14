package errors

// http return codes
const (
	HTTPSuccess                = 200 // Standard OK response.
	HTTPBadRequest             = 400 // Error executing request
	HTTPNotAuthenticated       = 401 // Not authenticated.
	HTTPLostAccess             = 402 // Lost access to the service
	HTTPEndpointForbidden      = 403 // No access to this endpoint.
	HTTPNotFound               = 404 // Entity not found or Endpoint not found. A message is sent to distinct these cases.
	HTTPSecurityKeyMissing     = 407 // Missing microservice security key from the headers.
	HTTPInvalidData            = 422 // Invalid data
	HTTPThrottling             = 429 // Too Many Attempts.
	HTTPServerError            = 500 // Unexpected error. Check ELK logs for the stack error.
	HTTPIncompleteRegistration = 501 // Uncompleted registration process.
)

// Form Errors
var (
	NotFoundErrorCode = NewErrorCode("GENERIC", 100, HTTPNotFound)
	InvalidFormData   = NewErrorCode("GENERIC", 422, HTTPInvalidData)
)

// Access Errors
var (
	TooManyAttemptsErrorCode = NewErrorCode("GENERIC", 429, HTTPThrottling)
	UnAuthorized             = NewErrorCode("GENERIC", 401, HTTPNotAuthenticated)
	Forbidden                = NewErrorCode("GENERIC", 403, HTTPEndpointForbidden)
)

// Internal Errors
var (
	FailedToWriteDataErrorCode = NewErrorCode("GENERIC", 997, HTTPServerError)
	FailedToReadDataErrorCode  = NewErrorCode("GENERIC", 998, HTTPServerError)
	DBErrorCode                = NewErrorCode("GENERIC", 999, HTTPServerError)
)

// Generic Errors
var (
	NoErrorCode                       = NewErrorCode("ERROR", 0, HTTPServerError)
	ErrorCreatingDependencyErrorCode  = NewErrorCode("ERROR", 100, HTTPServerError)
	ErrorLoadingStructConsulErrorCode = NewErrorCode("ERROR", 101, HTTPServerError)
)

// HTTP default error codes
var (
	UnknownErrorCode                     = NewErrorCode("HTTP_SERVER", 100, HTTPServerError)
	InvalidProcessHandlerErrorCode       = NewErrorCode("HTTP_SERVER", 101, HTTPServerError)
	InvalidCtxMetricErrorCode            = NewErrorCode("HTTP_SERVER", 103, HTTPServerError)
	ErrorCreatingMetricErrorCode         = NewErrorCode("HTTP_SERVER", 104, HTTPServerError)
	EventSourceMappingDontExistErrorCode = NewErrorCode("HTTP_SERVER", 105, HTTPServerError)
)

// Gates default error codes
var (
	InvalidJWTErrorCode       = NewErrorCode("GATES", 100, HTTPNotAuthenticated)
	InvalidAuthTokenErrorCode = NewErrorCode("GATES", 101, HTTPSecurityKeyMissing)
)

// Rest default errors
var (
	ErrorPerformingRequestErrorCode = NewErrorCode("REST", 101, HTTPBadRequest)
	ErrorUnmarshallBodyErrorCode    = NewErrorCode("REST", 102, HTTPInvalidData)
)
