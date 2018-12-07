package db

import (
	"github.com/cloudflare/mitmengine/fputil"
	"io"
)

type Database interface {
	Load(io.Reader) error
	Add(Record)
	Clear()
	Dump(io.Writer) error
	GetByRequestFingerprint(fp.RequestFingerprint) []Record
	GetByUAFingerprint(fp.UAFingerprint) []Record
}
