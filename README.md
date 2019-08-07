# cfmt

The cfmt program formats a comment in source code.  It reads from standard
input and writes to standard output.  It expects standard intput to start
with a comment.  It formats the comment and then passes the rest of the input
read through without any formatting.

To use it from vi, place your cursor on the first line of the comment you
wish to reformat and then call it with ```!}cfmt<ENTER<```.

If cfmt is not given a delimiter it will assume one of "-- ", "//" or "#" is
the delimiter.  The delimiter may be preceeded white space (tabs or spaces).

