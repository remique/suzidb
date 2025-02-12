package bitcask

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Bitcask struct {
	Options    *Options
	KeyDir     KeyDir
	ActiveFile *DataFile
	StaleFiles []*DataFile
}

func NewBitcask(opts ...Config) (*Bitcask, error) {
	b := &Bitcask{
		Options: DefaultOptions(),
		KeyDir:  KeyDir{},
	}

	for _, opt := range opts {
		opt(b.Options)
	}

	// Assing ActiveFile
	err := b.buildActiveFile()
	if err != nil {
		return nil, fmt.Errorf("Error while building activeFile: %s", err.Error())
	}

	// Assing StaleFiles
	err = b.buildStaleFiles()
	if err != nil {
		return nil, fmt.Errorf("Error while building staleFiles: %s", err.Error())
	}

	// Assing KeyDir
	err = b.buildKeydir()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Opens a new ActiveFile.
func (b *Bitcask) buildActiveFile() error {
	newActiveId, err := generateNewActiveFileId(b.Options.dir)
	if err != nil {
		return err
	}

	af, err := NewDataFile(b.Options.dir, newActiveId)
	if err != nil {
		return err
	}

	b.ActiveFile = af

	return nil
}

// Opens all other StaleFiles.
func (b *Bitcask) buildStaleFiles() error {
	allFilesGlob, err := glob(b.Options.dir)
	if err != nil {
		return err
	}

	for _, file := range allFilesGlob {
		fileBase := filepath.Base(file)
		asInt, err := strconv.Atoi(strings.Trim(fileBase, ".db"))
		if err != nil {
			return err
		}

		sf, err := NewDataFile(b.Options.dir, asInt)

		b.StaleFiles = append(b.StaleFiles, sf)
	}

	return nil
}

// TODO: Implement this
// Goes through the files and builds a KeyDir. In the future, it will be
// generated based on the 'hints file'.
func (b *Bitcask) buildKeydir() error {
	for _, file := range b.StaleFiles {
		decoder := json.NewDecoder(file.Fd)

		for {
			var dr DiskRecord
			if err := decoder.Decode(&dr); err != nil {
				if err.Error() == "EOF" {
					break
				}

				return fmt.Errorf("Error while building keydir: %s", err.Error())
			}

			// Convert DiskRecord into KeyDirRecord
			kdr := KeyDirRecord{
				FileId:    file.Id,
				ValueSize: len(dr.Value),
				// TODO: ValuePos
			}

			b.KeyDir[dr.Key] = kdr
		}
	}
	return nil
}
