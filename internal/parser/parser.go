package parser

import (
	"errors"
	"fmt"
	"os"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  *LoggerConfig  `yaml:"logger"`
	Netbox  *NetboxConfig  `yaml:"netbox"`
	Sources []SourceConfig `yaml:"source"`
}

type LoggerConfig struct {
	Level int    `yaml:"level"`
	Dest  string `yaml:"dest"`
}

func (l LoggerConfig) String() string {
	if l.Dest == "" {
		return fmt.Sprintf("LoggerConfig{Level: %d, Dest: stdout}", l.Level)
	}
	return fmt.Sprintf("LoggerConfig{Level: %d, Dest: %s}", l.Level, l.Dest)
}

type HTTPScheme string

const (
	HTTP  HTTPScheme = "http"
	HTTPS HTTPScheme = "https"
)

type NetboxConfig struct {
	APIToken string `yaml:"apiToken"`
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	// Can be http or https (default)
	HTTPScheme     HTTPScheme `yaml:"httpScheme"`
	ValidateCert   bool       `yaml:"validateCert"`
	Timeout        int        `yaml:"timeout"`
	Tag            string     `yaml:"tag"`
	TagColor       string     `yaml:"tagColor"`
	RemoveOrphans  bool       `yaml:"removeOrphans"`
	SourcePriority []string   `yaml:"sourcePriority"`
}

func (n NetboxConfig) String() string {
	return fmt.Sprintf("NetboxConfig{ApiToken: %s, Hostname: %s, Port: %d, HTTPScheme: %s, ValidateCert: %t, Timeout: %d, Tag: %s, TagColor: %s, RemoveOrphans: %t}", n.APIToken, n.Hostname, n.Port, n.HTTPScheme, n.ValidateCert, n.Timeout, n.Tag, n.TagColor, n.RemoveOrphans)
}

type SourceConfig struct {
	Name             string               `yaml:"name"`
	Type             constants.SourceType `yaml:"type"`
	HTTPScheme       HTTPScheme           `yaml:"httpScheme"`
	Hostname         string               `yaml:"hostname"`
	Port             int                  `yaml:"port"`
	Username         string               `yaml:"username"`
	Password         string               `yaml:"password"`
	PermittedSubnets []string             `yaml:"permittedSubnets"`
	ValidateCert     bool                 `yaml:"validateCert"`
	Tag              string               `yaml:"tag"`
	TagColor         string               `yaml:"tagColor"`

	// Relations
	HostSiteRelations      []string `yaml:"hostSiteRelations"`
	ClusterSiteRelations   []string `yaml:"clusterSiteRelations"`
	ClusterTenantRelations []string `yaml:"clusterTenantRelations"`
	HostTenantRelations    []string `yaml:"hostTenantRelations"`
	VMTenantRelations      []string `yaml:"vmTenantRelations"`
	VlanGroupRelations     []string `yaml:"vlanGroupRelations"`
	VlanTenantRelations    []string `yaml:"vlanTenantRelations"`

	// Vmware specific relations
	CustomFieldMappings []string `yaml:"customFieldMappings"`
}

func (s SourceConfig) String() string {
	return fmt.Sprintf("SourceConfig{Name: %s, Type: %s, HTTPScheme: %s, Hostname: %s, Port: %d, Username: %s, Password: %s, PermittedSubnets: %v, ValidateCert: %t, Tag: %s, TagColor: %s, HostSiteRelations: %v, ClusterSiteRelations: %v, clusterTenantRelations: %v, HostTenantRelations: %v, VmTenantRelations %v, VlanGroupRelations: %v, VlanTenantRelations: %v}", s.Name, s.Type, s.HTTPScheme, s.Hostname, s.Port, s.Username, s.Password, s.PermittedSubnets, s.ValidateCert, s.Tag, s.TagColor, s.HostSiteRelations, s.ClusterSiteRelations, s.ClusterTenantRelations, s.HostTenantRelations, s.VMTenantRelations, s.VlanGroupRelations, s.VlanTenantRelations)
}

// Validates the user's config for limits and required fields.
func validateConfig(config *Config) error {
	err := validateLoggerConfig(config)
	if err != nil {
		return err
	}

	err = validateNetboxConfig(config)
	if err != nil {
		return err
	}

	err = validateSourceConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func validateLoggerConfig(config *Config) error {
	if config.Logger.Level < 0 || config.Logger.Level > 3 {
		return errors.New("logger.level: must be between 0 and 3")
	}
	return nil
}

// Function that validates NetboxConfig.
func validateNetboxConfig(config *Config) error {
	// Validate Netbox config
	if config.Netbox.APIToken == "" {
		return errors.New("netbox.apiToken: cannot be empty")
	}
	if config.Netbox.HTTPScheme != HTTP && config.Netbox.HTTPScheme != HTTPS {
		return errors.New("netbox.httpScheme: must be either http or https. Is " + string(config.Netbox.HTTPScheme))
	}
	if config.Netbox.Hostname == "" {
		return errors.New("netbox.hostname: cannot be empty")
	}
	if config.Netbox.Port < 0 || config.Netbox.Port > 65535 {
		return errors.New("netbox.port: must be between 0 and 65535. Is " + fmt.Sprintf("%d", config.Netbox.Port))
	}
	if config.Netbox.Timeout < 0 {
		return errors.New("netbox.timeout: cannot be negative")
	}
	if config.Netbox.Tag == "" {
		config.Netbox.Tag = "netbox-ssot"
	}
	if config.Netbox.TagColor == "" {
		config.Netbox.TagColor = "00add8"
	} else {
		// Ensure that TagColor is a string of 6 hexadecimal characters
		if len(config.Netbox.TagColor) != len("ffffff") {
			return errors.New("netbox.tagColor: must be a string of 6 hexadecimal characters")
		}
		for _, c := range config.Netbox.TagColor {
			if c < '0' || c > '9' && c < 'a' || c > 'f' {
				return errors.New("netbox.tagColor: must be a string of 6 lowercase hexadecimal characters")
			}
		}
	}
	if len(config.Netbox.SourcePriority) > 0 {
		if len(config.Netbox.SourcePriority) != len(config.Sources) {
			return fmt.Errorf("netbox.sourcePriority: len(config.Netbox.SourcePriority != len(config.Sources))")
		}
		for _, sourceName := range config.Netbox.SourcePriority {
			contains := false
			for _, source := range config.Sources {
				if source.Name == sourceName {
					contains = true
					break
				}
			}
			if !contains {
				return fmt.Errorf("netbox.sourcePriority: source[%s] doesn't exist in sources array", sourceName)
			}
		}
	}
	return nil
}

func validateSourceConfig(config *Config) error {
	// Validate Sources
	for i := range config.Sources {
		externalSource := &config.Sources[i]
		externalSourceStr := "source[" + externalSource.Name + "]"
		if externalSource.Name == "" {
			return fmt.Errorf("%s: name cannot be empty", externalSourceStr)
		}
		if externalSource.HTTPScheme == "" {
			externalSource.HTTPScheme = "https"
		} else if externalSource.HTTPScheme != HTTP && externalSource.HTTPScheme != HTTPS {
			return fmt.Errorf("%s.httpScheme must be either http or https. Is %s", externalSourceStr, string(externalSource.HTTPScheme))
		}
		if externalSource.Hostname == "" {
			return fmt.Errorf("%s: hostname cannot be empty", externalSourceStr)
		}
		if externalSource.Port == 0 {
			externalSource.Port = 443
		} else if externalSource.Port < 0 || externalSource.Port > 65535 {
			return fmt.Errorf("%s: port must be between 0 and 65535. Is %d", externalSourceStr, externalSource.Port)
		}
		if externalSource.Username == "" {
			return fmt.Errorf("%s: username cannot be empty", externalSourceStr)
		}
		if externalSource.Password == "" {
			return fmt.Errorf("%s: password cannot be empty", externalSourceStr)
		}
		if externalSource.Tag == "" {
			externalSource.Tag = fmt.Sprintf("Source: %s", externalSource.Name)
		}
		if externalSource.TagColor == "" {
			externalSource.TagColor = constants.DefaultSourceToTagColorMap[externalSource.Type]
		}
		switch externalSource.Type {
		case constants.Ovirt:
		case constants.Vmware:
		case constants.Dnac:
		default:
			return fmt.Errorf("%s.type is not valid", externalSourceStr)
		}
		err := validateSourceConfigRelations(externalSource, externalSourceStr)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateSourceConfigRelations(externalSource *SourceConfig, externalSourceStr string) error {
	if len(externalSource.HostSiteRelations) > 0 {
		err := utils.ValidateRegexRelations(externalSource.HostSiteRelations)
		if err != nil {
			return fmt.Errorf("%s.hostSiteRelations: %s", externalSourceStr, err)
		}
	}
	if len(externalSource.ClusterSiteRelations) > 0 {
		err := utils.ValidateRegexRelations(externalSource.ClusterSiteRelations)
		if err != nil {
			return fmt.Errorf("%s.clusterSiteRelations: %s", externalSourceStr, err)
		}
	}
	if len(externalSource.ClusterTenantRelations) > 0 {
		err := utils.ValidateRegexRelations(externalSource.ClusterTenantRelations)
		if err != nil {
			return fmt.Errorf("%s.clusterTenantRelations: %s", externalSourceStr, err)
		}
	}
	if len(externalSource.HostTenantRelations) > 0 {
		err := utils.ValidateRegexRelations(externalSource.HostTenantRelations)
		if err != nil {
			return fmt.Errorf("%s.hostTenantRelations: %s", externalSourceStr, err)
		}
	}
	if len(externalSource.VMTenantRelations) > 0 {
		err := utils.ValidateRegexRelations(externalSource.VMTenantRelations)
		if err != nil {
			return fmt.Errorf("%s.vmTenantRelations: %s", externalSourceStr, err)
		}
	}
	if len(externalSource.VlanGroupRelations) > 0 {
		err := utils.ValidateRegexRelations(externalSource.VlanGroupRelations)
		if err != nil {
			return fmt.Errorf("%s.vlanGroupRelations: %v", externalSourceStr, err)
		}
	}
	if len(externalSource.VlanTenantRelations) > 0 {
		err := utils.ValidateRegexRelations((externalSource.VlanTenantRelations))
		if err != nil {
			return fmt.Errorf("%s.vlanTenantRelations: %v", externalSourceStr, err)
		}
	}
	return nil
}

func ParseConfig(filename string) (*Config, error) {
	// First we read the config file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Define Config with default values
	config := &Config{
		Logger: &LoggerConfig{
			Level: 1,
			Dest:  "",
		},
		Netbox: &NetboxConfig{
			HTTPScheme:    "https",
			Port:          constants.HTTPSDefaultPort,
			Timeout:       constants.DefaultAPITimeout,
			RemoveOrphans: true,
		},
		Sources: []SourceConfig{},
	}

	// Parse the config file into a Config struct
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	// Validate the config for limits and required fields
	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
