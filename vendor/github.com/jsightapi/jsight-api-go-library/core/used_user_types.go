package core

import (
	"fmt"

	jschema "github.com/jsightapi/jsight-schema-go-library"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

func fetchUsedUserTypes(ut jschema.Schema, userTypes *catalog.UserSchemas) ([]string, error) {
	f := &usedUserTypeFetcher{
		userTypes:        userTypes,
		alreadyProcessed: map[string]struct{}{},
	}
	if err := f.fetch(ut); err != nil {
		return nil, err
	}
	return f.usedUserTypes, nil
}

type usedUserTypeFetcher struct {
	userTypes        *catalog.UserSchemas
	alreadyProcessed map[string]struct{}
	usedUserTypes    []string
}

func (f *usedUserTypeFetcher) fetch(ut jschema.Schema) error {
	if ut == nil {
		return nil
	}

	tt, err := ut.UsedUserTypes()
	if err != nil {
		return err
	}

	if len(tt) == 0 {
		return nil
	}

	for _, t := range tt {
		if _, ok := f.alreadyProcessed[t]; ok {
			continue
		}

		f.alreadyProcessed[t] = struct{}{}
		f.usedUserTypes = append(f.usedUserTypes, t)
		if err := f.fetch(f.userTypes.GetValue(t)); err != nil {
			return fmt.Errorf("process type %q: %w", t, err)
		}
	}
	return nil
}
