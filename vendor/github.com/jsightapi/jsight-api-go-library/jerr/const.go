package jerr

const (
	InternalServerError                                 = "internal server error" // should not occur
	RequiredParameterNotSpecified                       = "required parameter not specified"
	ParametersAreForbiddenForTheDirective               = "parameters are forbidden for the directive"
	AnnotationIsForbiddenForTheDirective                = "annotation is forbidden for the directive"
	EmptyDescription                                    = "empty description"
	EmptyBody                                           = "empty body"
	HttpResourceNotFound                                = "resource not found"
	ResponsesIsEmpty                                    = "responses is empty"
	RequestIsEmpty                                      = "request is empty"
	NotUniqueDirective                                  = "not a unique directive"
	IncorrectParameter                                  = "incorrect parameter"
	BodyMustBeObject                                    = "body must be object"
	IsNotHTTPRequestMethod                              = "directive is not a HTTP request method"
	HttpMethodNotFound                                  = "HTTP method not found"
	PathNotFound                                        = "path not found"
	IncorrectPath                                       = "incorrect path"
	CannotUseTheTypeAndSchemaNotationParametersTogether = "cannot use the Type and SchemaNotation parameters together"
	IncorrectContextOfDirective                         = "incorrect context of directive"
	ThereIsNoExplicitContextForClosure                  = "there is no explicit context for closure"
	DirectiveNotAllowed                                 = "directive not allowed"
	JsonRpcMethodNotFound                               = "JSON-RPC method not found"
	JsonRpcResourceNotFound                             = "resource not found"
	ApartFromTheOpeningParenthesis                      = "apart from the opening parenthesis, there should be nothing else on this line"
)
