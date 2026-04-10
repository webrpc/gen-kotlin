package kotlin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestSuccinctGeneration(t *testing.T) {
	schema := `
webrpc = v1

name = Demo
version = v1.0.0
basepath = /rpc

struct FlattenRequest
  - name: string
  - amount: uint64

struct FlattenResponse
  - id: uint64
  - counter: uint64

service Demo
  - Flatten(FlattenRequest) => (FlattenResponse)
  - SendPair(first: string, second: string) => (accepted: bool, note: string)
`

	output := generateKotlin(t, schema)

	requireRegexp(t, `suspend fun flatten\(request:\s*FlattenRequest\):\s*FlattenResponse`, output)
	requireRegexp(t, `fun encodeRequest\(request:\s*FlattenRequest,\s*json:\s*Json = WebRpcJson\):\s*String`, output)
	requireRegexp(t, `fun decodeResponse\(body:\s*String,\s*json:\s*Json = WebRpcJson\):\s*FlattenResponse`, output)

	flattenBlock := regexp.MustCompile(`object Flatten \{(?s:.*?)\n    \}\n\n    object SendPair`).FindString(output)
	if flattenBlock == "" {
		t.Fatalf("expected to find Flatten method block\n\n%s", output)
	}
	requireNotContains(t, flattenBlock, "data class Request")
	requireNotContains(t, flattenBlock, "data class Response")

	requireRegexp(t, `suspend fun sendPair\(request:\s*DemoApi\.SendPair\.Request\):\s*DemoApi\.SendPair\.Response`, output)
	requireRegexp(t, `object SendPair \{(?s:.*?)data class Request`, output)
	requireRegexp(t, `object SendPair \{(?s:.*?)data class Response`, output)
}

func TestTransportSplitGeneration(t *testing.T) {
	schema := `
webrpc = v1

name = Transport
version = v1.0.0
basepath = /rpc

service Transport
  - Ping()
`

	coreOutput := generateKotlin(t, schema)
	requireNotContains(t, coreOutput, "import okhttp3.OkHttpClient")
	requireNotContains(t, coreOutput, "class OkHttpWebRpcTransport(")
	requireContains(t, coreOutput, "private val transport: WebRpcTransport,")
	requireNotContains(t, coreOutput, "OkHttpWebRpcTransport()")

	coreProject := writeGradleProject(t, "transport-core", map[string]string{
		"src/main/kotlin/TransportClient.kt": coreOutput,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
	})
	runGradle(t, coreProject, "compileKotlin")

	okhttpOutput := generateKotlin(t, schema, "-okhttpTransport=true")
	requireContains(t, okhttpOutput, "import okhttp3.OkHttpClient")
	requireContains(t, okhttpOutput, "class OkHttpWebRpcTransport(")

	okhttpProject := writeGradleProject(t, "transport-okhttp", map[string]string{
		"src/main/kotlin/TransportClient.kt": okhttpOutput,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
		withOkHttp:        true,
	})
	runGradle(t, okhttpProject, "compileKotlin")
}

func TestHelperApiRuntime(t *testing.T) {
	schema := `
webrpc = v1

name = Helper
version = v1.0.0
basepath = /rpc

service Helper
  - GetUser(userId: uint64) => (code: uint32, username: string)

error 200 UserNotFound "user not found"
`

	output := generateKotlin(t, schema)

	project := writeGradleProject(t, "helper-api", map[string]string{
		"src/main/kotlin/HelperClient.kt": output,
		"src/test/kotlin/HelperApiRuntimeTest.kt": `
import io.webrpc.client.ErrorKind
import io.webrpc.client.HelperApi
import io.webrpc.client.WebRpcError
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import kotlin.test.Test
import kotlin.test.assertEquals

class HelperApiRuntimeTest {
    @Test
    fun helperRoundTripWorks() {
        val request = HelperApi.GetUser.Request(userId = 7UL)
        val body = HelperApi.GetUser.encodeRequest(request)
        val bodyJson = Json.parseToJsonElement(body).jsonObject
        assertEquals("7", bodyJson.getValue("userId").jsonPrimitive.content)

        val response = HelperApi.GetUser.decodeResponse("""{"code":200,"username":"alice"}""")
        assertEquals(200U, response.code)
        assertEquals("alice", response.username)

        assertEquals("/GetUser", HelperApi.GetUser.path)
        assertEquals("/rpc/Helper/GetUser", HelperApi.GetUser.urlPath)
    }

    @Test
    fun helperErrorDecodeWorks() {
        val error: WebRpcError = io.webrpc.client.decodeWebRpcError(
            statusCode = 404,
            body = """{"error":"UserNotFound","code":200,"msg":"user not found","cause":"","status":404}""",
        )

        assertEquals("UserNotFound", error.error)
        assertEquals(ErrorKind.USER_NOT_FOUND, error.errorKind)
        assertEquals("user not found", error.message)
        assertEquals(404, error.status)
    }
}
`,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
	})

	runGradle(t, project, "test")
}

func TestNullOptionalGeneration(t *testing.T) {
	schema := `
webrpc = v1

name = Nulls
version = v1.0.0
basepath = /rpc

struct MaybeNull
  - value?: null

service Nulls
  - Echo(value?: null) => (value?: null)
`

	output := generateKotlin(t, schema)
	requireNotContains(t, output, "??")

	project := writeGradleProject(t, "null-optional", map[string]string{
		"src/main/kotlin/NullsClient.kt": output,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
	})

	runGradle(t, project, "compileKotlin")
}

func TestServiceNameCollisionGeneration(t *testing.T) {
	schema := `
webrpc = v1

name = Collide
version = v1.0.0
basepath = /rpc

struct WalletApi
  - id: uint64

struct WalletServiceApi
  - id: uint64

struct WalletWebRpcApi
  - id: uint64

struct WalletClient
  - id: uint64

struct WalletServiceClient
  - id: uint64

struct WalletWebRpcClient
  - id: uint64

struct CollideWalletApi
  - id: uint64

struct CollideWalletServiceApi
  - id: uint64

struct CollideWalletWebRpcApi
  - id: uint64

struct CollideWalletClient
  - id: uint64

struct CollideWalletServiceClient
  - id: uint64

struct CollideWalletWebRpcClient
  - id: uint64

service Wallet
  - Ping()
`

	output := generateKotlin(t, schema)
	requireContains(t, output, "object CollideWalletGeneratedRpcApi")
	requireContains(t, output, "class CollideWalletGeneratedRpcClient(")

	project := writeGradleProject(t, "service-collision", map[string]string{
		"src/main/kotlin/CollideClient.kt": output,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
	})

	runGradle(t, project, "compileKotlin")
}

func TestSchemaAwareServiceNaming(t *testing.T) {
	waasSchema := `
webrpc = v1

name = waas
version = v1.0.0
basepath = /rpc

service Wallet
  - Ping()
`

	waasOutput := generateKotlin(t, waasSchema)
	requireContains(t, waasOutput, "object WaasWalletApi")
	requireContains(t, waasOutput, "class WaasWalletClient(")
	requireContains(t, waasOutput, `const val basePath: String = "/rpc/Wallet"`)

	testSchema := `
webrpc = v1

name = Test
version = v1.0.0
basepath = /rpc

service TestApi
  - Ping()
`

	testOutput := generateKotlin(t, testSchema)
	requireContains(t, testOutput, "object TestApi")
	requireContains(t, testOutput, "class TestApiClient(")
	requireNotContains(t, testOutput, "TestTestApi")
}

func TestMultiServiceGenerationSeparatesTopLevelDeclarations(t *testing.T) {
	schema := `
webrpc = v1

name = Foo
version = v1.0.0
basepath = /rpc

service Bar
  - Ping()

service Baz
  - Pong()
`

	output := generateKotlin(t, schema)
	requireContains(t, output, "object FooBarApi")
	requireContains(t, output, "class FooBarClient(")
	requireContains(t, output, "object FooBazApi")
	requireContains(t, output, "class FooBazClient(")
	requireNotContains(t, output, "}object ")

	project := writeGradleProject(t, "multi-service-formatting", map[string]string{
		"src/main/kotlin/FooClient.kt": output,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
	})

	runGradle(t, project, "compileKotlin")
}

func TestCrossServiceSchemaAwareNameCollisionGeneration(t *testing.T) {
	schema := `
webrpc = v1

name = Foo
version = v1.0.0
basepath = /rpc

service Bar
  - Ping()

service FooBar
  - Pong()
`

	output := generateKotlin(t, schema)
	requireContains(t, output, "object FooBarApi")
	requireContains(t, output, "class FooBarClient(")
	requireContains(t, output, "object FooBarServiceApi")
	requireContains(t, output, "class FooBarServiceClient(")

	project := writeGradleProject(t, "multi-service-collision", map[string]string{
		"src/main/kotlin/FooClient.kt": output,
	}, gradleDeps{
		withCoroutines:    true,
		withSerialization: true,
	})

	runGradle(t, project, "compileKotlin")
}

type gradleDeps struct {
	withCoroutines    bool
	withSerialization bool
	withOkHttp        bool
}

func generateKotlin(t *testing.T, schema string, extraArgs ...string) string {
	t.Helper()

	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "schema.ridl")
	outPath := filepath.Join(dir, "Client.kt")

	if err := os.WriteFile(schemaPath, []byte(strings.TrimSpace(schema)+"\n"), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	args := []string{
		"-schema=" + schemaPath,
		"-target=" + repoRoot(t),
		"-client",
	}
	args = append(args, extraArgs...)
	args = append(args, "-out="+outPath)

	runWebrpcGen(t, args...)

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read generated output: %v", err)
	}
	return string(content)
}

func writeGradleProject(t *testing.T, name string, files map[string]string, deps gradleDeps) string {
	t.Helper()

	dir := filepath.Join(t.TempDir(), name)
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatalf("mkdir project root: %v", err)
	}

	settings := fmt.Sprintf("rootProject.name = %q\n", name)
	if err := os.WriteFile(filepath.Join(dir, "settings.gradle.kts"), []byte(settings), 0o644); err != nil {
		t.Fatalf("write settings.gradle.kts: %v", err)
	}

	var build strings.Builder
	build.WriteString(`import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.8.21"
    kotlin("plugin.serialization") version "1.8.21"
}

repositories {
    mavenCentral()
}

dependencies {
`)
	if deps.withCoroutines {
		build.WriteString(`    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.7.3")
`)
	}
	if deps.withSerialization {
		build.WriteString(`    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.6.3")
`)
	}
	if deps.withOkHttp {
		build.WriteString(`    implementation("com.squareup.okhttp3:okhttp:4.12.0")
`)
	}
	build.WriteString(`    testImplementation(kotlin("test"))
}

tasks.test {
    useJUnitPlatform()
}

tasks.withType<JavaCompile> {
    options.encoding = "UTF-8"
    sourceCompatibility = "11"
    targetCompatibility = "11"
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "11"
}
`)

	if err := os.WriteFile(filepath.Join(dir, "build.gradle.kts"), []byte(build.String()), 0o644); err != nil {
		t.Fatalf("write build.gradle.kts: %v", err)
	}

	copyFile(t, filepath.Join(repoRoot(t), "Tests", "gradlew"), filepath.Join(dir, "gradlew"), 0o755)
	copyFile(t, filepath.Join(repoRoot(t), "Tests", "gradlew.bat"), filepath.Join(dir, "gradlew.bat"), 0o644)
	copyDir(t, filepath.Join(repoRoot(t), "Tests", "gradle"), filepath.Join(dir, "gradle"))

	for relPath, content := range files {
		target := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", target, err)
		}
		if err := os.WriteFile(target, []byte(strings.TrimSpace(content)+"\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", target, err)
		}
	}

	return dir
}

func runGradle(t *testing.T, projectDir string, task string) {
	t.Helper()

	cmd := exec.Command("./gradlew", "--no-daemon", task)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "GRADLE_USER_HOME="+filepath.Join(repoRoot(t), ".gradle-test-home"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gradle %s failed: %v\n%s", task, err, output)
	}
}

func runWebrpcGen(t *testing.T, args ...string) {
	t.Helper()

	cmdName, prefix, workDir := webrpcGenCommand(t)
	cmd := exec.Command(cmdName, append(prefix, args...)...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("webrpc-gen failed: %v\n%s", err, output)
	}
}

func webrpcGenCommand(t *testing.T) (string, []string, string) {
	t.Helper()

	if bin := os.Getenv("WEBRPC_GEN_BIN"); bin != "" {
		return bin, nil, repoRoot(t)
	}
	return "go", []string{"-C", filepath.Join(repoRoot(t), "tools"), "tool", "webrpc-gen"}, repoRoot(t)
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return dir
}

func copyFile(t *testing.T, src, dst string, mode os.FileMode) {
	t.Helper()

	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, mode); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	}); err != nil {
		t.Fatalf("copy dir %s -> %s: %v", src, dst, err)
	}
}

func requireContains(t *testing.T, text, needle string) {
	t.Helper()
	if !strings.Contains(text, needle) {
		t.Fatalf("expected output to contain %q\n\n%s", needle, text)
	}
}

func requireNotContains(t *testing.T, text, needle string) {
	t.Helper()
	if strings.Contains(text, needle) {
		t.Fatalf("expected output not to contain %q\n\n%s", needle, text)
	}
}

func requireRegexp(t *testing.T, expr, text string) {
	t.Helper()
	if !regexp.MustCompile(expr).MatchString(text) {
		t.Fatalf("expected output to match %q\n\n%s", expr, text)
	}
}

func requireNotRegexp(t *testing.T, expr, text string) {
	t.Helper()
	if regexp.MustCompile(expr).MatchString(text) {
		t.Fatalf("expected output not to match %q\n\n%s", expr, text)
	}
}
