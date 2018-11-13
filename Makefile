VERSION:=$(shell git describe --tags)
default:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/valerius_$(VERSION)_linux_x64
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o build/valerius_$(VERSION)_linux_x86
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o build/valerius_$(VERSION)_linux_arm
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/valerius_$(VERSION)_windows_x64.exe
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o build/valerius_$(VERSION)_windows_x86.exe
clean-exec:
	rm -r build/
