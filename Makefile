#Build

build-mac-apple:
	GOOS=darwin GOARCH=arm64 go build -o bin/fold # Apple Silicon

build-mac-apple-intel:
	GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin # Apple amd64

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux app.go # 64-bit

list-platforms:
	go tool dist list

#Init

init-public:
	bin/fold --init --template public --dir examples/new

init-blanc:
	bin/fold --init --template blanc --dir examples/blanc

init-default:
	bin/fold --init --dir examples/example1

#Serve

run-new-example:
	bin/fold --dir examples/new

#Record

record-localhost:
	bin/fold --record bin/record-plan.json --dir examples/records/new

record-init:
	bin/fold --init --template blanc --record bin/record-plan.json --dir examples/records/new

test-record-run:
	bin/fold --dir examples/records/new
