package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jsightapi/jsight-schema-go-library/fs"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/scanner"
)

func (core *JApiCore) processInclude(keyword *scanner.Lexeme) *jerr.JApiError {
	// We got the "INCLUDE" directive here.
	// This directive shouldn't be among core.directives, because we simply
	// "paste" included file content inside current file.

	path, je := core.getIncludedFilePath(keyword)
	if je != nil {
		return je
	}

	file, err := readFile(path)
	if err != nil {
		return japiErrorForLexeme(keyword, fmt.Sprintf("%s (%s) %s", jerr.IncorrectParameter, "Filename", err))
	}

	if err := core.scannersStack.Push(core.scanner, keyword.Begin()); err != nil {
		return japiErrorForLexeme(keyword, err.Error())
	}
	core.scanner = scanner.NewJApiScanner(file)

	return nil
}

func (core *JApiCore) getIncludedFilePath(keyword *scanner.Lexeme) (string, *jerr.JApiError) {
	parameter, je := core.scanner.Next()
	if je != nil {
		return "", je
	}

	if parameter == nil {
		return "", japiErrorForLexeme(keyword, fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Filename"))
	}

	if parameter.Type() != scanner.Parameter {
		return "", japiErrorForLexeme(keyword, fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Filename"))
	}

	path := parameter.Value().String()

	if err := validateIncludeFileName(path); err != nil {
		return "", japiErrorForLexeme(keyword, fmt.Sprintf("%s (%s) %s", jerr.IncorrectParameter, "Filename", err))
	}

	// We included file path is always will be relative to currently scanned file
	// directory.
	return filepath.Join(filepath.Dir(core.scanner.File().Name()), path), nil
}

func readFile(p string) (*fs.File, error) {
	c, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return fs.NewFile(p, c), nil
}

func japiErrorForLexeme(lex *scanner.Lexeme, msg string) *jerr.JApiError {
	return jerr.NewJApiError(msg, lex.File(), lex.Begin())
}

func isIncludeKeyword(lex *scanner.Lexeme) bool {
	return lex != nil &&
		lex.Type() == scanner.Keyword &&
		lex.Value().String() == directive.Include.String()
}

func validateIncludeFileName(s string) error {
	if s[0] == '/' {
		return errors.New("mustn't starts with '/'")
	}

	hasForbiddenParts := strings.Contains(s, "/./") ||
		strings.Contains(s, "./") ||
		strings.Contains(s, "/.") ||
		strings.Contains(s, "/../") ||
		strings.Contains(s, "../") ||
		strings.Contains(s, "/..")
	if hasForbiddenParts {
		return errors.New("mustn't include '..' or '.'")
	}

	if strings.ContainsRune(s, '\\') {
		return errors.New("the separator for directories and files should be the symbol '/'")
	}

	return nil
}
