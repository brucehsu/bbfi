#define HEAP_MAX 10000
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct Heap{
	char *ptr;
	int currentIndex;
	char obj[HEAP_MAX];
	struct Heap *prev, *next;
} Heap;

int main(int argc, char **argv);
int interpret();
Heap* expandHeap(Heap* currentHeap);
Heap* createHeap();
Heap* movePtrToNext(Heap* heap);
Heap* movePtrToPrev(Heap* heap);
Heap* addPtr(Heap* heap);
Heap* subPtr(Heap* heap);
Heap* printPtr(Heap* heap);
Heap* inputPtr(Heap* heap);
