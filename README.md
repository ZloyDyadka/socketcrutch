SocketCrutch
================================
A extendable Go WebSocket-framework 

## Quick start

### Installation
-------------------------
```
go get github.com/ZloyDyadka/socketcrutch
```

### Simple example
-------------------------

```go

import (
	"github.com/ZloyDyadka/socketcrutch"
	"github.com/ZloyDyadka/socketcrutch/codec"
	"log"
	"net/http"
)

func main() {
	server := socketcrutch.New()
	server.SetCodec(&codec.MsgPackCodec{})
	
	router := server.GetRouter()
	
	api := router.Group("v1")
	
	api.Route("test", func(session *socketcrutch.Session, data []byte) error {
	    log.Println("Received a test message")
	    
	    return nil
	})
	
	
	http.HandleFunc("/ws", server.ServeHTTP)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

## TODO
-------------------------
* Testing
* Ack message
* Client library
* Documentation