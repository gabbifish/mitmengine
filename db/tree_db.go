package db

import (
	"bufio"
	"fmt"
	"github.com/cloudflare/mitmengine/fputil"
	"io"
	"strings"
)

// A Database contains a collection of records containing software signatures.
type TreeDatabase struct {
	RecordMap map[fp.UABrowserName]map[fp.UAOSName][]Record
}

// NewDatabase returns a new Database initialized from the configuration.
func NewTreeDatabase(input io.Reader) (*TreeDatabase, error) {
	var a TreeDatabase
	a.RecordMap = make(map[fp.UABrowserName]map[fp.UAOSName][]Record)
	err := a.Load(input)
	return &a, err
}

// Load records from input into the database, and return an error on bad records.
func (a *TreeDatabase) Load(input io.Reader) error {
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

func (a *TreeDatabase) Add(record Record) {
	// First index by browser name
	if _, ok := a.RecordMap[record.UASignature.BrowserName]; !ok {
		a.RecordMap[record.UASignature.BrowserName] = make(map[fp.UAOSName][]Record)
	}
	recordList := a.RecordMap[record.UASignature.BrowserName][record.UASignature.OSName]
	recordList = append(recordList, record)
	a.RecordMap[record.UASignature.BrowserName][record.UASignature.OSName] = recordList
}

// Clear all records from the database.
func (a *TreeDatabase) Clear() {
	for browserNameKey := range a.RecordMap {
		delete(a.RecordMap, browserNameKey)
	}
}

// Dump records in the database to output.
func (a *TreeDatabase) Dump(output io.Writer) error {
	for _, osMap := range a.RecordMap {
		for _, recordList := range osMap {
			for record := range recordList {
				_, err := fmt.Fprintln(output, record)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetByUAFingerprint returns all records in the database matching the
// user agent fingerprint.
func (a *TreeDatabase) GetByUAFingerprint(uaFingerprint fp.UAFingerprint) []Record {
	var candidateRecords []Record
	candidateUAFingerprints := a.RecordMap[uaFingerprint.BrowserName][uaFingerprint.OSName]
	for _, record := range candidateUAFingerprints {
		if record.UASignature.Match(uaFingerprint) != fp.MatchImpossible {
			candidateRecords = append(candidateRecords, record)
		}
	}
	return candidateRecords
}