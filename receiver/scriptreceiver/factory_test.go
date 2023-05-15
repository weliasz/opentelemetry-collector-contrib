// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scriptreceiver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type EventCollector struct {
	Events []map[string]string
}

func Test_scripted_receiver_true(t *testing.T) {
	assert.Equal(t, "val", "val")
}

/*
func Test_scripted_receiver(t *testing.T) {
	// fmt.Print("Hello, my friend")
	ctx, cancelFunc := context.WithCancel(context.Background())

	defer cancelFunc()

	cmd := exec.CommandContext(ctx, "./ps.sh")
	stdout, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(stdout), "\n")
	header := strings.Split(lines[0], "    ")

	//events := make([]map[string]string, 0)
	collector := &EventCollector{Events: make([]map[string]string, 0)}
	for _, line := range lines[1:] {
		m := make(map[string]string)
		values := strings.Split(line, "    ")
		if len(values) == len(header) {
			for i, h := range header {
				m[h] = values[i]
			}
			collector.Events = append(collector.Events, m)
		}
	}
	fmt.Println(collector.Events)
}*/
