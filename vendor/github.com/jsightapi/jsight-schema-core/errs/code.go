//nolint:lll
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
	ErrInfiniteRecursionDetected           Code = 104
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
	ErrUserTypeNotFound                               Code = 1302
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
	ErrUnknownValueOfTheTypeRule:           `Type %q does not exist. See the list of possible types here: https://jsight.io/docs/jsight-schema-0-3#rule-type`,
	ErrUnknownJSchemaType:                  `Type %q does not exist. See the list of possible types here: https://jsight.io/docs/jsight-schema-0-3#rule-type`,
	ErrInfiniteRecursionDetected:           "The infinite type recursion has been detected: %s. Use rules `optional: false` or `nullable: true` to stop the recursion.",
	ErrNodeTypeCantBeGuessed:               `Cannot determine the node type of the value "%s"`,
	ErrUnableToDetermineTheTypeOfJsonValue: "Unable to determine the type of the JSON value",

	// validator
	ErrValidator:                       "Validator error",
	ErrEmptySchema:                     "Empty schema",
	ErrEmptyJson:                       "Empty JSON",
	ErrOrRuleSetValidation:             "The example value does not match any of the options in the `or` rule. Change the example value. Learn more about the `or` rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or",
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
	ErrInvalidCharacterInAnnotationObjectKey: "Invalid character %s in the object key. See the rules for annotations here: https://jsight.io/docs/jsight-schema-0-3#rules",
	ErrUnexpectedEOF:                         "Unexpected end of file",
	ErrAnnotationNotAllowed:                  "The annotation is not allowed here. The ANNOTATION cannot be placed on lines containing more than one EXAMPLE element to which the ANNOTATION may apply. For more information, please refer to: https://jsight.io/docs/jsight-schema-0-3#rules",
	ErrEmptySetOfLexicalEvents:               "Empty set of found lexical events",
	ErrIncorrectEndingOfTheLexicalEvent:      "Incorrect ending of the lexical event",

	// schema
	ErrNodeGrow:                 "Node grow error",
	ErrDuplicateKeysInSchema:    "Duplicate key \"%s\"",
	ErrDuplicationOfNameOfTypes: "Duplication of the name of the types (%s)",

	// node
	ErrDuplicateRule:          "Duplicate rule %q",
	ErrUnexpectedLexicalEvent: "Unexpected lexical event %q %s",

	// constraint
	ErrUnknownRule:                                 `Unknown rule "%s". See the list of all possible rules here: https://jsight.io/docs/jsight-schema-0-3#rules`,
	ErrConstraintValidation:                        "The value in the example violates the rule `%q: %s` %s",
	ErrConstraintStringLengthValidation:            "The length of the string in the example violates the rule `%q: %q`",
	ErrInvalidValueOfConstraint:                    "Invalid value in the %q rule. Learn about the rules here: https://jsight.io/docs/jsight-schema-0-3#rules",
	ErrZeroPrecision:                               "Precision can not be zero",
	ErrEmptyEmail:                                  "Empty email",
	ErrInvalidEmail:                                "Invalid email (%s)",
	ErrConstraintMinItemsValidation:                `The number of the array elements does not match the "minItems" rule`,
	ErrConstraintMaxItemsValidation:                `The number of the array elements does not match the "maxItems" rule`,
	ErrDoesNotMatchAnyOfTheEnumValues:              "The value in the example does not match any of the enumeration values.",
	ErrDoesNotMatchRegularExpression:               "The value in the example does not match the regular expression.",
	ErrInvalidURI:                                  "Invalid URI (%s)",
	ErrInvalidDateTime:                             "Date/Time parsing error",
	ErrInvalidUUID:                                 "UUID parsing error: %s",
	ErrInvalidConst:                                "Does not match expected value (%s)",
	ErrInvalidDate:                                 "Date parsing error (%s)",
	ErrValueOfOneConstraintGreaterThanAnother:      "The value of the rule %q should be less or equal to the value of the rule %q", //nolint:lll
	ErrValueOfOneConstraintGreaterOrEqualToAnother: "The value of the rule %q should be less than the value of the rule %q",

	// loader
	ErrInvalidSchemaName:                `The type name "%s" is not valid. Learn more about the user types here: https://jsight.io/docs/jsight-schema-0-3#user-types`,
	ErrInvalidSchemaNameInAllOfRule:     "The type name \"%s\" is not valid. Learn more about the `allOf` rule here: https://jsight.io/docs/jsight-schema-0-3#rule-allof",
	ErrUnacceptableRecursionInAllOfRule: "The unacceptable recursion in the `allOf` rule",
	ErrUnacceptableUserTypeInAllOfRule:  `Unacceptable type. The "%s" type in the "allOf" rule must be an object`,
	ErrConflictAdditionalProperties:     `Conflicting value in additionalProperties rules when inheriting from allOf`,
	ErrLoadError:                        "load error: %w",

	// rule loader
	ErrLoader:                           "Loader error", // error somewhere in the loader code
	ErrIncorrectRuleValueType:           "Invalid rule value. Learn more about rules here: https://jsight.io/docs/jsight-schema-0-3#rules",
	ErrIncorrectRuleWithoutExample:      "You cannot place a RULE on a line without an EXAMPLE",
	ErrIncorrectRuleForSeveralNode:      "You cannot place a RULE on a line that contain more than one EXAMPLE value. The only exception is when an object key and its value are found in one line. Learn more about rules here: https://jsight.io/docs/jsight-schema-0-3#rules", //nolint:lll
	ErrLiteralValueExpected:             "Scalar value expected",
	ErrInvalidValueInEnumRule:           `An array or rule name was expected as a value for the "enum"`,
	ErrIncorrectArrayItemTypeInEnumRule: `Enums cannot contain arrays.`,
	ErrUnacceptableValueInAllOfRule:     "Incorrect value in the `allOf` rule. A type name or a list of type names are expected. Learn more about the `allOf` rule here: https://jsight.io/docs/jsight-schema-0-3#rule-allof", //nolint:lll
	ErrTypeNameNotFoundInAllOfRule:      "The `allOf` rule must contain a type name or an array of type names. Learn more about the `allOf` rule here: https://jsight.io/docs/jsight-schema-0-3#rule-allof",
	ErrDuplicationInEnumRule:            `The value %s is repeated in the "enum" rule!`,
	ErrRuleIsAlreadyCompiled:            "The rule is already compiled",
	ErrRuleIsNil:                        "The rule is nil",

	// "or" rule loader
	ErrArrayWasExpectedInOrRule:       `The "or" rule must contain an array. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,
	ErrEmptyArrayInOrRule:             `The empty array in the "or" rule! The "or" rule must contain a non-empty array. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,
	ErrOneElementInArrayInOrRule:      `The rule "or" must have at least two elements in the array. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,
	ErrIncorrectArrayItemTypeInOrRule: `Incorrect value in the "or" rule. The "or" array must contain strings (names of types) or objects with other rules. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,
	ErrEmptyRuleSet:                   `The object with the rules is empty! The "or" array must contain strings (names of types) or non-empty objects with other rules. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,
	ErrTypIsRequiredInsideOr:          `The "type" rule is missed inside the "or" rule. Specify the "type" rule inside. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,

	// compiler
	ErrRuleOptionalAppliesOnlyToObjectProperties: `The rule "optional" can be applied only to object properties!`,
	ErrCannotSpecifyOtherRulesWithTypeReference:  `Some of the rules can not be applied to the user type reference. Learn more about type referencing here: https://jsight.io/docs/jsight-schema-0-3#reference-to-the-user-type-in-the-example-value`,
	ErrShouldBeNoOtherRulesInSetWithOr:           `Some of the rules are not compatible with the "or" rule. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,
	ErrShouldBeNoOtherRulesInSetWithEnum:         `Some of the rules are not compatible with the "enum" rule. Learn more about the "enum" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-enum`,
	ErrShouldBeNoOtherRulesInSetWithAny:          `Some of the rules are not compatible with the "any" type. Learn more about types and rules compatibility here: https://jsight.io/docs/jsight-schema-0-3#appendix-1-a-table-of-all-built-in-types-and-rules`,
	ErrInvalidNestedElementsFoundForTypeAny:      `Example value for the type "any" can not have nested elements! Learn more about the type "any" here: https://jsight.io/docs/jsight-schema-0-3#type-any`,
	ErrInvalidChildNodeTogetherWithTypeReference: `Only scalar types can be referenced in the "type" rule. Use type references right in the example for referencing objects or arrays. See the examples here: https://jsight.io/docs/jsight-schema-0-3#reference-to-the-user-type-in-the-example-value`,
	ErrInvalidChildNodeTogetherWithOrRule:        `Only scalar values can be in an example when using the "or" rule. Use type references right in the example for referencing objects or arrays. See the examples here: https://jsight.io/docs/jsight-schema-0-3#reference-to-several-user-types-in-the-value-of-the-example`,
	ErrConstraintMinNotFound:                     `The rule "min" seems to be forgotten.`,
	ErrConstraintMaxNotFound:                     `The rule "max" seems to be forgotten.`,
	ErrInvalidValueInTheTypeRule:                 `The value of the "type" rule (%s) is not compatible with the other rules. Try to just remove the "type" rule.`,
	ErrNotFoundRulePrecision:                     `The rule "precision" is not found (required for the "decimal" type)`,
	ErrNotFoundRuleEnum:                          `The rule "enum" is not found (required for the "enum" type)`,
	ErrNotFoundRuleOr:                            `The rule "or" is not found (required for the "mixed" type)`,
	ErrIncompatibleTypes:                         `The value in the example does not match the specified type in the "type" rule (%s)`,
	// ErrUnknownAdditionalPropertiesTypes:          "Unknown type of additionalProperties (%s)",
	ErrUnexpectedConstraint: "The rule %q is not compatible with the %q type. Learn more about the rules and types compatibility here: https://jsight.io/docs/jsight-schema-0-3#appendix-1-a-table-of-all-built-in-types-and-rules",

	// checker
	ErrChecker:                               `Checker error`,
	ErrElementNotFoundInArray:                `Element not found in schema array node`,
	ErrIncorrectConstraintValueForEmptyArray: `The empty array in the example is not compatible with some of the rules. Learn more about the errors here: https://jsight.io/docs/jsight-schema-0-3#type-array`,

	// link checker
	ErrIncorrectUserType: "The value in the example does not match the rules!",
	ErrUserTypeNotFound:  "Type %q not found",
	ErrImpossibleToDetermineTheJsonTypeDueToRecursion: `It is impossible to determine the type due to the recursion of the type %q`, //nolint:lll
	ErrInvalidKeyShortcutType:                         "Reference in the object key %q must be string, not %q. Learn more about referencing user types in object properties here: https://jsight.io/docs/jsight-schema-0-3#reference-to-the-user-type-in-the-property-key",

	// sdk
	ErrEmptyType:                          `Type "%s" must not be empty`,
	ErrUnnecessaryLexemeAfterTheEndOfEnum: `An unnecessary non-space character after the end of the enum`,
	ErrRegexUnexpectedStart:               "Regular expression should start with the '/' character, not with %s",
	ErrRegexUnexpectedEnd:                 "Regular expression should end with the '/' character, not with %s",
	ErrRegexInvalid:                       "The regular expression is invalid: %s",

	// enum
	ErrEnumArrayExpected:  `An array was expected as a value for the "enum"`,
	ErrEnumIsHoldRuleName: "Can't append specific value to enum initialized with rule name",
	ErrEnumRuleNotFound:   "Enum %q is not found",
	ErrNotAnEnumRule:      "Rule %q is not an Enum",
	ErrInvalidEnumValues:  "Invalid enum values %q: %s",

	// value
	ErrInvalidBoolValue:         "Invalid boolean value",
	ErrNotEnoughDataInParseUint: "Not enough data in ParseUint",
	ErrInvalidByteInParseUint:   "Invalid byte %q in ParseUint %q",
	ErrTooMuchDataForInt:        "The value exceeds the maximum integer value",
	ErrIncorrectNumberValue:     "Incorrect number value %q",
	ErrURNPrefix:                "Invalid URN prefix: %q",
	ErrUUIDLength:               "Invalid UUID length: %d",
	ErrUUIDFormat:               "Invalid UUID format",
	ErrUUIDPrefix:               "Invalid prefix: braces expected",
	ErrIncorrectExponentValue:   "Incorrect exponent value",

	// example & ast
	ErrRegexExample:          "generate example for Regex type: %w",
	ErrCantCollectRulesTypes: `Can't collect rules: "types" constraint is required with "or" constraint. Learn more about the "or" rule here: https://jsight.io/docs/jsight-schema-0-3#rule-or`,

	// tests
	ErrInTheTest: "Error in the test: %s",
}

func (c Code) Itoa() string {
	return strconv.Itoa(int(c))
}

func (c Code) F(args ...any) *Err {
	return f(c, args...)
}
