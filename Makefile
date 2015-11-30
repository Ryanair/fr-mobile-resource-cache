buildit:
	./build
clean:
	rm -fr bin/*
buildclean: clean buildit
cleanbuild: clean buildit
test:
