package parser

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestValidonfig(t *testing.T) {
	filename := filepath.Join("testdata", "valid_config.yaml")
	want := &Config{
		Logger: &LoggerConfig{
			Level: 2,
			Dest:  "test",
		},
		Netbox: &NetboxConfig{
			APIToken:      "netbox-token",
			Hostname:      "netbox.example.com",
			HTTPScheme:    "https",
			Port:          666,
			ValidateCert:  false, // Default
			Timeout:       constants.DefaultAPITimeout,
			Tag:           constants.DefaultSourceName, // Default
			TagColor:      "00add8",                    // Default
			RemoveOrphans: true,                        // Default
		},
		Sources: []SourceConfig{
			{
				Name:       "testolvm",
				Type:       "ovirt",
				HTTPScheme: "http",
				Port:       443,
				Hostname:   "testolvm.example.com",
				Username:   "admin@internal",
				Password:   "adminpass",
				IgnoredSubnets: []string{
					"172.16.0.0/12",
					"192.168.0.0/16",
					"fd00::/8",
				},
				ValidateCert: true,
				Tag:          "testing",
				TagColor:     "ff0000",
			},
			{
				Name:       "prodolvm",
				Type:       "ovirt",
				Port:       80,
				HTTPScheme: "https",
				Hostname:   "ovirt.example.com",
				Username:   "admin",
				Password:   "adminpass",
				IgnoredSubnets: []string{
					"172.16.0.0/12",
				},
				ValidateCert: false,
				Tag:          "Source: prodolvm", // Default
				TagColor:     "aa1409",           // Default
				ClusterSiteRelations: []string{
					"Cluster_NYC = New York",
					"Cluster_FFM.* = Frankfurt",
					"Datacenter_BERLIN/* = Berlin",
				},
				HostSiteRelations: []string{
					".* = Berlin",
				},
				ClusterTenantRelations: []string{
					".*Stark = Stark Industries",
					".* = Default",
				},
				HostTenantRelations: []string{
					".*Health = Health Department",
					".* = Default",
				},
				VMTenantRelations: []string{
					".*Health = Health Department",
					".* = Default",
				},
			},
		},
	}
	got, err := ParseConfig(filename)
	if err != nil {
		t.Errorf("ParseConfig() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v\nwant %v\n", got, want)
	}
}
