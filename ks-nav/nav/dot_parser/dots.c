#include <stdio.h>
#include "dot.tab.h"

extern int totparse(void);
extern void set_input_string(const char* in);
extern void end_lexical_scan(void);
int main() {
	char *in="dgraph G {\na1  b3;\n}";
        set_input_string(in);
        int rv = dotparse();
        end_lexical_scan();
        return rv;
}
