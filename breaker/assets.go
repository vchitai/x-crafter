package breaker

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type version struct {
	major string
	minor string
	patch string
}

func (v version) toString() []byte {
	return []byte(fmt.Sprintf("%s.%s.%s", v.major, v.minor, v.patch))
}
func readVersion(s string) (*version, error) {
	verSeg := strings.Split(s, ".")
	if len(verSeg) != 3 {
		return nil, fmt.Errorf("not enough part")
	}
	return &version{
		major: verSeg[0],
		minor: verSeg[1],
		patch: verSeg[2],
	}, nil
}

func versioning(dest string) {
	d, err := ioutil.ReadFile(filepath.Join(dest, "version"))
	if errors.Is(err, os.ErrNotExist) {
		if err := ioutil.WriteFile(filepath.Join(dest, "version"), version{"1", "0", "0"}.toString(), 0644); err != nil {
			log.Println("Creating assets: ", err)
			return
		}
	} else if ver, err := readVersion(string(d)); err != nil {
		if err := ioutil.WriteFile(filepath.Join(dest, "version"), version{"1", "0", "0"}.toString(), 0644); err != nil {
			log.Println("Creating assets: ", err)
			return
		}
	} else {
		patch, _ := strconv.Atoi(ver.patch)
		ver.patch = strconv.Itoa(patch + 1)
		if err := ioutil.WriteFile(filepath.Join(dest, "version"), ver.toString(), 0644); err != nil {
			log.Println("Creating assets: ", err)
			return
		}
	}
}
