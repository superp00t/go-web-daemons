# Go Webdaemons (WIP)

This library allows you to offload intensive and/or redundant Web page processing into a singular Go WebAssembly process, that can be shared between your page's tabs. It also aims to provide similar functionality in Node/Electron.

If that doesn't interest you, you can also just import the following isomorphic wasm packages to be used as you see fit.

- `github.com/superp00t/go-web-daemons/ws` for a WebSocket client
- `github.com/superp00t/go-web-daemons/db` for IndexedDB (WebAssembly) and BoltDB (native)
- `github.com/superp00t/go-web-daemons/service` for inter-process communication (Worker MessagePort in WebAssembly, local WebSocket server natively)


## Go service: example.go

`gowebd gen example.go -o dist`

```go

package main

import (
  "fmt"

  "github.com/superp00t/etc"
  "github.com/superp00t/go-web-daemons/service"
)

func main() {
  svc := service.New()

  // Generates a UUID
  svc.On("uuid-gen", func(q *service.Query) {
    q.Send(etc.GenerateRandomUUID().String())
  })

  svc.OnPort(func(port *service.Port) {
    fmt.Println("Port opened")
  })

  sv.Run()
}
```

## JS client

```js
import Webdaemon from "go-web-daemons";

var service = new Webdaemon("example");
service.on("load", onload);

async function onload() {
  var uuid = await service.q("uuid-gen"); 
  console.log(`Service returned ${str}`);
}

```