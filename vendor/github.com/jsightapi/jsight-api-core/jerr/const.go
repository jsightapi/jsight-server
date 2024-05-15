package jerr

const (
	RuntimeFailure = "runtime failure" // should not occur
	TestFakeError  = "fake error"

	HTTPResourceNotFound    = "resource not found"
	HTTPMethodNotFound      = "HTTP method not found"
	PathNotFound            = "path not found"
	TagNotFound             = "tag not found"
	UserTypeNotFound        = "user type not found"
	ParentNotFound          = "parent directive not found"
	MacroNotFound           = "macro not found"
	ServerNotFound          = "server not found"
	JsonRpcMethodNotFound   = "JSON-RPC method not found"
	JsonRpcResourceNotFound = "resource not found"

	ProtocolNotFound     = `the directive "Protocol" not found`
	ProtocolParameterErr = `the parameter value have to be "json-rpc-2.0"`

	InfoIsEmpty        = "the INFO directive cannot be empty"
	ResponsesIsEmpty   = "the response cannot be empty"
	RequestIsEmpty     = "the request cannot be empty"
	DescriptionIsEmpty = "the description cannot be empty, learn more about the Description directive here: https://jsight.io/docs/jsight-api-0-3#directive-description" //nolint:lll
	BodyIsEmpty        = "the body cannot be empty"
	MacroIsEmpty       = "the macros cannot be empty, learn more about the MACRO directive here: https://jsight.io/docs/jsight-api-0-3#directive-macro" //nolint:lll

	IncorrectPath             = "incorrect path"
	IncorrectRequest          = "incorrect request"
	IncorrectDirectiveContext = "incorrect context for the directive"
	IncorrectParameter        = "incorrect parameter"

	PathOrErr                          = "the root schema object cannot have the `or` rule in the Path directive"
	PathObjectErr                      = "the body of the Path directive must be an object, learn more about the Path directive here: https://jsight.io/docs/jsight-api-0-3#directive-path"                         //nolint:lll
	PathAdditionalPropertiesErr        = `the "additionalProperties" rule should not be used in the Path directive, learn more about the Path directive here: https://jsight.io/docs/jsight-api-0-3#directive-path` //nolint:lll
	PathNullableErr                    = `the "nullable" rule should not be used in the Path directive, learn more about the Path directive here: https://jsight.io/docs/jsight-api-0-3#directive-path`             //nolint:lll
	PathEmptyErr                       = "the object in the Path directive can not be empty, learn more about the Path directive here: https://jsight.io/docs/jsight-api-0-3#directive-path"                        //nolint:lll
	PathMultiLevelPropertyErr          = "the multi-level property is not allowed in the Path directive"
	PathEmptyParameter                 = "empty PATH parameter"
	PathParameterIsDuplicatedInThePath = "the parameter of the path is duplicated"
	PathsAreSimilar                    = "the ambiguous paths are not allowed: \"/%s\", \"/%s\", see the details here: https://jsight.io/docs/jsight-api-0-3#parameter-path"                    //nolint:lll
	PathParameterAlreadyDefined        = "The parameter %q has already been defined earlier, see more details about path parameters here: https://jsight.io/docs/jsight-api-0-3#parameter-path" //nolint:lll

	IncludeRootErr      = "cannot not start with `/`"
	IncludeUpErr        = "cannot contain `..` or `.`"
	IncludeSeparatorErr = "directories must be separated by slashes `/`"
	IncludeDirectiveErr = "the directive is not allowed in included files:"

	UnsupportedVersion                = "The specified JSight version is not supported"
	DirectiveJSIGHTShouldBeTheFirst   = "The first directive in the document must be JSIGHT"
	DirectiveJSIGHTGottaBeOnlyOneTime = "The directive JSIGHT has already been specified before"
	DirectiveINFOGottaBeOnlyOneTime   = "The directive INFO has already been specified before"
	DirectiveBaseURLAlreadyDefined    = "The directive BaseUrl has already been defined before"

	UnknownDirective = "unknown directive"
	UnknownNotation  = "unknown notation"

	RequiredParameterNotSpecified         = "required parameter(s) not specified"
	ParametersAreForbiddenForTheDirective = "the directive should not have parameters in this case"
	ParametersIsAlreadyDefined            = "the parameter %q is already defined for the directive"

	AnnotationIsForbiddenForTheDirective                = "the annotation is not allowed for this directive"
	NotUniqueDirective                                  = "the directive has already been defined"
	NotUniquePath                                       = "the path %q has already been defined"
	BodyMustBeObject                                    = "there must be an object or a reference to an object in the directive body"                                                                                                              //nolint:lll
	CannotUseTheTypeAndSchemaNotationParametersTogether = "directive parameters `Type` and `SchemaNotation` cannot be declared simultaneously"                                                                                                     //nolint:lll
	ThereIsNoExplicitContextForClosure                  = "nothing to close with this closing parenthesis, learn more about the explicit direcitve boundaries here: https://jsight.io/docs/jsight-api-0-3#boundaries-of-the-body-of-the-directive" //nolint:lll
	DirectiveNotAllowed                                 = "the directive is not allowed"
	ApartFromTheOpeningParenthesis                      = "apart from the opening parenthesis, there should be nothing else on this line, learn more about the explicit direcitve boundaries here: https://jsight.io/docs/jsight-api-0-3#boundaries-of-the-body-of-the-directive" //nolint:lll
	DuplicateNames                                      = "the name %q has already been declared before"
	NotAllowedToOverrideTheProperty                     = "it is not allowed to override the %q property from the user type %q"                                                                                                            //nolint:lll
	ContextNotClosed                                    = "this opening parenthesis is not closed, learn more about the explicit direcitve boundaries here: https://jsight.io/docs/jsight-api-0-3#boundaries-of-the-body-of-the-directive" //nolint:lll
	WrongDescriptionContext                             = "wrong description context"
	MethodIsAlreadyDefinedInResource                    = "this method has already been defined in the resource"
	UndefinedRequestBodyForResource                     = "undefined request body for resource"
	RecursionIsProhibited                               = "file dependency recursion is detected, learn more about the INCLUDE directive here: https://jsight.io/docs/jsight-api-0-3#directive-include" //nolint:lll
	UserTypeIsNotAnObject                               = "the user type is not an object"
	ProcessTypeErr                                      = "process type"
	FailedToComputeScannersHash                         = "failed to compute the scanner's hash"
)
