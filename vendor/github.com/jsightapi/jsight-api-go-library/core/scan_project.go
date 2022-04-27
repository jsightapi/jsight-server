package core

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/fs"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/scanner"
)

func (core *JApiCore) scanProject() *jerr.JAPIError {
	core.setScanner(core.file)

	for {
		lexeme, je := core.scanner.Next()
		if je != nil {
			return je
		}

		if lexeme == nil { // EOF
			break
		}

		if je := core.next(*lexeme); je != nil {
			return je
		}
	}

	return core.processEOF()
}

func (core *JApiCore) setScanner(file *fs.File) {
	core.scanner = scanner.NewJApiScanner(file)
	core.scanner.SetCurrentIndex(0)
}

// simply decides which function to call based on lexeme type
func (core *JApiCore) next(lexeme scanner.Lexeme) *jerr.JAPIError {
	switch lexeme.Type() {
	case scanner.Keyword:
		return core.processKeyword(lexeme)
	case scanner.Parameter:
		return core.processParameter(lexeme)
	case scanner.Annotation:
		core.processAnnotation(lexeme)
		return nil
	case scanner.Schema, scanner.Array, scanner.Text, scanner.Json:
		core.processBody(lexeme)
		return nil
	case scanner.ContextExplicitOpening:
		core.processContextBegin()
		return nil
	case scanner.ContextExplicitClosing:
		return core.processContextEnd()
	default:
		return core.japiError(`Unknown lexeme type (`+lexeme.Type().String()+`)`, lexeme.Begin())
	}
}

func (core *JApiCore) processKeyword(lexeme scanner.Lexeme) *jerr.JAPIError {
	// previous directive is ready to be processed
	if je := core.processCurrentDirective(); je != nil {
		return je
	}

	keyword := lexeme.Value().String()
	return core.setCurrentDirective(keyword, coordsFromLexeme(lexeme))
}

func (core *JApiCore) processParameter(lexeme scanner.Lexeme) *jerr.JAPIError {
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

func (core *JApiCore) processContextBegin( /*lexeme scanner.Lexeme*/ ) {
	core.currentDirective.HasExplicitContext = true
}

func (core *JApiCore) closeLastExplicitContext() *jerr.JAPIError {
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
	d := core.currentContextDirective

	for {
		if d == nil {
			return false
		}

		if d.HasExplicitContext {
			return true
		}

		d = d.Parent
	}
}

func (core *JApiCore) processContextEnd() *jerr.JAPIError {
	if je := core.processCurrentDirective(); je != nil {
		return je
	}
	return core.closeLastExplicitContext()
}

func (core *JApiCore) processEOF() *jerr.JAPIError {
	// previous directive should be processed
	if je := core.processCurrentDirective(); je != nil {
		return je
	}
	if core.HasUnclosedExplicitContext() {
		return core.japiError("not all explicit contexts are closed", core.scanner.CurrentIndex()-1)
	}
	return nil
}

func (core *JApiCore) setCurrentDirective(keyword string, keywordCoords directive.Coords) *jerr.JAPIError {
	de, err := directive.NewDirectiveType(keyword)
	if err != nil {
		return core.japiError(fmt.Sprintf("unknown directive %q", keyword), keywordCoords.B())
	}

	d := directive.New(de, keywordCoords)
	d.Keyword = keyword

	core.currentDirective = d

	return nil
}

func coordsFromLexeme(lex scanner.Lexeme) directive.Coords {
	return directive.NewCoords(lex.File(), lex.Begin(), lex.End())
}
