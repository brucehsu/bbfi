#include "bbfi.h"

int main(int argc, char **argv) {
	Heap** heap = (Heap**) malloc(sizeof(Heap*));
	*heap = 0;
	InstCache** cache = (InstCache**) malloc(sizeof(InstCache*));
	*cache = 0;
	return interpret(heap,cache,-1); 
}

int interpret(Heap** heap, InstCache** cache, int loopBegin) {
	char c;
	if(*heap==0) *heap = createHeap();
	if(*cache==0) *cache = createInstCache();
	while((c=fgetc(stdin))!=EOF) {
		switch(c) {
			case '>':
			case '<':
			case '+':
			case '-':
			case '.':
			case ',':
				if(execute(heap, &c, 1)) return 13;
				*cache = addInst(c, *cache);
				if(c==',') addInst(*((*heap)->ptr), *cache);
				break;
			case '[':
				*cache = addInst(c, *cache);

				if(loopBegin==-1) loopBegin = strlen((*cache)->inst)-1;
				if(interpret(heap,cache,strlen((*cache)->inst)-1)) return 13;
				break;
			case ']':
				*cache = addInst(c, *cache);
				execute(heap, (*cache)->inst+loopBegin, strlen((*cache)->inst)-loopBegin);
				return 0;
				break;
		}
	}
	return 0;
}

char* execute(Heap** heap, char* str, int len) {
	char* sptr = str;
	int slen = len;
	while(slen--) {
		switch(*sptr) {
			case '>':
				*heap = movePtrToNext(*heap);
				break;
			case '<':
				*heap = movePtrToPrev(*heap);
				if(*heap==0) return 0;
				break;
			case '+':
				*heap = addPtr(*heap);
				break;
			case '-':
				*heap = subPtr(*heap);
				break;
			case '.':
				*heap = printPtr(*heap);
				break;
			case ',':
				if(slen==0) *heap = inputPtr(*heap);
				break;
			case '[':
				sptr = execute(heap, sptr+1, strlen(sptr));
				slen = strlen(sptr);
				break;
			case ']':
				if(*((*heap)->ptr)!=0) {
					sptr = str;
					slen = len;
					continue;
				} else {
					return sptr;
				}
				break;
			default:
				*((*heap)->ptr) = *sptr;
		}
		++sptr;
	}
	return 0;
}

InstCache* createInstCache() {
	InstCache *cache = (InstCache*) malloc(sizeof(InstCache));
	memset(cache, 0, sizeof(InstCache));
	cache->inst = (char*) malloc(sizeof(char) * CACHE_MAX);
	cache->cache_size = CACHE_MAX;
	memset(cache->inst, 0, cache->cache_size);
	return cache;
}
InstCache* expandInstCache(InstCache *cache) {
	InstCache *expanded_cache = (InstCache*) malloc(sizeof(InstCache));
	memset(expanded_cache, 0, sizeof(InstCache));
	expanded_cache->inst = (char*) malloc(sizeof(char)* (cache->cache_size*2));
	expanded_cache->cache_size = cache->cache_size*2;
	memset(cache->inst, 0, cache->cache_size);

	strncpy(expanded_cache->inst, cache->inst, cache->cache_size);

	free(cache->inst);
	free(cache);
	return expanded_cache;	
}

InstCache* addInst(char c, InstCache *cache) {
	if(strlen(cache->inst)==cache->cache_size-1) {
		cache = expandInstCache(cache);
	}
	cache->inst[strlen(cache->inst)] = c;
	return cache;
}

Heap* createHeap() {
	Heap* heap = (Heap*) malloc(sizeof(Heap));
	memset(heap,0,sizeof(Heap));
	heap->ptr = heap->obj;
	return heap;
}

Heap* expandHeap(Heap *currentHeap) {
	Heap* heap = createHeap();
	currentHeap->next = heap;
	heap->prev = currentHeap;
	return heap;
}

Heap* movePtrToNext(Heap* heap) {
	if(heap->currentIndex<HEAP_MAX-1) {
		++(heap->ptr);
		++(heap->currentIndex);
	} else {
		if(heap->next) {
			return heap->next;
		} else {
			heap = expandHeap(heap);
		}
	}
	return heap;
}

Heap* movePtrToPrev(Heap* heap) {
	if(heap->currentIndex==0) {
		if(heap->prev==0) return 0;
		heap = heap->prev;
	} else {
		--(heap->ptr);
		--(heap->currentIndex);
	}
	return heap;
}

Heap* addPtr(Heap* heap) {
	++(*(heap->ptr));
	return heap;
}

Heap* subPtr(Heap* heap) {
	--(*(heap->ptr));
	return heap;	
}

Heap* printPtr(Heap* heap) {
	fputc((int) *(heap->ptr), stdout);
	return heap;
}

Heap* inputPtr(Heap* heap) {
	*(heap->ptr) = fgetc(stdin);
	return heap;
}