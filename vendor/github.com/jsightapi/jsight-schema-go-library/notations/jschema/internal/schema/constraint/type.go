package constraint

// Type available constraint types.
// gen:Stringer t Unknown constraint type
type Type int

const (
	MinLengthConstraintType            Type = iota // minLength
	MaxLengthConstraintType                        // maxLength
	MinConstraintType                              // min
	MaxConstraintType                              // max
	ExclusiveMinimumConstraintType                 // exclusiveMinimum
	ExclusiveMaximumConstraintType                 // exclusiveMaximum
	PrecisionConstraintType                        // precision
	TypeConstraintType                             // type
	TypesListConstraintType                        // types
	OptionalConstraintType                         // optional
	OrConstraintType                               // or
	RequiredKeysConstraintType                     // required-keys
	EmailConstraintType                            // email
	MinItemsConstraintType                         // minItems
	MaxItemsConstraintType                         // maxItems
	EnumConstraintType                             // enum
	AdditionalPropertiesConstraintType             // additionalProperties
	AllOfConstraintType                            // allOf
	AnyConstraintType                              // any
	NullableConstraintType                         // nullable
	RegexConstraintType                            // regex
	UriConstraintType                              // uri
	DateConstraintType                             // date
	DateTimeConstraintType                         // datetime
	UuidConstraintType                             // uuid
	ConstConstraintType                            // const
)
