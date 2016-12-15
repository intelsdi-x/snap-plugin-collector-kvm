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
	"os"
	"path/filepath"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKvmPlugin(t *testing.T) {
	sysFs, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	sysFs = filepath.Join(sysFs, "sys")
	config := plugin.Config{
		"sys_path": "/sys",
	}
	Convey("Create Kvm Collector", t, func() {
		kvmCol := KvmCollector{}
		Convey("So Kvm should not be nil", func() {
			So(kvmCol, ShouldNotBeNil)
		})
		Convey("So Kvm should be of Kvm type", func() {
			So(kvmCol, ShouldHaveSameTypeAs, KvmCollector{})
		})
		Convey("kvmCol.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := kvmCol.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a plugin.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})

	Convey("Get Metric Kvm Types", t, func() {
		kvmCol := KvmCollector{}
		var cfg = plugin.Config{}
		metrics, err := kvmCol.GetMetricTypes(cfg)
		So(err, ShouldBeNil)
		So(len(metrics), ShouldResemble, 34)
	})

	Convey("Collect kvm debug Metrics", t, func() {
		kvmCol := KvmCollector{}
		mts := []plugin.Metric{}
		for _, v := range nsTypes.kvmMetricTypes {
			mts = append(mts, plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin).AddStaticElements(v), Config: config})
		}
		metrics, err := kvmCol.CollectMetrics(mts)
		So(err, ShouldBeNil)
		// this test should use mocked sys_path, because current the result depends on system where these tests are running
		// what is more, it will fail on Travis CLI
		So(len(metrics), ShouldResemble, 29)
		So(metrics[0].Data, ShouldNotBeNil)
	})
}
