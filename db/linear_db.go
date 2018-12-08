package db

import (
	"bufio"
	"fmt"
	"github.com/cloudflare/mitmengine/fputil"
	"io"
	"strings"
)

type LinearDatabase struct {
	Records   []Record
}

// NewDatabase returns a new Database initialized from the configuration.
func NewLinearDatabase(input io.Reader) (*LinearDatabase, error) {
	var a LinearDatabase
	err := a.Load(input)
	return &a, err
}

// Load records from input into the database, and return an error on bad records.
func (a *LinearDatabase) Load(input io.Reader) error {
	var record Record
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		recordString := scanner.Text()
		if idx := strings.IndexRune(recordString, '\t'); idx != -1 {
			// remove anything before a tab
			recordString = recordString[idx+1:]
		}
		if idx := strings.IndexRune(recordString, '#'); idx != -1 {
			// remove comments at end of lines
			recordString = recordString[:idx]
		}
		// remove any whitespace or quotes
		recordString = strings.Trim(strings.TrimSpace(recordString), "\"")
		if len(recordString) == 0 {
			continue // skip empty lines
		}
		if err := record.Parse(recordString); err != nil {
			return fmt.Errorf("unable to parse record: %s, %s", recordString, err)
		}
		a.Add(record)
	}
	return nil
}

// USE MAP FOR FASTER LOOKUP; map of browser type?
// Add a single record to the database.
func (a *LinearDatabase) Add(record Record) {
	a.Records = append(a.Records, record)
}

// Clear all records from the linear database.
func (a *LinearDatabase) Clear() {
	a.Records = nil
}

// Dump records in the database to output.
func (a *LinearDatabase) Dump(output io.Writer) error {
	for _, record := range a.Records {
		_, err := fmt.Fprintln(output, record)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetByRequestFingerprint returns all records in the database matching the
// request fingerprint.
func (a *LinearDatabase) GetByRequestFingerprint(requestFingerprint fp.RequestFingerprint) []Record {
	return a.GetBy(func(r Record) bool {
		match, _ := r.RequestSignature.Match(requestFingerprint)
		return match != fp.MatchImpossible
	})
}

// GetByUAFingerprint returns all records in the database matching the
// user agent fingerprint.
func (a *LinearDatabase) GetByUAFingerprint(uaFingerprint fp.UAFingerprint) []Record {
	return a.GetBy(func(r Record) bool { return r.UASignature.Match(uaFingerprint) != fp.MatchImpossible })
}

// GetBy returns a list of records for which GetBy returns true.
func (a *LinearDatabase) GetBy(getFunc func(Record) bool) []Record {
	var matchedRecords []Record
	for _, record := range a.Records {
		if getFunc(record) {
			matchedRecords = append(matchedRecords, record)
		}
	}
	return matchedRecords
}

// DeleteBy deletes records for which rmFunc returns true.
func (a *LinearDatabase) DeleteBy(deleteFunc func(Record) bool) {
	for idx, record := range a.Records {
		if deleteFunc(record) {
			a.Records = append(a.Records[:idx], a.Records[idx+1:]...)
		}
	}
}

// MergeBy merges records for which mergeFunc returns true.
func (a *LinearDatabase) MergeBy(mergeFunc func(Record, Record) bool) (int, int) {
	before := len(a.Records)
	for id1, record1 := range a.Records {
		for id2, record2 := range a.Records {
			if record1.String() == record2.String() {
				continue
			}
			if mergeFunc(record1, record2) {
				a.Records[id1] = record1.Merge(record2)
				// If elements are deleted from the map during the iteration, they will not be produced.
				// https://golang.org/ref/spec#For_statements
				//delete(a.RecordMap, id2)
				a.Records[id2] = Record{} // Treat the merged record as deleted.
			}
		}
	}
	return before, len(a.Records)
}

