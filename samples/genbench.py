#!/usr/bin/env python
# encoding: utf-8

import sys
import os.path
from random import Random
import shutil
random = Random(0) # initialise with seed to have reproductible benches

HELP_USAGE = """Usage: generate_libs.py root libs classes internal external.
    root      - Root directory where to create libs.
    libs      - Number of libraries (libraries only depend on those with smaller numbers)
    classes   - Number of classes per library
    internal  - Number of includes per file referring to that same library
    external  - Number of includes per file pointing to other libraries
    code_size - Number of code blocks in functions

For example:
./genbench.py /tmp/build 200 100 15 5 1000

To try the waf part, do:
waf configure build -p -j5

To test the autotools part, do:
touch README AUTHORS NEWS ChangeLog &&
autoreconf --install --symlink --verbose &&
mkdir autotools-build-dir &&
cd autotools-build-dir &&
../configure --disable-shared CXXFLAGS=-Wall &&
time make -j4 --silent &&
time make -j4 --silent
"""

def lib_name(i):
    return "lib_" + str(i)

def createHeader(name):
    filename = name + ".h"
    handle = open(filename, "w" )

    guard = name + '_h_'
    handle.write ('#ifndef ' + guard + '\n');
    handle.write ('#define ' + guard + '\n\n');

    handle.write ('class ' + name + ' {\n');
    handle.write ("""    int i;
    int arr[1024];
 """);
    handle.write ('public:\n');
    handle.write ("""    int ip;
    int arrp[1024];
 """);
    handle.write ('    ' + name + '();\n');
    handle.write ('    ~' + name + '();\n');
    handle.write ('};\n\n');

    handle.write ('#endif\n');


def createCPP(name, lib_number, classes_per_lib, internal_includes, external_includes, code_size):
    global CPPCONTENT
    filename = name + ".cpp"
    handle = open(filename, "w" )

    handle.write ('#include <stdio.h>\n')
    header= name + ".h"
    handle.write ('#include "' + header + '"\n');

    includes = random.sample(range(classes_per_lib), internal_includes)
    for i in includes:
        handle.write ('#include "class_' + str(i) + '.h"\n')

    if (lib_number > 0):
        includes = random.sample(range(classes_per_lib), external_includes)
        lib_list = range(lib_number)
        for i in includes:
            libname = 'lib_' + str(random.choice(lib_list))
            handle.write ('#include <' + libname + '/' + 'class_' + str(i) + '.h>\n')

    handle.write ('\n');
    handle.write (name + '::' + name + "() {\n");
    handle.write ('    int i;\n');
    for i in range(code_size):
        handle.write("""
    for(i=0; i<1024; i++) {
        printf("Hello World: %d", i);
    }
    \n""")
    handle.write ("}\n");
    handle.write (name + '::~' + name + "() {\n");
    handle.write ('    int i;\n');
    for i in range(code_size):
        handle.write("""
    for(i=0; i<1024; i++) {
        printf("Hello World: %d", i);
    }
    \n""")
    handle.write ("}\n");


def createSConscript(lib_number, classes):
    handle = open("SConscript", "w");
    handle.write("Import('env')\n")
    handle.write('list = Split("""\n');
    for i in range(classes):
        handle.write('    class_' + str(i) + '.cpp\n')
    handle.write('    """)\n\n')
    handle.write('env.StaticLibrary("lib_' + str(lib_number) + '", list)\n\n')

def createLibCMakeLists(lib_number, classes):
    handle = open("CMakeLists.txt", "w")
    handle.write("""add_library(lib_%s STATIC %s)\n""" % (str(lib_number), ' '.join(('class_%s' % str(i) for i in range(classes)))))

def createLibMakefile(lib_number, classes):
    handle = open("Makefile", "w");
    handle.write ("""COMPILER = g++
INC = -I..
CCFLAGS = -g -Wall $(INC)
ARCHIVE = ar
DEPEND = makedepend
.SUFFIXES: .o .cpp

""")
    handle.write ("lib = lib_" + str(lib_number) + ".a\n")
    handle.write ("src = \\\n")
    for i in range(classes):
        handle.write('class_' + str(i) + '.cpp \\\n')
    handle.write ("""

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

""")

def createLibJamFile(lib_number, classes):
    handle = open("Jamfile", "w")
    handle.write ("SubDir TOP lib_" + str(lib_number) + " ;\n\n")
    handle.write ("SubDirHdrs $(INCLUDES) ;\n\n")
    handle.write ("Library lib_" + str(lib_number) + " :\n")
    for i in range(classes):
        handle.write('    class_' + str(i) + '.cpp\n')
    handle.write ('    ;\n')

def createVCProjFile(lib_number, classes):
    handle = open("lib_" + str(lib_number) + ".vcproj", "w")
    handle.write("""<?xml version="1.0" encoding="Windows-1252"?>
<VisualStudioProject
	ProjectType="Visual C++"
	Version="7.10"
	Name=""" + '"' + lib_name(lib_number) + '"' + """
	ProjectGUID="{CF495178-8865-4D20-939D-AAA""" + str(lib_number) + """}"
	Keyword="Win32Proj">
	<Platforms>
		<Platform
			Name="Win32"/>
	</Platforms>
	<Configurations>
		<Configuration
			Name="Debug|Win32"
			OutputDirectory="Debug"
			IntermediateDirectory="Debug"
			ConfigurationType="4"
			CharacterSet="2">
			<Tool
				Name="VCCLCompilerTool"
				Optimization="0"
				PreprocessorDefinitions="WIN32;_DEBUG;_LIB"
                AdditionalIncludeDirectories=".."
				MinimalRebuild="TRUE"
				BasicRuntimeChecks="3"
				RuntimeLibrary="5"
				UsePrecompiledHeader="0"
				WarningLevel="3"
				Detect64BitPortabilityProblems="TRUE"
				DebugInformationFormat="4"/>
			<Tool
				Name="VCCustomBuildTool"/>
			<Tool
				Name="VCLibrarianTool"
				OutputFile="$(OutDir)/""" + lib_name(lib_number) + """.lib"/>
		</Configuration>
	</Configurations>
	<References>
	</References>
	<Files>
""")

    for i in range(classes):
        handle.write(r'  <File RelativePath=".\class_' + str(i) + '.cpp"/>\n')

    handle.write("""
	</Files>
	<Globals>
	</Globals>
</VisualStudioProject>
""")

def createLibrary(lib_number, classes, internal_includes, external_includes, code_size):
    name = "lib_" + str(lib_number)
    setDir(name)
    for i in range(classes):
        classname = "class_" + str(i)
        createHeader(classname)
        createCPP(classname, lib_number, classes, internal_includes, external_includes, code_size)
    createSConscript(lib_number, classes)
    createLibCMakeLists(lib_number, classes)
    createLibMakefile(lib_number, classes)
    createAutotools(lib_number, classes)

    os.chdir("..")

def createCMakeLists(libs):
    handle = open("CMakeLists.txt", "w")
    handle.write("""project('profiling-test')
cmake_minimum_required(VERSION 2.8)

include_directories(${CMAKE_SOURCE_DIR})
""")

    for i in range(libs):
        handle.write("""add_subdirectory(lib_%s)\n""" % str(i))

def createSConstruct(libs):
    handle = open("SConstruct", "w");
    handle.write("""env = Environment(CPPFLAGS=['-Wall'], CPPDEFINES=['LINUX'], CPPPATH=[Dir('#')])\n""")
    handle.write("""env.Decider('timestamp-newer')\n""")
    handle.write("""env.SetOption('implicit_cache', True)\n""")
    handle.write("""env.SourceCode('.', None)\n""")

    for i in range(libs):
        handle.write("""env.SConscript("lib_%s/SConscript", exports=['env'])\n""" % str(i))

def createFullMakefile(libs):
    handle = open("Makefile", "w")

    handle.write('subdirs = \\\n')
    for i in range(libs):
        handle.write('lib_' + str(i) + '\\\n')
    handle.write("""

.PHONY: all clean wiskcache ccache
all wiskcache ccache: $(subdirs)

clean:
	@for i in $(subdirs); do 	(cd $$i; $(MAKE) clean); done

.PHONY: $(subdirs)
$(subdirs):
	PREFIX=$(PREFIX) $(MAKE) -C $@ $(MAKECMDGOALS)

depend:
	@echo Skip depend

demo_clean:
	@rm -rf /nobackup/$(USER)/wisk-nfscache/*

demo:  demo_clean
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo Regular build
	@/usr/bin/time  --format "%U;%S;%E;%PCPU" $(MAKE) all
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Wiskcache build #1"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) wiskcache
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Wiskcache build #2"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) wiskcache
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Ccache build #1"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) ccache > /dev/null
	@echo ============================================
	@echo
	-@$(MAKE) clean >/dev/null 2>&1
	@echo
	@echo "Ccache build #2"
	@/usr/bin/time --format "%U;%S;%E;%PCPU"  $(MAKE) ccache > /dev/null
	@echo ============================================
""")

def createFullJamfile(libs):
    handle = open("Jamfile", "w")
    handle.write ("SubDir TOP ;\n\n")

    for i in range(libs):
        handle.write('SubInclude TOP ' + lib_name(i) + ' ;\n')

    handle = open("Jamrules", "w")
    handle.write('INCLUDES = $(TOP) ;\n')

WT = """#! /usr/bin/env python
# encoding: utf-8

VERSION = '0.0.2'
APPNAME = 'build_bench'
top  = '.'
out  = 'out'

def options(opt):
	opt.load('compiler_cxx')

def configure(conf):
	conf.load('compiler_cxx')

def build(bld):
	for i in range(%d):
		filez = ' '.join(['lib_%%d/class_%%d.cpp' %% (i, j) for j in range(%d)])
		bld.stlib(
			source = filez,
			target = 'lib_%%d' %% i,
			includes = '.', # include the top-level
		)
"""

def createWtop(libs, classes):
	f = open('wscript', 'w')
	f.write(WT % (libs, classes))
	f.close()

def createFullSolution(libs):
    handle = open("solution.sln", "w")
    handle.write("Microsoft Visual Studio Solution File, Format Version 8.00\n")

    for i in range(libs):
        project_name = lib_name(i) + '\\' + lib_name(i) + '.vcproj'
        handle.write('Project("{8BC9CEB8-8B4A-11D0-8D11-00A0C91BC942}") = "' + lib_name(i) +
                      '", "' + project_name + '", "{CF495178-8865-4D20-939D-AAA' + str(i) + '}"\n')
        handle.write('EndProject\n')

def createAutotoolsTop(libs):
    handle = open("configure.ac", "w")
    handle.write('''\
AC_INIT([bench], [1.0.0])
AC_CONFIG_AUX_DIR([autotools-aux])
AM_INIT_AUTOMAKE([subdir-objects nostdinc no-define tar-pax dist-bzip2])
AM_PROG_LIBTOOL
AC_CONFIG_HEADERS([config.h])
AC_CONFIG_FILES([Makefile])
AC_OUTPUT
''')

    handle = open("Makefile.am", "w")
    handle.write('''\
AM_CPPFLAGS = -I$(srcdir)
lib_LTLIBRARIES =
''')
    for i in range(libs): handle.write('include lib_%s/Makefile.am\n' % str(i))

def createAutotools(lib_number, classes):

    handle = open("Makefile.am", "w")
    handle.write('''\
lib_LTLIBRARIES += lib%s.la
lib%s_la_SOURCES =''' % (str(lib_number), str(lib_number)))
    for i in range(classes): handle.write(' lib_%s/class_%s.cpp' % (str(lib_number), str(i)))
    handle.write('\n')

def createExtraMakefiles():
    handle_pre = open("Makefile.pre", "w")
    handle_pre.write('''\
wiskcache:
	PREFIX=wiskcache  $(MAKE) all
ccache:
	PREFIX=ccache $(MAKE) all
    ''')
    handle_ccache = open("Makefile.ccache", "w")
    handle_ccache.write('''\
CCACHE=ccache  # Will updated as required
CACHE_PRE=$(CCACHE)
    ''')
    handle_wisk = open("Makefile.wiskcache", "w")
    handle_wisk.write('''\
WCACHE=/router/bin/wiskcache
CACHE_PRE=$(WCACHE) --
    ''')
    print "Be sure to confirm settings in"
    print "\to Makefile.ccache"
    print "\to Makefile.wiskcache"

def setDir(dir):
    if (not os.path.exists(dir)):
        os.mkdir(dir)
    os.chdir(dir)

def createWiskConfig(root_dir):
    wiskcacheroot = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
    os.makedirs(os.path.join(root_dir, "wisk/config"))
    shutil.copyfile(
        os.path.join(wiskcacheroot, "wisk/config/wiskcache_config.yaml"),
        os.path.join(root_dir, "wisk/config/wiskcache_config.yaml"))
    shutil.copyfile(
        os.path.join(wiskcacheroot, "wisk/config/wisktrack.ini"),
        os.path.join(root_dir, "wisk/config/wisktrack.ini"))
    print("Creating Config directories", wiskcacheroot)

def main(argv):
    global CPPCONTENT
    if len(argv) != 7:
        print(HELP_USAGE)
        return

    root_dir = argv[1]
    libs = int(argv[2])
    classes = int(argv[3])
    internal_includes = int(argv[4])
    external_includes = int(argv[5])
    code_size = int(argv[6])

    CPPCONTENT = open(os.path.join(os.path.dirname(sys.argv[0]), "cppcontent.template")).read()

    setDir(root_dir)
    for i in range(libs):
        createLibrary(i, classes, internal_includes, external_includes, code_size)

    createSConstruct(libs)
    createCMakeLists(libs)
    createFullMakefile(libs)
    createWtop(libs, classes)
    createAutotoolsTop(libs)
    createExtraMakefiles()
    createWiskConfig(root_dir)

if __name__ == "__main__":
    main( sys.argv )


