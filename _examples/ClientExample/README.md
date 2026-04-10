webrpc node-ts
==============

* Server: Nodejs (TypeScript)
* Client: CLI (Kotlin JVM)

example of generating a webrpc client from [service.ridl](./service.ridl) schema.

## Usage

1. Install nodejs, yarn and Go 1.24+
2. $ `make bootstrap` -- runs yarn on ./server
3. $ `make generate` -- generates both server and client code via the pinned tool module in `../../tools`
4. $ `make run-server`
5. $ `make run-client`
