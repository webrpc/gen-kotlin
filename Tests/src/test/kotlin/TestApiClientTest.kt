import io.webrpc.client.*
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertThrows
import org.junit.jupiter.api.Test
import kotlinx.serialization.json.JsonPrimitive

class TestApiClientTest {
    private val client = TestApiClient(
        baseUrl = "http://localhost:9988",
        transport = OkHttpWebRpcTransport(),
    )

    @Test
    fun `Test GetEmpty`() {
        runBlocking {
            client.getEmpty()
        }
    }

    @Test
    fun `Test GetError`() {
        val error = assertThrows(WebRpcError::class.java) {
            runBlocking {
                client.getError()
            }
        }

        assertEquals(
            WebRpcError(
                code = ErrorKind.WEBRPC_ENDPOINT.code,
                error = "WebrpcEndpoint",
                message = "endpoint error",
                causeString = "internal error",
                status = 400
            ), error
        )
    }

    @Test
    fun `Test GetOne`() {
        runBlocking {
            assertEquals(
                TestApiApi.GetOne.Response(one = simple),
                client.getOne(),
            )
        }
    }

    @Test
    fun `Test SendOne`() {
        runBlocking {
            client.sendOne(
                request = TestApiApi.SendOne.Request(one = simple),
            )
        }
    }

    @Test
    fun `Test GetMulti`() {
        runBlocking {
            assertEquals(
                multi,
                client.getMulti(),
            )
        }
    }

    @Test
    fun `Test SendMulti`() {
        runBlocking {
            client.sendMulti(
                request = TestApiApi.SendMulti.Request(
                    one = multi.one,
                    two = multi.two,
                    three = multi.three,
                ),
            )
        }
    }

    @Test
    fun `Test GetComplex`() {
        runBlocking {
            assertEquals(
                TestApiApi.GetComplex.Response(complex = complex),
                client.getComplex(),
            )
        }
    }

    @Test
    fun `Test SendComplex`() {
        runBlocking {
            client.sendComplex(
                request = TestApiApi.SendComplex.Request(complex = complex),
            )
        }
    }

    @Test
    fun `Test GetSchemaError`() {
        runBlocking {
            client.getSchemaError(
                request = TestApiApi.GetSchemaError.Request(code = -999),
            )
        }
    }

    companion object {
        private val simple = Simple(
            id = 1, name = "one"
        )
        private val multi = TestApiApi.GetMulti.Response(
            one = Simple(
                id = 1, name = "one"
            ), two = Simple(
                id = 2, name = "two"
            ), three = Simple(
                id = 3, name = "three"
            )
        )
        private val complex = Complex(
            meta = mapOf("1" to JsonPrimitive("23"), "2" to JsonPrimitive(24)),
            metaNestedExample = mapOf(
                "1" to mapOf(
                    "2" to 1U
                )
            ),
            namesList = listOf("John", "Alice", "Jakob"),
            numsList = listOf(1L, 2L, 3L, 4534643543L),
            doubleArray = listOf(listOf("testing"), listOf("api")),
            listOfMaps = listOf(
                mapOf(
                    "Jakob" to 251U, "alice" to 2U, "john" to 1U
                )
            ),
            listOfUsers = listOf(
                User(
                    id = 1UL, username = "John-Doe", role = "admin"
                )
            ),
            mapOfUsers = mapOf(
                "admin" to User(
                    id = 1UL, username = "John-Doe", role = "admin"
                )
            ),
            user = User(
                id = 1UL, username = "John-Doe", role = "admin"
            ),
            status = Status.AVAILABLE
        )
    }

}
