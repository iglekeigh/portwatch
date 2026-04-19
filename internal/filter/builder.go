package filter

import "github.com/yourorg/portwatch/internal/scanner"

// Config holds raw filter configuration.
type Config struct {
	IncludePorts string `yaml:"include_ports"`
	ExcludePorts string `yaml:"exclude_ports"`
}

// FromConfig builds a Filter from a Config, parsing port ranges.
func FromConfig(cfg Config) (*Filter, error) {
	var rules []Rule

	if cfg.IncludePorts != "" {
		ports, err := scanner.ParsePortRange(cfg.IncludePorts)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Ports: ports, Exclude: false})
	}

	if cfg.ExcludePorts != "" {
		ports, err := scanner.ParsePortRange(cfg.ExcludePorts)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Ports: ports, Exclude: true})
	}

	return New(rules), nil
}
