package catalog

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-go-library/bytes"

	"github.com/jsightapi/jsight-api-go-library/directive"
)

const JDocExchangeVersion = "2.0.0"

// Catalog is the main result of processing .jst file, represents all accumulated Api data more closely to API terms.
// Has no direct connection to Directives but operates with models with similar structure, but in API terms.
// Has flat URL structure, each URL is unique and indicates a separate Resource.
// Each Resource has Requests — available http methods and their properties.
// Each Resource or Request has it its own complex properties, which can be added through dedicated methods.
type Catalog struct {
	rawUserTypes *directive.Directives

	// Info represent "INFO" directive structure.
	Info *Info

	// Servers contains information about servers.
	Servers *Servers

	// UserTypes contains information about defined user types.
	UserTypes *UserTypes

	// UserEnums contains information about defined user enums.
	UserEnums *UserRules

	Interactions *Interactions

	Tags *Tags

	JSightVersion string
}

var _ json.Marshaler = &Catalog{}

func (c *Catalog) ToJson() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Catalog) ToJsonIndent() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func NewCatalog() *Catalog {
	return &Catalog{
		rawUserTypes: &directive.Directives{},
		Servers:      &Servers{},
		UserTypes:    &UserTypes{},
		UserEnums:    &UserRules{},
		Interactions: &Interactions{},
		Tags:         &Tags{},
	}
}

func (c *Catalog) AddRawUserType(d *directive.Directive) {
	c.rawUserTypes.Set(d.Parameter("Name"), d)
}

func (c *Catalog) GetRawUserTypes() *directive.Directives {
	return c.rawUserTypes
}

func (*Catalog) Read(coords directive.Coords) bytes.Bytes {
	return coords.Read()
}

func (c *Catalog) MarshalJSON() ([]byte, error) {
	var data struct {
		Tags                *Tags         `json:"tags"`
		Info                *Info         `json:"info,omitempty"`
		Servers             *Servers      `json:"servers,omitempty"`
		UserTypes           *UserTypes    `json:"userTypes,omitempty"`
		UserEnums           *UserRules    `json:"userEnums,omitempty"`
		Interactions        *Interactions `json:"interactions"`
		JSightVersion       string        `json:"jsight"`
		JdocExchangeVersion string        `json:"jdocExchangeVersion"`
	}

	data.JdocExchangeVersion = JDocExchangeVersion
	data.JSightVersion = c.JSightVersion
	data.Info = c.Info
	if c.Servers.Len() > 0 {
		data.Servers = c.Servers
	}
	if c.UserTypes.Len() > 0 {
		data.UserTypes = c.UserTypes
	}
	if c.UserEnums.Len() > 0 {
		data.UserEnums = c.UserEnums
	}
	data.Interactions = c.Interactions
	data.Tags = c.Tags

	return json.Marshal(data)
}
