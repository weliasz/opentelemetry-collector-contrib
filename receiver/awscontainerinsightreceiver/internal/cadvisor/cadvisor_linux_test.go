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

//go:build linux
// +build linux

package cadvisor

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/cadvisor/cache/memory"
	"github.com/google/cadvisor/container"
	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver/internal/cadvisor/extractors"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver/internal/cadvisor/testutils"
)

type mockCadvisorManager struct {
	t *testing.T
}

// Start the manager. Calling other manager methods before this returns
// may produce undefined behavior.
func (m *mockCadvisorManager) Start() error {
	return nil
}

// Get information about all subcontainers of the specified container (includes self).
func (m *mockCadvisorManager) SubcontainersInfo(containerName string, query *info.ContainerInfoRequest) ([]*info.ContainerInfo, error) {
	containerInfos := testutils.LoadContainerInfo(m.t, "./extractors/testdata/CurInfoContainer.json")
	return containerInfos, nil
}

type mockCadvisorManager2 struct {
}

func (m *mockCadvisorManager2) Start() error {
	return errors.New("new error")
}

func (m *mockCadvisorManager2) SubcontainersInfo(containerName string, query *info.ContainerInfoRequest) ([]*info.ContainerInfo, error) {
	return nil, nil
}

func newMockCreateManager(t *testing.T) createCadvisorManager {
	return func(memoryCache *memory.InMemoryCache, sysfs sysfs.SysFs, houskeepingConfig manager.HouskeepingConfig,
		includedMetricsSet container.MetricSet, collectorHTTPClient *http.Client, rawContainerCgroupPathPrefixWhiteList []string,
		perfEventsFile string) (cadvisorManager, error) {
		return &mockCadvisorManager{t: t}, nil
	}
}

var mockCreateManager2 = func(memoryCache *memory.InMemoryCache, sysfs sysfs.SysFs, houskeepingConfig manager.HouskeepingConfig,
	includedMetricsSet container.MetricSet, collectorHTTPClient *http.Client, rawContainerCgroupPathPrefixWhiteList []string,
	perfEventsFile string) (cadvisorManager, error) {
	return &mockCadvisorManager2{}, nil
}

var mockCreateManagerWithError = func(memoryCache *memory.InMemoryCache, sysfs sysfs.SysFs, houskeepingConfig manager.HouskeepingConfig,
	includedMetricsSet container.MetricSet, collectorHTTPClient *http.Client, rawContainerCgroupPathPrefixWhiteList []string,
	perfEventsFile string) (cadvisorManager, error) {
	return nil, errors.New("error")
}

type MockK8sDecorator struct {
}

func (m *MockK8sDecorator) Decorate(metric *extractors.CAdvisorMetric) *extractors.CAdvisorMetric {
	return metric
}

func TestGetMetrics(t *testing.T) {
	t.Setenv("HOST_NAME", "host")
	hostInfo := testutils.MockHostInfo{ClusterName: "cluster"}
	k8sdecoratorOption := WithDecorator(&MockK8sDecorator{})

	c, err := New("eks", hostInfo, zap.NewNop(), cadvisorManagerCreator(newMockCreateManager(t)), k8sdecoratorOption)
	assert.NotNil(t, c)
	assert.NoError(t, err)
	assert.NotNil(t, c.GetMetrics())
}

func TestGetMetricsNoEnv(t *testing.T) {
	hostInfo := testutils.MockHostInfo{ClusterName: "cluster"}
	k8sdecoratorOption := WithDecorator(&MockK8sDecorator{})

	c, err := New("eks", hostInfo, zap.NewNop(), cadvisorManagerCreator(newMockCreateManager(t)), k8sdecoratorOption)
	assert.Nil(t, c)
	assert.Error(t, err)
}

func TestGetMetricsNoClusterName(t *testing.T) {
	t.Setenv("HOST_NAME", "host")
	hostInfo := testutils.MockHostInfo{}
	k8sdecoratorOption := WithDecorator(&MockK8sDecorator{})

	c, err := New("eks", hostInfo, zap.NewNop(), cadvisorManagerCreator(newMockCreateManager(t)), k8sdecoratorOption)
	assert.NotNil(t, c)
	assert.NoError(t, err)
	assert.Nil(t, c.GetMetrics())
}

func TestGetMetricsErrorWhenCreatingManager(t *testing.T) {
	t.Setenv("HOST_NAME", "host")
	hostInfo := testutils.MockHostInfo{ClusterName: "cluster"}
	k8sdecoratorOption := WithDecorator(&MockK8sDecorator{})

	c, err := New("eks", hostInfo, zap.NewNop(), cadvisorManagerCreator(mockCreateManagerWithError), k8sdecoratorOption)
	assert.Nil(t, c)
	assert.Error(t, err)
}

func TestGetMetricsErrorWhenCallingManagerStart(t *testing.T) {
	t.Setenv("HOST_NAME", "host")
	hostInfo := testutils.MockHostInfo{ClusterName: "cluster"}
	k8sdecoratorOption := WithDecorator(&MockK8sDecorator{})

	c, err := New("eks", hostInfo, zap.NewNop(), cadvisorManagerCreator(mockCreateManager2), k8sdecoratorOption)
	assert.Nil(t, c)
	assert.Error(t, err)
}
