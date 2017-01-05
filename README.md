[![Build Status](https://api.travis-ci.org/intelsdi-x/snap-plugin-collector-kvm.svg)](https://travis-ci.org/intelsdi-x/snap-plugin-collector-kvm )
[![Go Report Card](http://goreportcard.com/badge/intelsdi-x/snap-plugin-collector-kvm)](http://goreportcard.com/report/intelsdi-x/snap-plugin-collector-kvm)

# Snap collector plugin - kvm

This plugin collects metrics from Linux KVM debug statistics.  

It's used in the [Snap framework](http://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
  Plugin collects specified metrics from linux kvm

### System Requirements
* [golang 1.7+](https://golang.org/dl/)  - needed only for building
* This Plugin compatible with kernel > 2.6
* Linux/x86_64

### Operating systems
* Linux/amd64


### Installation

#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture under the plugin's [release](https://github.com/intelsdi-x/snap-plugin-collector-kvm/releases) page. For Snap, check [here](https://github.com/intelsdi-x/snap/releases).

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-kvm
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-kvm.git
```

Build the snap kvm plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage

Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).

Load the plugin and create a task, see example in [Examples](#examples).

## Documentation

This collector gathers metrics from kvm.


### Collected Metrics

List of collected metrics is described in [METRICS.md](METRICS.md).

### Examples 

Example of running snap kvm collector and writing data to file.

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `snapteld -l 1 -t 0 &`

Download and load snap plugins:

```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-kvm/latest/linux/x86_64/snap-plugin-collector-kvm
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snaptel plugin load snap-plugin-collector-kvm
$ snaptel plugin load snap-plugin-publisher-file
```

Create a task manifest file  (exemplary files in [examples/tasks/] (examples/tasks/)):
```yaml
---
  version: 1
  schedule:
    type: "simple"
    interval: "1s"
  max-failures: 10
  workflow:
    collect:
      metrics:
        /intel/kvm/insn_emulation: {}
        /intel/kvm/insn_emulation_fail: {}
        /intel/kvm/invlpq: {}
        /intel/kvm/io_exits: {}
        /intel/kvm/irq_exits: {}
        /intel/kvm/irq_injections: {}
        /intel/kvm/irq_window: {}
        /intel/kvm/largepages: {}
        /intel/kvm/mmio_exits: {}
        /intel/kvm/mmu_cache_miss: {}
        /intel/kvm/mmu_flooded: {}
        /intel/kvm/mmu_pde_zapped: {}
        /intel/kvm/mmu_pte_updated: {}
        /intel/kvm/mmu_pte_write: {}
        /intel/kvm/mmu_recycled: {}
        /intel/kvm/mmu_shadow_zapped: {}
        /intel/kvm/mmu_unsync: {}
        /intel/kvm/nmi_injections: {}
        /intel/kvm/nmi_window: {}
        /intel/kvm/pf_fixed: {}
        /intel/kvm/pf_quest: {}
        /intel/kvm/remote_tlb_flush: {}
        /intel/kvm/request_irq: {}
        /intel/kvm/signal_exits: {}
        /intel/kvm/tlb_flush: {}
        /intel/kvm/efer_reload: {}
        /intel/kvm/exits: {}
        /intel/kvm/fpu_reload: {}
        /intel/kvm/halt_attempted_poll: {}
        /intel/kvm/halt_exits: {}
        /intel/kvm/halt_successful_poll: {}
        /intel/kvm/halt_wakeup: {}
        /intel/kvm/host_state_reload: {}
        /intel/kvm/hypercalls: {}

      publish:
        - plugin_name: "file"
          config:
            file: "/tmp/kvm_metrics"
```
Download an [example task file](https://github.com/intelsdi-x/snap-plugin-collector-kvm/blob/master/examples/tasks/) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-kvm/master/examples/tasks/kvm-file.yml
$ snaptel task create -t kvm-file.yml
Using task manifest to create task
Task created
ID: 480323af-15b0-4af8-a526-eb2ca6d8ae67
Name: Task-480323af-15b0-4af8-a526-eb2ca6d8ae67
State: Running
```

See realtime output from `snaptel task watch <task_id>` (CTRL+C to exit)
```
$ snaptel task watch 480323af-15b0-4af8-a526-eb2ca6d8ae67
```

This data is published to a file `/tmp/kvm_metrics` per task specification

Stop task:
```
$ snaptel task stop 480323af-15b0-4af8-a526-eb2ca6d8ae67
Task stopped:
ID: 480323af-15b0-4af8-a526-eb2ca6d8ae67
```



### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-kvm/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-kvm/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@Ramesh Raju](https://github.com/rraju2/)

This software has been contributed by MIKELANGELO, a Horizon 2020 project co-funded by the European Union. https://www.mikelangelo-project.eu/
## Thank You
And **thank you!** Your contribution, through code and participation, is incredibly important to us.
