package test

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger-labs/cc-tools/mock"
)

func TestRichQueryIntegration(t *testing.T) {
	stub := mock.NewMockStub("richQueryTest", nil)

	alice := map[string]interface{}{
		"@assetType": "person",
		"name":       "alice",
		"age":        25,
		"city":       "NYC",
	}
	bob := map[string]interface{}{
		"@assetType": "person",
		"name":       "bob",
		"age":        30,
		"city":       "LA",
	}
	charlie := map[string]interface{}{
		"@assetType": "person",
		"name":       "charlie",
		"age":        35,
		"city":       "NYC",
	}

	aliceData, _ := json.Marshal(alice)
	bobData, _ := json.Marshal(bob)
	charlieData, _ := json.Marshal(charlie)

	stub.MockTransactionStart("setup")
	stub.PutState("a", aliceData)
	stub.PutState("b", bobData)
	stub.PutState("c", charlieData)
	stub.MockTransactionEnd("setup")

	t.Run("basic selector", func(t *testing.T) {
		query := `{"selector": {"@assetType": "person"}}`
		iter, err := stub.GetQueryResult(query)
		if err != nil {
			t.Fatalf("GetQueryResult failed: %v", err)
		}
		defer iter.Close()

		count := 0
		for iter.HasNext() {
			_, err := iter.Next()
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			count++
		}
		if count != 3 {
			t.Errorf("Expected 3 person documents, got %d", count)
		}
	})

	t.Run("comparison operator", func(t *testing.T) {
		query := `{"selector": {"@assetType": "person", "age": {"$gte": 30}}}`
		iter, err := stub.GetQueryResult(query)
		if err != nil {
			t.Fatalf("GetQueryResult failed: %v", err)
		}
		defer iter.Close()

		count := 0
		for iter.HasNext() {
			_, err := iter.Next()
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			count++
		}
		if count != 2 {
			t.Errorf("Expected 2 people aged >= 30, got %d", count)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		query := `{"selector": {"@assetType": "person"}}`
		iter, meta, err := stub.GetQueryResultWithPagination(query, 2, "")
		if err != nil {
			t.Fatalf("GetQueryResultWithPagination failed: %v", err)
		}
		defer iter.Close()

		count := 0
		for iter.HasNext() {
			_, err := iter.Next()
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			count++
		}
		if count != 2 {
			t.Errorf("Expected 2 results in first page, got %d", count)
		}
		if meta.FetchedRecordsCount != 2 {
			t.Errorf("Expected FetchedRecordsCount=2, got %d", meta.FetchedRecordsCount)
		}
		if meta.Bookmark == "" {
			t.Error("Expected non-empty bookmark")
		}
	})

	t.Run("registry override", func(t *testing.T) {
		query := `{"selector": {"city": "NYC"}}`
		stub.RegisterQueryResult(query, []string{"c"})

		iter, err := stub.GetQueryResult(query)
		if err != nil {
			t.Fatalf("GetQueryResult failed: %v", err)
		}
		defer iter.Close()

		results := make([]string, 0)
		for iter.HasNext() {
			kv, err := iter.Next()
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			results = append(results, kv.Key)
		}
		if len(results) != 1 {
			t.Fatalf("Expected 1 result from override, got %d", len(results))
		}
		if results[0] != "c" {
			t.Errorf("Expected key 'c', got %s", results[0])
		}
	})
}
