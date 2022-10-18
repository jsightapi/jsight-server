package core

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/scanner"
)

func (core *JApiCore) scanProject() (je *jerr.JApiError) {
	defer func() {
		// We might get an error during scanning included file, and we should return
		// correct error in that case.
		core.scannersStack.AddIncludeTraceToError(je)
	}()

	for {
		if je := core.drainCurrentScanner(); je != nil {
			return je
		}

		if je := core.processEOF(); je != nil {
			return je
		}

		if core.isScanningFinished() {
			return nil
		}
	}
}

func (core *JApiCore) drainCurrentScanner() *jerr.JApiError {
	for {
		lexeme, je := core.scanner.Next()
		if je != nil {
			return je
		}

		if lexeme == nil { // EOF
			break
		}

		if isIncludeKeyword(lexeme) {
			je = core.processInclude(lexeme)
		} else {
			je = core.next(*lexeme)
		}
		if je != nil {
			return je
		}
	}
	return nil
}

// simply decides which function to call based on lexeme type
func (core *JApiCore) next(lexeme scanner.Lexeme) *jerr.JApiError {
	switch lexeme.Type() {
	case scanner.Keyword:
		return core.processKeyword(lexeme)

	case scanner.Parameter:
		return core.processParameter(lexeme)

	case scanner.Annotation:
		core.processAnnotation(lexeme)
		return nil

	case scanner.Schema, scanner.Text, scanner.Json, scanner.Enum:
		core.processBody(lexeme)
		return nil

	case scanner.ContextExplicitOpening:
		core.processContextBegin()
		return nil

	case scanner.ContextExplicitClosing:
		return core.processContextEnd()

	default:
		return core.japiError("Unknown lexeme type ("+lexeme.Type().String()+")", lexeme.Begin())
	}
}

func (core *JApiCore) processKeyword(lexeme scanner.Lexeme) *jerr.JApiError {
	// previous directive is ready to be processed
	if je := core.processCurrentDirective(); je != nil {
		return je
	}

	keyword := lexeme.Value().String()
	coords := coordsFromLexeme(lexeme)
	if !core.scannersStack.Empty() && keyword == directive.Jsight.String() {
		return core.japiError(fmt.Sprintf("directive %q not allowed in included file", keyword), coords.Begin())
	}

	return core.setCurrentDirective(keyword, coords)
}

func (core *JApiCore) processParameter(lexeme scanner.Lexeme) *jerr.JApiError {
	if err := core.currentDirective.AppendParameter(lexeme.Value()); err != nil {
		return core.japiError(err.Error(), lexeme.Begin())
	}
	return nil
}

func (core *JApiCore) processAnnotation(lexeme scanner.Lexeme) {
	core.currentDirective.Annotation = catalog.Annotation(lexeme.Value().String())
}

func (core *JApiCore) processBody(lexeme scanner.Lexeme) {
	core.currentDirective.BodyCoords = coordsFromLexeme(lexeme)
}

func (core *JApiCore) processContextBegin() {
	core.currentDirective.HasExplicitContext = true
}

func (core *JApiCore) closeLastExplicitContext() *jerr.JApiError {
	for {
		if core.currentContextDirective == nil {
			return core.japiError(jerr.ThereIsNoExplicitContextForClosure, core.scanner.CurrentIndex()-1)
		}

		if core.currentContextDirective.HasExplicitContext {
			core.currentContextDirective = core.currentContextDirective.Parent // Parent can be nil
			return nil
		}

		core.currentContextDirective = core.currentContextDirective.Parent
	}
}

func (core *JApiCore) HasUnclosedExplicitContext() bool {
	for d := core.currentContextDirective; d != nil; d = d.Parent {
		if d.HasExplicitContext {
			return true
		}
	}
	return false
}

func (core *JApiCore) processContextEnd() *jerr.JApiError {
	if je := core.processCurrentDirective(); je != nil {
		return je
	}
	return core.closeLastExplicitContext()
}

func (core *JApiCore) processEOF() *jerr.JApiError {
	// previous directive should be processed
	if je := core.processCurrentDirective(); je != nil {
		return je
	}
	if core.HasUnclosedExplicitContext() {
		return core.japiError("not all explicit contexts are closed", core.scanner.CurrentIndex()-1)
	}
	return nil
}

func (core *JApiCore) setCurrentDirective(keyword string, keywordCoords directive.Coords) *jerr.JApiError {
	de, err := directive.NewDirectiveType(keyword)
	if err != nil {
		return core.japiError(fmt.Sprintf("unknown directive %q", keyword), keywordCoords.Begin())
	}

	d := directive.NewWithCallStack(de, keywordCoords, core.scannersStack.ToDirectiveIncludeTracer())
	d.Keyword = keyword

	core.currentDirective = d

	return nil
}

func (core *JApiCore) isScanningFinished() bool {
	s := core.scannersStack.Pop()
	if s == nil {
		return true
	}
	core.scanner = s
	return false
}

func coordsFromLexeme(lex scanner.Lexeme) directive.Coords {
	return directive.NewCoords(lex.File(), lex.Begin(), lex.End())
}
