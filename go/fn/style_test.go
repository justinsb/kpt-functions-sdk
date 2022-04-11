package fn

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStyle(t *testing.T) {
	j := `{"apiVersion": "v1", "kind": "SomeKind", "metadata": {"name": "somename"}, "spec": {"foo": "bar"}}`

	obj, err := parseOneKubeObject([]byte(j))
	if err != nil {
		t.Fatalf("error from parseOneKubeObject: %v", err)
	}

	obj.NormalizeStyle()

	b, err := obj.ToYAML()
	if err != nil {
		t.Fatalf("error from ToYAML: %v", err)
	}

	got := string(b)
	want := `
apiVersion: v1
kind: SomeKind
metadata:
  name: somename
spec:
  foo: bar
`

	got = strings.TrimSpace(got)
	want = strings.TrimSpace(want)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected result (-want, +got) %s", diff)
	}
}
