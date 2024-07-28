# promviz

`promviz` is a command-line tool to visualize the Prometheus metrics in CLI. We have forked prometheus and reused the code to create this new tool. The tool is available in the `cmd/promviz` directory alongside the `cmd/prometheus` (TSDB) and `cmd/promtool`. The tool is in `alpha` stage and is not yet ready for production use. The project was created by [@chtushar](https://github.com/chtushar) and [@pushkar-anand](https://github.com/pushkar-anand) as a part of the [FOSSHack 2024 hackathon](https://fossunited.org/fosshack/2024).

## Table of Contents

- [promviz](#promviz)
  - [Table of Contents](#table-of-contents)
- [Usage](#usage)
  - [Example](#example)
  - [Installation](#installation)
    - [From Binary](#from-binary)
- [Features](#features)
- [Roadmap](#roadmap)

# Usage

You can download the binary from the [releases](https://github.com/chtushar/prometheues-with-viz/releases) page. The tool is available for Linux, MacOS, and Windows. You can run the tool by executing the binary in the terminal.

```bash
./promviz [prometheus-url] [path-to-dashboard-json]
```

## Example

```bash
./promviz http://localhost:9090 /path/to/dashboard.json
```

## Installation

### From Binary

Download the binary for your operating system from the releases page.

Make the binary executable:

Run the binary:

```bash
./promviz [prometheus-url] [path-to-dashboard-json]
```

# Features

Currently, the tool supports the following panel types:

- [x] Gauge
- [x] Timeseries
- [x] Stat

Many more are in the roadmap

# Roadmap

Features that are planned to be added in the future:

- [ ] Feature to directly query the local storage for metrics.
- [ ] Add support to control the variables for the panels.
- [ ] Add support for changing the time and date range for the panels.

Panel type support:

- [ ] Bar graph
- [ ] Heatmap
- [ ] Table
- [ ] Text
