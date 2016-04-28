package mntent

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Entry is an entry in a filesystem table.
type Entry struct {
	Fsname string
	Dir    string
	Type   string
	Opts   string
	Freq   int
	Passno int
}

type Entries []*Entry

func (e *Entry) String() string {
	return fmt.Sprintf("%s %s %s %s %d %d", e.Fsname, e.Dir, e.Type, e.Opts, e.Freq, e.Passno)
}

func parseError(line, msg string) error {
	return fmt.Errorf("Error parsing %q: %s", strings.TrimSpace(line), msg)
}

func parseEntry(line string) (*Entry, error) {
	var entry Entry
	var err error
	fields := strings.Fields(line)
	if len(fields) <= 3 {
		return nil, parseError(line, "less than 3 fields")
	}
	entry.Fsname = fields[0]
	entry.Dir = fields[1]
	entry.Type = fields[2]
	if len(fields) >= 4 {
		entry.Opts = fields[3]
	}
	if len(fields) >= 5 {
		entry.Freq, err = strconv.Atoi(fields[4])
		if err != nil {
			return nil, err
		}
	}
	if len(fields) >= 5 {
		entry.Passno, err = strconv.Atoi(fields[5])
		if err != nil {
			return nil, err
		}

	}
	return &entry, nil
}

func getEntries(fp io.Reader) (Entries, error) {
	var entries Entries
	input := bufio.NewScanner(fp)
	ignoreRe := regexp.MustCompile(`^\s*$|^\s*#.*$`)
	for input.Scan() {
		line := input.Text()
		if ignoreRe.MatchString(line) {
			continue
		}
		entry, err := parseEntry(line)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if input.Err() != nil {
		return nil, input.Err()
	}
	return entries, nil
}

func (entries Entries) ByDir(dir string) (*Entry, error) {
	for _, mnt := range entries {
		if mnt.Dir == dir {
			return mnt, nil
		}
	}
	return nil, fmt.Errorf("%q: mount point not found", dir)

}

func (entries Entries) ByType(fstype string) ([]*Entry, error) {
	var selected []*Entry
	for _, mnt := range entries {
		if mnt.Type == fstype {
			selected = append(selected, mnt)
		}
	}
	return selected, nil
}

// GetEntryByDir returns the mounted filesystem entry for
// the provided mount point.
func GetEntryByDir(dir string) (*Entry, error) {
	dir = filepath.Clean(dir)
	entries, err := GetMounted()
	if err != nil {
		return nil, err
	}
	return entries.ByDir(dir)
}

// GetEntriesByType returns a slice of mounted filesystem
// entries that match the provided type.
func GetEntriesByType(fstype string) ([]*Entry, error) {
	entries, err := GetMounted()
	if err != nil {
		return nil, err
	}
	return entries.ByType(fstype)
}
