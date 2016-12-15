/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kvm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func getEvent(sysPath string, eventName string, ns plugin.Namespace) (plugin.Metric, error) {
	filePath := filepath.Join(sysPath, eventName)

	value, err := getValue(filePath)
	if err != nil {
		return plugin.Metric{}, err
	}

	metric := plugin.Metric{
		Namespace: ns,
		Data:      value,
		Timestamp: time.Now(),
		Version:   Version,
	}
	return metric, nil
}

func getValue(filename string) (int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("Cannot open file, err: %v", err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	// The file should contain only one line.
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("Cannot read the content of file %s, err: %v", filename, err)
	}
	// trim white spaces
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, hexPrefix) {
		// check if value is a hex value and do appropriate parsing
		line = strings.TrimPrefix(line, hexPrefix)
		return strconv.ParseInt(line, 16, 0)
	}
	return strconv.ParseInt(line, 10, 0)
}
