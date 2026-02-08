package mock

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// bookmark represents a pagination bookmark
type bookmark struct {
	Key string `json:"key"`
}

// encodeBookmark creates a bookmark string from a key
func encodeBookmark(key string) string {
	if key == "" {
		return ""
	}

	bm := bookmark{Key: key}
	data, _ := json.Marshal(bm)
	return base64.StdEncoding.EncodeToString(data)
}

// decodeBookmark parses a bookmark string and returns the key
func decodeBookmark(bookmarkStr string) (string, error) {
	if bookmarkStr == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(bookmarkStr)
	if err != nil {
		return "", fmt.Errorf("invalid bookmark: %w", err)
	}

	var bm bookmark
	if err := json.Unmarshal(data, &bm); err != nil {
		return "", fmt.Errorf("invalid bookmark format: %w", err)
	}

	return bm.Key, nil
}

// paginateKeys applies pagination to a list of keys
// Returns the paginated keys and the next bookmark
func paginateKeys(keys []string, pageSize int32, bookmarkStr string) ([]string, string, error) {
	if pageSize <= 0 {
		return nil, "", fmt.Errorf("invalid pageSize: must be greater than 0")
	}

	// Decode bookmark to find starting position
	startKey, err := decodeBookmark(bookmarkStr)
	if err != nil {
		return nil, "", err
	}

	// Find the starting index
	startIdx := 0
	if startKey != "" {
		found := false
		for i, key := range keys {
			if key == startKey {
				startIdx = i + 1 // Start after the bookmark
				found = true
				break
			}
		}
		if !found {
			return nil, "", fmt.Errorf("bookmark not found in result set")
		}
	}

	// Calculate end index
	endIdx := startIdx + int(pageSize)
	if endIdx > len(keys) {
		endIdx = len(keys)
	}

	// Extract page
	pageKeys := keys[startIdx:endIdx]

	// Create next bookmark
	var nextBookmark string
	if endIdx < len(keys) {
		nextBookmark = encodeBookmark(keys[endIdx-1])
	}

	return pageKeys, nextBookmark, nil
}

// createQueryMetadata creates pagination metadata
func createQueryMetadata(fetchedRecordsCount int32, bookmark string) *pb.QueryResponseMetadata {
	return &pb.QueryResponseMetadata{
		FetchedRecordsCount: fetchedRecordsCount,
		Bookmark:            bookmark,
	}
}
