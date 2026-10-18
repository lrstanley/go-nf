.DEFAULT_GOAL := generate

export NF_VERSION := 3.5.1
export NVIM_DEVICONS_VERSION := ec309a0eaff6ea967e58c275996155ab5c16938d

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

up:
	cd ./cmd/codegen && go get -u -t ./... && go mod tidy
	# cd ./_examples && go get -u -t ./... && go mod tidy
	go get -u -t ./... && go mod tidy

prepare:
	cd ./cmd/codegen && go mod tidy
	cd ./_examples && go mod tidy
	go mod tidy

generate: license
	rm -rf glyphs/ constants.gen.go
	cd ./cmd/codegen && go run . ../../
	gofmt -e -s -w glyphs/**/*.go constants.gen.go
	go test -v ./...
