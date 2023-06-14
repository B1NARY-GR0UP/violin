# VIOLIN

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/violin)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/violin)

> Please ... listen.

![VIOLIN](images/VIOLIN.png)

VIOLIN worker pool, rich APIs and configuration options are provided, inspired by [gammazero/workerpool](https://github.com/gammazero/workerpool).

## Install

```shell
go get github.com/B1NARY-GR0UP/violin
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/B1NARY-GR0UP/violin"
)

func main() {
	v := violin.New()
	defer v.Shutdown()
	v.Submit(func() {
		fmt.Println("Hello, VIOLIN!")
	})
}
```

## Configuration

| Option                  | Default           | Description                               |
|-------------------------|-------------------|-------------------------------------------|
| `WithMaxWorkers`        | `5`               | Set the maximum number of workers         |
| `WithWaitingQueueSize`  | `64`              | Set the size of the waiting queue         |
| `WithWorkerIdleTimeout` | `time.Second * 3` | Set the destroyed timeout of idle workers |

## License

VIOLIN is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of VIOLIN are explained [here](./licenses).

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BINARY-WEB-ECO.png" alt="BINARY-WEB-ECO"/>
<br/><br/>
VIOLIN is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">BINARY WEB ECOLOGY</a>
</p>