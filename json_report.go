package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedPageData := make([]PageData, len(pages))
	for _, val := range keys {
		sortedPageData = append(sortedPageData, pages[val])
	}
	data, err := json.MarshalIndent(sortedPageData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal the pages: %v", err)
	}
	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write '%v':%v", filename, err)
	}
	return nil
}
