FLAGS = -I include/
PROG  = BOA

main.out: objects/chunk.o objects/main.o
	gcc $^ ${FLAGS} -o bin/${PROG}

objects/chunk.o: src/chunk.c
	gcc -c src/chunk.c ${FLAGS} -o $@

objects/main.o: src/main.c
	gcc -c src/main.c ${FLAGS} -o $@

.PHONY: clean
clean:
	rm -rf objects/*.o bin/*

.PHONY: run
run:
	./bin/${PROG}

