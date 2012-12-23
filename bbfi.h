#define HEAP_MAX 10240
#define CACHE_MAX 1024
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct Heap{
	char *ptr;
	int currentIndex;
	char obj[HEAP_MAX];
	struct Heap *prev, *next;
} Heap;

typedef struct  {
	char *inst;
	int cache_size;
} InstCache;

int main(int argc, char **argv);
int interpret(Heap** heap, InstCache** cache, int loopBegin);
char* execute(Heap **heap, char *str, int len);

InstCache* createInstCache();
InstCache* expandInstCache(InstCache *cache);
InstCache* addInst(char c, InstCache *cache);

Heap* createHeap();
Heap* expandHeap(Heap* currentHeap);

Heap* movePtrToNext(Heap* heap);
Heap* movePtrToPrev(Heap* heap);
Heap* addPtr(Heap* heap);
Heap* subPtr(Heap* heap);
Heap* printPtr(Heap* heap);
Heap* inputPtr(Heap* heap);
