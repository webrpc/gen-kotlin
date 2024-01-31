webrpc node-ts
==============

* Server: Nodejs (TypeScript)
* Client: CLI (Kotlin JVM)

example of generating a webrpc client from [service.ridl](./service.ridl) schema.

## Usage

1. Install nodejs, yarn and webrpc-gen
2. $ `make bootstrap` -- runs yarn on ./server
3. $ `make generate` -- generates both server and client code
4. $ `make run-server`
5. $ `make run-client`
