package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dotcloud/docker/pkg/cgroups"
)

type blkioGroup struct {
}

func (s *blkioGroup) Set(d *data) error {
	// we just want to join this group even though we don't set anything
	if _, err := d.join("blkio"); err != nil && err != cgroups.ErrNotFound {
		return err
	}
	fmt.Println(s.Stats(d))
	return nil
}

func (s *blkioGroup) Remove(d *data) error {
	return removePath(d.path("blkio"))
}

func (s *blkioGroup) Stats(d *data) (map[string]float64, error) {
	paramData := make(map[string]float64)
	path, err := d.path("blkio")
	if err != nil {
		return paramData, fmt.Errorf("Unable to read %s cgroup param: %s", path, err)
	}
	params := []string{
		"sectors",
		"io_service_bytes",
		"io_serviced",
		"io_queued",
	}
	for _, param := range params {
		p := fmt.Sprintf("blkio.%s", param)
		paramPath := filepath.Join(path, p)
		f, err := os.Open(paramPath)
		if err != nil {
			return paramData, err
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			fields := strings.Fields(sc.Text())
			v, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				fmt.Printf("Error parsing %s stats: %s", param, err)
				continue
			}
			paramData[param] = v
		}
	}
	return paramData, nil
}
