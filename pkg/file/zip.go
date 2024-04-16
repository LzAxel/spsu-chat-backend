package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type ByNumericalFilename []*zip.File

func (nf ByNumericalFilename) Len() int      { return len(nf) }
func (nf ByNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilename) Less(i, j int) bool {

	pathA := nf[i].Name
	pathB := nf[j].Name

	a, err1 := strconv.ParseInt(pathA[0:strings.LastIndex(pathA, ".")], 10, 64)
	b, err2 := strconv.ParseInt(pathB[0:strings.LastIndex(pathB, ".")], 10, 64)

	if err1 != nil || err2 != nil {
		return pathA < pathB
	}

	return a < b
}

func GetFilesFromZip(fileReader io.Reader) ([]*zip.File, error) {
	zipBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return []*zip.File{}, fmt.Errorf("GetFilesFromZip: failed to read file: %w", err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return []*zip.File{}, fmt.Errorf("GetFilesFromZip: failed to create zip reader: %w", err)
	}

	return zipReader.File, nil
}

func SortZipFilesNumerically(files []*zip.File) {
	sort.Sort(ByNumericalFilename(files))
}

func IsOnlyNumericalFiles(files []*zip.File) bool {
	for _, file := range files {
		filename := strings.TrimSuffix(filepath.Base(file.Name), filepath.Ext(file.Name))
		_, err := strconv.Atoi(filename)
		if err != nil {
			return false
		}
	}

	return true
}
