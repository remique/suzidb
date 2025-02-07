package bitcask

import (
	// "fmt"
	// "path/filepath"
	"sort"
	"strconv"
	"strings"
)

func getLastFileId(matches []string) (int, error) {
	if len(matches) == 0 {
		return 1, nil
	}

	sort.Strings(matches)

	lastString := matches[len(matches)-1]
	lastInt, err := strconv.Atoi(strings.Trim(lastString, ".db"))
	if err != nil {
		return -1, err
	}

	return lastInt, nil
}

// func glob(dirName string) (int, error) {
// 	matches, err := filepath.Glob(dirPath + "/*.db")
// 	fmt.Println(matches)
// 	if err != nil {
// 		return -1, err
// 	}

// 	if len(matches) == 0 {
// 		return 1, nil
// 	}

// 	sort.Strings(matches)

// 	lastString := matches[len(matches)-1]
// 	lastInt, err := strconv.Atoi(lastString)
// 	if err != nil {
// 		return -1, err
// 	}

// 	return lastInt, nil
// }
