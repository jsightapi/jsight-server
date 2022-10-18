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
		return "", requiredParameterNotSpecified(keyword)
	}

	if parameter.Type() != scanner.Parameter {
		return "", requiredParameterNotSpecified(keyword)
	}

	path := parameter.Value().String()

	if err := validateIncludeFileName(path); err != nil {
		return "", incorrectParameter(keyword, path, err.Error())
	}

	// We included file path is always will be relative to currently scanned file
	// directory.
	absolutePath := filepath.Join(filepath.Dir(core.scanner.File().Name()), path)
	info, err := os.Stat(absolutePath)
	if err == nil {
		if info.IsDir() {
			return "", incorrectParameter(keyword, path, "is a directory")
		}
		return absolutePath, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return "", incorrectParameter(keyword, path, "isn't exists")
	}
	return "", incorrectParameter(keyword, path, err.Error())
}

func readFile(p string) (*fs.File, error) {
	c, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return fs.NewFile(p, c), nil
}

func incorrectParameter(lex *scanner.Lexeme, path, msg string) *jerr.JApiError {
	return japiErrorForLexeme(lex, fmt.Sprintf("%s (%s) %q: %s", jerr.IncorrectParameter, "Filename", path, msg))
}

func requiredParameterNotSpecified(lex *scanner.Lexeme) *jerr.JApiError {
	return japiErrorForLexeme(lex, fmt.Sprintf("%s (%s)", jerr.RequiredParameterNotSpecified, "Filename"))
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
