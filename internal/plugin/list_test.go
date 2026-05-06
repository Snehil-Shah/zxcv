// Package plugin_test tests the `plugin` package.
package plugin_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestListAll(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{
		"install":  "#!/bin/sh\n",
		"list-all": "#!/bin/sh\necho \"1.0.0 1.1.0 2.0.0\"\n",
	})
	got, err := p.ListAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"1.0.0", "1.1.0", "2.0.0"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
