CC=gcc
all:
	$(CC) -g -O2 -o bbfi bbfi.c

clean:
	rm bbfi
