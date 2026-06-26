package compose

import "gopkg.in/yaml.v3"

// Config represents the top-level docker-compose.yml structure.
type Config struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

// Service represents a single service definition.
type Service struct {
	Image       string            `yaml:"image"`
	Build       interface{}       `yaml:"build"`
	DependsOn   DependsOn         `yaml:"depends_on"`
	Healthcheck *Healthcheck      `yaml:"healthcheck"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
}

// DependsOn supports both forms of the depends_on key:
//
//	depends_on: [db, redis]           ← list form
//	depends_on:
//	  db:
//	    condition: service_healthy    ← map form
type DependsOn map[string]DependsOnCondition

// UnmarshalYAML handles both the list and map forms of depends_on.
func (d *DependsOn) UnmarshalYAML(value *yaml.Node) error {
	*d = make(DependsOn)
	switch value.Kind {
	case yaml.SequenceNode:
		var names []string
		if err := value.Decode(&names); err != nil {
			return err
		}
		for _, name := range names {
			(*d)[name] = DependsOnCondition{Condition: "service_started"}
		}
	case yaml.MappingNode:
		type plain map[string]DependsOnCondition
		var m plain
		if err := value.Decode(&m); err != nil {
			return err
		}
		*d = DependsOn(m)
	}
	return nil
}

// DependsOnCondition holds the optional condition for a single dependency.
type DependsOnCondition struct {
	Condition string `yaml:"condition"`
}

// Healthcheck mirrors the docker-compose healthcheck block.
type Healthcheck struct {
	Test     []string `yaml:"test"`
	Interval string   `yaml:"interval"`
	Timeout  string   `yaml:"timeout"`
	Retries  int      `yaml:"retries"`
}
