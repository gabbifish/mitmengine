package db_test

import (
	"bytes"
	"testing"

	"github.com/cloudflare/mitmengine/db"
	fp "github.com/cloudflare/mitmengine/fputil"
	"github.com/cloudflare/mitmengine/testutil"
)

func TestNewLinearDatabase(t *testing.T) {
	_, err := db.NewLinearDatabase(bytes.NewReader(nil))
	testutil.Ok(t, err)
}

func TestDatabaseLoad(t *testing.T) {
	a, _ := db.NewLinearDatabase(bytes.NewReader(nil))
	err := a.Load(bytes.NewReader(nil))
	testutil.Ok(t, err)
}

func TestDatabaseAdd(t *testing.T) {
	a, _ := db.NewLinearDatabase(bytes.NewReader(nil))
	testutil.Equals(t, 0, len(a.Records))
	a.Add(db.Record{})
	testutil.Equals(t, 1, len(a.Records))
	a.Add(db.Record{})
	testutil.Equals(t, 2, len(a.Records))
}

func TestDatabaseClear(t *testing.T) {
	a, _ := db.NewLinearDatabase(bytes.NewReader(nil))
	a.Add(db.Record{})
	a.Clear()
	testutil.Equals(t, 0, len(a.Records))
	a.Add(db.Record{})
	testutil.Equals(t, 1, len(a.Records))
}

func TestDatabaseGetByUAFingerprint(t *testing.T) {
	var tests = []struct {
		in  fp.UAFingerprint
		out []db.Record
	}{
		{fp.UAFingerprint{}, []db.Record(nil)},
		{fp.UAFingerprint{BrowserName: 1}, []db.Record{{
			UASignature: fp.UASignature{
				BrowserName: 1,
			},
		}}},
		{fp.UAFingerprint{BrowserName: 2}, []db.Record(nil)},
	}
	a, _ := db.NewLinearDatabase(bytes.NewReader(nil))
	a.Add(db.Record{UASignature: fp.UASignature{
		BrowserName: 1,
	}})
	for _, test := range tests {
		testutil.Equals(t, test.out, a.GetByUAFingerprint(test.in))
	}
}

func TestDatabaseGetByRequestFingerprint(t *testing.T) {
	var tests = []struct {
		in  fp.RequestFingerprint
		out []db.Record
	}{
		{fp.RequestFingerprint{}, []db.Record(nil)},
		{fp.RequestFingerprint{Version: 1}, []db.Record{{
			RequestSignature: fp.RequestSignature{
				Version: fp.VersionSignature{
					Exp: 1,
					Min: 1,
					Max: 1,
				},
			},
		}}},
		{fp.RequestFingerprint{Version: 2}, []db.Record(nil)},
	}
	a, _ := db.NewLinearDatabase(bytes.NewReader(nil))
	a.Add(db.Record{RequestSignature: fp.RequestSignature{
		Version: fp.VersionSignature{Exp: 1, Min: 1, Max: 1},
	}})
	for _, test := range tests {
		testutil.Equals(t, test.out, a.GetByRequestFingerprint(test.in))
	}
}
