package parser

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func findDifferentFields(firstConf Config, secondConf Config) map[string]string {
	field2diff := make(map[string]string)

	// Compare logger fields
	if firstConf.Logger.Level != secondConf.Logger.Level {
		field2diff["Logger.Level"] = fmt.Sprintf("%d != %d", firstConf.Logger.Level, secondConf.Logger.Level)
	}
	if firstConf.Logger.Dest != secondConf.Logger.Dest {
		field2diff["Logger.Dest"] = fmt.Sprintf("%s != %s", firstConf.Logger.Dest, secondConf.Logger.Dest)
	}

	// Compare netbox fields
	if firstConf.Netbox.ApiToken != secondConf.Netbox.ApiToken {
		field2diff["Netbox.ApiToken"] = fmt.Sprintf("%s != %s", firstConf.Netbox.ApiToken, secondConf.Netbox.ApiToken)
	}
	if firstConf.Netbox.Hostname != secondConf.Netbox.Hostname {
		field2diff["Netbox.Hostname"] = fmt.Sprintf("%s != %s", firstConf.Netbox.Hostname, secondConf.Netbox.Hostname)
	}
	if firstConf.Netbox.HTTPScheme != secondConf.Netbox.HTTPScheme {
		field2diff["Netbox.HTTPScheme"] = fmt.Sprintf("%s != %s", firstConf.Netbox.HTTPScheme, secondConf.Netbox.HTTPScheme)
	}
	if firstConf.Netbox.Port != secondConf.Netbox.Port {
		field2diff["Netbox.Port"] = fmt.Sprintf("%d != %d", firstConf.Netbox.Port, secondConf.Netbox.Port)
	}
	if firstConf.Netbox.ValidateCert != secondConf.Netbox.ValidateCert {
		field2diff["Netbox.ValidateCert"] = fmt.Sprintf("%t != %t", firstConf.Netbox.ValidateCert, secondConf.Netbox.ValidateCert)
	}
	if firstConf.Netbox.Timeout != secondConf.Netbox.Timeout {
		field2diff["Netbox.Timeout"] = fmt.Sprintf("%d != %d", firstConf.Netbox.Timeout, secondConf.Netbox.Timeout)
	}
	if firstConf.Netbox.Tag != secondConf.Netbox.Tag {
		field2diff["Netbox.Tag"] = fmt.Sprintf("%s != %s", firstConf.Netbox.Tag, secondConf.Netbox.Tag)
	}
	if firstConf.Netbox.TagColor != secondConf.Netbox.TagColor {
		field2diff["Netbox.TagColor"] = fmt.Sprintf("%s != %s", firstConf.Netbox.TagColor, secondConf.Netbox.TagColor)
	}

	// Compare sources list
	if len(firstConf.Sources) != len(secondConf.Sources) {
		field2diff["Sources"] = fmt.Sprintf("len(%d) != len(%d)", len(firstConf.Sources), len(secondConf.Sources))
	}
	for i := 0; i < len(firstConf.Sources); i++ {
		firstSource := firstConf.Sources[i]
		secondSource := secondConf.Sources[i]
		if firstSource.Name != secondSource.Name {
			field2diff[fmt.Sprintf("Sources[%d].Name", i)] = fmt.Sprintf("%s != %s", firstSource.Name, secondSource.Name)
		}
		if firstSource.Type != secondSource.Type {
			field2diff[fmt.Sprintf("Sources[%d].Type", i)] = fmt.Sprintf("%s != %s", firstSource.Type, secondSource.Type)
		}
		if firstSource.HTTPScheme != secondSource.HTTPScheme {
			field2diff[fmt.Sprintf("Sources[%d].HTTPScheme", i)] = fmt.Sprintf("%s != %s", firstSource.HTTPScheme, secondSource.HTTPScheme)
		}
		if firstSource.Hostname != secondSource.Hostname {
			field2diff[fmt.Sprintf("Sources[%d].Hostname", i)] = fmt.Sprintf("%s != %s", firstSource.Hostname, secondSource.Hostname)
		}
		if firstSource.Port != secondSource.Port {
			field2diff[fmt.Sprintf("Sources[%d].Port", i)] = fmt.Sprintf("%d != %d", firstSource.Port, secondSource.Port)
		}
		if firstSource.Username != secondSource.Username {
			field2diff[fmt.Sprintf("Sources[%d].Username", i)] = fmt.Sprintf("%s != %s", firstSource.Username, secondSource.Username)
		}
		if firstSource.Password != secondSource.Password {
			field2diff[fmt.Sprintf("Sources[%d].Password", i)] = fmt.Sprintf("%s != %s", firstSource.Password, secondSource.Password)
		}
		if firstSource.Tag != secondSource.Tag {
			field2diff[fmt.Sprintf("Sources[%d].Tag", i)] = fmt.Sprintf("%s != %s", firstSource.Tag, secondSource.Tag)
		}
		if firstSource.TagColor != secondSource.TagColor {
			field2diff[fmt.Sprintf("Sources[%d].TagColor", i)] = fmt.Sprintf("%s != %s", firstSource.TagColor, secondSource.TagColor)
		}
		if len(firstSource.PermittedSubnets) != len(secondSource.PermittedSubnets) {
			field2diff[fmt.Sprintf("Sources[%d].PermittedSubnets", i)] = fmt.Sprintf("len(%d) != len(%d)", len(firstSource.PermittedSubnets), len(secondSource.PermittedSubnets))
		}
		for j := 0; j < len(firstSource.PermittedSubnets); j++ {
			if firstSource.PermittedSubnets[j] != secondSource.PermittedSubnets[j] {
				field2diff[fmt.Sprintf("Sources[%d].PermittedSubnets[%d]", i, j)] = fmt.Sprintf("%s != %s", firstSource.PermittedSubnets[j], secondSource.PermittedSubnets[j])
			}
		}
		if len(firstSource.ClusterSiteRelations) != len(secondSource.ClusterSiteRelations) {
			field2diff[fmt.Sprintf("Sources[%d].ClusterSiteRelations", i)] = fmt.Sprintf("len(%d) != len(%d)", len(firstSource.ClusterSiteRelations), len(secondSource.ClusterSiteRelations))
		}
		for j := 0; j < len(firstSource.ClusterSiteRelations); j++ {
			if firstSource.ClusterSiteRelations[j] != secondSource.ClusterSiteRelations[j] {
				field2diff[fmt.Sprintf("Sources[%d].ClusterSiteRelations[%d]", i, j)] = fmt.Sprintf("%s != %s", firstSource.ClusterSiteRelations[j], secondSource.ClusterSiteRelations[j])
			}
		}
		if len(firstSource.HostSiteRelations) != len(secondSource.HostSiteRelations) {
			field2diff[fmt.Sprintf("Sources[%d].HostSiteRelations", i)] = fmt.Sprintf("len(%d) != len(%d)", len(firstSource.HostSiteRelations), len(secondSource.HostSiteRelations))
		}
		for j := 0; j < len(firstSource.HostSiteRelations); j++ {
			if firstSource.HostSiteRelations[j] != secondSource.HostSiteRelations[j] {
				field2diff[fmt.Sprintf("Sources[%d].HostSiteRelations[%d]", i, j)] = fmt.Sprintf("%s != %s", firstSource.HostSiteRelations[j], secondSource.HostSiteRelations[j])
			}
		}
		if len(firstSource.ClusterTenantRelations) != len(secondSource.ClusterTenantRelations) {
			field2diff[fmt.Sprintf("Sources[%d].ClusterTenantRelations", i)] = fmt.Sprintf("len(%d) != len(%d)", len(firstSource.ClusterTenantRelations), len(secondSource.ClusterTenantRelations))
		}
		for j := 0; j < len(firstSource.ClusterTenantRelations); j++ {
			if firstSource.ClusterTenantRelations[j] != secondSource.ClusterTenantRelations[j] {
				field2diff[fmt.Sprintf("Sources[%d].ClusterTenantRelations[%d]", i, j)] = fmt.Sprintf("%s != %s", firstSource.ClusterTenantRelations[j], secondSource.ClusterTenantRelations[j])
			}
		}
		if len(firstSource.HostTenantRelations) != len(secondSource.HostTenantRelations) {
			field2diff[fmt.Sprintf("Sources[%d].HostTenantRelations", i)] = fmt.Sprintf("len(%d) != len(%d)", len(firstSource.HostTenantRelations), len(secondSource.HostTenantRelations))
		}
		for j := 0; j < len(firstSource.HostTenantRelations); j++ {
			if firstSource.HostTenantRelations[j] != secondSource.HostTenantRelations[j] {
				field2diff[fmt.Sprintf("Sources[%d].HostTenantRelations[%d]", i, j)] = fmt.Sprintf("%s != %s", firstSource.HostTenantRelations[j], secondSource.HostTenantRelations[j])
			}
		}
		if len(firstSource.VmTenantRelations) != len(secondSource.VmTenantRelations) {
			field2diff[fmt.Sprintf("Sources[%d].VmTenantRelations", i)] = fmt.Sprintf("len(%d) != len(%d)", len(firstSource.VmTenantRelations), len(secondSource.VmTenantRelations))
		}
		for j := 0; j < len(firstSource.VmTenantRelations); j++ {
			if firstSource.VmTenantRelations[j] != secondSource.VmTenantRelations[j] {
				field2diff[fmt.Sprintf("Sources[%d].VmTenantRelations[%d]", i, j)] = fmt.Sprintf("%s != %s", firstSource.VmTenantRelations[j], secondSource.VmTenantRelations[j])
			}
		}
	}
	return field2diff
}

func TestValidConfig(t *testing.T) {
	filename := filepath.Join("testdata", "valid_config.yaml")
	want := &Config{
		Logger: &LoggerConfig{
			Level: 2,
			Dest:  "test",
		},
		Netbox: &NetboxConfig{
			ApiToken:      "netbox-token",
			Hostname:      "netbox.example.com",
			HTTPScheme:    "https",
			Port:          666,
			ValidateCert:  false,         // Default
			Timeout:       30,            // Default
			Tag:           "netbox-ssot", // Default
			TagColor:      "00add8",      // Default
			RemoveOrphans: true,          // Default
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
				PermittedSubnets: []string{
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
				PermittedSubnets: []string{
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
				VmTenantRelations: []string{
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
		field2diff := findDifferentFields(*got, *want)
		t.Errorf("got = %v\nwant %v\nDifferent fields: %v", got, want, field2diff)
	}
}
