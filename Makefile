default:
	go build -ldflags="-s -w" -o valerius_linux
	GOOS=windows go build -ldflags="-s -w" -o valerius_windows.exe
