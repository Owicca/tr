bin=/web.exe
dev_reload=CompileDaemon \
					 -build="go build -o /web.exe ./main.go" \
					 -command="/web.exe" \
					 -pattern="(.*)\.(go|json|sql)" \
					 -exclude-dir=./log/ \
					 -exclude="./web.exe" \
					 -polling=true \
					 -polling-interval=300 \
					 -color \
					 -verbose

dev_reload_no_pool=CompileDaemon \
					 -build="go build -o /web.exe ./main.go" \
					 -command="/web.exe" \
					 -pattern="(.*)\.(go|json|sql)" \
					 -exclude-dir=./log/ \
					 -exclude="./web.exe" \
					 -color \
					 -verbose

all: run

run: build
	$(bin)

clean:
	rm $(bin)

build:
	go build -o $(bin) main.go

dev_reload:
	$(dev_reload)

clean:
	-rm ./log/*
	-rm ./web.exe
