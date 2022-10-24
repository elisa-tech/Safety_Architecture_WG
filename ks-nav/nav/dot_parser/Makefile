#DEBUG=-g
CC=gcc
#CC=clang-13

all: dotf
dot.tab.h: dot.y
	bison -t -v -d -g --report=all -p "dot" dot.y
dot.tab.c: dot.y
	bison -t -v -d -g --report=all -p "dot" dot.y
lex.dot.c: dot.l
	flex -Pdot dot.l
dotf: lex.dot.c dot.tab.c dot.tab.h dotf.c
	$(CC) $(DEBUG) -o dotf lex.dot.c dot.tab.c dotf.c
dots: lex.dot.c dot.tab.c dot.tab.h dots.c
	$(CC) $(DEBUG) -o dots lex.dot.c dot.tab.c dots.c

libdotparser.so: lex.dot.c dot.tab.c dot.tab.h
	$(CC)  $(DEBUG) -shared -fPIC -o libdotparser.so lex.dot.c dot.tab.c
clean:
	rm -f dot.output dot.tab.c dot.tab.h dot.tab.o lex.dot.o dot.dot lex.dot.c dots dotf libdotparser.so
