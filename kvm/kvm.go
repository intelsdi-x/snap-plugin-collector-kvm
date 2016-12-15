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

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	// Name of plugin
	Name = "kvm"
	// Vendor  prefix
	Vendor = "intel"
	// Plugin plugin name
	Plugin = "kvm"
	// Version of plugin
	Version = 1
	// Class of the collector
	Class = "sys"
)

const (
	hexPrefix = "0x"
	// sysPath source of data for metrics
	kvmPath = "kernel/debug/kvm"
)

var nsTypes = struct {
	kvmMetricTypes []string
}{
	// name of available metrics
	kvmMetricTypes: []string{"insn_emulation", "insn_emulation_fail", "invlpq", "io_exits", "irq_exits", "irq_injections", "irq_window",
		"largepages", "mmio_exits", "mmu_cache_miss", "mmu_flooded", "mmu_pde_zapped", "mmu_pte_updated", "mmu_pte_write", "mmu_recycled", "mmu_shadow_zapped", "mmu_unsync",
		"nmi_injections", "nmi_window", "pf_fixed", "pf_quest", "remote_tlb_flush", "request_irq", "signal_exits", "tlb_flush", "efer_reload", "exits", "fpu_reload", "halt_attempted_poll", "halt_exits", "halt_successful_poll", "halt_wakeup", "host_state_reload", "hypercalls"},
}

// KvmCollector type
type KvmCollector struct {
}

// GetConfigPolicy returns a config policy
func (KvmCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{""}, "sys_path", false, plugin.SetDefaultString("/sys"))

	return *policy, nil
}

// CollectMetrics returns collected metrics
func (KvmCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	sysPathConf, err := mts[0].Config.GetString("sys_path")
	if err != nil {
		return nil, err
	}

	sysPath := filepath.Join(sysPathConf, kvmPath)
	if _, err := os.Stat(sysPath); err != nil {
		return nil, err
	}

	for _, m := range mts {
		ns := m.Namespace
		eventName := m.Namespace.Strings()[len(ns.Strings())-1]
		metric, err := getEvent(sysPath, eventName, ns)
		if err != nil {
			// data for some kvm event might be unavailable, so log about that
			log.Warningf("Data for event `%s` is not available, err: %v", eventName, err)
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

//GetMetricTypes method makes values from Global config available
//returns namespaces for my metrics.
func (KvmCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {

	mts := []plugin.Metric{}

	for _, metric := range nsTypes.kvmMetricTypes {
		metric := plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin).AddStaticElement(metric)}
		mts = append(mts, metric)
	}
	return mts, nil
}

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
