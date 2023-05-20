package errs

import (
	"strconv"
)

type Code int

const (
	// Basic

	ErrGeneric        Code = 0
	ErrRuntimeFailure Code = 1

	// Common

	ErrUserTypeFound                       Code = 101
	ErrUnknownValueOfTheTypeRule           Code = 102
	ErrUnknownJSchemaType                  Code = 103
	ErrInfinityRecursionDetected           Code = 104
	ErrNodeTypeCantBeGuessed               Code = 105
	ErrUnableToDetermineTheTypeOfJsonValue Code = 106

	// Validator

	ErrValidator                       Code = 201
	ErrEmptySchema                     Code = 202
	ErrEmptyJson                       Code = 203
	ErrOrRuleSetValidation             Code = 204
	ErrRequiredKeyNotFound             Code = 205
	ErrSchemaDoesNotSupportKey         Code = 206
	ErrUnexpectedLexInLiteralValidator Code = 207
	ErrUnexpectedLexInObjectValidator  Code = 208
	ErrUnexpectedLexInArrayValidator   Code = 209
	ErrInvalidValueType                Code = 210
	ErrInvalidKeyType                  Code = 211
	ErrUnexpectedLexInMixedValidator   Code = 212
	ErrObjectExpected                  Code = 213
	ErrPropertyNotFound                Code = 214

	// Scanner

	ErrInvalidCharacter                      Code = 301
	ErrInvalidCharacterInAnnotationObjectKey Code = 302
	ErrUnexpectedEOF                         Code = 303
	ErrAnnotationNotAllowed                  Code = 304
	ErrEmptySetOfLexicalEvents               Code = 305
	ErrIncorrectEndingOfTheLexicalEvent      Code = 306

	// Schema

	ErrNodeGrow                 Code = 401
	ErrDuplicateKeysInSchema    Code = 402
	ErrDuplicationOfNameOfTypes Code = 403

	// Node

	ErrDuplicateRule          Code = 501
	ErrUnexpectedLexicalEvent Code = 502

	// Constraint

	ErrUnknownRule                                 Code = 601
	ErrConstraintValidation                        Code = 602
	ErrConstraintStringLengthValidation            Code = 603
	ErrInvalidValueOfConstraint                    Code = 604
	ErrZeroPrecision                               Code = 605
	ErrEmptyEmail                                  Code = 606
	ErrInvalidEmail                                Code = 607
	ErrConstraintMinItemsValidation                Code = 608
	ErrConstraintMaxItemsValidation                Code = 609
	ErrDoesNotMatchAnyOfTheEnumValues              Code = 610
	ErrDoesNotMatchRegularExpression               Code = 611
	ErrInvalidURI                                  Code = 612
	ErrInvalidDateTime                             Code = 613
	ErrInvalidUUID                                 Code = 614
	ErrInvalidConst                                Code = 615
	ErrInvalidDate                                 Code = 616
	ErrValueOfOneConstraintGreaterThanAnother      Code = 617
	ErrValueOfOneConstraintGreaterOrEqualToAnother Code = 618

	// Loader

	ErrInvalidSchemaName                Code = 701
	ErrInvalidSchemaNameInAllOfRule     Code = 702
	ErrUnacceptableRecursionInAllOfRule Code = 703
	ErrUnacceptableUserTypeInAllOfRule  Code = 704
	ErrConflictAdditionalProperties     Code = 705
	ErrLoadError                        Code = 706

	// Rule loader

	ErrLoader                           Code = 801
	ErrIncorrectRuleValueType           Code = 802
	ErrIncorrectRuleWithoutExample      Code = 803
	ErrIncorrectRuleForSeveralNode      Code = 804
	ErrLiteralValueExpected             Code = 805
	ErrInvalidValueInEnumRule           Code = 806
	ErrIncorrectArrayItemTypeInEnumRule Code = 807
	ErrUnacceptableValueInAllOfRule     Code = 808
	ErrTypeNameNotFoundInAllOfRule      Code = 809
	ErrDuplicationInEnumRule            Code = 810
	ErrRuleIsAlreadyCompiled            Code = 811
	ErrRuleIsNil                        Code = 812

	// "or" rule loader

	ErrArrayWasExpectedInOrRule       Code = 901
	ErrEmptyArrayInOrRule             Code = 902
	ErrOneElementInArrayInOrRule      Code = 903
	ErrIncorrectArrayItemTypeInOrRule Code = 904
	ErrEmptyRuleSet                   Code = 905
	ErrTypIsRequiredInsideOr          Code = 906

	// Compiler

	ErrRuleOptionalAppliesOnlyToObjectProperties Code = 1101
	ErrCannotSpecifyOtherRulesWithTypeReference  Code = 1102
	ErrShouldBeNoOtherRulesInSetWithOr           Code = 1103
	ErrShouldBeNoOtherRulesInSetWithEnum         Code = 1104
	ErrShouldBeNoOtherRulesInSetWithAny          Code = 1105
	ErrInvalidNestedElementsFoundForTypeAny      Code = 1106
	ErrInvalidChildNodeTogetherWithTypeReference Code = 1107
	ErrInvalidChildNodeTogetherWithOrRule        Code = 1108
	ErrConstraintMinNotFound                     Code = 1109
	ErrConstraintMaxNotFound                     Code = 1110
	ErrInvalidValueInTheTypeRule                 Code = 1111
	ErrNotFoundRulePrecision                     Code = 1112
	ErrNotFoundRuleEnum                          Code = 1113
	ErrNotFoundRuleOr                            Code = 1114
	ErrIncompatibleTypes                         Code = 1115

	ErrUnexpectedConstraint Code = 1117

	// Checker

	ErrChecker                               Code = 1201
	ErrElementNotFoundInArray                Code = 1203
	ErrIncorrectConstraintValueForEmptyArray Code = 1204

	// Link checker

	ErrIncorrectUserType                              Code = 1301
	ErrTypeNotFound                                   Code = 1302
	ErrImpossibleToDetermineTheJsonTypeDueToRecursion Code = 1303
	ErrInvalidKeyShortcutType                         Code = 1304

	// SDK

	ErrEmptyType                          Code = 1401
	ErrUnnecessaryLexemeAfterTheEndOfEnum Code = 1402

	// Regex

	ErrRegexUnexpectedStart Code = 1500
	ErrRegexUnexpectedEnd   Code = 1501
	ErrRegexInvalid         Code = 1502

	// Enum

	ErrEnumArrayExpected  Code = 1600
	ErrEnumIsHoldRuleName Code = 1601
	ErrEnumRuleNotFound   Code = 1602
	ErrNotAnEnumRule      Code = 1603
	ErrInvalidEnumValues  Code = 1604

	// Value

	ErrInvalidBoolValue         Code = 1701
	ErrNotEnoughDataInParseUint Code = 1702
	ErrInvalidByteInParseUint   Code = 1703
	ErrTooMuchDataForInt        Code = 1704
	ErrIncorrectNumberValue     Code = 1705
	ErrURNPrefix                Code = 1706
	ErrUUIDLength               Code = 1708
	ErrUUIDFormat               Code = 1709
	ErrUUIDPrefix               Code = 1710
	ErrIncorrectExponentValue   Code = 1711

	// Example & AST

	ErrRegexExample          Code = 1801
	ErrCantCollectRulesTypes Code = 1802

	// Tests

	ErrInTheTest Code = 9901
)

var errorFormat = map[Code]string{
	ErrGeneric:        "%s",
	ErrRuntimeFailure: "Runtime Failure",

	// main & common
	ErrUserTypeFound:                       "Found an invalid reference to the type",
	ErrUnknownValueOfTheTypeRule:           "Unknown value of the type rule %q",
	ErrUnknownJSchemaType:                  "Unknown JSchema type %q",
	ErrInfinityRecursionDetected:           "Infinity recursion detected %s",
	ErrNodeTypeCantBeGuessed:               "Node type can't be guessed by value (%s)",
	ErrUnableToDetermineTheTypeOfJsonValue: "Unable to determine the type of JSON value",

	// validator
	ErrValidator:                       "Validator error",
	ErrEmptySchema:                     "Empty schema",
	ErrEmptyJson:                       "Empty JSON",
	ErrOrRuleSetValidation:             `None of the rules in the "OR" set has been validated`,
	ErrRequiredKeyNotFound:             `Required key(s) %q not found`,
	ErrSchemaDoesNotSupportKey:         `Schema does not support key %q`,
	ErrUnexpectedLexInLiteralValidator: `Invalid value, scalar expected`,
	ErrUnexpectedLexInObjectValidator:  `Invalid value, object expected`,
	ErrUnexpectedLexInArrayValidator:   `Invalid value, array expected`,
	ErrUnexpectedLexInMixedValidator:   `Invalid value, scalar, array, or object expected`,
	ErrInvalidValueType:                "Invalid value type `%s`, expected `%s`",
	ErrInvalidKeyType:                  `Incorrect key type "%s"`,
	ErrObjectExpected:                  `An object is expected to validate the property`,
	ErrPropertyNotFound:                `The %q property was not found`,

	// scanner
	ErrInvalidCharacter:                      "Invalid character %q %s",
	ErrInvalidCharacterInAnnotationObjectKey: "Invalid character %s in object key (inside comment)",
	ErrUnexpectedEOF:                         "Unexpected end of file",
	ErrAnnotationNotAllowed:                  "Annotation not allowed here",
	ErrEmptySetOfLexicalEvents:               "Empty set of found lexical events",
	ErrIncorrectEndingOfTheLexicalEvent:      "Incorrect ending of the lexical event",

	// schema
	ErrNodeGrow:                 "Node grow error",
	ErrDuplicateKeysInSchema:    "Duplicate keys (%s) in the schema",
	ErrDuplicationOfNameOfTypes: "Duplication of the name of the types (%s)",

	// node
	ErrDuplicateRule:          "Duplicate %q rule",
	ErrUnexpectedLexicalEvent: "Unexpected lexical event %q %s",

	// constraint
	ErrUnknownRule:                                 `Unknown rule "%s"`,
	ErrConstraintValidation:                        "Invalid value for %q = %s constraint %s",
	ErrConstraintStringLengthValidation:            "Invalid string length for %q = %q constraint",
	ErrInvalidValueOfConstraint:                    "Invalid value of %q constraint",
	ErrZeroPrecision:                               "Precision can't be zero",
	ErrEmptyEmail:                                  "Empty email",
	ErrInvalidEmail:                                "Invalid email (%s)",
	ErrConstraintMinItemsValidation:                `The number of array elements does not match the "minItems" rule`,
	ErrConstraintMaxItemsValidation:                `The number of array elements does not match the "maxItems" rule`,
	ErrDoesNotMatchAnyOfTheEnumValues:              "Does not match any of the enumeration values",
	ErrDoesNotMatchRegularExpression:               "Does not match regular expression",
	ErrInvalidURI:                                  "Invalid URI (%s)",
	ErrInvalidDateTime:                             "Date/Time parsing error",
	ErrInvalidUUID:                                 "UUID parsing error: %s",
	ErrInvalidConst:                                "Does not match expected value (%s)",
	ErrInvalidDate:                                 "Date parsing error (%s)",
	ErrValueOfOneConstraintGreaterThanAnother:      "Value of constraint %q should be less or equal to value of %q constraint", //nolint:lll
	ErrValueOfOneConstraintGreaterOrEqualToAnother: "Value of constraint %q should be less than value of %q constraint",

	// loader
	ErrInvalidSchemaName:                "Invalid schema name (%s)",
	ErrInvalidSchemaNameInAllOfRule:     `Invalid schema name (%s) in "allOf" rule`,
	ErrUnacceptableRecursionInAllOfRule: `Unacceptable recursion in "allOf" rule`,
	ErrUnacceptableUserTypeInAllOfRule:  `Unacceptable type. The "%s" type in the "allOf" rule must be an object`,
	ErrConflictAdditionalProperties:     `Conflicting value in AdditionalProperties rules when inheriting from allOf`,
	ErrLoadError:                        "load error: %w",

	// rule loader
	ErrLoader:                           "Loader error", // error somewhere in the loader code
	ErrIncorrectRuleValueType:           "Incorrect rule value type",
	ErrIncorrectRuleWithoutExample:      "You cannot place a RULE on line without EXAMPLE",
	ErrIncorrectRuleForSeveralNode:      "You cannot place a RULE on lines that contain more than one EXAMPLE node to which any RULES can apply. The only exception is when an object key and its value are found in one line.", //nolint:lll
	ErrLiteralValueExpected:             "Literal value expected",
	ErrInvalidValueInEnumRule:           `An array or rule name was expected as a value for the "enum"`,
	ErrIncorrectArrayItemTypeInEnumRule: `Incorrect array item type in "enum". Only literals are allowed.`,
	ErrUnacceptableValueInAllOfRule:     `Incorrect value in "allOf" rule. A type name, or list of type names, is expected.`, //nolint:lll
	ErrTypeNameNotFoundInAllOfRule:      `Type name not found in "allOf" rule`,
	ErrDuplicationInEnumRule:            `%s value duplicates in "enum"`,
	ErrRuleIsAlreadyCompiled:            "Rule is already compiled",
	ErrRuleIsNil:                        "Rule is nil",

	// "or" rule loader
	ErrArrayWasExpectedInOrRule:       `An array was expected as a value for the "or" rule`,
	ErrEmptyArrayInOrRule:             `Empty array in "or" rule`,
	ErrOneElementInArrayInOrRule:      `Array rule "or" must have at least two elements`,
	ErrIncorrectArrayItemTypeInOrRule: `Incorrect array item type in "or" rule`,
	ErrEmptyRuleSet:                   `Empty rule set`,
	ErrTypIsRequiredInsideOr:          `Type is required inside the "or" rule `,

	// compiler
	ErrRuleOptionalAppliesOnlyToObjectProperties: `The rule "optional" applies only to object properties`,
	ErrCannotSpecifyOtherRulesWithTypeReference:  `Invalid rule set shared with a type reference`,
	ErrShouldBeNoOtherRulesInSetWithOr:           `Invalid rule set shared with "or"`,
	ErrShouldBeNoOtherRulesInSetWithEnum:         `Invalid rule set shared with "enum"`,
	ErrShouldBeNoOtherRulesInSetWithAny:          `Invalid rule set shared with "any"`,
	ErrInvalidNestedElementsFoundForTypeAny:      `Invalid nested elements found for an element of type "any"`,
	ErrInvalidChildNodeTogetherWithTypeReference: `You cannot specify child node if you use a type reference`,
	ErrInvalidChildNodeTogetherWithOrRule:        `You cannot specify child node if you use a "or" rule`,
	ErrConstraintMinNotFound:                     `Constraint "min" not found`,
	ErrConstraintMaxNotFound:                     `Constraint "max" not found`,
	ErrInvalidValueInTheTypeRule:                 `Invalid value in the "type" rule (%s)`,
	ErrNotFoundRulePrecision:                     `Not found the rule "precision" for the "decimal" type`,
	ErrNotFoundRuleEnum:                          `Not found the rule "enum" for the "enum" type`,
	ErrNotFoundRuleOr:                            `Not found the rule "or" for the "mixed" type`,
	ErrIncompatibleTypes:                         `Incompatible value of example and "type" rule (%s)`,
	// ErrUnknownAdditionalPropertiesTypes:          "Unknown type of additionalProperties (%s)",
	ErrUnexpectedConstraint: "The %q constraint can't be used for the %q type",

	// checker
	ErrChecker:                               `Checker error`,
	ErrElementNotFoundInArray:                `Element not found in schema array node`,
	ErrIncorrectConstraintValueForEmptyArray: `Incorrect constraint value for empty array`,

	// link checker
	ErrIncorrectUserType: "Incorrect type of user type",
	ErrTypeNotFound:      "Type %q not found",
	ErrImpossibleToDetermineTheJsonTypeDueToRecursion: `It is impossible to determine the json type due to recursion of type %q`, //nolint:lll
	ErrInvalidKeyShortcutType:                         "Key shortcut %q should be string but %q given",

	// sdk
	ErrEmptyType:                          `Type "%s" must not be empty`,
	ErrUnnecessaryLexemeAfterTheEndOfEnum: `An unnecessary non-space character after the end of the enum`,
	ErrRegexUnexpectedStart:               "Regex should starts with '/' character, but found %s",
	ErrRegexUnexpectedEnd:                 "Regex should ends with '/' character, but found %s",
	ErrRegexInvalid:                       "Invalid regex %s",

	// enum
	ErrEnumArrayExpected:  `An array was expected as a value for the "enum"`,
	ErrEnumIsHoldRuleName: "Can't append specific value to enum initialized with rule name",
	ErrEnumRuleNotFound:   "Enum rule %q not found",
	ErrNotAnEnumRule:      "Rule %q not an Enum",
	ErrInvalidEnumValues:  "Invalid enum values %q: %s",

	// value
	ErrInvalidBoolValue:         "Invalid bool value",
	ErrNotEnoughDataInParseUint: "Not enough data in ParseUint",
	ErrInvalidByteInParseUint:   "Invalid byte %q in ParseUint %q",
	ErrTooMuchDataForInt:        "Too much data for int",
	ErrIncorrectNumberValue:     "Incorrect number value %q",
	ErrURNPrefix:                "Invalid URN prefix: %q",
	ErrUUIDLength:               "Invalid UUID length: %d",
	ErrUUIDFormat:               "Invalid UUID format",
	ErrUUIDPrefix:               "Invalid prefix: braces expected",
	ErrIncorrectExponentValue:   "Incorrect exponent value",

	// example & ast
	ErrRegexExample:          "generate example for Regex type: %w",
	ErrCantCollectRulesTypes: `Can't collect rules: "types" constraint is required with "or" constraint`,

	// tests
	ErrInTheTest: "Error in the test: %s",
}

func (c Code) Itoa() string {
	return strconv.Itoa(int(c))
}

func (c Code) F(args ...any) *Err {
	return f(c, args...)
}
