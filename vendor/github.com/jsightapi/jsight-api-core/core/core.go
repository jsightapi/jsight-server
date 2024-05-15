package core

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/fs"

	"github.com/jsightapi/jsight-api-core/catalog"
	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/scanner"
)

// JApiCore the Brain and heart of jApi. Collects lexemes from scanner, validates document logic and structure,
// builds catalog, renders documentation
type JApiCore struct {
	rawUserTypes *directive.Directives

	// processedUserTypes a "set" of already build user types.
	processedUserTypes map[string]struct{}

	// userTypes represent all user types.
	userTypes *catalog.UserSchemas

	// onlyOneProtocolIntoURL to verify the uniqueness of a directive Protocol within a directive URL
	onlyOneProtocolIntoURL map[*directive.Directive]struct{}

	// uniqURLPath used for checking the uniqueness URL paths.
	uniqURLPath map[catalog.Path]struct{}

	// similarPaths to check the forbidding of "similar" paths.
	similarPaths map[string]string

	// macro contains list of all project macros.
	macro map[string]*directive.Directive

	// directiveFunctions map between available directives and function which
	// should be used for processing.
	directiveFunctions map[directive.Enumeration]func(*directive.Directive) *jerr.JApiError

	// bannedDirectives a set of directives which didn't allow to use in schema.
	bannedDirectives map[directive.Enumeration]struct{}

	// rules all defined rules.
	rules map[string]schema.Rule

	// representation of Api data.
	catalog *catalog.Catalog

	// scanner for parsing .jst files.
	scanner *scanner.Scanner

	// scannersStack holds all saved scanners with position in that scanner when
	// it was pushed onto stack.
	// This property is used for processing INCLUDE keywords.
	scannersStack *scanner.Stack

	// currentContextDirective is current context for adding a child directive.
	currentContextDirective *directive.Directive

	// to accumulate directive data pieces from scanner.
	currentDirective *directive.Directive

	// rawPathVariables contains properties of the Path directives.
	rawPathVariables []rawPathVariable

	piecesOfPathVariables map[PathParameter]PieceOfPathVariable

	// directives from loaded from project.
	directives []*directive.Directive

	// directivesWithPastes directives after processing the MACRO and PASTE directives.
	directivesWithPastes []*directive.Directive
}

type Option func(*JApiCore)

// WithBannedDirectives ban specified directives.
// Core will emit an error if found one of them.
func WithBannedDirectives(dd ...directive.Enumeration) Option {
	return func(c *JApiCore) {
		if len(dd) == 0 {
			return
		}

		if c.bannedDirectives == nil {
			c.bannedDirectives = make(map[directive.Enumeration]struct{}, len(dd))
		}

		for _, v := range dd {
			c.bannedDirectives[v] = struct{}{}
		}
	}
}

func NewJApiCore(file *fs.File, oo ...Option) *JApiCore {
	core := &JApiCore{
		rawUserTypes:           &directive.Directives{},
		userTypes:              &catalog.UserSchemas{},
		processedUserTypes:     make(map[string]struct{}, 30),
		scanner:                scanner.NewJApiScanner(file),
		catalog:                catalog.NewCatalog(),
		currentDirective:       nil,
		directives:             make([]*directive.Directive, 0, 200),
		onlyOneProtocolIntoURL: make(map[*directive.Directive]struct{}, 20),
		uniqURLPath:            make(map[catalog.Path]struct{}, 20),
		similarPaths:           make(map[string]string, 20),
		rawPathVariables:       make([]rawPathVariable, 0, 40),
		piecesOfPathVariables:  make(map[PathParameter]PieceOfPathVariable, 20),
		macro:                  make(map[string]*directive.Directive, 20),
		scannersStack:          &scanner.Stack{},
		rules:                  map[string]schema.Rule{},
	}
	core.directiveFunctions = map[directive.Enumeration]func(*directive.Directive) *jerr.JApiError{
		directive.Jsight:           core.addJSight,
		directive.Info:             core.addInfo,
		directive.Title:            core.addTitle,
		directive.Version:          core.addVersion,
		directive.Description:      core.addDescription,
		directive.Server:           core.addServer,
		directive.BaseURL:          core.addBaseUrl,
		directive.Type:             core.addType,
		directive.URL:              core.addURL,
		directive.Get:              core.addHTTPMethod,
		directive.Post:             core.addHTTPMethod,
		directive.Put:              core.addHTTPMethod,
		directive.Patch:            core.addHTTPMethod,
		directive.Delete:           core.addHTTPMethod,
		directive.Query:            core.addQuery,
		directive.Request:          core.addRequest,
		directive.HTTPResponseCode: core.addResponse,
		directive.Headers:          core.addHeaders,
		directive.Body:             core.addBody,
		directive.Protocol:         core.addProtocol,
		directive.Method:           core.addJsonRpcMethod,
		directive.Params:           core.addJsonRpcParams,
		directive.Result:           core.addJsonRpcResult,
	}

	for _, o := range oo {
		o(core)
	}

	return core
}

func (core *JApiCore) BuildCatalog() *jerr.JApiError {
	return core.processJApiProject()
}

func (core *JApiCore) Catalog() *catalog.Catalog {
	return core.catalog
}

func (core *JApiCore) AddRawUserType(d *directive.Directive) {
	core.rawUserTypes.Set(d.NamedParameter("Name"), d)
}

// processJApiProject the main internal method of Core. Only when done, core is
// ready to handle data validation or doc rendering.
func (core *JApiCore) processJApiProject() *jerr.JApiError {
	if je := core.scanProject(); je != nil {
		return je
	}

	if je := core.compileCore(); je != nil {
		return je
	}

	if je := core.buildCatalog(); je != nil {
		return je
	}

	if je := core.compileCatalog(); je != nil {
		return je
	}

	return core.validateCatalog()
}
