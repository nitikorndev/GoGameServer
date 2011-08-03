include $(GOROOT)/src/Make.inc

TARG=gogameserver
GOFILES=$(wildcard src/*.go)

include $(GOROOT)/src/Make.cmd
