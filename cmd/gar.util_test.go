package cmd

import "testing"

func TestPackageToMicroserviceName(t *testing.T) {
	previousConfig := config
	t.Cleanup(func() {
		config = previousConfig
	})

	config = configuration{
		GAR: GAR{
			PackagePrefix: "company-",
		},
		Mappings: map[string]ServiceMappings{
			"api-events": {
				GAR: "events",
			},
		},
	}

	tests := []struct {
		name        string
		packageName string
		want        string
	}{
		{
			name:        "prefix is stripped when mapping is absent",
			packageName: "company-api-users",
			want:        "api-users",
		},
		{
			name:        "package name is returned as is when no mapping matches",
			packageName: "standalone-service",
			want:        "standalone-service",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := packageToMicroserviceName(test.packageName); got != test.want {
				t.Fatalf("resolveMicroserviceNameFromGarPackage() = %q, want %q", got, test.want)
			}
		})
	}
}
