package cmd

import (
	"github.com/google/go-containerregistry/pkg/v1/google"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestSliceMapping(t *testing.T) {
	inputEntries := []Entry{
		EntryOf("key1", "val1"),
		EntryOf("key2", "val2"),
		EntryOf("key2", "val2"),
		EntryOf("compl1", "complval"),
		EntryOf("", "where is key?"),
	}
	wantOutput := []string{"key1", "key2", "key2", "compl1", ""}
	gotOutput := SliceMapping(inputEntries, func(e Entry) string { return e.Key })
	if !reflect.DeepEqual(gotOutput, wantOutput) {
		t.Errorf("SliceMapping() = %v, want %v", gotOutput, wantOutput)
	}
}

type result struct {
	Created time.Time
	Tags    []string
	Digest  string
}

func TestMapToSliceMapping(t *testing.T) {
	now := time.Now()
	inputs := map[string]google.ManifestInfo{
		"3sef9j3k": {
			Created: now,
			Tags:    []string{"2022.12"},
		},
		"d9vmq-d9": {
			Created: now.Truncate(40 * time.Hour),
			Tags:    []string{"JIRA-3942"},
		},
		"91ndiv35": {
			Created: now.Truncate(10 * time.Minute),
			Tags:    []string{},
		},
		"5fdw3ms": {
			Created: now.Truncate(5 * time.Second),
			Tags:    nil,
		},
	}
	wantOutput := []result{
		{
			Created: now,
			Tags:    []string{"2022.12"},
			Digest:  "3sef9j3k",
		},
		{
			Created: now.Truncate(40 * time.Hour),
			Tags:    []string{"JIRA-3942"},
			Digest:  "d9vmq-d9",
		},
		{
			Created: now.Truncate(10 * time.Minute),
			Tags:    []string{},
			Digest:  "91ndiv35",
		},
		{
			Created: now.Truncate(5 * time.Second),
			Digest:  "5fdw3ms",
		},
	}
	gotOutput := MapToSliceMapping(inputs, func(k string, v google.ManifestInfo) result {
		return result{v.Created, v.Tags, k}
	})

	sortResults(wantOutput)
	sortResults(gotOutput)
	if !reflect.DeepEqual(gotOutput, wantOutput) {
		t.Errorf("MapToSliceMapping() = %v, want %v", gotOutput, wantOutput)
	}
}

func sortResults(s []result) {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Created.After(s[j].Created)
	})
}
