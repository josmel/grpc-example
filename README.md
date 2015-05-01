# Simple gRPC service example

## Install

    $ go get github.com/sikang99/grpc-example

## Development

refer `Makefile` if you know how to build, use, and test.
	
	$ make 
	Makefile for grpc-example, by Stoney Kang, sikang99@gmail.com

	make [proto|build|run|kill|test]
   	- proto : compile interface spec
   	- build : compile client/server
   	- run   : start the server and exec client
   	- kill  : stop the server


compile IDL proto of gRPC.
	
	$ make proto

build client and server programs

	$ make build

run and test the service
	
	$ make run
	$ make test

## Usage

start to run the server and use any number of clients

	$ server/server &
	$ client/client add [name] [age]
	$ client/client get [id]
	$ client/client update [id]
	$ client/client delete [id]
	$ client/client list [age]

## History

- 2015/05/01 : list support search with condition optionally
- 2015/04/30 : start to code with mattn/grpc-example


## References

- [Bolt — an embedded key/value database for Go](https://www.progville.com/go/bolt-embedded-db-golang/)
- [Golang : How to reverse elements order in map?](https://www.socketloop.com/tutorials/golang-how-to-reverse-elements-order-in-map)
- [gRPC-JSON Proxy](http://yugui.jp/articles/889)
- [Protocol Buffers を利用した RPC、gRPC を golang から試してみ](http://mattn.kaoriya.net/software/lang/go/20150227144125.htm) 
- [mattn/grpc-example](https://github.com/mattn/grpc-example)


## License

MIT

