COMPILER = g++
INC = -I..
CCFLAGS = -g -Wall $(INC)
ARCHIVE = ar
DEPEND = echo
.SUFFIXES: .o .cpp

lib = libX.a
# generated list
# src = \

objects = $(patsubst %.cpp, %.o, $(src))

all: depend $(lib)

$(lib): $(objects)
	@rm -f $@
	$(CACHE_PRE) $(ARCHIVE) cr $@ $^
	touch $@

.cpp.o:
	$(CACHE_PRE) $(COMPILER) $(CCFLAGS) -o $@ -c $<

clean:
	-@rm $(objects) $(lib) 2>/dev/null

depend:
	@echo skip depend

include ../Makefile.pre
-include ../Makefile.$(PREFIX)

# DO NOT DELETE
# list of generated dependencies
#
