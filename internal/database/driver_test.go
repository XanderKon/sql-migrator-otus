package database

import (
	"io"
	"testing"
)

type testDriver struct {
	url       string
	tablename string
}

func (t *testDriver) Open(url string, tablename string) (Driver, error) {
	return &testDriver{
		url:       url,
		tablename: tablename,
	}, nil
}

func (t *testDriver) Close() error {
	return nil
}

func (t *testDriver) Lock() error {
	return nil
}

func (t *testDriver) Unlock() error {
	return nil
}

func (t *testDriver) Run(migration io.Reader) error {
	return nil
}

func (t *testDriver) SetVersion(version int) error {
	return nil
}

func (t *testDriver) DeleteVersion(version int) error {
	return nil
}

func (t *testDriver) Version() (version int, err error) {
	return 0, nil
}

func (t *testDriver) List() (versions []int, err error) {
	return make([]int, 0), nil
}

func (t *testDriver) PrepareTable() error {
	return nil
}

func TestOpen(t *testing.T) {
	// Make sure the driver is registered.
	// But if the previous test already registered it just ignore the panic.
	// If we don't do this it will be impossible to run this test standalone.
	func() {
		defer func() {
			_ = recover()
		}()
		Register("test", &testDriver{})
	}()

	cases := []struct {
		url string
		err bool
	}{
		{
			"test://app:!ChangeMe!@pgsql:5432/app?serverVersion=15&charset=utf8",
			false,
		},
		{
			"postgresql://app:!ChangeMe!@pgsql:5432/app?serverVersion=15&charset=utf8",
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.url, func(t *testing.T) {
			d, err := Open(c.url, "migrations")

			if err == nil {
				if c.err {
					t.Fatal("should be error for wrong driver")
				} else {
					if md, ok := d.(*testDriver); !ok {
						t.Fatalf("expected *testDriver got %T", d)
					} else if md.url != c.url {
						t.Fatalf("expected %q got %q", c.url, md.url)
					}
				}
			} else if !c.err {
				t.Fatalf("did not expect %q", err)
			}
		})
	}
}
