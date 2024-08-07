![VIOLIN](images/VIOLIN.png)

VIOLIN worker pool / connection pool, rich APIs and configuration options are provided.

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/violin)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/violin)

## Install

```shell
go get github.com/B1NARY-GR0UP/violin
```

## Quick Start

### Worker Pool

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

### Connection Pool

```go
package main

import (
	"net"
	"time"

	"github.com/B1NARY-GR0UP/violin/cool"
)

func main() {
	producer := func() (net.Conn, error) {
		return net.Dial("your-network", "your-address")
	}
	c, _ := cool.New(5, 30, producer, cool.WithConnIdleTimeout(30*time.Second))
	defer c.Close()
	_ = c.Len()
	conn, _ := c.Get()
	_ = conn.Close()
	if cc, ok := conn.(*cool.Conn); ok {
		cc.MarkUnusable()
		if cc.IsUnusable() {
			_ = cc.Close()
		}
	}
}
```

## Configuration

### Worker Pool

| Option                  | Default           | Description                               |
|-------------------------|-------------------|-------------------------------------------|
| `WithMinWorkers`        | `0`               | Set the minimum number of workers         |
| `WithMaxWorkers`        | `5`               | Set the maximum number of workers         |
| `WithWorkerIdleTimeout` | `3 * time.Second` | Set the destroyed timeout of idle workers |

### Connection Pool

| Option                | Default | Description                     |
|-----------------------|---------|---------------------------------|
| `WithConnIdleTimeout` | `0`     | Set the connection idle timeout |

## Blogs

- [GO: How to Write a Worker Pool](https://dev.to/justlorain/go-how-to-write-a-worker-pool-1h3b) | [中文](https://juejin.cn/post/7244733519948333111)

## Credits

Sincere appreciation to the following repositories that made the development of VIOLIN possible.

- [gammazero/workerpool](https://github.com/gammazero/workerpool)
- [fatih/pool](https://github.com/fatih/pool)

## License

VIOLIN is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of VIOLIN are explained [here](./licenses).

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BMS.png" alt="BMS"/>
<br/><br/>
VIOLIN is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">Basic Middleware Service</a>
</p>