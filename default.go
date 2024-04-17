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

// Errors 40000
var (
	SystemClientErrorCodeBase       = 40000
	InvalidFormDataCode             = NewErrorCode("InvalidFormDataError", SystemClientErrorCodeBase+HTTPInvalidData, HTTPInvalidData)
	NotFoundErrorCode               = NewErrorCode("NotFoundError", SystemClientErrorCodeBase+HTTPNotFound, HTTPNotFound)
	TooManyAttemptsErrorCode        = NewErrorCode("TooManyAttemptsError", SystemClientErrorCodeBase+HTTPThrottling, HTTPThrottling)
	UnauthorizedErrorCode           = NewErrorCode("UnauthorizedError", SystemClientErrorCodeBase+HTTPNotAuthenticated, HTTPNotAuthenticated)
	ForbiddenErrorCode              = NewErrorCode("ForbiddenError", SystemClientErrorCodeBase+HTTPEndpointForbidden, HTTPEndpointForbidden)
	InvalidJWTErrorCode             = NewErrorCode("InvalidJWTError", SystemClientErrorCodeBase+HTTPNotAuthenticated, HTTPNotAuthenticated)
	InvalidAuthTokenErrorCode       = NewErrorCode("InvalidAuthTokenError", SystemClientErrorCodeBase+HTTPSecurityKeyMissing, HTTPSecurityKeyMissing)
	ErrorPerformingRequestErrorCode = NewErrorCode("ErrorPerformingRequestError", SystemClientErrorCodeBase+HTTPBadRequest, HTTPBadRequest)
	ErrorUnmarshallBodyErrorCode    = NewErrorCode("ErrorUnmarshallBodyError", SystemClientErrorCodeBase+HTTPInvalidData, HTTPInvalidData)
)

// Errors 50000
var (
	SystemSystemErrorCodeBase            = 50000
	FailedToWriteDataErrorCode           = NewErrorCode("FailedToWriteDataError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	FailedToReadDataErrorCode            = NewErrorCode("FailedToReadDataError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	DBErrorCode                          = NewErrorCode("DBError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	UnknownErrorCode                     = NewErrorCode("UnknownError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	InvalidProcessHandlerErrorCode       = NewErrorCode("InvalidProcessHandlerError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	InvalidCtxMetricErrorCode            = NewErrorCode("InvalidCtxMetricError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	ErrorCreatingMetricErrorCode         = NewErrorCode("ErrorCreatingMetricError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
	EventSourceMappingDontExistErrorCode = NewErrorCode("EventSourceMappingDontExistError", SystemSystemErrorCodeBase+HTTPServerError, HTTPServerError)
)

// Generic Errors
var (
	SystemNoCodeCodeBase              = 90000
	GenericErrorCode                  = NewErrorCode("GenericErrorCode", SystemNoCodeCodeBase+HTTPServerError, HTTPServerError)
	ErrorCreatingDependencyErrorCode  = NewErrorCode("ErrorCreatingDependencyError", SystemNoCodeCodeBase+HTTPServerError, HTTPServerError)
	ErrorLoadingStructConsulErrorCode = NewErrorCode("ErrorLoadingStructConsulError", SystemNoCodeCodeBase+HTTPServerError, HTTPServerError)
)
