all:
	scripts/build.sh

#dist:
#	scripts/dist.sh

qa: vet fmt lint test

lint:
	go get github.com/golang/lint/golint
	$(GOPATH)/bin/golint ./...

fmt:
	gofmt -s -w .

clean:
	rm -f bin/* || true
	rm -rf .gopath || true

test:
	go test -v ./...

vet:
	go vet ./...

here: build qa

build: deps
	go build -o bin/app-presets

deps: version-deps

version-deps: version.go ninjapack/root/opt/ninjablocks/apps/app-presets/package.json

version.go: pkgversion
	sed -i.bak "s/\"\([^\"]*\)\"/\"$$(cat pkgversion)\"/" version.go

ninjapack/root/opt/ninjablocks/apps/app-presets/package.json: pkgversion
	jq ".version = \"$$(cat pkgversion)\"" < ninjapack/root/opt/ninjablocks/apps/app-presets/package.json > ninjapack/root/opt/ninjablocks/apps/app-presets/package.json.tmp && \
	mv ninjapack/root/opt/ninjablocks/apps/app-presets/package.json.tmp ninjapack/root/opt/ninjablocks/apps/app-presets/package.json

.PHONY: all	dist clean test

