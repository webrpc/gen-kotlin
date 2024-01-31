import io.webrpc.client.*
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class TestApiClientTest {
    private val client = TestApiClient(
        baseUrl = "http://localhost:9988",
        httpClientBuilder = { HttpClient(CIO) },
    )

    @Test
    fun `Test GetEmpty`() {
        runBlocking {
            client.getEmpty()
        }
    }

    @Test
    fun `Test GetError`() {
        runBlocking {
            try {
                client.getError()
            } catch (e: WebRpcError) {
                assertEquals(
                    WebRpcError(
                        code = ErrorKind.WEBRPC_ENDPOINT.code,
                        error = "WebrpcEndpoint",
                        message = "endpoint error",
                        causeString = "internal error",
                        status = 400
                    ), e
                )
            }
        }
    }

    @Test
    fun `Test GetOne`() {
        runBlocking {
            assertEquals(
                simple, client.getOne()
            )
        }
    }

    @Test
    fun `Test SendOne`() {
        runBlocking {
            client.sendOne(
                one = simple,
            )
        }
    }

    @Test
    fun `Test GetMulti`() {
        runBlocking {
            assertEquals(
                multi, client.getMulti()
            )
        }
    }

    @Test
    fun `Test SendMulti`() {
        runBlocking {
            client.sendMulti(
                one = multi.one,
                two = multi.two,
                three = multi.three,
            )
        }
    }

    @Test
    fun `Test GetComplex`() {
        runBlocking {
            assertEquals(complex, client.getComplex())
        }
    }

    @Test
    fun `Test SendComplex`() {
        runBlocking {
            client.sendComplex(complex = complex)
        }
    }

    @Test
    fun `Test GetSchemaError`() {
        runBlocking {
            client.getSchemaError(code = -999)
        }
    }

    companion object {
        private val simple = Simple(
            id = 1, name = "one"
        )
        private val multi = TestApi.GetMultiResponse(
            one = Simple(
                id = 1, name = "one"
            ), two = Simple(
                id = 2, name = "two"
            ), three = Simple(
                id = 3, name = "three"
            )
        )
        private val complex = Complex(
            meta = mapOf("1" to "23", "2" to 24.0),
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
            enum = Status.AVAILABLE
        )
    }

}