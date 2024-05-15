package jschema

import (
	"strings"

	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

func collectUserTypes(node ischema.Node) []string {
	c := &userTypesCollector{
		alreadyProcessed: map[string]struct{}{},
	}
	c.collect(node)
	return c.userTypes
}

type userTypesCollector struct {
	alreadyProcessed map[string]struct{}
	userTypes        []string
}

func (c *userTypesCollector) collect(node ischema.Node) {
	c.collectUserTypesFromTypesListConstraint(node)
	c.collectUserTypesFromTypeConstraint(node)
	c.collectUserTypesFromAllOfConstraint(node)

	switch n := node.(type) {
	case *ischema.ObjectNode:
		c.collectUserTypesFromAdditionalPropertiesOfConstraint(node)
		c.collectUserTypesObjectNode(n)

	case *ischema.ArrayNode:
		for _, child := range n.Children() {
			c.collect(child)
		}

	case *ischema.MixedValueNode:
		for _, ut := range strings.Split(n.Value().String(), "|") {
			s := strings.TrimSpace(ut)
			if s[0] == '@' {
				c.addType(s)
			}
		}
	}
}

func (c *userTypesCollector) collectUserTypesFromTypesListConstraint(node ischema.Node) {
	for _, name := range UserTypeNamesFromTypesListConstraint(node) {
		c.addType(name)
	}
}

func (c *userTypesCollector) collectUserTypesFromTypeConstraint(node ischema.Node) {
	for _, name := range UserTypeNamesFromTypeConstraint(node) {
		c.addType(name)
	}
}

func (c *userTypesCollector) collectUserTypesFromAllOfConstraint(node ischema.Node) {
	cnstr := node.Constraint(constraint.AllOfConstraintType)
	if c == nil {
		return
	}

	allOf, ok := cnstr.(*constraint.AllOf)
	if !ok {
		return
	}

	for _, name := range allOf.SchemaNames() {
		if name[0] == '@' {
			c.addType(name)
		}
	}
}

func (c *userTypesCollector) collectUserTypesFromAdditionalPropertiesOfConstraint(node ischema.Node) {
	cnstr := node.Constraint(constraint.AdditionalPropertiesConstraintType)
	if c == nil {
		return
	}

	ap, ok := cnstr.(*constraint.AdditionalProperties)
	if !ok {
		return
	}

	if ap.Mode() == constraint.AdditionalPropertiesMustBeUserType {
		c.addType(ap.TypeName().String())
	}
}

func (c *userTypesCollector) collectUserTypesObjectNode(node *ischema.ObjectNode) {
	for _, v := range node.Keys().Data {
		k := v.Key

		if v.IsShortcut {
			if k[0] == '@' {
				c.addType(k)
			}
		}

		child, ok := node.Child(k, v.IsShortcut)
		if ok {
			c.collect(child)
		}
	}
}

func (c *userTypesCollector) addType(n string) {
	if _, ok := c.alreadyProcessed[n]; ok {
		return
	}
	c.alreadyProcessed[n] = struct{}{}
	c.userTypes = append(c.userTypes, n)
}
