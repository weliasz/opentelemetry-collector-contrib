// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dockerstatsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver"

import (
	dtypes "github.com/docker/docker/api/types"
)

// Following functions has been copied from: calculateCPUPercentUnix(), calculateMemUsageUnixNoCache(), calculateMemPercentUnixNoCache()
// https://github.com/docker/cli/blob/a2e9ed3b874fccc177b9349f3b0277612403934f/cli/command/container/stats_helpers.go

// Copyright 2012-2017 Docker, Inc.
// This product includes software developed at Docker, Inc. (https://www.docker.com).
// The following is courtesy of our legal counsel:
// Use and transfer of Docker may be subject to certain restrictions by the
// United States and other governments.
// It is your responsibility to ensure that your use and/or transfer does not
// violate applicable laws.
// For more information, please see https://www.bis.doc.gov
// See also https://www.apache.org/dev/crypto.html and/or seek legal counsel.

func calculateCPUPercent(previous *dtypes.CPUStats, v *dtypes.CPUStats) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUUsage.TotalUsage) - float64(previous.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.SystemUsage) - float64(previous.SystemUsage)
		onlineCPUs  = float64(v.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(v.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

// calculateMemUsageNoCache calculate memory usage of the container.
// Cache is intentionally excluded to avoid misinterpretation of the output.
//
// On cgroup v1 host, the result is `mem.Usage - mem.Stats["total_inactive_file"]` .
// On cgroup v2 host, the result is `mem.Usage - mem.Stats["inactive_file"] `.
//
// This definition is consistent with cadvisor and containerd/CRI.
// * https://github.com/google/cadvisor/commit/307d1b1cb320fef66fab02db749f07a459245451
// * https://github.com/containerd/cri/commit/6b8846cdf8b8c98c1d965313d66bc8489166059a
//
// On Docker 19.03 and older, the result was `mem.Usage - mem.Stats["cache"]`.
// See https://github.com/moby/moby/issues/40727 for the background.
func calculateMemUsageNoCache(memoryStats *dtypes.MemoryStats) uint64 {
	// cgroup v1
	if v, isCgroup1 := memoryStats.Stats["total_inactive_file"]; isCgroup1 && v < memoryStats.Usage {
		return memoryStats.Usage - v
	}
	// cgroup v2
	if v := memoryStats.Stats["inactive_file"]; v < memoryStats.Usage {
		return memoryStats.Usage - v
	}
	return memoryStats.Usage
}

func calculateMemoryPercent(limit uint64, usedNoCache uint64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return float64(usedNoCache) / float64(limit) * 100.0
	}
	return 0.0
}
