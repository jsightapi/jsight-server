package constraint

type Type int

const (
	MinLengthConstraintType Type = iota
	MaxLengthConstraintType
	MinConstraintType
	MaxConstraintType
	ExclusiveMinimumConstraintType
	ExclusiveMaximumConstraintType
	PrecisionConstraintType
	TypeConstraintType
	TypesListConstraintType
	OptionalConstraintType
	OrConstraintType
	RequiredKeysConstraintType
	EmailConstraintType
	MinItemsConstraintType
	MaxItemsConstraintType
	EnumConstraintType
	AdditionalPropertiesConstraintType
	AllOfConstraintType
	AnyConstraintType
	NullableConstraintType
	RegexConstraintType
	UriConstraintType
	DateConstraintType
	DateTimeConstraintType
	UuidConstraintType
	ConstType
)

func (t Type) String() string { //nolint:gocyclo // For now it's okay.
	switch t {
	case MinLengthConstraintType:
		return "minLength"
	case MaxLengthConstraintType:
		return "maxLength"
	case MinConstraintType:
		return "min"
	case MaxConstraintType:
		return "max"
	case ExclusiveMinimumConstraintType:
		return "exclusiveMinimum"
	case ExclusiveMaximumConstraintType:
		return "exclusiveMaximum"
	case PrecisionConstraintType:
		return "precision"
	case TypeConstraintType:
		return "type"
	case TypesListConstraintType:
		return "types"
	case OptionalConstraintType:
		return "optional"
	case OrConstraintType:
		return "or"
	case RequiredKeysConstraintType:
		return "required-keys"
	case EmailConstraintType:
		return "email"
	case MinItemsConstraintType:
		return "minItems"
	case MaxItemsConstraintType:
		return "maxItems"
	case EnumConstraintType:
		return "enum"
	case AdditionalPropertiesConstraintType:
		return "additionalProperties"
	case AllOfConstraintType:
		return "allOf"
	case AnyConstraintType:
		return "any"
	case NullableConstraintType:
		return "nullable"
	case RegexConstraintType:
		return "regex"
	case UriConstraintType:
		return "uri"
	case DateConstraintType:
		return "date"
	case DateTimeConstraintType:
		return "datetime"
	case UuidConstraintType:
		return "uuid"
	case ConstType:
		return "const"
	}
	panic("Unknown constraint type")
}
