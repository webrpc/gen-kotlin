all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

bootstrap:
	rm -rf server/node_modules
	cd server && yarn

generate: generate-server generate-client

generate-server: bootstrap
	go -C ../../tools tool webrpc-gen -schema=../_examples/ClientExample/service.ridl -target=typescript -server -out=../_examples/ClientExample/server/server.gen.ts

generate-client:
	go -C ../../tools tool webrpc-gen -schema=../_examples/ClientExample/service.ridl -target=.. -client -okhttpTransport=true -out=../_examples/ClientExample/client/src/main/kotlin/ExampleServiceClient.kt
	cd client/ && ./gradlew build --no-daemon

run-server:
	yarn --cwd ./server start
	
run-client:
	cd client/ && ./gradlew run --no-daemon
