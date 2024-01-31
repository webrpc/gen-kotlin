package io.webrpc.client

import io.ktor.client.*
import io.ktor.client.engine.cio.*
import kotlinx.coroutines.runBlocking

fun main(args: Array<String>) {
    val client = ExampleServiceClient(
        baseUrl = "http://localhost:3000",
        httpClientBuilder = { HttpClient(CIO) },
    )

    runBlocking {
        println("Sending ping…")
        client.ping()
    }

    runBlocking {
        println("Getting user 1…")
        client.getUser(1UL).let {
            println(it.toString())
        }
    }

    println("Done.")
}