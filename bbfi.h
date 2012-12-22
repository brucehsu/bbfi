#define HEAP_MAX 10000
#include <stdio.h>
#include <stdlib.h>

typedef struct heap{
	char *ptr;
	char obj[HEAP_MAX];
	struct heap *prev, *next;
} heap;

int main(int argc, char **argv);
int interpret();
heap* expandHeap(heap* currentHeap);
heap* createHeap();
char* movePtrToNext(heap* heap);
char* movePtrToPrev(heap* heap);
char* addPtr(heap* heap);
char* subPtr(heap* heap);
char* printPtr(heap* heap);
char* inputPtr(heap* heap);
