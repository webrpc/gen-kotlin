all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

bootstrap:
	rm -rf server/node_modules
	cd server && yarn

generate: generate-server generate-client

generate-server: bootstrap
	webrpc-gen -schema=service.ridl -target=typescript -server -out=./server/server.gen.ts

generate-client:
	webrpc-gen -schema=service.ridl -target=../../ -client -out=client/src/main/kotlin/ExampleServiceClient.kt
	cd client/ && ./gradlew build --no-daemon

run-server:
	yarn --cwd ./server start
	
run-client:
	cd client/ && ./gradlew run --no-daemon