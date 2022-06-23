package errors

import (
	"strconv"
	"strings"
)

type ErrorCode int //nolint:errname // This is okay.

const (
	ErrGeneric    ErrorCode = 0
	ErrImpossible ErrorCode = 1

	// main & common
	ErrUserTypeFound      ErrorCode = 101
	ErrUnknownType        ErrorCode = 102
	ErrUnknownJSchemaType ErrorCode = 103

	// validator
	ErrValidator                       ErrorCode = 201
	ErrEmptySchema                     ErrorCode = 202
	ErrEmptyJson                       ErrorCode = 203
	ErrOrRuleSetValidation             ErrorCode = 204
	ErrRequiredKeyNotFound             ErrorCode = 205
	ErrSchemaDoesNotSupportKey         ErrorCode = 206
	ErrUnexpectedLexInLiteralValidator ErrorCode = 207
	ErrUnexpectedLexInObjectValidator  ErrorCode = 208
	ErrUnexpectedLexInArrayValidator   ErrorCode = 209
	ErrInvalidValueType                ErrorCode = 210
	ErrInvalidKeyType                  ErrorCode = 211
	ErrUnexpectedLexInMixedValidator   ErrorCode = 212

	// scanner
	ErrInvalidCharacter                      ErrorCode = 301
	ErrInvalidCharacterInAnnotationObjectKey ErrorCode = 302
	ErrUnexpectedEOF                         ErrorCode = 303
	ErrAnnotationNotAllowed                  ErrorCode = 304

	// schema
	ErrNodeGrow                 ErrorCode = 401
	ErrDuplicateKeysInSchema    ErrorCode = 402
	ErrDuplicationOfNameOfTypes ErrorCode = 403

	// node
	ErrDuplicateRule ErrorCode = 501

	// constraint
	ErrUnknownRule                      ErrorCode = 601
	ErrConstraintValidation             ErrorCode = 602
	ErrConstraintStringLengthValidation ErrorCode = 603
	ErrInvalidValueOfConstraint         ErrorCode = 604
	ErrZeroPrecision                    ErrorCode = 605
	ErrEmptyEmail                       ErrorCode = 606
	ErrInvalidEmail                     ErrorCode = 607
	ErrConstraintMinItemsValidation     ErrorCode = 608
	ErrConstraintMaxItemsValidation     ErrorCode = 609
	ErrDoesNotMatchAnyOfTheEnumValues   ErrorCode = 610
	ErrDoesNotMatchRegularExpression    ErrorCode = 611
	ErrInvalidUri                       ErrorCode = 612
	ErrInvalidDateTime                  ErrorCode = 613
	ErrInvalidUuid                      ErrorCode = 614
	ErrInvalidConst                     ErrorCode = 615
	ErrInvalidDate                      ErrorCode = 616

	// loader
	ErrInvalidSchemaName                ErrorCode = 701
	ErrInvalidSchemaNameInAllOfRule     ErrorCode = 702
	ErrUnacceptableRecursionInAllOfRule ErrorCode = 703
	ErrUnacceptableUserTypeInAllOfRule  ErrorCode = 704
	ErrConflictAdditionalProperties     ErrorCode = 705

	// rule loader
	ErrLoader                           ErrorCode = 801
	ErrIncorrectRuleValueType           ErrorCode = 802
	ErrIncorrectRuleWithoutExample      ErrorCode = 803
	ErrIncorrectRuleForSeveralNode      ErrorCode = 804
	ErrLiteralValueExpected             ErrorCode = 805
	ErrInvalidValueInEnumRule           ErrorCode = 806
	ErrIncorrectArrayItemTypeInEnumRule ErrorCode = 807
	ErrUnacceptableValueInAllOfRule     ErrorCode = 808
	ErrTypeNameNotFoundInAllOfRule      ErrorCode = 809
	ErrDuplicationInEnumRule            ErrorCode = 810

	// "or" rule loader
	ErrArrayWasExpectedInOrRule       ErrorCode = 901
	ErrEmptyArrayInOrRule             ErrorCode = 902
	ErrOneElementInArrayInOrRule      ErrorCode = 903
	ErrIncorrectArrayItemTypeInOrRule ErrorCode = 904
	ErrEmptyRuleSet                   ErrorCode = 905

	// compiler
	ErrRuleOptionalAppliesOnlyToObjectProperties ErrorCode = 1101
	ErrCannotSpecifyOtherRulesWithTypeReference  ErrorCode = 1102
	ErrShouldBeNoOtherRulesInSetWithOr           ErrorCode = 1103
	ErrShouldBeNoOtherRulesInSetWithEnum         ErrorCode = 1104
	ErrShouldBeNoOtherRulesInSetWithAny          ErrorCode = 1105
	ErrInvalidNestedElementsFoundForTypeAny      ErrorCode = 1106
	ErrInvalidChildNodeTogetherWithTypeReference ErrorCode = 1107
	ErrInvalidChildNodeTogetherWithOrRule        ErrorCode = 1108
	ErrConstraintMinNotFound                     ErrorCode = 1109
	ErrConstraintMaxNotFound                     ErrorCode = 1110
	ErrInvalidValueInTheTypeRule                 ErrorCode = 1111
	ErrNotFoundRulePrecision                     ErrorCode = 1112
	ErrNotFoundRuleEnum                          ErrorCode = 1113
	ErrNotFoundRuleOr                            ErrorCode = 1114
	ErrIncompatibleTypes                         ErrorCode = 1115
	ErrUnknownAdditionalPropertiesTypes          ErrorCode = 1116
	ErrUnexpectedConstraint                      ErrorCode = 1117

	// checker
	ErrChecker                               ErrorCode = 1201
	ErrElementNotFoundInArray                ErrorCode = 1203
	ErrIncorrectConstraintValueForEmptyArray ErrorCode = 1204

	// link checker
	ErrIncorrectUserType                              ErrorCode = 1301
	ErrTypeNotFound                                   ErrorCode = 1302
	ErrImpossibleToDetermineTheJsonTypeDueToRecursion ErrorCode = 1303

	// sdk
	ErrEmptyType                          ErrorCode = 1401
	ErrUnnecessaryLexemeAfterTheEndOfEnum ErrorCode = 1402

	// regex
	ErrRegexUnexpectedStart ErrorCode = 1500
	ErrRegexUnexpectedEnd   ErrorCode = 1501

	// enum
	ErrEnumArrayExpected  ErrorCode = 1600
	ErrEnumIsHoldRuleName ErrorCode = 1601
	ErrEnumRuleNotFound   ErrorCode = 1602
	ErrNotAnEnumRule      ErrorCode = 1603
)

var errorFormat = map[ErrorCode]string{
	// old error format
	ErrGeneric: "%s",

	ErrImpossible: "The error should not occur during regular operation. May appear only in the process of unfinished code refactoring.",

	// main & common
	ErrUserTypeFound:      "Found an invalid reference to the type",
	ErrUnknownType:        "Unknown type %q",
	ErrUnknownJSchemaType: "Unknown JSchema type %q",

	// validator
	ErrValidator:                       "Validator error",
	ErrEmptySchema:                     "Empty schema",
	ErrEmptyJson:                       "Empty JSON",
	ErrOrRuleSetValidation:             `None of the rules in the "OR" set has been validated`,
	ErrRequiredKeyNotFound:             `Required key "%s" not found`,
	ErrSchemaDoesNotSupportKey:         `Schema does not support key "%s"`,
	ErrUnexpectedLexInLiteralValidator: `Invalid value, scalar expected`,
	ErrUnexpectedLexInObjectValidator:  `Invalid value, object expected`,
	ErrUnexpectedLexInArrayValidator:   `Invalid value, array expected`,
	ErrUnexpectedLexInMixedValidator:   `Invalid value, scalar, array, or object expected`,
	ErrInvalidValueType:                `Invalid value type "%s", expected "%s"`,
	ErrInvalidKeyType:                  `Incorrect key type "%s"`,

	// scanner
	ErrInvalidCharacter:                      "Invalid character %q %s",
	ErrInvalidCharacterInAnnotationObjectKey: "Invalid character %s in object key (inside comment)",
	ErrUnexpectedEOF:                         "Unexpected end of file",
	ErrAnnotationNotAllowed:                  "Annotation not allowed here",

	// schema
	ErrNodeGrow:                 "Node grow error",
	ErrDuplicateKeysInSchema:    "Duplicate keys (%s) in the schema",
	ErrDuplicationOfNameOfTypes: "Duplication of the name of the types (%s)",

	// node
	ErrDuplicateRule: `Duplicate "%s" rule`,

	// constraint
	ErrUnknownRule:                      `Unknown rule "%s"`,
	ErrConstraintValidation:             `Invalid value for "%s" = %s constraint %s`,
	ErrConstraintStringLengthValidation: `Invalid string length for "%s" = "%s" constraint`,
	ErrInvalidValueOfConstraint:         `Invalid value of "%s" constraint`,
	ErrZeroPrecision:                    `Precision can't be zero`,
	ErrEmptyEmail:                       `Empty email`,
	ErrInvalidEmail:                     `Invalid email (%s)`,
	ErrConstraintMinItemsValidation:     `The number of array elements does not match the "minItems" rule`,
	ErrConstraintMaxItemsValidation:     `The number of array elements does not match the "maxItems" rule`,
	ErrDoesNotMatchAnyOfTheEnumValues:   `Does not match any of the enumeration values`,
	ErrDoesNotMatchRegularExpression:    `Does not match regular expression`,
	ErrInvalidUri:                       `Invalid URI (%s)`,
	ErrInvalidDateTime:                  `Date/Time parsing error`,
	ErrInvalidUuid:                      `UUID parsing error: %s`,
	ErrInvalidConst:                     "Does not match expected value (%s)",
	ErrInvalidDate:                      `Date parsing error (%s)`,

	// loader
	ErrInvalidSchemaName:                "Invalid schema name (%s)",
	ErrInvalidSchemaNameInAllOfRule:     `Invalid schema name (%s) in "allOf" rule`,
	ErrUnacceptableRecursionInAllOfRule: `Unacceptable recursion in "allOf" rule`,
	ErrUnacceptableUserTypeInAllOfRule:  `Unacceptable type. The "%s" type in the "allOf" rule must be an object`,
	ErrConflictAdditionalProperties:     `Conflicting value in AdditionalProperties rules when inheriting from allOf`,

	// rule loader
	ErrLoader:                           "Loader error", // error somewhere in the loader code
	ErrIncorrectRuleValueType:           "Incorrect rule value type",
	ErrIncorrectRuleWithoutExample:      "You cannot place a RULE on line without EXAMPLE",
	ErrIncorrectRuleForSeveralNode:      "You cannot place a RULE on lines that contain more than one EXAMPLE node to which any RULES can apply. The only exception is when an object key and its value are found in one line.",
	ErrLiteralValueExpected:             "Literal value expected",
	ErrInvalidValueInEnumRule:           `An array or rule name was expected as a value for the "enum"`,
	ErrIncorrectArrayItemTypeInEnumRule: `Incorrect array item type in "enum". Only literals are allowed.`,
	ErrUnacceptableValueInAllOfRule:     `Incorrect value in "allOf" rule. A type name, or list of type names, is expected.`,
	ErrTypeNameNotFoundInAllOfRule:      `Type name not found in "allOf" rule`,
	ErrDuplicationInEnumRule:            `%s value duplicates in "enum"`,

	// "or" rule loader
	ErrArrayWasExpectedInOrRule:       `An array was expected as a value for the "or" rule`,
	ErrEmptyArrayInOrRule:             `Empty array in "or" rule`,
	ErrOneElementInArrayInOrRule:      `Array rule "or" must have at least two elements`,
	ErrIncorrectArrayItemTypeInOrRule: `Incorrect array item type in "or" rule`,
	ErrEmptyRuleSet:                   `Empty rule set`,

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
	ErrUnknownAdditionalPropertiesTypes:          "Unknown type of additionalProperties (%s)",
	ErrUnexpectedConstraint:                      "The %q constraint can't be used for the %q type",

	// checker
	ErrChecker:                               `Checker error`,
	ErrElementNotFoundInArray:                `Element not found in schema array node`,
	ErrIncorrectConstraintValueForEmptyArray: `Incorrect constraint value for empty array`,

	// link checker
	ErrIncorrectUserType: `Incorrect type of user type`,
	ErrTypeNotFound:      `Type "%s" not found`,
	ErrImpossibleToDetermineTheJsonTypeDueToRecursion: `It is impossible to determine the json type due to recursion of type "%s"`,

	// sdk
	ErrEmptyType:                          `Type "%s" must not be empty`,
	ErrUnnecessaryLexemeAfterTheEndOfEnum: `An unnecessary non-space character after the end of the enum`,
	ErrRegexUnexpectedStart:               "Regex should starts with '/' character, but found %s",
	ErrRegexUnexpectedEnd:                 "Regex should ends with '/' character, but found %s",

	// enum
	ErrEnumArrayExpected:  `An array was expected as a value for the "enum"`,
	ErrEnumIsHoldRuleName: "Can't append specific value to enum initialized with rule name",
	ErrEnumRuleNotFound:   "Enum rule %q not found",
	ErrNotAnEnumRule:      "Rule %q not an Enum",
}

func (c ErrorCode) Code() ErrorCode {
	return c
}

func (c ErrorCode) Itoa() string {
	return strconv.Itoa(int(c))
}

func (c ErrorCode) Error() string {
	if format, ok := errorFormat[c]; ok {
		cnt := strings.Count(format, "%s")
		cnt += strings.Count(format, "%q")
		if cnt == 0 {
			return format
		} else {
			panic("Not enough data to generate an error message from template: " + format)
		}
	}
	panic("Unknown error code")
}
