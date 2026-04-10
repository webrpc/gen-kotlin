package io.webrpc.client

import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val client = NodeTsExampleServiceClient(
        baseUrl = "http://localhost:3000",
        transport = OkHttpWebRpcTransport(),
    )

    println("Sending ping...")
    client.ping()

    println("Getting user 1...")
    val response = client.getUser(
        NodeTsExampleServiceApi.GetUser.Request(userId = 1UL),
    )
    println(response)

    println("Done.")
}
