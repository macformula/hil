package flow

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/macformula/hil/config" // Assuming config has SetupFirebase()
)

type Store struct {
	*config.FireDB
}

// NewStore returns a Store with the Realtime Database client
func NewStore() *Store {
	d := config.FirebaseDB()
	return &Store{FireDB: d}
}

// Creates the "report" that stores the testID as a root
func (s *Store) CreateReport(reportID uuid.UUID, report *Report) error {
	report.ID = reportID         // Use the provided report ID
	report.DateTime = time.Now() // Set the current date/time
	reportRef := s.FireDB.NewRef(fmt.Sprintf("report/%s", reportID))
	return reportRef.Set(context.Background(), report)
}

// Updates a test under the report tree
func (s *Store) AddTest(reportID string, testNumber int, test *Test) error {
	testRef := s.FireDB.NewRef(fmt.Sprintf("report/%s/test%d", reportID, testNumber))
	return testRef.Set(context.Background(), test)
}

// // Example of adding a sequence of tests (Test1 to TestN)
// func (s *Store) AddSequenceTests(reportID string, tests []Test) error {
// 	for i, test := range tests {
// 		if err := s.AddTest(reportID, i+1, &test); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
