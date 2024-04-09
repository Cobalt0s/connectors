package msdsales

import (
	"errors"
	"strings"
)

var ErrInvalidParentRef = errors.New("referenced parent type doesn't exist for schema")

// EntitySet organised by Entity.Name.
type EntitySet map[string]*Entity

func NewEntitySet() EntitySet {
	return make(map[string]*Entity)
}

func (s EntitySet) GetNames() []string {
	names := make([]string, len(s))
	i := 0

	for name := range s {
		names[i] = name
		i++
	}

	return names
}

func (s EntitySet) GetOrCreate(name, parentName string) *Entity {
	if _, ok := s[name]; !ok {
		s[name] = &Entity{
			Name:       name,
			properties: make([]string, 0),
			parentName: parentName,
		}
	}

	return s[name]
}

type Entity struct {
	Name       string
	properties []string
	parentName string
	parent     *Entity
}

func (e *Entity) AddProperty(property string) {
	e.properties = append(e.properties, property)
}

// GetRawParentName parents that are defined under schema are prefixed with its alias
// this strips the prefix.
func (e *Entity) GetRawParentName(schemaAlias string) string {
	// strip prefix if we are using local schema
	name, _ := strings.CutPrefix(e.parentName, schemaAlias+".")
	return name
}

func (s EntitySet) MatchParentsWithChildren(schemaAlias string) error {
	for _, entity := range s {
		name := entity.GetRawParentName(schemaAlias)
		// if it has parent then match
		if len(name) != 0 {
			parent, ok := s[name]
			if !ok {
				return ErrInvalidParentRef
			}

			entity.parent = parent
		}
	}

	return nil
}

// GetAllProperties recursive function that includes inherited fields from parents.
func (e *Entity) GetAllProperties() []string {
	if e.parent == nil {
		// this is root
		return e.properties
	}

	parentProperties := e.parent.GetAllProperties()

	return append(e.properties, parentProperties...)
}
