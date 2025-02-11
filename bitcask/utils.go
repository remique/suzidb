package bitcask

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Based on array of matched file names, returns the last index of currently
// used file.
func getLastFileId(matches []string) (int, error) {
	if len(matches) == 0 {
		return 0, nil
	}

	sort.Strings(matches)

	lastString := matches[len(matches)-1]
	lastInt, err := strconv.Atoi(strings.Trim(lastString, ".db"))
	if err != nil {
		return -1, err
	}

	return lastInt, nil
}

func glob(dirName string) ([]string, error) {
	matches, err := filepath.Glob(dirName + "/*.db")
	if err != nil {
		return []string{}, err
	}

	return matches, nil
}

func generateNewActiveFileId(dir string) (int, error) {
	matches, err := glob(dir)
	fmt.Println("matches", matches)
	if err != nil {
		return -1, err
	}

	lastId, err := getLastFileId(matches)
	if err != nil {
		return -1, err
	}

	return lastId + 1, nil
}
