package openstack

import (
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/parser"
)

func TestResolveDomainConfig(t *testing.T) {
	tests := []struct {
		name                      string
		cfg                       *parser.SourceConfig
		expectedDomainName        string
		expectedDomainID          string
		expectedProjectDomainName string
		expectedProjectDomainID   string
	}{
		{
			name: "Only domainName - should be used for both user and project domains",
			cfg: &parser.SourceConfig{
				DomainName: "MyDomain",
			},
			expectedDomainName:        "MyDomain",
			expectedDomainID:          "",
			expectedProjectDomainName: "MyDomain",
			expectedProjectDomainID:   "",
		},
		{
			name: "Only domainID - should clear Name and be used for both domains",
			cfg: &parser.SourceConfig{
				DomainID: "abc123",
			},
			expectedDomainName:        "",
			expectedDomainID:          "abc123",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "abc123",
		},
		{
			name: "Both domainName and domainID - ID takes precedence, Name cleared",
			cfg: &parser.SourceConfig{
				DomainName: "MyDomain",
				DomainID:   "abc123",
			},
			expectedDomainName:        "",
			expectedDomainID:          "abc123",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "abc123",
		},
		{
			name: "domainName with id| prefix - should convert to domainID",
			cfg: &parser.SourceConfig{
				DomainName: "id|xyz789",
			},
			expectedDomainName:        "",
			expectedDomainID:          "xyz789",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "xyz789",
		},
		{
			name:                      "Neither domainName nor domainID - should default to 'Default'",
			cfg:                       &parser.SourceConfig{},
			expectedDomainName:        "Default",
			expectedDomainID:          "",
			expectedProjectDomainName: "Default",
			expectedProjectDomainID:   "",
		},
		{
			name: "Explicit projectDomainName - should override domainName fallback",
			cfg: &parser.SourceConfig{
				DomainName:        "UserDomain",
				ProjectDomainName: "ProjectDomain",
			},
			expectedDomainName:        "UserDomain",
			expectedDomainID:          "",
			expectedProjectDomainName: "ProjectDomain",
			expectedProjectDomainID:   "",
		},
		{
			name: "Explicit projectDomainID - should override fallback and clear projectDomainName",
			cfg: &parser.SourceConfig{
				DomainName:      "UserDomain",
				ProjectDomainID: "proj123",
			},
			expectedDomainName:        "UserDomain",
			expectedDomainID:          "",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "proj123",
		},
		{
			name: "Both projectDomainName and projectDomainID - ID precedence",
			cfg: &parser.SourceConfig{
				DomainName:        "UserDomain",
				ProjectDomainName: "ProjectDomain",
				ProjectDomainID:   "proj123",
			},
			expectedDomainName:        "UserDomain",
			expectedDomainID:          "",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "proj123",
		},
		{
			name: "Complex: domainID set, projectDomainName explicit",
			cfg: &parser.SourceConfig{
				DomainID:          "user123",
				ProjectDomainName: "ProjectDomain",
			},
			expectedDomainName:        "",
			expectedDomainID:          "user123",
			expectedProjectDomainName: "ProjectDomain",
			expectedProjectDomainID:   "", // Do not inherit user domainID when projectDomainName is explicit
		},
		{
			name: "Complex: domainName set, projectDomainID explicit",
			cfg: &parser.SourceConfig{
				DomainName:      "UserDomain",
				ProjectDomainID: "proj456",
			},
			expectedDomainName:        "UserDomain",
			expectedDomainID:          "",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "proj456",
		},
		{
			name: "All fields set - IDs take precedence, Names cleared",
			cfg: &parser.SourceConfig{
				DomainName:        "UserDomain",
				DomainID:          "user123",
				ProjectDomainName: "ProjectDomain",
				ProjectDomainID:   "proj456",
			},
			expectedDomainName:        "",
			expectedDomainID:          "user123",
			expectedProjectDomainName: "",
			expectedProjectDomainID:   "proj456",
		},
		{
			name: "id| prefix with explicit projectDomainName",
			cfg: &parser.SourceConfig{
				DomainName:        "id|user789",
				ProjectDomainName: "ProjectDomain",
			},
			expectedDomainName:        "",
			expectedDomainID:          "user789",
			expectedProjectDomainName: "ProjectDomain",
			expectedProjectDomainID:   "", // Do not inherit user domainID when projectDomainName is explicit
		},
		{
			name: "Empty strings treated as not set",
			cfg: &parser.SourceConfig{
				DomainName:        "",
				DomainID:          "",
				ProjectDomainName: "",
				ProjectDomainID:   "",
			},
			expectedDomainName:        "Default",
			expectedDomainID:          "",
			expectedProjectDomainName: "Default",
			expectedProjectDomainID:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveDomainConfig(tt.cfg)

			if result.domainName != tt.expectedDomainName {
				t.Errorf("domainName = %q, want %q", result.domainName, tt.expectedDomainName)
			}
			if result.domainID != tt.expectedDomainID {
				t.Errorf("domainID = %q, want %q", result.domainID, tt.expectedDomainID)
			}
			if result.projectDomainName != tt.expectedProjectDomainName {
				t.Errorf("projectDomainName = %q, want %q", result.projectDomainName, tt.expectedProjectDomainName)
			}
			if result.projectDomainID != tt.expectedProjectDomainID {
				t.Errorf("projectDomainID = %q, want %q", result.projectDomainID, tt.expectedProjectDomainID)
			}
		})
	}
}

// TestResolveDomainConfig_IDPrecedence explicitly tests ID precedence rules.
func TestResolveDomainConfig_IDPrecedence(t *testing.T) {
	t.Run("domainID clears domainName", func(t *testing.T) {
		cfg := &parser.SourceConfig{
			DomainName: "ShouldBeCleared",
			DomainID:   "123",
		}
		result := resolveDomainConfig(cfg)

		if result.domainName != "" {
			t.Errorf("Expected domainName to be cleared when domainID is set, got %q", result.domainName)
		}
		if result.domainID != "123" {
			t.Errorf("Expected domainID = '123', got %q", result.domainID)
		}
	})

	t.Run("explicit domainID overrides id| in domainName", func(t *testing.T) {
		cfg := &parser.SourceConfig{
			DomainName: "id|fromName",
			DomainID:   "fromField",
		}
		result := resolveDomainConfig(cfg)

		if result.domainID != "fromField" {
			t.Errorf("Expected domainID to prefer explicit DomainID, got %q", result.domainID)
		}
	})

	t.Run("projectDomainID clears projectDomainName", func(t *testing.T) {
		cfg := &parser.SourceConfig{
			DomainName:        "UserDomain",
			ProjectDomainName: "ShouldBeCleared",
			ProjectDomainID:   "456",
		}
		result := resolveDomainConfig(cfg)

		if result.projectDomainName != "" {
			t.Errorf("Expected projectDomainName to be cleared when projectDomainID is set, got %q", result.projectDomainName)
		}
		if result.projectDomainID != "456" {
			t.Errorf("Expected projectDomainID = '456', got %q", result.projectDomainID)
		}
	})
}

// TestResolveDomainConfig_Fallbacks tests fallback behavior.
func TestResolveDomainConfig_Fallbacks(t *testing.T) {
	t.Run("projectDomainName falls back to domainName", func(t *testing.T) {
		cfg := &parser.SourceConfig{
			DomainName: "MyDomain",
		}
		result := resolveDomainConfig(cfg)

		if result.projectDomainName != "MyDomain" {
			t.Errorf("Expected projectDomainName to fall back to domainName, got %q", result.projectDomainName)
		}
	})

	t.Run("projectDomainID falls back to domainID", func(t *testing.T) {
		cfg := &parser.SourceConfig{
			DomainID: "abc123",
		}
		result := resolveDomainConfig(cfg)

		if result.projectDomainID != "abc123" {
			t.Errorf("Expected projectDomainID to fall back to domainID, got %q", result.projectDomainID)
		}
	})

	t.Run("explicit projectDomain overrides fallback", func(t *testing.T) {
		cfg := &parser.SourceConfig{
			DomainName:        "UserDomain",
			DomainID:          "user123",
			ProjectDomainName: "ExplicitProject",
			ProjectDomainID:   "proj456",
		}
		result := resolveDomainConfig(cfg)

		// Even though domainID would clear domainName, projectDomain* are explicit
		if result.projectDomainName != "" {
			t.Errorf("Expected projectDomainName to be cleared due to projectDomainID precedence, got %q", result.projectDomainName)
		}
		if result.projectDomainID != "proj456" {
			t.Errorf("Expected projectDomainID = 'proj456', got %q", result.projectDomainID)
		}
	})
}

// TestResolveDomainConfig_DefaultFallback tests the 'Default' domain fallback.
func TestResolveDomainConfig_DefaultFallback(t *testing.T) {
	testCases := []struct {
		name string
		cfg  *parser.SourceConfig
	}{
		{
			name: "empty config",
			cfg:  &parser.SourceConfig{},
		},
		{
			name: "nil strings",
			cfg: &parser.SourceConfig{
				DomainName: "",
				DomainID:   "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := resolveDomainConfig(tc.cfg)

			if result.domainName != "Default" {
				t.Errorf("Expected domainName to default to 'Default', got %q", result.domainName)
			}
			if result.domainID != "" {
				t.Errorf("Expected domainID to remain empty, got %q", result.domainID)
			}
			if result.projectDomainName != "Default" {
				t.Errorf("Expected projectDomainName to fall back to 'Default', got %q", result.projectDomainName)
			}
			if result.projectDomainID != "" {
				t.Errorf("Expected projectDomainID to remain empty, got %q", result.projectDomainID)
			}
		})
	}
}
