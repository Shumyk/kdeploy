package cmd

import "testing"

func TestComposeImagePath(t *testing.T) {
	type args struct {
		registry string
		repo     string
		service  string
		tag      string
		digest   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"full image path",
			args{
				"us.gcr.io",
				"company-infra/company/company-",
				"api-composer",
				"2022.12",
				"ds324knf2foi43kkl",
			},
			"us.gcr.io/company-infra/company/company-api-composer:2022.12@sha256:ds324knf2foi43kkl",
		},
		{
			"image path without tag",
			args{
				"gcr.io",
				"company-infra/company-",
				"api-events",
				"",
				"3kfo3mcs98dvj",
			},
			"gcr.io/company-infra/company-api-events@sha256:3kfo3mcs98dvj",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ComposeImagePath(test.args.registry, test.args.repo, test.args.service, test.args.tag, test.args.digest); got != test.want {
				t.Errorf("ComposeImagePath() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestParseImagePath(t *testing.T) {
	tests := []struct {
		name       string
		imagePath  string
		wantTag    string
		wantDigest string
	}{
		{
			"image path with tag",
			"eu.gcr.io/company-infra/company-data-manager:2023.01@sha256:30d9vj3kmf2oi39u93n",
			"2023.01",
			"30d9vj3kmf2oi39u93n",
		},
		{
			"image path without tag",
			"us.gcr.io/company-infra/company-datastore@sha256:3fm1i0s3d9v32osk",
			"",
			"3fm1i0s3d9v32osk",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotTag, gotDigest := ParseImagePath(test.imagePath)
			if gotTag != test.wantTag {
				t.Errorf("ParseImagePath() gotTag = %v, want %v", gotTag, test.wantTag)
			}
			if gotDigest != test.wantDigest {
				t.Errorf("ParseImagePath() gotDigest = %v, want %v", gotDigest, test.wantDigest)
			}
		})
	}
}
