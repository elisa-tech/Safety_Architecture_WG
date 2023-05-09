#include <stdio.h>
#include "dot.tab.h"

extern FILE *dotin;
extern int totparse(void);
int main () {
	dotin = stdin;
	return dotparse();
}
