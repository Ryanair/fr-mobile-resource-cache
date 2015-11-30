buildit:
	go get && go build -ldflags "-w" -o bin/resource_sync
clean:
	rm -fr bin/*
buildclean: clean buildit
cleanbuild: clean buildit
test:
