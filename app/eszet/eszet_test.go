package eszet

import (
	"testing"

	"github.com/spf13/afero"
)

func TestWrite(t *testing.T) {
	var AppFs = afero.NewOsFs()
	if err := AppFs.MkdirAll("/tmp/eszet", 0777); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := AppFs.Create("/tmp/eszet/test.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := New(file)
	if err := e.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e.Write("test", []byte("test"))
	if value, ok := e.Get("test"); ok {
		if string(value) != "test" {
			t.Fatalf("unexpected value: %s", value)
		}
	} else {
		t.Fatalf("unexpected value: %v", ok)
	}

	e.Write("test", []byte("test2"))
	if value, ok := e.Get("test"); ok {
		if string(value) != "test2" {
			t.Fatalf("unexpected value: %s", value)
		}
	} else {
		t.Fatalf("unexpected value: %v", ok)
	}

	t.Logf("%v", e.content)

	e.Close()

	file, err = AppFs.Open("/tmp/eszet/test.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e = New(file)
	if err := e.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("content: %+v\n", e.content)
	if value, ok := e.Get("test"); ok {
		if string(value) != "test2" {
			t.Fatalf("unexpected value: %s", value)
		}
	} else {
		t.Fatalf("unexpected value: %v", ok)
	}
}
