CC=gcc
all:
	$(CC) -O2 -o bbfi bbfi.c

clean:
	rm bbfi