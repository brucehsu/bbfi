#include "bbfi.h"

int main(int argc, char **argv) {
	return interpret(); 
}

int interpret() {
	char c;
	while((c=fgetc(stdin))) {
		switch(c) {
			case '>':
				break;
			case '<':
				break;
			case '+':
				break;
			case '-':
				break;
			case '.':
				break;
			case ',':
				break;
			case '[':
				break;
		}
	}
	return 0;
}