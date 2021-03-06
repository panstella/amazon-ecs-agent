//+build unit

// Copyright 2014-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package stats

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestIsNetworkStatsError(t *testing.T) {
	isNetStatsErr := isNetworkStatsError(fmt.Errorf("no such file or directory"))
	if isNetStatsErr {
		// Expect it to not be a net stats error
		t.Error("Error incorrectly reported as network stats error")
	}

	isNetStatsErr = isNetworkStatsError(fmt.Errorf("open /sys/class/net/veth2f5f3e4/statistics/tx_bytes: no such file or directory"))
	if !isNetStatsErr {
		// Expect this to be a net stats error
		t.Error("Error incorrectly reported as non network stats error")
	}
}

func TestDockerStatsToContainerStatsMemUsage(t *testing.T) {
	jsonStat := fmt.Sprintf(`
		{
			"cpu_stats":{
				"cpu_usage":{
					"percpu_usage":[%d, %d, %d, %d],
					"total_usage":%d
				}
			},
			"memory_stats":{
				"usage": %d,
				"max_usage": %d,
				"stats": {
					"cache": %d,
					"rss": %d
				},
				"privateworkingset": %d
			}
		}`, 1, 2, 3, 4, 100, 30, 100, 20, 10, 10)
	dockerStat := &types.StatsJSON{}
	json.Unmarshal([]byte(jsonStat), dockerStat)
	containerStats, err := dockerStatsToContainerStats(dockerStat)
	if err != nil {
		t.Errorf("Error converting container stats: %v", err)
	}
	if containerStats == nil {
		t.Fatal("containerStats should not be nil")
	}
	if containerStats.memoryUsage != 10 {
		t.Error("Unexpected value for memoryUsage", containerStats.memoryUsage)
	}
}
