package errors

import (
	"net/http"
)

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
	UserInputErrorCode                    = 40000
	InvalidFormDataCode                   = NewErrorCode("InvalidFormDataError", UserInputErrorCode+HTTPInvalidData)
	NotFoundErrorCode                     = NewErrorCode("NotFoundError", UserInputErrorCode+HTTPNotFound)
	TooManyAttemptsErrorCode              = NewErrorCode("TooManyAttemptsError", UserInputErrorCode+HTTPThrottling)
	UnauthorizedErrorCode                 = NewErrorCode("UnauthorizedError", UserInputErrorCode+HTTPNotAuthenticated)
	ForbiddenErrorCode                    = NewErrorCode("ForbiddenError", UserInputErrorCode+HTTPEndpointForbidden)
	InvalidJWTErrorCode                   = NewErrorCode("InvalidJWTError", UserInputErrorCode+HTTPNotAuthenticated)
	InvalidAuthTokenErrorCode             = NewErrorCode("InvalidAuthTokenError", UserInputErrorCode+HTTPSecurityKeyMissing)
	ErrorPerformingRequestErrorCode       = NewErrorCode("ErrorPerformingRequestError", UserInputErrorCode+HTTPBadRequest)
	ErrorUnmarshallBodyErrorCode          = NewErrorCode("ErrorUnmarshallBodyError", UserInputErrorCode+HTTPInvalidData)
	EntityNotFoundErrorCode               = NewErrorCode("EntityNotFoundErrorCode", UserInputErrorCode+HTTPNotFound)
	UserNotFoundErrorCode                 = NewErrorCode("UserNotFoundErrorCode", UserInputErrorCode+HTTPNotFound)
	UserNotActiveErrorCode                = NewErrorCode("UserNotActiveErrorCode", UserInputErrorCode+HTTPInvalidData)
	SessionChannelNotSupportedErrorCode   = NewErrorCode("SessionChannelNotSupportedErrorCode", UserInputErrorCode+HTTPInvalidData)
	APIValidationErrorCode                = NewErrorCode("APIValidationErrorCode", UserInputErrorCode+HTTPInvalidData)
	EntitiesInactiveUnauthorizedErrorCode = NewErrorCode("EntitiesInactiveUnauthorizedErrorCode", UserInputErrorCode+HTTPNotAuthenticated)
)

var (
	StreamsErrorCode                   = 55000
	ServerErrorErrorCode               = NewErrorCode("ServerErrorErrorCode", StreamsErrorCode+http.StatusInternalServerError)
	ConnectionNotActive                = NewErrorCode("ConnectionNotActive", StreamsErrorCode+http.StatusGone)
	ProducerErrorCode                  = NewErrorCode("ProducerErrorCode", StreamsErrorCode+http.StatusServiceUnavailable)
	RateLimitErrorCode                 = NewErrorCode("RateLimitErrorCode", StreamsErrorCode+http.StatusForbidden)
	ProcessFailedDoNotRequeueErrorCode = NewErrorCode("ProcessFailedDoNotRequeueErrorCode", StreamsErrorCode+HTTPServerError)
	InvalidScopeRequeueErrorCode       = NewErrorCode("InvalidScopeRequeueErrorCode", StreamsErrorCode+HTTPServerError)
	InvalidRecordsListErrorCode        = NewErrorCode("InvalidRecordsListErrorCode", StreamsErrorCode+HTTPServerError)
)

var (
	SystemErrorCode                      = 50000
	JoinedErrorCode                      = NewErrorCode("JoinedError", SystemErrorCode+http.StatusMultipleChoices)
	FailedToWriteDataErrorCode           = NewErrorCode("FailedToWriteDataError", SystemErrorCode+HTTPServerError)
	FailedToReadDataErrorCode            = NewErrorCode("FailedToReadDataError", SystemErrorCode+HTTPServerError)
	DBErrorCode                          = NewErrorCode("DBError", SystemErrorCode+HTTPServerError)
	UnknownErrorCode                     = NewErrorCode("UnknownError", SystemErrorCode+HTTPServerError)
	InvalidProcessHandlerErrorCode       = NewErrorCode("InvalidProcessHandlerError", SystemErrorCode+HTTPServerError)
	InvalidCtxMetricErrorCode            = NewErrorCode("InvalidCtxMetricError", SystemErrorCode+HTTPServerError)
	ErrorCreatingMetricErrorCode         = NewErrorCode("ErrorCreatingMetricError", SystemErrorCode+HTTPServerError)
	EventSourceMappingDontExistErrorCode = NewErrorCode("EventSourceMappingDontExistError", SystemErrorCode+HTTPServerError)
	LambdaInitFailedErrorCode            = NewErrorCode("LambdaInitFailedErrorCode", SystemErrorCode+HTTPServerError)
	LambdaPanicErrorCode                 = NewErrorCode("LambdaPanicErrorCode", SystemErrorCode+HTTPServerError)
	FailedToAcquireLockErrorCode         = NewErrorCode("FailedToAcquireLockErrorCode", SystemErrorCode+HTTPServerError)
	NoRetryErrorCode                     = NewErrorCode("NoRetryErrorCode", SystemErrorCode+HTTPServerError)
	InvalidTypeErrorCode                 = NewErrorCode("InvalidTypeErrorCode", SystemErrorCode+HTTPServerError)
)

// Generic Errors
var (
	SystemNoCodeCodeBase              = 90000
	GenericErrorCode                  = NewErrorCode("GenericErrorCode", SystemNoCodeCodeBase+HTTPServerError)
	ErrorCreatingDependencyErrorCode  = NewErrorCode("ErrorCreatingDependencyError", SystemNoCodeCodeBase+HTTPServerError)
	ErrorLoadingStructConsulErrorCode = NewErrorCode("ErrorLoadingStructConsulError", SystemNoCodeCodeBase+HTTPServerError)
)

// State Machine errors
var (
	StateMachineErrorCode                  = 60000
	StateMachineInvalidTransitionErrorCode = NewErrorCode("StateMachineInvalidTransitionErrorCode", StateMachineErrorCode+HTTPEndpointForbidden)
	StateMachineInvalidStateErrorCode      = NewErrorCode("StateMachineInvalidStateErrorCode", StateMachineErrorCode+HTTPInvalidData)
	StateMachineStateNotVisitedErrorCode   = NewErrorCode("StateMachineStateNotVisitedErrorCode", StateMachineErrorCode+HTTPEndpointForbidden)
)
