bin=/app.exe
dev_reload=CompileDaemon \
					 -build="go build -o /app.exe ./cmd/tr/main.go" \
					 -command="/app.exe" \
					 -pattern="(.*)\.(go|json|sql)" \
					 -exclude-dir=./log/ \
					 -exclude="./app.exe" \
					 -polling=true \
					 -polling-interval=300 \
					 -color \
					 -verbose

dev_reload_no_pool=CompileDaemon \
					 -build="go build -o /app.exe ./cmd/tr/main.go" \
					 -command="/app.exe" \
					 -pattern="(.*)\.(go|json|sql)" \
					 -exclude-dir=./log/ \
					 -exclude="./app.exe" \
					 -color \
					 -verbose

all: run

run: build
	$(bin) ./config/config.json

clean:
	rm $(bin)

build:
	go build -o $(bin) ./cmd/tr/main.go

dev_reload:
	$(dev_reload)

clean:
	-rm ./log/*
	-rm ./app.exe
