default:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o valerius_linux_x64
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o valerius_linux_x86
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o valerius_linux_arm
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o valerius_windows_x64.exe
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o valerius_windows_x86.exe
clean:
	rm valerius_linux_x64
	rm valerius_linux_arm
	rm valerius_linux_x86
	rm valerius_windows_x64
	rm valerius_windows_x86
