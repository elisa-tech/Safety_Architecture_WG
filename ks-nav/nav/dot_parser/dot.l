%{
#include <stdio.h>
#include <string.h>
#include "dot.tab.h"

#define MAXSIZE  1024
char name_s[MAXSIZE];

%}

%option noinput
%option nounput

ws				[ \t\n]
anum				[a-zA-Z0-9_][a-zA-Z0-9_.]*
digit				[0-9]
A				[aA]
B				[bB]
C				[cC]
D				[dD]
E				[eE]
F				[fF]
G				[gG]
H				[hH]
I				[iI]
J				[jJ]
K				[kK]
L				[lL]
M				[mM]
N				[nN]
O				[oO]
P				[pP]
Q				[qQ]
R				[rR]
S				[sS]
T				[tT]
U				[uU]
V				[vV]
W				[wW]
X				[xX]
Y				[yY]
Z				[zZ]

%%
"//".*				;
\"(\\.|[^\\"])*\"		{ strcpy(name_s,yytext);return(T_STRING); }
{ws}+				;
"["				{ return(T_LBRACKET); }
"]"				{ return(T_RBRACKET); }
"("				{ return(T_LPAREN); }
")"				{ return(T_RPAREN); }
"{"				{ return(T_LBRACE); }
"}"				{ return(T_RBRACE); }
","				{ return(T_COMMA); }
":"				{ return(T_COLON); }
";"				{ return(T_SEMI); }
"->"				{ return(T_DEDGE); }
"--"				{ return(T_UEDGE); }
"="				{ return(T_EQ); }
{D}{I}{G}{R}{A}{P}{H}		{ return(T_DIGRAPH); }
{G}{R}{A}{P}{H}			{ return(T_GRAPH); }
{S}{U}{B}{G}{R}{A}{P}{H}	{ return(T_SUBGRAPH); }
{A}{T}				{ return(T_AT); }
{S}{T}{R}{I}{C}{T}		{ return(T_STRICT); }
{N}{O}{D}{E}			{ return(T_NODE); }
{E}{D}{G}{E}			{ return(T_EDGE); }
{anum}				{ strcpy(name_s, yytext);return(T_ID); }

%%


int yywrap(){
	return(1);
}
void set_input_string(const char* in) {
	yy_scan_string(in);
}

void end_lexical_scan(void) {
	yy_delete_buffer(YY_CURRENT_BUFFER);
}
