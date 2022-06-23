package core

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/fs"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/scanner"
)

// JApiCore the Brain and heart of jApi. Collects lexemes from scanner, validates document logic and structure,
// builds catalog, renders documentation
type JApiCore struct {
	// builtUserTypes a "set" of already build user types.
	processedUserTypes map[string]struct{}

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
	rules map[string]jschema.Rule

	// userTypes represent all user types.
	userTypes *catalog.UserSchemas

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

	// rawPathVariables contains properties of the Path directive.
	rawPathVariables []rawPathVariable

	// directives from loaded from project.
	directives []*directive.Directive

	// directivesWithPastes directives after processing the MACRO and PASTE directives.
	directivesWithPastes []*directive.Directive
}

type Option func(*JApiCore)

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
		userTypes:          &catalog.UserSchemas{},
		processedUserTypes: make(map[string]struct{}, 30),
		scanner:            scanner.NewJApiScanner(file),
		catalog:            catalog.NewCatalog(),
		currentDirective:   nil,
		directives:         make([]*directive.Directive, 0, 200),
		uniqURLPath:        make(map[catalog.Path]struct{}, 20),
		similarPaths:       make(map[string]string, 20),
		rawPathVariables:   make([]rawPathVariable, 0, 40),
		macro:              make(map[string]*directive.Directive, 20),
		scannersStack:      &scanner.Stack{},
		rules:              map[string]jschema.Rule{},
	}
	core.directiveFunctions = map[directive.Enumeration]func(*directive.Directive) *jerr.JApiError{
		directive.Jsight:           core.addJSight,
		directive.Info:             core.addInfo,
		directive.Title:            core.addTitle,
		directive.Version:          core.addVersion,
		directive.Description:      core.addDescription,
		directive.Server:           core.addServer,
		directive.BaseUrl:          core.addBaseUrl,
		directive.Type:             core.addType,
		directive.Url:              core.addURL,
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
		// TODO directive.Method:           core.addMethod,
		// TODO directive.Params:           core.addParams,
		// TODO directive.Result:           core.addResult,
	}

	for _, o := range oo {
		o(core)
	}

	return core
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
