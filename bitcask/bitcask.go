package bitcask

import (
	// "bytes"
	// "encoding/json"
	"fmt"
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
		asInt, err := strconv.Atoi(strings.Trim(file, ".db"))
		if err != nil {
			return err
		}

		sf, err := NewDataFile(".", asInt)

		b.StaleFiles = append(b.StaleFiles, sf)
	}

	return nil
}

// TODO: Implement this
// Goes through the files and builds a KeyDir. In the future, it will be
// generated based on the 'hints file'.
func (b *Bitcask) buildKeydir() error {
	return nil
}
