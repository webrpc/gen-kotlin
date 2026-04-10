# gen-kotlin

This repo contains the templates used by the `webrpc-gen` cli to code-generate
webrpc Kotlin client code.

This generator, from a webrpc schema/design file will code-generate:

1. Client -- a Kotlin client (via an injected transport, with optional provided OkHttp transport support and
`kotlinx.serialization`) to speak to a webrpc server using the
provided schema. This client is compatible with any webrpc server language (ie. Go, nodejs, etc.).

## Dependencies

In order to support communication with server, dependencies to few libraries must be provided.
This is a dependency of the generated code, so you must add it to your project.

Add this to `build.gradle.kts`:
```kotlin
plugins {
    kotlin("plugin.serialization") version "<your-kotlin-version>"
}

dependencies {
    val coroutinesVersion = "1.7.3"
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:$coroutinesVersion")

    val serializationVersion = "1.6.3"
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:$serializationVersion")
}
```

Generated clients depend on a `WebRpcTransport` abstraction.

If you want the generated file to also include the provided
`OkHttpWebRpcTransport`, generate with:

```sh
webrpc-gen -schema=example.ridl -target=kotlin -client -okhttpTransport=true -out=./example.gen.kt
```

and add:

```kotlin
dependencies {
    val okhttpVersion = "4.12.0"
    implementation("com.squareup.okhttp3:okhttp:$okhttpVersion")
}
```

Generated output also exposes low-level method helpers for custom flows:

- schema-aware service symbols, for example `WaasWalletApi` / `WaasWalletClient`
- `SchemaServiceApi.basePath`
- `SchemaServiceApi.Method.path`
- `SchemaServiceApi.Method.urlPath`
- `SchemaServiceApi.Method.encodeRequest(...)`
- `SchemaServiceApi.Method.decodeResponse(...)`

## Usage

```
webrpc-gen -schema=example.ridl -target=kotlin -client -out=./example.gen.kt
```

or 

```
webrpc-gen -schema=example.ridl -target=github.com/webrpc/gen-kotlin@latest -client -out=./example.gen.kt
```

or

```
webrpc-gen -schema=example.ridl -target=./local-templates-on-disk -client -out=./example.gen.kt
```

As you can see, the `-target` supports default `kotlin`, any git URI, or a local folder :)

## Tooling

This repo pins the published webrpc tool module in `tools/go.mod` using Go tool dependencies:

- `github.com/webrpc/webrpc v0.37.1`
- `tool github.com/webrpc/webrpc/cmd/webrpc-gen`
- `tool github.com/webrpc/webrpc/cmd/webrpc-test`

Use the pinned tools from this repo with:

```sh
go -C tools tool webrpc-gen
go -C tools tool webrpc-test
```

To update both to the latest published version and write that version back into `tools/go.mod`:

```sh
go -C tools get -tool github.com/webrpc/webrpc/cmd/webrpc-gen@latest github.com/webrpc/webrpc/cmd/webrpc-test@latest
```

### Set custom template variables
Change any of the following values by passing `-option="Value"` CLI flag to `webrpc-gen`.

| webrpc-gen -option              | Description                | Default value              |
|---------------------------------|----------------------------|----------------------------|
| `-client`                       | generate client code       | unset (`false`)            |
| `-okhttpTransport=%bool%`       | include optional OkHttp transport | `false`            |
| `-packageName=%package name%`   | define package name        | `io.webrpc.client`         |

## LICENSE

[MIT LICENSE](./LICENSE)
