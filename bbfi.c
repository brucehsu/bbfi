#include "bbfi.h"

int main(int argc, char **argv) {
	return interpret(); 
}

int interpret() {
	char c;
	Heap* heap = createHeap();
	while((c=fgetc(stdin))!=EOF) {
		switch(c) {
			case '>':
				heap = movePtrToNext(heap);
				break;
			case '<':
				heap = movePtrToPrev(heap);
				if(heap==0) return 13;
				break;
			case '+':
				heap = addPtr(heap);
				break;
			case '-':
				heap = subPtr(heap);
				break;
			case '.':
				heap = printPtr(heap);
				break;
			case ',':
				heap = inputPtr(heap);
				break;
			case '[':
				break;
		}
	}
	return 0;
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