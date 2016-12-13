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
	//Class of the collector
	Class = "sys"
)

var (
	//sysPath source of data for metrics
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
	configKey := []string{Vendor, Plugin}
	//policy.AddNewStringRule([]string{"intel", "kvm"},
	policy.AddNewStringRule(configKey, "sysPath", false, plugin.SetDefaultString("/sys"))

	//	"sysPath", false, plugin.SetDefaultString("/sys"))
	return *policy, nil
}

// CollectMetrics returns collected metrics
func (KvmCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	sysPathConf, err := mts[0].Config.GetString("sysPath")
	if err != nil {
		return nil, err
	}
	sysPath := filepath.Join(sysPathConf, kvmPath)
	if _, err := os.Stat(sysPath); os.IsNotExist(err) {
		return nil, err
	}
	for _, m := range mts {
		lastElement := len(m.Namespace.Strings()) - 1
		cnt := m.Namespace.Strings()[lastElement]
		metric, err := getEvents(sysPathConf, cnt, m.Namespace)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric...)
	}
	return metrics, nil
}

//GetMetricTypes method makes values from Global config available
//returns namespaes for my metrics.
func (KvmCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {

	mts := []plugin.Metric{}

	for _, metric := range nsTypes.kvmMetricTypes {
		metric := plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin).AddStaticElement(metric)}
		mts = append(mts, metric)
	}
	return mts, nil
}

func getEvents(sysPath string, eventName string, ns plugin.Namespace) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	filePath := filepath.Join(sysPath, kvmPath, eventName)
	cnt, err := readHex(filePath)
	if err != nil {
		return metrics, nil
	}
	metric := plugin.Metric{
		Namespace: ns,
		Data:      cnt,
	}
	metrics = append(metrics, metric)
	return metrics, nil
}

// Read Hex value
func readHex(filename string) (int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	// The int files that this is concerned with should only be one liners.
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return 0, err // if you return error
		}
		i := strings.TrimSpace(line)[2:] // 0x removed from Hex value
		number, err := strconv.ParseInt(i, 16, 0)
		if err != nil {
			return 0, nil
		}
		fmt.Println(number)
		return number, nil

	}
}
