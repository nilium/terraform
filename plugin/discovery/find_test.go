package discovery

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFindPluginPaths(t *testing.T) {
	got := findPluginPaths(
		"foo",
		"mockos_mockarch",
		[]string{
			"test-fixtures/current-style-plugins",
			"test-fixtures/legacy-style-plugins",
			"test-fixtures/non-existent",
			"test-fixtures/not-a-dir",
		},
	)
	want := []string{
		filepath.Join("test-fixtures", "current-style-plugins", "mockos_mockarch", "terraform-foo-bar_v0.0.1"),
		filepath.Join("test-fixtures", "current-style-plugins", "mockos_mockarch", "terraform-foo-bar_v1.0.0"),
		filepath.Join("test-fixtures", "legacy-style-plugins", "terraform-foo-bar"),
		filepath.Join("test-fixtures", "legacy-style-plugins", "terraform-foo-baz"),
	}

	// Turn the paths back into relative paths, since we don't care exactly
	// where this code is present on the system for the sake of this test.
	baseDir, err := os.Getwd()
	if err != nil {
		// Should never happen
		panic(err)
	}
	for i, absPath := range got {
		if !filepath.IsAbs(absPath) {
			t.Errorf("got non-absolute path %s", absPath)
		}

		got[i], err = filepath.Rel(baseDir, absPath)
		if err != nil {
			t.Fatalf("Can't make %s relative to current directory %s", absPath, baseDir)
		}
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestResolvePluginPaths(t *testing.T) {
	got := ResolvePluginPaths([]string{
		"/example/mockos_mockarch/terraform-foo-bar_v0.0.1",
		"/example/mockos_mockarch/terraform-foo-baz_v0.0.1",
		"/example/mockos_mockarch/terraform-foo-baz_v1.0.0",
		"/example/mockos_mockarch/terraform-foo-baz_v2.0.0_x4",
		"/example/mockos_mockarch/terraform-foo-upper_V2.0.0_X4",
		"/example/terraform-foo-bar",
		"/example/mockos_mockarch/terraform-foo-bar_vbananas",
		"/example/mockos_mockarch/terraform-foo-bar_v",
		"/example2/mockos_mockarch/terraform-foo-bar_v0.0.1",
	})

	want := []PluginMeta{
		{
			Name:    "bar",
			Version: "0.0.1",
			Path:    "/example/mockos_mockarch/terraform-foo-bar_v0.0.1",
		},
		{
			Name:    "baz",
			Version: "0.0.1",
			Path:    "/example/mockos_mockarch/terraform-foo-baz_v0.0.1",
		},
		{
			Name:    "baz",
			Version: "1.0.0",
			Path:    "/example/mockos_mockarch/terraform-foo-baz_v1.0.0",
		},
		{
			Name:    "baz",
			Version: "2.0.0",
			Path:    "/example/mockos_mockarch/terraform-foo-baz_v2.0.0_x4",
		},
		{
			Name:    "upper",
			Version: "2.0.0",
			Path:    "/example/mockos_mockarch/terraform-foo-upper_V2.0.0_X4",
		},
		{
			Name:    "bar",
			Version: "0.0.0",
			Path:    "/example/terraform-foo-bar",
		},
		{
			Name:    "bar",
			Version: "bananas",
			Path:    "/example/mockos_mockarch/terraform-foo-bar_vbananas",
		},
		{
			Name:    "bar",
			Version: "",
			Path:    "/example/mockos_mockarch/terraform-foo-bar_v",
		},
	}

	for p := range got {
		t.Logf("got %#v", p)
	}

	if got, want := got.Count(), len(want); got != want {
		t.Errorf("got %d items; want %d", got, want)
	}

	for _, wantMeta := range want {
		if !got.Has(wantMeta) {
			t.Errorf("missing meta %#v", wantMeta)
		}
	}
}
