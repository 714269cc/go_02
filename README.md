# The Go Programming Language

Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.

![Gopher image](https://golang.org/doc/gopher/fiveyears.jpg)
*Gopher image by [Renee French][rf], licensed under [Creative Commons 3.0 Attributions license][cc3-by].*

Our canonical Git repository is located at https://go.googlesource.com/go.
There is a mirror of the repository at https://github.com/golang/go.

Unless otherwise noted, the Go source files are distributed under the
BSD-style license found in the LICENSE file.

### Download and Install

#### Binary Distributions

Official binary distributions are available at https://golang.org/dl/.

After downloading a binary release, visit https://golang.org/doc/install
for installation instructions.

#### Install From Source

If a binary distribution is not available for your combination of
operating system and architecture, visit
https://golang.org/doc/install/source
for source installation instructions.

### Contributing

Go is the work of thousands of contributors. We appreciate your help!

To contribute, please read the contribution guidelines at https://golang.org/doc/contribute.html.

Note that the Go project uses the issue tracker for bug reports and
proposals only. See https://golang.org/wiki/Questions for a list of
places to ask questions about the Go language.

[rf]: https://reneefrench.blogspot.com/
[cc3-by]: https://creativecommons.org/licenses/by/3.0/
Go
Command go
Start a bug report
Compile packages and dependencies
Remove object files and cached files
Show documentation for package or symbol
Print Go environment information
Update packages to use new APIs
Gofmt (reformat) package sources
Generate Go files by processing source
Add dependencies to current module and install them
Compile and install packages and dependencies
List packages or modules
Module maintenance
Download modules to local cache
Edit go.mod from tools or scripts
Print module requirement graph
Initialize new module in current directory
Add missing and remove unused modules
Make vendored copy of dependencies
Verify dependencies have expected content
Explain why packages or modules are needed
Compile and run Go program
Test packages
Run specified go tool
Print Go version
Report likely mistakes in packages
Build constraints
Build modes
Calling between Go and C
Build and test caching
Environment variables
File types
The go.mod file
GOPATH environment variable
GOPATH and Modules
Internal Directories
Vendor Directories
Legacy GOPATH go get
Module proxy protocol
Import path syntax
Relative import paths
Remote import paths
Import path checking
Modules, module versions, and more
Module authentication using go.sum
Package lists and patterns
Configuration for downloading non-public code
Testing flags
Testing functions
Controlling version control with GOVCS
Go is a tool for managing Go source code.

Usage:

go <command> [arguments]
The commands are:

bug         start a bug report
build       compile packages and dependencies
clean       remove object files and cached files
doc         show documentation for package or symbol
env         print Go environment information
fix         update packages to use new APIs
fmt         gofmt (reformat) package sources
generate    generate Go files by processing source
get         add dependencies to current module and install them
install     compile and install packages and dependencies
list        list packages or modules
mod         module maintenance
run         compile and run Go program
test        test packages
tool        run specified go tool
version     print Go version
vet         report likely mistakes in packages
Use "go help <command>" for more information about a command.

Additional help topics:

buildconstraint build constraints
buildmode       build modes
c               calling between Go and C
cache           build and test caching
environment     environment variables
filetype        file types
go.mod          the go.mod file
gopath          GOPATH environment variable
gopath-get      legacy GOPATH go get
goproxy         module proxy protocol
importpath      import path syntax
modules         modules, module versions, and more
module-get      module-aware go get
module-auth     module authentication using go.sum
packages        package lists and patterns
private         configuration for downloading non-public code
testflag        testing flags
testfunc        testing functions
vcs             controlling version control with GOVCS
Use "go help <topic>" for more information about that topic.

Start a bug report
Usage:

go bug
Bug opens the default browser and starts a new bug report. The report includes useful system information.

Compile packages and dependencies
Usage:

go build [-o output] [build flags] [packages]
Build compiles the packages named by the import paths, along with their dependencies, but it does not install the results.

If the arguments to build are a list of .go files from a single directory, build treats them as a list of source files specifying a single package.

When compiling packages, build ignores files that end in '_test.go'.

When compiling a single main package, build writes the resulting executable to an output file named after the first source file ('go build ed.go rx.go' writes 'ed' or 'ed.exe') or the source code directory ('go build unix/sam' writes 'sam' or 'sam.exe'). The '.exe' suffix is added when writing a Windows executable.

When compiling multiple packages or a single non-main package, build compiles the packages but discards the resulting object, serving only as a check that the packages can be built.

The -o flag forces build to write the resulting executable or object to the named output file or directory, instead of the default behavior described in the last two paragraphs. If the named output is an existing directory or ends with a slash or backslash, then any resulting executables will be written to that directory.

The -i flag installs the packages that are dependencies of the target. The -i flag is deprecated. Compiled packages are cached automatically.

The build flags are shared by the build, clean, get, install, list, run, and test commands:

-a
	force rebuilding of packages that are already up-to-date.
-n
	print the commands but do not run them.
-p n
	the number of programs, such as build commands or
	test binaries, that can be run in parallel.
	The default is the number of CPUs available.
-race
	enable data race detection.
	Supported only on linux/amd64, freebsd/amd64, darwin/amd64, windows/amd64,
	linux/ppc64le and linux/arm64 (only for 48-bit VMA).
-msan
	enable interoperation with memory sanitizer.
	Supported only on linux/amd64, linux/arm64
	and only with Clang/LLVM as the host C compiler.
	On linux/arm64, pie build mode will be used.
-v
	print the names of packages as they are compiled.
-work
	print the name of the temporary work directory and
	do not delete it when exiting.
-x
	print the commands.

-asmflags '[pattern=]arg list'
	arguments to pass on each go tool asm invocation.
-buildmode mode
	build mode to use. See 'go help buildmode' for more.
-compiler name
	name of compiler to use, as in runtime.Compiler (gccgo or gc).
-gccgoflags '[pattern=]arg list'
	arguments to pass on each gccgo compiler/linker invocation.
-gcflags '[pattern=]arg list'
	arguments to pass on each go tool compile invocation.
-installsuffix suffix
	a suffix to use in the name of the package installation directory,
	in order to keep output separate from default builds.
	If using the -race flag, the install suffix is automatically set to race
	or, if set explicitly, has _race appended to it. Likewise for the -msan
	flag. Using a -buildmode option that requires non-default compile flags
	has a similar effect.
-ldflags '[pattern=]arg list'
	arguments to pass on each go tool link invocation.
-linkshared
	build code that will be linked against shared libraries previously
	created with -buildmode=shared.
-mod mode
	module download mode to use: readonly, vendor, or mod.
	By default, if a vendor directory is present and the go version in go.mod
	is 1.14 or higher, the go command acts as if -mod=vendor were set.
	Otherwise, the go command acts as if -mod=readonly were set.
	See https://golang.org/ref/mod#build-commands for details.
-modcacherw
	leave newly-created directories in the module cache read-write
	instead of making them read-only.
-modfile file
	in module aware mode, read (and possibly write) an alternate go.mod
	file instead of the one in the module root directory. A file named
	"go.mod" must still be present in order to determine the module root
	directory, but it is not accessed. When -modfile is specified, an
	alternate go.sum file is also used: its path is derived from the
	-modfile flag by trimming the ".mod" extension and appending ".sum".
-overlay file
	read a JSON config file that provides an overlay for build operations.
	The file is a JSON struct with a single field, named 'Replace', that
	maps each disk file path (a string) to its backing file path, so that
	a build will run as if the disk file path exists with the contents
	given by the backing file paths, or as if the disk file path does not
	exist if its backing file path is empty. Support for the -overlay flag
	has some limitations:importantly, cgo files included from outside the
	include path must be  in the same directory as the Go package they are
	included from, and overlays will not appear when binaries and tests are
	run through go run and go test respectively.
-pkgdir dir
	install and load all packages from dir instead of the usual locations.
	For example, when building with a non-standard configuration,
	use -pkgdir to keep generated packages in a separate location.
-tags tag,list
	a comma-separated list of build tags to consider satisfied during the
	build. For more information about build tags, see the description of
	build constraints in the documentation for the go/build package.
	(Earlier versions of Go used a space-separated list, and that form
	is deprecated but still recognized.)
-trimpath
	remove all file system paths from the resulting executable.
	Instead of absolute file system paths, the recorded file names
	will begin with either "go" (for the standard library),
	or a module path@version (when using modules),
	or a plain import path (when using GOPATH).
-toolexec 'cmd args'
	a program to use to invoke toolchain programs like vet and asm.
	For example, instead of running asm, the go command will run
	'cmd args /path/to/asm <arguments for asm>'.
The -asmflags, -gccgoflags, -gcflags, and -ldflags flags accept a space-separated list of arguments to pass to an underlying tool during the build. To embed spaces in an element in the list, surround it with either single or double quotes. The argument list may be preceded by a package pattern and an equal sign, which restricts the use of that argument list to the building of packages matching that pattern (see 'go help packages' for a description of package patterns). Without a pattern, the argument list applies only to the packages named on the command line. The flags may be repeated with different patterns in order to specify different arguments for different sets of packages. If a package matches patterns given in multiple flags, the latest match on the command line wins. For example, 'go build -gcflags=-S fmt' prints the disassembly only for package fmt, while 'go build -gcflags=all=-S fmt' prints the disassembly for fmt and all its dependencies.

For more about specifying packages, see 'go help packages'. For more about where packages and binaries are installed, run 'go help gopath'. For more about calling between Go and C/C++, run 'go help c'.

Note: Build adheres to certain conventions such as those described by 'go help gopath'. Not all projects can follow these conventions, however. Installations that have their own conventions or that use a separate software build system may choose to use lower-level invocations such as 'go tool compile' and 'go tool link' to avoid some of the overheads and design decisions of the build tool.

See also: go install, go get, go clean.

Remove object files and cached files
Usage:

go clean [clean flags] [build flags] [packages]
Clean removes object files from package source directories. The go command builds most objects in a temporary directory, so go clean is mainly concerned with object files left by other tools or by manual invocations of go build.

If a package argument is given or the -i or -r flag is set, clean removes the following files from each of the source directories corresponding to the import paths:

_obj/            old object directory, left from Makefiles
_test/           old test directory, left from Makefiles
_testmain.go     old gotest file, left from Makefiles
test.out         old test log, left from Makefiles
build.out        old test log, left from Makefiles
*.[568ao]        object files, left from Makefiles

DIR(.exe)        from go build
DIR.test(.exe)   from go test -c
MAINFILE(.exe)   from go build MAINFILE.go
*.so             from SWIG
In the list, DIR represents the final path element of the directory, and MAINFILE is the base name of any Go source file in the directory that is not included when building the package.

The -i flag causes clean to remove the corresponding installed archive or binary (what 'go install' would create).

The -n flag causes clean to print the remove commands it would execute, but not run them.

The -r flag causes clean to be applied recursively to all the dependencies of the packages named by the import paths.

The -x flag causes clean to print remove commands as it executes them.

The -cache flag causes clean to remove the entire go build cache.

The -testcache flag causes clean to expire all test results in the go build cache.

The -modcache flag causes clean to remove the entire module download cache, including unpacked source code of versioned dependencies.

For more about build flags, see 'go help build'.

For more about specifying packages, see 'go help packages'.

Show documentation for package or symbol
Usage:

go doc [-u] [-c] [package|[package.]symbol[.methodOrField]]
Doc prints the documentation comments associated with the item identified by its arguments (a package, const, func, type, var, method, or struct field) followed by a one-line summary of each of the first-level items "under" that item (package-level declarations for a package, methods for a type, etc.).

Doc accepts zero, one, or two arguments.

Given no arguments, that is, when run as

go doc
it prints the package documentation for the package in the current directory. If the package is a command (package main), the exported symbols of the package are elided from the presentation unless the -cmd flag is provided.

When run with one argument, the argument is treated as a Go-syntax-like representation of the item to be documented. What the argument selects depends on what is installed in GOROOT and GOPATH, as well as the form of the argument, which is schematically one of these:

go doc <pkg>
go doc <sym>[.<methodOrField>]
go doc [<pkg>.]<sym>[.<methodOrField>]
go doc [<pkg>.][<sym>.]<methodOrField>
The first item in this list matched by the argument is the one whose documentation is printed. (See the examples below.) However, if the argument starts with a capital letter it is assumed to identify a symbol or method in the current directory.

For packages, the order of scanning is determined lexically in breadth-first order. That is, the package presented is the one that matches the search and is nearest the root and lexically first at its level of the hierarchy. The GOROOT tree is always scanned in its entirety before GOPATH.

If there is no package specified or matched, the package in the current directory is selected, so "go doc Foo" shows the documentation for symbol Foo in the current package.

The package path must be either a qualified path or a proper suffix of a path. The go tool's usual package mechanism does not apply: package path elements like . and ... are not implemented by go doc.

When run with two arguments, the first must be a full package path (not just a suffix), and the second is a symbol, or symbol with method or struct field. This is similar to the syntax accepted by godoc:

go doc <pkg> <sym>[.<methodOrField>]
In all forms, when matching symbols, lower-case letters in the argument match either case but upper-case letters match exactly. This means that there may be multiple matches of a lower-case argument in a package if different symbols have different cases. If this occurs, documentation for all matches is printed.

Examples:

go doc
	Show documentation for current package.
go doc Foo
	Show documentation for Foo in the current package.
	(Foo starts with a capital letter so it cannot match
	a package path.)
go doc encoding/json
	Show documentation for the encoding/json package.
go doc json
	Shorthand for encoding/json.
go doc json.Number (or go doc json.number)
	Show documentation and method summary for json.Number.
go doc json.Number.Int64 (or go doc json.number.int64)
	Show documentation for json.Number's Int64 method.
go doc cmd/doc
	Show package docs for the doc command.
go doc -cmd cmd/doc
	Show package docs and exported symbols within the doc command.
go doc template.new
	Show documentation for html/template's New function.
	(html/template is lexically before text/template)
go doc text/template.new # One argument
	Show documentation for text/template's New function.
go doc text/template new # Two arguments
	Show documentation for text/template's New function.

At least in the current tree, these invocations all print the
documentation for json.Decoder's Decode method:

go doc json.Decoder.Decode
go doc json.decoder.decode
go doc json.decode
cd go/src/encoding/json; go doc decode
Flags:

-all
	Show all the documentation for the package.
-c
	Respect case when matching symbols.
-cmd
	Treat a command (package main) like a regular package.
	Otherwise package main's exported symbols are hidden
	when showing the package's top-level documentation.
-short
	One-line representation for each symbol.
-src
	Show the full source code for the symbol. This will
	display the full Go source of its declaration and
	definition, such as a function definition (including
	the body), type declaration or enclosing const
	block. The output may therefore include unexported
	details.
-u
	Show documentation for unexported as well as exported
	symbols, methods, and fields.
Print Go environment information
Usage:

go env [-json] [-u] [-w] [var ...]
Env prints Go environment information.

By default env prints information as a shell script (on Windows, a batch file). If one or more variable names is given as arguments, env prints the value of each named variable on its own line.

The -json flag prints the environment in JSON format instead of as a shell script.

The -u flag requires one or more arguments and unsets the default setting for the named environment variables, if one has been set with 'go env -w'.

The -w flag requires one or more arguments of the form NAME=VALUE and changes the default settings of the named environment variables to the given values.

For more about environment variables, see 'go help environment'.

Update packages to use new APIs
Usage:

go fix [packages]
Fix runs the Go fix command on the packages named by the import paths.

For more about fix, see 'go doc cmd/fix'. For more about specifying packages, see 'go help packages'.

To run fix with specific options, run 'go tool fix'.

See also: go fmt, go vet.

Gofmt (reformat) package sources
Usage:

go fmt [-n] [-x] [packages]
Fmt runs the command 'gofmt -l -w' on the packages named by the import paths. It prints the names of the files that are modified.

For more about gofmt, see 'go doc cmd/gofmt'. For more about specifying packages, see 'go help packages'.

The -n flag prints commands that would be executed. The -x flag prints commands as they are executed.

The -mod flag's value sets which module download mode to use: readonly or vendor. See 'go help modules' for more.

To run gofmt with specific options, run gofmt itself.

See also: go fix, go vet.

Generate Go files by processing source
Usage:

go generate [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
Generate runs commands described by directives within existing files. Those commands can run any process but the intent is to create or update Go source files.

Go generate is never run automatically by go build, go get, go test, and so on. It must be run explicitly.

Go generate scans the file for directives, which are lines of the form,

//go:generate command argument...
(note: no leading spaces and no space in "//go") where command is the generator to be run, corresponding to an executable file that can be run locally. It must either be in the shell path (gofmt), a fully qualified path (/usr/you/bin/mytool), or a command alias, described below.

Note that go generate does not parse the file, so lines that look like directives in comments or multiline strings will be treated as directives.

The arguments to the directive are space-separated tokens or double-quoted strings passed to the generator as individual arguments when it is run.

Quoted strings use Go syntax and are evaluated before execution; a quoted string appears as a single argument to the generator.

To convey to humans and machine tools that code is generated, generated source should have a line that matches the following regular expression (in Go syntax):

^// Code generated .* DO NOT EDIT\.$
This line must appear before the first non-comment, non-blank text in the file.

Go generate sets several variables when it runs the generator:

$GOARCH
	The execution architecture (arm, amd64, etc.)
$GOOS
	The execution operating system (linux, windows, etc.)
$GOFILE
	The base name of the file.
$GOLINE
	The line number of the directive in the source file.
$GOPACKAGE
	The name of the package of the file containing the directive.
$DOLLAR
	A dollar sign.
Other than variable substitution and quoted-string evaluation, no special processing such as "globbing" is performed on the command line.

As a last step before running the command, any invocations of any environment variables with alphanumeric names, such as $GOFILE or $HOME, are expanded throughout the command line. The syntax for variable expansion is $NAME on all operating systems. Due to the order of evaluation, variables are expanded even inside quoted strings. If the variable NAME is not set, $NAME expands to the empty string.

A directive of the form,

//go:generate -command xxx args...
specifies, for the remainder of this source file only, that the string xxx represents the command identified by the arguments. This can be used to create aliases or to handle multiword generators. For example,

//go:generate -command foo go tool foo
specifies that the command "foo" represents the generator "go tool foo".

Generate processes packages in the order given on the command line, one at a time. If the command line lists .go files from a single directory, they are treated as a single package. Within a package, generate processes the source files in a package in file name order, one at a time. Within a source file, generate runs generators in the order they appear in the file, one at a time. The go generate tool also sets the build tag "generate" so that files may be examined by go generate but ignored during build.

For packages with invalid code, generate processes only source files with a valid package clause.

If any generator returns an error exit status, "go generate" skips all further processing for that package.

The generator is run in the package's source directory.

Go generate accepts one specific flag:

-run=""
	if non-empty, specifies a regular expression to select
	directives whose full original source text (excluding
	any trailing spaces and final newline) matches the
	expression.
It also accepts the standard build flags including -v, -n, and -x. The -v flag prints the names of packages and files as they are processed. The -n flag prints commands that would be executed. The -x flag prints commands as they are executed.

For more about build flags, see 'go help build'.

For more about specifying packages, see 'go help packages'.

Add dependencies to current module and install them
Usage:

go get [-d] [-t] [-u] [-v] [-insecure] [build flags] [packages]
Get resolves its command-line arguments to packages at specific module versions, updates go.mod to require those versions, downloads source code into the module cache, then builds and installs the named packages.

To add a dependency for a package or upgrade it to its latest version:

go get example.com/pkg
To upgrade or downgrade a package to a specific version:

go get example.com/pkg@v1.2.3
To remove a dependency on a module and downgrade modules that require it:

go get example.com/mod@none
See https://golang.org/ref/mod#go-get for details.

The 'go install' command may be used to build and install packages. When a version is specified, 'go install' runs in module-aware mode and ignores the go.mod file in the current directory. For example:

go install example.com/pkg@v1.2.3
go install example.com/pkg@latest
See 'go help install' or https://golang.org/ref/mod#go-install for details.

In addition to build flags (listed in 'go help build') 'go get' accepts the following flags.

The -t flag instructs get to consider modules needed to build tests of packages specified on the command line.

The -u flag instructs get to update modules providing dependencies of packages named on the command line to use newer minor or patch releases when available.

The -u=patch flag (not -u patch) also instructs get to update dependencies, but changes the default to select patch releases.

When the -t and -u flags are used together, get will update test dependencies as well.

The -insecure flag permits fetching from repositories and resolving custom domains using insecure schemes such as HTTP, and also bypassess module sum validation using the checksum database. Use with caution. This flag is deprecated and will be removed in a future version of go. To permit the use of insecure schemes, use the GOINSECURE environment variable instead. To bypass module sum validation, use GOPRIVATE or GONOSUMDB. See 'go help environment' for details.

The -d flag instructs get not to build or install packages. get will only update go.mod and download source code needed to build packages.

Building and installing packages with get is deprecated. In a future release, the -d flag will be enabled by default, and 'go get' will be only be used to adjust dependencies of the current module. To install a package using dependencies from the current module, use 'go install'. To install a package ignoring the current module, use 'go install' with an @version suffix like "@latest" after each argument.

For more about modules, see https://golang.org/ref/mod.

For more about specifying packages, see 'go help packages'.

This text describes the behavior of get using modules to manage source code and dependencies. If instead the go command is running in GOPATH mode, the details of get's flags and effects change, as does 'go help get'. See 'go help gopath-get'.

See also: go build, go install, go clean, go mod.

Compile and install packages and dependencies
Usage:

go install [build flags] [packages]
Install compiles and installs the packages named by the import paths.

Executables are installed in the directory named by the GOBIN environment variable, which defaults to $GOPATH/bin or $HOME/go/bin if the GOPATH environment variable is not set. Executables in $GOROOT are installed in $GOROOT/bin or $GOTOOLDIR instead of $GOBIN.

If the arguments have version suffixes (like @latest or @v1.0.0), "go install" builds packages in module-aware mode, ignoring the go.mod file in the current directory or any parent directory, if there is one. This is useful for installing executables without affecting the dependencies of the main module. To eliminate ambiguity about which module versions are used in the build, the arguments must satisfy the following constraints:

- Arguments must be package paths or package patterns (with "..." wildcards). They must not be standard packages (like fmt), meta-patterns (std, cmd, all), or relative or absolute file paths.

- All arguments must have the same version suffix. Different queries are not allowed, even if they refer to the same version.

- All arguments must refer to packages in the same module at the same version.

- No module is considered the "main" module. If the module containing packages named on the command line has a go.mod file, it must not contain directives (replace and exclude) that would cause it to be interpreted differently than if it were the main module. The module must not require a higher version of itself.

- Package path arguments must refer to main packages. Pattern arguments will only match main packages.

If the arguments don't have version suffixes, "go install" may run in module-aware mode or GOPATH mode, depending on the GO111MODULE environment variable and the presence of a go.mod file. See 'go help modules' for details. If module-aware mode is enabled, "go install" runs in the context of the main module.

When module-aware mode is disabled, other packages are installed in the directory $GOPATH/pkg/$GOOS_$GOARCH. When module-aware mode is enabled, other packages are built and cached but not installed.

The -i flag installs the dependencies of the named packages as well. The -i flag is deprecated. Compiled packages are cached automatically.

For more about the build flags, see 'go help build'. For more about specifying packages, see 'go help packages'.

See also: go build, go get, go clean.

List packages or modules
Usage:

go list [-f format] [-json] [-m] [list flags] [build flags] [packages]
List lists the named packages, one per line. The most commonly-used flags are -f and -json, which control the form of the output printed for each package. Other list flags, documented below, control more specific details.

The default output shows the package import path:

bytes
encoding/json
github.com/gorilla/mux
golang.org/x/net/html
The -f flag specifies an alternate format for the list, using the syntax of package template. The default output is equivalent to -f '{{.ImportPath}}'. The struct being passed to the template is:

type Package struct {
    Dir           string   // directory containing package sources
    ImportPath    string   // import path of package in dir
    ImportComment string   // path in import comment on package statement
    Name          string   // package name
    Doc           string   // package documentation string
    Target        string   // install path
    Shlib         string   // the shared library that contains this package (only set when -linkshared)
    Goroot        bool     // is this package in the Go root?
    Standard      bool     // is this package part of the standard Go library?
    Stale         bool     // would 'go install' do anything for this package?
    StaleReason   string   // explanation for Stale==true
    Root          string   // Go root or Go path dir containing this package
    ConflictDir   string   // this directory shadows Dir in $GOPATH
    BinaryOnly    bool     // binary-only package (no longer supported)
    ForTest       string   // package is only for use in named test
    Export        string   // file containing export data (when using -export)
    BuildID       string   // build ID of the compiled package (when using -export)
    Module        *Module  // info about package's containing module, if any (can be nil)
    Match         []string // command-line patterns matching this package
    DepOnly       bool     // package is only a dependency, not explicitly listed

    // Source files
    GoFiles         []string   // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
    CgoFiles        []string   // .go source files that import "C"
    CompiledGoFiles []string   // .go files presented to compiler (when using -compiled)
    IgnoredGoFiles  []string   // .go source files ignored due to build constraints
    IgnoredOtherFiles []string // non-.go source files ignored due to build constraints
    CFiles          []string   // .c source files
    CXXFiles        []string   // .cc, .cxx and .cpp source files
    MFiles          []string   // .m source files
    HFiles          []string   // .h, .hh, .hpp and .hxx source files
    FFiles          []string   // .f, .F, .for and .f90 Fortran source files
    SFiles          []string   // .s source files
    SwigFiles       []string   // .swig files
    SwigCXXFiles    []string   // .swigcxx files
    SysoFiles       []string   // .syso object files to add to archive
    TestGoFiles     []string   // _test.go files in package
    XTestGoFiles    []string   // _test.go files outside package

    // Embedded files
    EmbedPatterns      []string // //go:embed patterns
    EmbedFiles         []string // files matched by EmbedPatterns
    TestEmbedPatterns  []string // //go:embed patterns in TestGoFiles
    TestEmbedFiles     []string // files matched by TestEmbedPatterns
    XTestEmbedPatterns []string // //go:embed patterns in XTestGoFiles
    XTestEmbedFiles    []string // files matched by XTestEmbedPatterns

    // Cgo directives
    CgoCFLAGS    []string // cgo: flags for C compiler
    CgoCPPFLAGS  []string // cgo: flags for C preprocessor
    CgoCXXFLAGS  []string // cgo: flags for C++ compiler
    CgoFFLAGS    []string // cgo: flags for Fortran compiler
    CgoLDFLAGS   []string // cgo: flags for linker
    CgoPkgConfig []string // cgo: pkg-config names

    // Dependency information
    Imports      []string          // import paths used by this package
    ImportMap    map[string]string // map from source import to ImportPath (identity entries omitted)
    Deps         []string          // all (recursively) imported dependencies
    TestImports  []string          // imports from TestGoFiles
    XTestImports []string          // imports from XTestGoFiles

    // Error information
    Incomplete bool            // this package or a dependency has an error
    Error      *PackageError   // error loading package
    DepsErrors []*PackageError // errors loading dependencies
}
Packages stored in vendor directories report an ImportPath that includes the path to the vendor directory (for example, "d/vendor/p" instead of "p"), so that the ImportPath uniquely identifies a given copy of a package. The Imports, Deps, TestImports, and XTestImports lists also contain these expanded import paths. See golang.org/s/go15vendor for more about vendoring.

The error information, if any, is

type PackageError struct {
    ImportStack   []string // shortest path from package named on command line to this one
    Pos           string   // position of error (if present, file:line:col)
    Err           string   // the error itself
}
The module information is a Module struct, defined in the discussion of list -m below.

The template function "join" calls strings.Join.

The template function "context" returns the build context, defined as:

type Context struct {
    GOARCH        string   // target architecture
    GOOS          string   // target operating system
    GOROOT        string   // Go root
    GOPATH        string   // Go path
    CgoEnabled    bool     // whether cgo can be used
    UseAllFiles   bool     // use files regardless of +build lines, file names
    Compiler      string   // compiler to assume when computing target paths
    BuildTags     []string // build constraints to match in +build lines
    ReleaseTags   []string // releases the current release is compatible with
    InstallSuffix string   // suffix to use in the name of the install dir
}
For more information about the meaning of these fields see the documentation for the go/build package's Context type.

The -json flag causes the package data to be printed in JSON format instead of using the template format.

The -compiled flag causes list to set CompiledGoFiles to the Go source files presented to the compiler. Typically this means that it repeats the files listed in GoFiles and then also adds the Go code generated by processing CgoFiles and SwigFiles. The Imports list contains the union of all imports from both GoFiles and CompiledGoFiles.

The -deps flag causes list to iterate over not just the named packages but also all their dependencies. It visits them in a depth-first post-order traversal, so that a package is listed only after all its dependencies. Packages not explicitly listed on the command line will have the DepOnly field set to true.

The -e flag changes the handling of erroneous packages, those that cannot be found or are malformed. By default, the list command prints an error to standard error for each erroneous package and omits the packages from consideration during the usual printing. With the -e flag, the list command never prints errors to standard error and instead processes the erroneous packages with the usual printing. Erroneous packages will have a non-empty ImportPath and a non-nil Error field; other information may or may not be missing (zeroed).

The -export flag causes list to set the Export field to the name of a file containing up-to-date export information for the given package.

The -find flag causes list to identify the named packages but not resolve their dependencies: the Imports and Deps lists will be empty.

The -test flag causes list to report not only the named packages but also their test binaries (for packages with tests), to convey to source code analysis tools exactly how test binaries are constructed. The reported import path for a test binary is the import path of the package followed by a ".test" suffix, as in "math/rand.test". When building a test, it is sometimes necessary to rebuild certain dependencies specially for that test (most commonly the tested package itself). The reported import path of a package recompiled for a particular test binary is followed by a space and the name of the test binary in brackets, as in "math/rand [math/rand.test]" or "regexp [sort.test]". The ForTest field is also set to the name of the package being tested ("math/rand" or "sort" in the previous examples).

The Dir, Target, Shlib, Root, ConflictDir, and Export file paths are all absolute paths.

By default, the lists GoFiles, CgoFiles, and so on hold names of files in Dir (that is, paths relative to Dir, not absolute paths). The generated files added when using the -compiled and -test flags are absolute paths referring to cached copies of generated Go source files. Although they are Go source files, the paths may not end in ".go".

The -m flag causes list to list modules instead of packages.

When listing modules, the -f flag still specifies a format template applied to a Go struct, but now a Module struct:

type Module struct {
    Path      string       // module path
    Version   string       // module version
    Versions  []string     // available module versions (with -versions)
    Replace   *Module      // replaced by this module
    Time      *time.Time   // time version was created
    Update    *Module      // available update, if any (with -u)
    Main      bool         // is this the main module?
    Indirect  bool         // is this module only an indirect dependency of main module?
    Dir       string       // directory holding files for this module, if any
    GoMod     string       // path to go.mod file used when loading this module, if any
    GoVersion string       // go version used in module
    Retracted string       // retraction information, if any (with -retracted or -u)
    Error     *ModuleError // error loading module
}

type ModuleError struct {
    Err string // the error itself
}
The file GoMod refers to may be outside the module directory if the module is in the module cache or if the -modfile flag is used.

The default output is to print the module path and then information about the version and replacement if any. For example, 'go list -m all' might print:

my/main/module
golang.org/x/text v0.3.0 => /tmp/text
rsc.io/pdf v0.1.1
The Module struct has a String method that formats this line of output, so that the default format is equivalent to -f '{{.String}}'.

Note that when a module has been replaced, its Replace field describes the replacement module, and its Dir field is set to the replacement's source code, if present. (That is, if Replace is non-nil, then Dir is set to Replace.Dir, with no access to the replaced source code.)

The -u flag adds information about available upgrades. When the latest version of a given module is newer than the current one, list -u sets the Module's Update field to information about the newer module. list -u will also set the module's Retracted field if the current version is retracted. The Module's String method indicates an available upgrade by formatting the newer version in brackets after the current version. If a version is retracted, the string "(retracted)" will follow it. For example, 'go list -m -u all' might print:

my/main/module
golang.org/x/text v0.3.0 [v0.4.0] => /tmp/text
rsc.io/pdf v0.1.1 (retracted) [v0.1.2]
(For tools, 'go list -m -u -json all' may be more convenient to parse.)

The -versions flag causes list to set the Module's Versions field to a list of all known versions of that module, ordered according to semantic versioning, earliest to latest. The flag also changes the default output format to display the module path followed by the space-separated version list.

The -retracted flag causes list to report information about retracted module versions. When -retracted is used with -f or -json, the Retracted field will be set to a string explaining why the version was retracted. The string is taken from comments on the retract directive in the module's go.mod file. When -retracted is used with -versions, retracted versions are listed together with unretracted versions. The -retracted flag may be used with or without -m.

The arguments to list -m are interpreted as a list of modules, not packages. The main module is the module containing the current directory. The active modules are the main module and its dependencies. With no arguments, list -m shows the main module. With arguments, list -m shows the modules specified by the arguments. Any of the active modules can be specified by its module path. The special pattern "all" specifies all the active modules, first the main module and then dependencies sorted by module path. A pattern containing "..." specifies the active modules whose module paths match the pattern. A query of the form path@version specifies the result of that query, which is not limited to active modules. See 'go help modules' for more about module queries.

The template function "module" takes a single string argument that must be a module path or query and returns the specified module as a Module struct. If an error occurs, the result will be a Module struct with a non-nil Error field.

For more about build flags, see 'go help build'.

For more about specifying packages, see 'go help packages'.

For more about modules, see https://golang.org/ref/mod.

Module maintenance
Go mod provides access to operations on modules.

Note that support for modules is built into all the go commands, not just 'go mod'. For example, day-to-day adding, removing, upgrading, and downgrading of dependencies should be done using 'go get'. See 'go help modules' for an overview of module functionality.

Usage:

go mod <command> [arguments]
The commands are:

download    download modules to local cache
edit        edit go.mod from tools or scripts
graph       print module requirement graph
init        initialize new module in current directory
tidy        add missing and remove unused modules
vendor      make vendored copy of dependencies
verify      verify dependencies have expected content
why         explain why packages or modules are needed
Use "go help mod <command>" for more information about a command.

Download modules to local cache
Usage:

go mod download [-x] [-json] [modules]
Download downloads the named modules, which can be module patterns selecting dependencies of the main module or module queries of the form path@version. With no arguments, download applies to all dependencies of the main module (equivalent to 'go mod download all').

The go command will automatically download modules as needed during ordinary execution. The "go mod download" command is useful mainly for pre-filling the local cache or to compute the answers for a Go module proxy.

By default, download writes nothing to standard output. It may print progress messages and errors to standard error.

The -json flag causes download to print a sequence of JSON objects to standard output, describing each downloaded module (or failure), corresponding to this Go struct:

type Module struct {
    Path     string // module path
    Version  string // module version
    Error    string // error loading module
    Info     string // absolute path to cached .info file
    GoMod    string // absolute path to cached .mod file
    Zip      string // absolute path to cached .zip file
    Dir      string // absolute path to cached source root directory
    Sum      string // checksum for path, version (as in go.sum)
    GoModSum string // checksum for go.mod (as in go.sum)
}
The -x flag causes download to print the commands download executes.

See https://golang.org/ref/mod#go-mod-download for more about 'go mod download'.

See https://golang.org/ref/mod#version-queries for more about version queries.

Edit go.mod from tools or scripts
Usage:

go mod edit [editing flags] [go.mod]
Edit provides a command-line interface for editing go.mod, for use primarily by tools or scripts. It reads only go.mod; it does not look up information about the modules involved. By default, edit reads and writes the go.mod file of the main module, but a different target file can be specified after the editing flags.

The editing flags specify a sequence of editing operations.

The -fmt flag reformats the go.mod file without making other changes. This reformatting is also implied by any other modifications that use or rewrite the go.mod file. The only time this flag is needed is if no other flags are specified, as in 'go mod edit -fmt'.

The -module flag changes the module's path (the go.mod file's module line).

The -require=path@version and -droprequire=path flags add and drop a requirement on the given module path and version. Note that -require overrides any existing requirements on path. These flags are mainly for tools that understand the module graph. Users should prefer 'go get path@version' or 'go get path@none', which make other go.mod adjustments as needed to satisfy constraints imposed by other modules.

The -exclude=path@version and -dropexclude=path@version flags add and drop an exclusion for the given module path and version. Note that -exclude=path@version is a no-op if that exclusion already exists.

The -replace=old[@v]=new[@v] flag adds a replacement of the given module path and version pair. If the @v in old@v is omitted, a replacement without a version on the left side is added, which applies to all versions of the old module path. If the @v in new@v is omitted, the new path should be a local module root directory, not a module path. Note that -replace overrides any redundant replacements for old[@v], so omitting @v will drop existing replacements for specific versions.

The -dropreplace=old[@v] flag drops a replacement of the given module path and version pair. If the @v is omitted, a replacement without a version on the left side is dropped.

The -retract=version and -dropretract=version flags add and drop a retraction on the given version. The version may be a single version like "v1.2.3" or a closed interval like "[v1.1.0,v1.1.9]". Note that -retract=version is a no-op if that retraction already exists.

The -require, -droprequire, -exclude, -dropexclude, -replace, -dropreplace, -retract, and -dropretract editing flags may be repeated, and the changes are applied in the order given.

The -go=version flag sets the expected Go language version.

The -print flag prints the final go.mod in its text format instead of writing it back to go.mod.

The -json flag prints the final go.mod file in JSON format instead of writing it back to go.mod. The JSON output corresponds to these Go types:

type Module struct {
	Path string
	Version string
}

type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
	Replace []Replace
	Retract []Retract
}

type Require struct {
	Path string
	Version string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

type Retract struct {
	Low       string
	High      string
	Rationale string
}
Retract entries representing a single version (not an interval) will have the "Low" and "High" fields set to the same value.

Note that this only describes the go.mod file itself, not other modules referred to indirectly. For the full set of modules available to a build, use 'go list -m -json all'.

See https://golang.org/ref/mod#go-mod-edit for more about 'go mod edit'.

Print module requirement graph
Usage:

go mod graph
Graph prints the module requirement graph (with replacements applied) in text form. Each line in the output has two space-separated fields: a module and one of its requirements. Each module is identified as a string of the form path@version, except for the main module, which has no @version suffix.

See https://golang.org/ref/mod#go-mod-graph for more about 'go mod graph'.

Initialize new module in current directory
Usage:

go mod init [module]
Init initializes and writes a new go.mod file in the current directory, in effect creating a new module rooted at the current directory. The go.mod file must not already exist.

Init accepts one optional argument, the module path for the new module. If the module path argument is omitted, init will attempt to infer the module path using import comments in .go files, vendoring tool configuration files (like Gopkg.lock), and the current directory (if in GOPATH).

If a configuration file for a vendoring tool is present, init will attempt to import module requirements from it.

See https://golang.org/ref/mod#go-mod-init for more about 'go mod init'.

Add missing and remove unused modules
Usage:

go mod tidy [-e] [-v]
Tidy makes sure go.mod matches the source code in the module. It adds any missing modules necessary to build the current module's packages and dependencies, and it removes unused modules that don't provide any relevant packages. It also adds any missing entries to go.sum and removes any unnecessary ones.

The -v flag causes tidy to print information about removed modules to standard error.

The -e flag causes tidy to attempt to proceed despite errors encountered while loading packages.

See https://golang.org/ref/mod#go-mod-tidy for more about 'go mod tidy'.

Make vendored copy of dependencies
Usage:

go mod vendor [-e] [-v]
Vendor resets the main module's vendor directory to include all packages needed to build and test all the main module's packages. It does not include test code for vendored packages.

The -v flag causes vendor to print the names of vendored modules and packages to standard error.

The -e flag causes vendor to attempt to proceed despite errors encountered while loading packages.

See https://golang.org/ref/mod#go-mod-vendor for more about 'go mod vendor'.

Verify dependencies have expected content
Usage:

go mod verify
Verify checks that the dependencies of the current module, which are stored in a local downloaded source cache, have not been modified since being downloaded. If all the modules are unmodified, verify prints "all modules verified." Otherwise it reports which modules have been changed and causes 'go mod' to exit with a non-zero status.

See https://golang.org/ref/mod#go-mod-verify for more about 'go mod verify'.

Explain why packages or modules are needed
Usage:

go mod why [-m] [-vendor] packages...
Why shows a shortest path in the import graph from the main module to each of the listed packages. If the -m flag is given, why treats the arguments as a list of modules and finds a path to any package in each of the modules.

By default, why queries the graph of packages matched by "go list all", which includes tests for reachable packages. The -vendor flag causes why to exclude tests of dependencies.

The output is a sequence of stanzas, one for each package or module name on the command line, separated by blank lines. Each stanza begins with a comment line "# package" or "# module" giving the target package or module. Subsequent lines give a path through the import graph, one package per line. If the package or module is not referenced from the main module, the stanza will display a single parenthesized note indicating that fact.

For example:

$ go mod why golang.org/x/text/language golang.org/x/text/encoding
# golang.org/x/text/language
rsc.io/quote
rsc.io/sampler
golang.org/x/text/language

# golang.org/x/text/encoding
(main module does not need package golang.org/x/text/encoding)
$
See https://golang.org/ref/mod#go-mod-why for more about 'go mod why'.

Compile and run Go program
Usage:

go run [build flags] [-exec xprog] package [arguments...]
Run compiles and runs the named main Go package. Typically the package is specified as a list of .go source files from a single directory, but it may also be an import path, file system path, or pattern matching a single known package, as in 'go run .' or 'go run my/cmd'.

By default, 'go run' runs the compiled binary directly: 'a.out arguments...'. If the -exec flag is given, 'go run' invokes the binary using xprog:

'xprog a.out arguments...'.
If the -exec flag is not given, GOOS or GOARCH is different from the system default, and a program named go_$GOOS_$GOARCH_exec can be found on the current search path, 'go run' invokes the binary using that program, for example 'go_js_wasm_exec a.out arguments...'. This allows execution of cross-compiled programs when a simulator or other execution method is available.

The exit status of Run is not the exit status of the compiled binary.

For more about build flags, see 'go help build'. For more about specifying packages, see 'go help packages'.

See also: go build.

Test packages
Usage:

go test [build/test flags] [packages] [build/test flags & test binary flags]
'Go test' automates testing the packages named by the import paths. It prints a summary of the test results in the format:

ok   archive/tar   0.011s
FAIL archive/zip   0.022s
ok   compress/gzip 0.033s
...
followed by detailed output for each failed package.

'Go test' recompiles each package along with any files with names matching the file pattern "*_test.go". These additional files can contain test functions, benchmark functions, and example functions. See 'go help testfunc' for more. Each listed package causes the execution of a separate test binary. Files whose names begin with "_" (including "_test.go") or "." are ignored.

Test files that declare a package with the suffix "_test" will be compiled as a separate package, and then linked and run with the main test binary.

The go tool will ignore a directory named "testdata", making it available to hold ancillary data needed by the tests.

As part of building a test binary, go test runs go vet on the package and its test source files to identify significant problems. If go vet finds any problems, go test reports those and does not run the test binary. Only a high-confidence subset of the default go vet checks are used. That subset is: 'atomic', 'bool', 'buildtags', 'errorsas', 'ifaceassert', 'nilfunc', 'printf', and 'stringintconv'. You can see the documentation for these and other vet tests via "go doc cmd/vet". To disable the running of go vet, use the -vet=off flag.

All test output and summary lines are printed to the go command's standard output, even if the test printed them to its own standard error. (The go command's standard error is reserved for printing errors building the tests.)

Go test runs in two different modes:

The first, called local directory mode, occurs when go test is invoked with no package arguments (for example, 'go test' or 'go test -v'). In this mode, go test compiles the package sources and tests found in the current directory and then runs the resulting test binary. In this mode, caching (discussed below) is disabled. After the package test finishes, go test prints a summary line showing the test status ('ok' or 'FAIL'), package name, and elapsed time.

The second, called package list mode, occurs when go test is invoked with explicit package arguments (for example 'go test math', 'go test ./...', and even 'go test .'). In this mode, go test compiles and tests each of the packages listed on the command line. If a package test passes, go test prints only the final 'ok' summary line. If a package test fails, go test prints the full test output. If invoked with the -bench or -v flag, go test prints the full output even for passing package tests, in order to display the requested benchmark results or verbose logging. After the package tests for all of the listed packages finish, and their output is printed, go test prints a final 'FAIL' status if any package test has failed.

In package list mode only, go test caches successful package test results to avoid unnecessary repeated running of tests. When the result of a test can be recovered from the cache, go test will redisplay the previous output instead of running the test binary again. When this happens, go test prints '(cached)' in place of the elapsed time in the summary line.

The rule for a match in the cache is that the run involves the same test binary and the flags on the command line come entirely from a restricted set of 'cacheable' test flags, defined as -cpu, -list, -parallel, -run, -short, and -v. If a run of go test has any test or non-test flags outside this set, the result is not cached. To disable test caching, use any test flag or argument other than the cacheable flags. The idiomatic way to disable test caching explicitly is to use -count=1. Tests that open files within the package's source root (usually $GOPATH) or that consult environment variables only match future runs in which the files and environment variables are unchanged. A cached test result is treated as executing in no time at all, so a successful package test result will be cached and reused regardless of -timeout setting.

In addition to the build flags, the flags handled by 'go test' itself are:

-args
    Pass the remainder of the command line (everything after -args)
    to the test binary, uninterpreted and unchanged.
    Because this flag consumes the remainder of the command line,
    the package list (if present) must appear before this flag.

-c
    Compile the test binary to pkg.test but do not run it
    (where pkg is the last element of the package's import path).
    The file name can be changed with the -o flag.

-exec xprog
    Run the test binary using xprog. The behavior is the same as
    in 'go run'. See 'go help run' for details.

-i
    Install packages that are dependencies of the test.
    Do not run the test.
    The -i flag is deprecated. Compiled packages are cached automatically.

-json
    Convert test output to JSON suitable for automated processing.
    See 'go doc test2json' for the encoding details.

-o file
    Compile the test binary to the named file.
    The test still runs (unless -c or -i is specified).
The test binary also accepts flags that control execution of the test; these flags are also accessible by 'go test'. See 'go help testflag' for details.

For more about build flags, see 'go help build'. For more about specifying packages, see 'go help packages'.

See also: go build, go vet.

Run specified go tool
Usage:

go tool [-n] command [args...]
Tool runs the go tool command identified by the arguments. With no arguments it prints the list of known tools.

The -n flag causes tool to print the command that would be executed but not execute it.

For more about each tool command, see 'go doc cmd/<command>'.

Print Go version
Usage:

go version [-m] [-v] [file ...]
Version prints the build information for Go executables.

Go version reports the Go version used to build each of the named executable files.

If no files are named on the command line, go version prints its own version information.

If a directory is named, go version walks that directory, recursively, looking for recognized Go binaries and reporting their versions. By default, go version does not report unrecognized files found during a directory scan. The -v flag causes it to report unrecognized files.

The -m flag causes go version to print each executable's embedded module version information, when available. In the output, the module information consists of multiple lines following the version line, each indented by a leading tab character.

See also: go doc runtime/debug.BuildInfo.

Report likely mistakes in packages
Usage:

go vet [-n] [-x] [-vettool prog] [build flags] [vet flags] [packages]
Vet runs the Go vet command on the packages named by the import paths.

For more about vet and its flags, see 'go doc cmd/vet'. For more about specifying packages, see 'go help packages'. For a list of checkers and their flags, see 'go tool vet help'. For details of a specific checker such as 'printf', see 'go tool vet help printf'.

The -n flag prints commands that would be executed. The -x flag prints commands as they are executed.

The -vettool=prog flag selects a different analysis tool with alternative or additional checks. For example, the 'shadow' analyzer can be built and run using these commands:

go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
go vet -vettool=$(which shadow)
The build flags supported by go vet are those that control package resolution and execution, such as -n, -x, -v, -tags, and -toolexec. For more about these flags, see 'go help build'.

See also: go fmt, go fix.

Build constraints
A build constraint, also known as a build tag, is a line comment that begins

// +build
that lists the conditions under which a file should be included in the package. Constraints may appear in any kind of source file (not just Go), but they must appear near the top of the file, preceded only by blank lines and other line comments. These rules mean that in Go files a build constraint must appear before the package clause.

To distinguish build constraints from package documentation, a series of build constraints must be followed by a blank line.

A build constraint is evaluated as the OR of space-separated options. Each option evaluates as the AND of its comma-separated terms. Each term consists of letters, digits, underscores, and dots. A term may be negated with a preceding !. For example, the build constraint:

// +build linux,386 darwin,!cgo
corresponds to the boolean formula:

(linux AND 386) OR (darwin AND (NOT cgo))
A file may have multiple build constraints. The overall constraint is the AND of the individual constraints. That is, the build constraints:

// +build linux darwin
// +build amd64
corresponds to the boolean formula:

(linux OR darwin) AND amd64
During a particular build, the following words are satisfied:

- the target operating system, as spelled by runtime.GOOS, set with the
  GOOS environment variable.
- the target architecture, as spelled by runtime.GOARCH, set with the
  GOARCH environment variable.
- the compiler being used, either "gc" or "gccgo"
- "cgo", if the cgo command is supported (see CGO_ENABLED in
  'go help environment').
- a term for each Go major release, through the current version:
  "go1.1" from Go version 1.1 onward, "go1.12" from Go 1.12, and so on.
- any additional tags given by the -tags flag (see 'go help build').
There are no separate build tags for beta or minor releases.

If a file's name, after stripping the extension and a possible _test suffix, matches any of the following patterns:

*_GOOS
*_GOARCH
*_GOOS_GOARCH
(example: source_windows_amd64.go) where GOOS and GOARCH represent any known operating system and architecture values respectively, then the file is considered to have an implicit build constraint requiring those terms (in addition to any explicit constraints in the file).

Using GOOS=android matches build tags and files as for GOOS=linux in addition to android tags and files.

Using GOOS=illumos matches build tags and files as for GOOS=solaris in addition to illumos tags and files.

Using GOOS=ios matches build tags and files as for GOOS=darwin in addition to ios tags and files.

To keep a file from being considered for the build:

// +build ignore
(any other unsatisfied word will work as well, but "ignore" is conventional.)

To build a file only when using cgo, and only on Linux and OS X:

// +build linux,cgo darwin,cgo
Such a file is usually paired with another file implementing the default functionality for other systems, which in this case would carry the constraint:

// +build !linux,!darwin !cgo
Naming a file dns_windows.go will cause it to be included only when building the package for Windows; similarly, math_386.s will be included only when building the package for 32-bit x86.

Build modes
The 'go build' and 'go install' commands take a -buildmode argument which indicates which kind of object file is to be built. Currently supported values are:

-buildmode=archive
	Build the listed non-main packages into .a files. Packages named
	main are ignored.

-buildmode=c-archive
	Build the listed main package, plus all packages it imports,
	into a C archive file. The only callable symbols will be those
	functions exported using a cgo //export comment. Requires
	exactly one main package to be listed.

-buildmode=c-shared
	Build the listed main package, plus all packages it imports,
	into a C shared library. The only callable symbols will
	be those functions exported using a cgo //export comment.
	Requires exactly one main package to be listed.

-buildmode=default
	Listed main packages are built into executables and listed
	non-main packages are built into .a files (the default
	behavior).

-buildmode=shared
	Combine all the listed non-main packages into a single shared
	library that will be used when building with the -linkshared
	option. Packages named main are ignored.

-buildmode=exe
	Build the listed main packages and everything they import into
	executables. Packages not named main are ignored.

-buildmode=pie
	Build the listed main packages and everything they import into
	position independent executables (PIE). Packages not named
	main are ignored.

-buildmode=plugin
	Build the listed main packages, plus all packages that they
	import, into a Go plugin. Packages not named main are ignored.
On AIX, when linking a C program that uses a Go archive built with -buildmode=c-archive, you must pass -Wl,-bnoobjreorder to the C compiler.

Calling between Go and C
There are two different ways to call between Go and C/C++ code.

The first is the cgo tool, which is part of the Go distribution. For information on how to use it see the cgo documentation (go doc cmd/cgo).

The second is the SWIG program, which is a general tool for interfacing between languages. For information on SWIG see http://swig.org/. When running go build, any file with a .swig extension will be passed to SWIG. Any file with a .swigcxx extension will be passed to SWIG with the -c++ option.

When either cgo or SWIG is used, go build will pass any .c, .m, .s, .S or .sx files to the C compiler, and any .cc, .cpp, .cxx files to the C++ compiler. The CC or CXX environment variables may be set to determine the C or C++ compiler, respectively, to use.

Build and test caching
The go command caches build outputs for reuse in future builds. The default location for cache data is a subdirectory named go-build in the standard user cache directory for the current operating system. Setting the GOCACHE environment variable overrides this default, and running 'go env GOCACHE' prints the current cache directory.

The go command periodically deletes cached data that has not been used recently. Running 'go clean -cache' deletes all cached data.

The build cache correctly accounts for changes to Go source files, compilers, compiler options, and so on: cleaning the cache explicitly should not be necessary in typical use. However, the build cache does not detect changes to C libraries imported with cgo. If you have made changes to the C libraries on your system, you will need to clean the cache explicitly or else use the -a build flag (see 'go help build') to force rebuilding of packages that depend on the updated C libraries.

The go command also caches successful package test results. See 'go help test' for details. Running 'go clean -testcache' removes all cached test results (but not cached build results).

The GODEBUG environment variable can enable printing of debugging information about the state of the cache:

GODEBUG=gocacheverify=1 causes the go command to bypass the use of any cache entries and instead rebuild everything and check that the results match existing cache entries.

GODEBUG=gocachehash=1 causes the go command to print the inputs for all of the content hashes it uses to construct cache lookup keys. The output is voluminous but can be useful for debugging the cache.

GODEBUG=gocachetest=1 causes the go command to print details of its decisions about whether to reuse a cached test result.

Environment variables
The go command and the tools it invokes consult environment variables for configuration. If an environment variable is unset, the go command uses a sensible default setting. To see the effective setting of the variable <NAME>, run 'go env <NAME>'. To change the default setting, run 'go env -w <NAME>=<VALUE>'. Defaults changed using 'go env -w' are recorded in a Go environment configuration file stored in the per-user configuration directory, as reported by os.UserConfigDir. The location of the configuration file can be changed by setting the environment variable GOENV, and 'go env GOENV' prints the effective location, but 'go env -w' cannot change the default location. See 'go help env' for details.

General-purpose environment variables:

GO111MODULE
	Controls whether the go command runs in module-aware mode or GOPATH mode.
	May be "off", "on", or "auto".
	See https://golang.org/ref/mod#mod-commands.
GCCGO
	The gccgo command to run for 'go build -compiler=gccgo'.
GOARCH
	The architecture, or processor, for which to compile code.
	Examples are amd64, 386, arm, ppc64.
GOBIN
	The directory where 'go install' will install a command.
GOCACHE
	The directory where the go command will store cached
	information for reuse in future builds.
GOMODCACHE
	The directory where the go command will store downloaded modules.
GODEBUG
	Enable various debugging facilities. See 'go doc runtime'
	for details.
GOENV
	The location of the Go environment configuration file.
	Cannot be set using 'go env -w'.
GOFLAGS
	A space-separated list of -flag=value settings to apply
	to go commands by default, when the given flag is known by
	the current command. Each entry must be a standalone flag.
	Because the entries are space-separated, flag values must
	not contain spaces. Flags listed on the command line
	are applied after this list and therefore override it.
GOINSECURE
	Comma-separated list of glob patterns (in the syntax of Go's path.Match)
	of module path prefixes that should always be fetched in an insecure
	manner. Only applies to dependencies that are being fetched directly.
	Unlike the -insecure flag on 'go get', GOINSECURE does not disable
	checksum database validation. GOPRIVATE or GONOSUMDB may be used
	to achieve that.
GOOS
	The operating system for which to compile code.
	Examples are linux, darwin, windows, netbsd.
GOPATH
	For more details see: 'go help gopath'.
GOPROXY
	URL of Go module proxy. See https://golang.org/ref/mod#environment-variables
	and https://golang.org/ref/mod#module-proxy for details.
GOPRIVATE, GONOPROXY, GONOSUMDB
	Comma-separated list of glob patterns (in the syntax of Go's path.Match)
	of module path prefixes that should always be fetched directly
	or that should not be compared against the checksum database.
	See https://golang.org/ref/mod#private-modules.
GOROOT
	The root of the go tree.
GOSUMDB
	The name of checksum database to use and optionally its public key and
	URL. See https://golang.org/ref/mod#authenticating.
GOTMPDIR
	The directory where the go command will write
	temporary source files, packages, and binaries.
GOVCS
	Lists version control commands that may be used with matching servers.
	See 'go help vcs'.
Environment variables for use with cgo:

AR
	The command to use to manipulate library archives when
	building with the gccgo compiler.
	The default is 'ar'.
CC
	The command to use to compile C code.
CGO_ENABLED
	Whether the cgo command is supported. Either 0 or 1.
CGO_CFLAGS
	Flags that cgo will pass to the compiler when compiling
	C code.
CGO_CFLAGS_ALLOW
	A regular expression specifying additional flags to allow
	to appear in #cgo CFLAGS source code directives.
	Does not apply to the CGO_CFLAGS environment variable.
CGO_CFLAGS_DISALLOW
	A regular expression specifying flags that must be disallowed
	from appearing in #cgo CFLAGS source code directives.
	Does not apply to the CGO_CFLAGS environment variable.
CGO_CPPFLAGS, CGO_CPPFLAGS_ALLOW, CGO_CPPFLAGS_DISALLOW
	Like CGO_CFLAGS, CGO_CFLAGS_ALLOW, and CGO_CFLAGS_DISALLOW,
	but for the C preprocessor.
CGO_CXXFLAGS, CGO_CXXFLAGS_ALLOW, CGO_CXXFLAGS_DISALLOW
	Like CGO_CFLAGS, CGO_CFLAGS_ALLOW, and CGO_CFLAGS_DISALLOW,
	but for the C++ compiler.
CGO_FFLAGS, CGO_FFLAGS_ALLOW, CGO_FFLAGS_DISALLOW
	Like CGO_CFLAGS, CGO_CFLAGS_ALLOW, and CGO_CFLAGS_DISALLOW,
	but for the Fortran compiler.
CGO_LDFLAGS, CGO_LDFLAGS_ALLOW, CGO_LDFLAGS_DISALLOW
	Like CGO_CFLAGS, CGO_CFLAGS_ALLOW, and CGO_CFLAGS_DISALLOW,
	but for the linker.
CXX
	The command to use to compile C++ code.
FC
	The command to use to compile Fortran code.
PKG_CONFIG
	Path to pkg-config tool.
Architecture-specific environment variables:

GOARM
	For GOARCH=arm, the ARM architecture for which to compile.
	Valid values are 5, 6, 7.
GO386
	For GOARCH=386, how to implement floating point instructions.
	Valid values are sse2 (default), softfloat.
GOMIPS
	For GOARCH=mips{,le}, whether to use floating point instructions.
	Valid values are hardfloat (default), softfloat.
GOMIPS64
	For GOARCH=mips64{,le}, whether to use floating point instructions.
	Valid values are hardfloat (default), softfloat.
GOWASM
	For GOARCH=wasm, comma-separated list of experimental WebAssembly features to use.
	Valid values are satconv, signext.
Special-purpose environment variables:

GCCGOTOOLDIR
	If set, where to find gccgo tools, such as cgo.
	The default is based on how gccgo was configured.
GOROOT_FINAL
	The root of the installed Go tree, when it is
	installed in a location other than where it is built.
	File names in stack traces are rewritten from GOROOT to
	GOROOT_FINAL.
GO_EXTLINK_ENABLED
	Whether the linker should use external linking mode
	when using -linkmode=auto with code that uses cgo.
	Set to 0 to disable external linking mode, 1 to enable it.
GIT_ALLOW_PROTOCOL
	Defined by Git. A colon-separated list of schemes that are allowed
	to be used with git fetch/clone. If set, any scheme not explicitly
	mentioned will be considered insecure by 'go get'.
	Because the variable is defined by Git, the default value cannot
	be set using 'go env -w'.
Additional information available from 'go env' but not read from the environment:

GOEXE
	The executable file name suffix (".exe" on Windows, "" on other systems).
GOGCCFLAGS
	A space-separated list of arguments supplied to the CC command.
GOHOSTARCH
	The architecture (GOARCH) of the Go toolchain binaries.
GOHOSTOS
	The operating system (GOOS) of the Go toolchain binaries.
GOMOD
	The absolute path to the go.mod of the main module.
	If module-aware mode is enabled, but there is no go.mod, GOMOD will be
	os.DevNull ("/dev/null" on Unix-like systems, "NUL" on Windows).
	If module-aware mode is disabled, GOMOD will be the empty string.
GOTOOLDIR
	The directory where the go tools (compile, cover, doc, etc...) are installed.
GOVERSION
	The version of the installed Go tree, as reported by runtime.Version.
File types
The go command examines the contents of a restricted set of files in each directory. It identifies which files to examine based on the extension of the file name. These extensions are:

.go
	Go source files.
.c, .h
	C source files.
	If the package uses cgo or SWIG, these will be compiled with the
	OS-native compiler (typically gcc); otherwise they will
	trigger an error.
.cc, .cpp, .cxx, .hh, .hpp, .hxx
	C++ source files. Only useful with cgo or SWIG, and always
	compiled with the OS-native compiler.
.m
	Objective-C source files. Only useful with cgo, and always
	compiled with the OS-native compiler.
.s, .S, .sx
	Assembler source files.
	If the package uses cgo or SWIG, these will be assembled with the
	OS-native assembler (typically gcc (sic)); otherwise they
	will be assembled with the Go assembler.
.swig, .swigcxx
	SWIG definition files.
.syso
	System object files.
Files of each of these types except .syso may contain build constraints, but the go command stops scanning for build constraints at the first item in the file that is not a blank line or //-style line comment. See the go/build package documentation for more details.

The go.mod file
A module version is defined by a tree of source files, with a go.mod file in its root. When the go command is run, it looks in the current directory and then successive parent directories to find the go.mod marking the root of the main (current) module.

The go.mod file format is described in detail at https://golang.org/ref/mod#go-mod-file.

To create a new go.mod file, use 'go help init'. For details see 'go help mod init' or https://golang.org/ref/mod#go-mod-init.

To add missing module requirements or remove unneeded requirements, use 'go mod tidy'. For details, see 'go help mod tidy' or https://golang.org/ref/mod#go-mod-tidy.

To add, upgrade, downgrade, or remove a specific module requirement, use 'go get'. For details, see 'go help module-get' or https://golang.org/ref/mod#go-get.

To make other changes or to parse go.mod as JSON for use by other tools, use 'go mod edit'. See 'go help mod edit' or https://golang.org/ref/mod#go-mod-edit.

GOPATH environment variable
The Go path is used to resolve import statements. It is implemented by and documented in the go/build package.

The GOPATH environment variable lists places to look for Go code. On Unix, the value is a colon-separated string. On Windows, the value is a semicolon-separated string. On Plan 9, the value is a list.

If the environment variable is unset, GOPATH defaults to a subdirectory named "go" in the user's home directory ($HOME/go on Unix, %USERPROFILE%\go on Windows), unless that directory holds a Go distribution. Run "go env GOPATH" to see the current GOPATH.

See https://golang.org/wiki/SettingGOPATH to set a custom GOPATH.

Each directory listed in GOPATH must have a prescribed structure:

The src directory holds source code. The path below src determines the import path or executable name.

The pkg directory holds installed package objects. As in the Go tree, each target operating system and architecture pair has its own subdirectory of pkg (pkg/GOOS_GOARCH).

If DIR is a directory listed in the GOPATH, a package with source in DIR/src/foo/bar can be imported as "foo/bar" and has its compiled form installed to "DIR/pkg/GOOS_GOARCH/foo/bar.a".

The bin directory holds compiled commands. Each command is named for its source directory, but only the final element, not the entire path. That is, the command with source in DIR/src/foo/quux is installed into DIR/bin/quux, not DIR/bin/foo/quux. The "foo/" prefix is stripped so that you can add DIR/bin to your PATH to get at the installed commands. If the GOBIN environment variable is set, commands are installed to the directory it names instead of DIR/bin. GOBIN must be an absolute path.

Here's an example directory layout:

GOPATH=/home/user/go

/home/user/go/
    src/
        foo/
            bar/               (go code in package bar)
                x.go
            quux/              (go code in package main)
                y.go
    bin/
        quux                   (installed command)
    pkg/
        linux_amd64/
            foo/
                bar.a          (installed package object)
Go searches each directory listed in GOPATH to find source code, but new packages are always downloaded into the first directory in the list.

See https://golang.org/doc/code.html for an example.

GOPATH and Modules
When using modules, GOPATH is no longer used for resolving imports. However, it is still used to store downloaded source code (in GOPATH/pkg/mod) and compiled commands (in GOPATH/bin).

Internal Directories
Code in or below a directory named "internal" is importable only by code in the directory tree rooted at the parent of "internal". Here's an extended version of the directory layout above:

/home/user/go/
    src/
        crash/
            bang/              (go code in package bang)
                b.go
        foo/                   (go code in package foo)
            f.go
            bar/               (go code in package bar)
                x.go
            internal/
                baz/           (go code in package baz)
                    z.go
            quux/              (go code in package main)
                y.go
The code in z.go is imported as "foo/internal/baz", but that import statement can only appear in source files in the subtree rooted at foo. The source files foo/f.go, foo/bar/x.go, and foo/quux/y.go can all import "foo/internal/baz", but the source file crash/bang/b.go cannot.

See https://golang.org/s/go14internal for details.

Vendor Directories
Go 1.6 includes support for using local copies of external dependencies to satisfy imports of those dependencies, often referred to as vendoring.

Code below a directory named "vendor" is importable only by code in the directory tree rooted at the parent of "vendor", and only using an import path that omits the prefix up to and including the vendor element.

Here's the example from the previous section, but with the "internal" directory renamed to "vendor" and a new foo/vendor/crash/bang directory added:

/home/user/go/
    src/
        crash/
            bang/              (go code in package bang)
                b.go
        foo/                   (go code in package foo)
            f.go
            bar/               (go code in package bar)
                x.go
            vendor/
                crash/
                    bang/      (go code in package bang)
                        b.go
                baz/           (go code in package baz)
                    z.go
            quux/              (go code in package main)
                y.go
The same visibility rules apply as for internal, but the code in z.go is imported as "baz", not as "foo/vendor/baz".

Code in vendor directories deeper in the source tree shadows code in higher directories. Within the subtree rooted at foo, an import of "crash/bang" resolves to "foo/vendor/crash/bang", not the top-level "crash/bang".

Code in vendor directories is not subject to import path checking (see 'go help importpath').

When 'go get' checks out or updates a git repository, it now also updates submodules.

Vendor directories do not affect the placement of new repositories being checked out for the first time by 'go get': those are always placed in the main GOPATH, never in a vendor subtree.

See https://golang.org/s/go15vendor for details.

Legacy GOPATH go get
The 'go get' command changes behavior depending on whether the go command is running in module-aware mode or legacy GOPATH mode. This help text, accessible as 'go help gopath-get' even in module-aware mode, describes 'go get' as it operates in legacy GOPATH mode.

Usage: go get [-d] [-f] [-t] [-u] [-v] [-fix] [-insecure] [build flags] [packages]

Get downloads the packages named by the import paths, along with their dependencies. It then installs the named packages, like 'go install'.

The -d flag instructs get to stop after downloading the packages; that is, it instructs get not to install the packages.

The -f flag, valid only when -u is set, forces get -u not to verify that each package has been checked out from the source control repository implied by its import path. This can be useful if the source is a local fork of the original.

The -fix flag instructs get to run the fix tool on the downloaded packages before resolving dependencies or building the code.

The -insecure flag permits fetching from repositories and resolving custom domains using insecure schemes such as HTTP. Use with caution. This flag is deprecated and will be removed in a future version of go. The GOINSECURE environment variable should be used instead, since it provides control over which packages may be retrieved using an insecure scheme. See 'go help environment' for details.

The -t flag instructs get to also download the packages required to build the tests for the specified packages.

The -u flag instructs get to use the network to update the named packages and their dependencies. By default, get uses the network to check out missing packages but does not use it to look for updates to existing packages.

The -v flag enables verbose progress and debug output.

Get also accepts build flags to control the installation. See 'go help build'.

When checking out a new package, get creates the target directory GOPATH/src/<import-path>. If the GOPATH contains multiple entries, get uses the first one. For more details see: 'go help gopath'.

When checking out or updating a package, get looks for a branch or tag that matches the locally installed version of Go. The most important rule is that if the local installation is running version "go1", get searches for a branch or tag named "go1". If no such version exists it retrieves the default branch of the package.

When go get checks out or updates a Git repository, it also updates any git submodules referenced by the repository.

Get never checks out or updates code stored in vendor directories.

For more about specifying packages, see 'go help packages'.

For more about how 'go get' finds source code to download, see 'go help importpath'.

This text describes the behavior of get when using GOPATH to manage source code and dependencies. If instead the go command is running in module-aware mode, the details of get's flags and effects change, as does 'go help get'. See 'go help modules' and 'go help module-get'.

See also: go build, go install, go clean.

Module proxy protocol
A Go module proxy is any web server that can respond to GET requests for URLs of a specified form. The requests have no query parameters, so even a site serving from a fixed file system (including a file:/// URL) can be a module proxy.

For details on the GOPROXY protocol, see https://golang.org/ref/mod#goproxy-protocol.

Import path syntax
An import path (see 'go help packages') denotes a package stored in the local file system. In general, an import path denotes either a standard package (such as "unicode/utf8") or a package found in one of the work spaces (For more details see: 'go help gopath').

Relative import paths
An import path beginning with ./ or ../ is called a relative path. The toolchain supports relative import paths as a shortcut in two ways.

First, a relative path can be used as a shorthand on the command line. If you are working in the directory containing the code imported as "unicode" and want to run the tests for "unicode/utf8", you can type "go test ./utf8" instead of needing to specify the full path. Similarly, in the reverse situation, "go test .." will test "unicode" from the "unicode/utf8" directory. Relative patterns are also allowed, like "go test ./..." to test all subdirectories. See 'go help packages' for details on the pattern syntax.

Second, if you are compiling a Go program not in a work space, you can use a relative path in an import statement in that program to refer to nearby code also not in a work space. This makes it easy to experiment with small multipackage programs outside of the usual work spaces, but such programs cannot be installed with "go install" (there is no work space in which to install them), so they are rebuilt from scratch each time they are built. To avoid ambiguity, Go programs cannot use relative import paths within a work space.

Remote import paths
Certain import paths also describe how to obtain the source code for the package using a revision control system.

A few common code hosting sites have special syntax:

Bitbucket (Git, Mercurial)

	import "bitbucket.org/user/project"
	import "bitbucket.org/user/project/sub/directory"

GitHub (Git)

	import "github.com/user/project"
	import "github.com/user/project/sub/directory"

Launchpad (Bazaar)

	import "launchpad.net/project"
	import "launchpad.net/project/series"
	import "launchpad.net/project/series/sub/directory"

	import "launchpad.net/~user/project/branch"
	import "launchpad.net/~user/project/branch/sub/directory"

IBM DevOps Services (Git)

	import "hub.jazz.net/git/user/project"
	import "hub.jazz.net/git/user/project/sub/directory"
For code hosted on other servers, import paths may either be qualified with the version control type, or the go tool can dynamically fetch the import path over https/http and discover where the code resides from a <meta> tag in the HTML.

To declare the code location, an import path of the form

repository.vcs/path
specifies the given repository, with or without the .vcs suffix, using the named version control system, and then the path inside that repository. The supported version control systems are:

Bazaar      .bzr
Fossil      .fossil
Git         .git
Mercurial   .hg
Subversion  .svn
For example,

import "example.org/user/foo.hg"
denotes the root directory of the Mercurial repository at example.org/user/foo or foo.hg, and

import "example.org/repo.git/foo/bar"
denotes the foo/bar directory of the Git repository at example.org/repo or repo.git.

When a version control system supports multiple protocols, each is tried in turn when downloading. For example, a Git download tries https://, then git+ssh://.

By default, downloads are restricted to known secure protocols (e.g. https, ssh). To override this setting for Git downloads, the GIT_ALLOW_PROTOCOL environment variable can be set (For more details see: 'go help environment').

If the import path is not a known code hosting site and also lacks a version control qualifier, the go tool attempts to fetch the import over https/http and looks for a <meta> tag in the document's HTML <head>.

The meta tag has the form:

<meta name="go-import" content="import-prefix vcs repo-root">
The import-prefix is the import path corresponding to the repository root. It must be a prefix or an exact match of the package being fetched with "go get". If it's not an exact match, another http request is made at the prefix to verify the <meta> tags match.

The meta tag should appear as early in the file as possible. In particular, it should appear before any raw JavaScript or CSS, to avoid confusing the go command's restricted parser.

The vcs is one of "bzr", "fossil", "git", "hg", "svn".

The repo-root is the root of the version control system containing a scheme and not containing a .vcs qualifier.

For example,

import "example.org/pkg/foo"
will result in the following requests:

https://example.org/pkg/foo?go-get=1 (preferred)
http://example.org/pkg/foo?go-get=1  (fallback, only with -insecure)
If that page contains the meta tag

<meta name="go-import" content="example.org git https://code.org/r/p/exproj">
the go tool will verify that https://example.org/?go-get=1 contains the same meta tag and then git clone https://code.org/r/p/exproj into GOPATH/src/example.org.

When using GOPATH, downloaded packages are written to the first directory listed in the GOPATH environment variable. (See 'go help gopath-get' and 'go help gopath'.)

When using modules, downloaded packages are stored in the module cache. See https://golang.org/ref/mod#module-cache.

When using modules, an additional variant of the go-import meta tag is recognized and is preferred over those listing version control systems. That variant uses "mod" as the vcs in the content value, as in:

<meta name="go-import" content="example.org mod https://code.org/moduleproxy">
This tag means to fetch modules with paths beginning with example.org from the module proxy available at the URL https://code.org/moduleproxy. See https://golang.org/ref/mod#goproxy-protocol for details about the proxy protocol.

Import path checking
When the custom import path feature described above redirects to a known code hosting site, each of the resulting packages has two possible import paths, using the custom domain or the known hosting site.

A package statement is said to have an "import comment" if it is immediately followed (before the next newline) by a comment of one of these two forms:

package math // import "path"
package math /* import "path" */
The go command will refuse to install a package with an import comment unless it is being referred to by that import path. In this way, import comments let package authors make sure the custom import path is used and not a direct path to the underlying code hosting site.

Import path checking is disabled for code found within vendor trees. This makes it possible to copy code into alternate locations in vendor trees without needing to update import comments.

Import path checking is also disabled when using modules. Import path comments are obsoleted by the go.mod file's module statement.

See https://golang.org/s/go14customimport for details.

Modules, module versions, and more
Modules are how Go manages dependencies.

A module is a collection of packages that are released, versioned, and distributed together. Modules may be downloaded directly from version control repositories or from module proxy servers.

For a series of tutorials on modules, see https://golang.org/doc/tutorial/create-module.

For a detailed reference on modules, see https://golang.org/ref/mod.

By default, the go command may download modules from https://proxy.golang.org. It may authenticate modules using the checksum database at https://sum.golang.org. Both services are operated by the Go team at Google. The privacy policies for these services are available at https://proxy.golang.org/privacy and https://sum.golang.org/privacy, respectively.

The go command's download behavior may be configured using GOPROXY, GOSUMDB, GOPRIVATE, and other environment variables. See 'go help environment' and https://golang.org/ref/mod#private-module-privacy for more information.

Module authentication using go.sum
When the go command downloads a module zip file or go.mod file into the module cache, it computes a cryptographic hash and compares it with a known value to verify the file hasn't changed since it was first downloaded. Known hashes are stored in a file in the module root directory named go.sum. Hashes may also be downloaded from the checksum database depending on the values of GOSUMDB, GOPRIVATE, and GONOSUMDB.

For details, see https://golang.org/ref/mod#authenticating.

Package lists and patterns
Many commands apply to a set of packages:

go action [packages]
Usually, [packages] is a list of import paths.

An import path that is a rooted path or that begins with a . or .. element is interpreted as a file system path and denotes the package in that directory.

Otherwise, the import path P denotes the package found in the directory DIR/src/P for some DIR listed in the GOPATH environment variable (For more details see: 'go help gopath').

If no import paths are given, the action applies to the package in the current directory.

There are four reserved names for paths that should not be used for packages to be built with the go tool:

- "main" denotes the top-level package in a stand-alone executable.

- "all" expands to all packages found in all the GOPATH trees. For example, 'go list all' lists all the packages on the local system. When using modules, "all" expands to all packages in the main module and their dependencies, including dependencies needed by tests of any of those.

- "std" is like all but expands to just the packages in the standard Go library.

- "cmd" expands to the Go repository's commands and their internal libraries.

Import paths beginning with "cmd/" only match source code in the Go repository.

An import path is a pattern if it includes one or more "..." wildcards, each of which can match any string, including the empty string and strings containing slashes. Such a pattern expands to all package directories found in the GOPATH trees with names matching the patterns.

To make common patterns more convenient, there are two special cases. First, /... at the end of the pattern can match an empty string, so that net/... matches both net and packages in its subdirectories, like net/http. Second, any slash-separated pattern element containing a wildcard never participates in a match of the "vendor" element in the path of a vendored package, so that ./... does not match packages in subdirectories of ./vendor or ./mycode/vendor, but ./vendor/... and ./mycode/vendor/... do. Note, however, that a directory named vendor that itself contains code is not a vendored package: cmd/vendor would be a command named vendor, and the pattern cmd/... matches it. See golang.org/s/go15vendor for more about vendoring.

An import path can also name a package to be downloaded from a remote repository. Run 'go help importpath' for details.

Every package in a program must have a unique import path. By convention, this is arranged by starting each path with a unique prefix that belongs to you. For example, paths used internally at Google all begin with 'google', and paths denoting remote repositories begin with the path to the code, such as 'github.com/user/repo'.

Packages in a program need not have unique package names, but there are two reserved package names with special meaning. The name main indicates a command, not a library. Commands are built into binaries and cannot be imported. The name documentation indicates documentation for a non-Go program in the directory. Files in package documentation are ignored by the go command.

As a special case, if the package list is a list of .go files from a single directory, the command is applied to a single synthesized package made up of exactly those files, ignoring any build constraints in those files and ignoring any other files in the directory.

Directory and file names that begin with "." or "_" are ignored by the go tool, as are directories named "testdata".

Configuration for downloading non-public code
The go command defaults to downloading modules from the public Go module mirror at proxy.golang.org. It also defaults to validating downloaded modules, regardless of source, against the public Go checksum database at sum.golang.org. These defaults work well for publicly available source code.

The GOPRIVATE environment variable controls which modules the go command considers to be private (not available publicly) and should therefore not use the proxy or checksum database. The variable is a comma-separated list of glob patterns (in the syntax of Go's path.Match) of module path prefixes. For example,

GOPRIVATE=*.corp.example.com,rsc.io/private
causes the go command to treat as private any module with a path prefix matching either pattern, including git.corp.example.com/xyzzy, rsc.io/private, and rsc.io/private/quux.

For fine-grained control over module download and validation, the GONOPROXY and GONOSUMDB environment variables accept the same kind of glob list and override GOPRIVATE for the specific decision of whether to use the proxy and checksum database, respectively.

For example, if a company ran a module proxy serving private modules, users would configure go using:

GOPRIVATE=*.corp.example.com
GOPROXY=proxy.example.com
GONOPROXY=none
The GOPRIVATE variable is also used to define the "public" and "private" patterns for the GOVCS variable; see 'go help vcs'. For that usage, GOPRIVATE applies even in GOPATH mode. In that case, it matches import paths instead of module paths.

The 'go env -w' command (see 'go help env') can be used to set these variables for future go command invocations.

For more details, see https://golang.org/ref/mod#private-modules.

Testing flags
The 'go test' command takes both flags that apply to 'go test' itself and flags that apply to the resulting test binary.

Several of the flags control profiling and write an execution profile suitable for "go tool pprof"; run "go tool pprof -h" for more information. The --alloc_space, --alloc_objects, and --show_bytes options of pprof control how the information is presented.

The following flags are recognized by the 'go test' command and control the execution of any test:

-bench regexp
    Run only those benchmarks matching a regular expression.
    By default, no benchmarks are run.
    To run all benchmarks, use '-bench .' or '-bench=.'.
    The regular expression is split by unbracketed slash (/)
    characters into a sequence of regular expressions, and each
    part of a benchmark's identifier must match the corresponding
    element in the sequence, if any. Possible parents of matches
    are run with b.N=1 to identify sub-benchmarks. For example,
    given -bench=X/Y, top-level benchmarks matching X are run
    with b.N=1 to find any sub-benchmarks matching Y, which are
    then run in full.

-benchtime t
    Run enough iterations of each benchmark to take t, specified
    as a time.Duration (for example, -benchtime 1h30s).
    The default is 1 second (1s).
    The special syntax Nx means to run the benchmark N times
    (for example, -benchtime 100x).

-count n
    Run each test and benchmark n times (default 1).
    If -cpu is set, run n times for each GOMAXPROCS value.
    Examples are always run once.

-cover
    Enable coverage analysis.
    Note that because coverage works by annotating the source
    code before compilation, compilation and test failures with
    coverage enabled may report line numbers that don't correspond
    to the original sources.

-covermode set,count,atomic
    Set the mode for coverage analysis for the package[s]
    being tested. The default is "set" unless -race is enabled,
    in which case it is "atomic".
    The values:
	set: bool: does this statement run?
	count: int: how many times does this statement run?
	atomic: int: count, but correct in multithreaded tests;
		significantly more expensive.
    Sets -cover.

-coverpkg pattern1,pattern2,pattern3
    Apply coverage analysis in each test to packages matching the patterns.
    The default is for each test to analyze only the package being tested.
    See 'go help packages' for a description of package patterns.
    Sets -cover.

-cpu 1,2,4
    Specify a list of GOMAXPROCS values for which the tests or
    benchmarks should be executed. The default is the current value
    of GOMAXPROCS.

-failfast
    Do not start new tests after the first test failure.

-list regexp
    List tests, benchmarks, or examples matching the regular expression.
    No tests, benchmarks or examples will be run. This will only
    list top-level tests. No subtest or subbenchmarks will be shown.

-parallel n
    Allow parallel execution of test functions that call t.Parallel.
    The value of this flag is the maximum number of tests to run
    simultaneously; by default, it is set to the value of GOMAXPROCS.
    Note that -parallel only applies within a single test binary.
    The 'go test' command may run tests for different packages
    in parallel as well, according to the setting of the -p flag
    (see 'go help build').

-run regexp
    Run only those tests and examples matching the regular expression.
    For tests, the regular expression is split by unbracketed slash (/)
    characters into a sequence of regular expressions, and each part
    of a test's identifier must match the corresponding element in
    the sequence, if any. Note that possible parents of matches are
    run too, so that -run=X/Y matches and runs and reports the result
    of all tests matching X, even those without sub-tests matching Y,
    because it must run them to look for those sub-tests.

-short
    Tell long-running tests to shorten their run time.
    It is off by default but set during all.bash so that installing
    the Go tree can run a sanity check but not spend time running
    exhaustive tests.

-timeout d
    If a test binary runs longer than duration d, panic.
    If d is 0, the timeout is disabled.
    The default is 10 minutes (10m).

-v
    Verbose output: log all tests as they are run. Also print all
    text from Log and Logf calls even if the test succeeds.

-vet list
    Configure the invocation of "go vet" during "go test"
    to use the comma-separated list of vet checks.
    If list is empty, "go test" runs "go vet" with a curated list of
    checks believed to be always worth addressing.
    If list is "off", "go test" does not run "go vet" at all.
The following flags are also recognized by 'go test' and can be used to profile the tests during execution:

-benchmem
    Print memory allocation statistics for benchmarks.

-blockprofile block.out
    Write a goroutine blocking profile to the specified file
    when all tests are complete.
    Writes test binary as -c would.

-blockprofilerate n
    Control the detail provided in goroutine blocking profiles by
    calling runtime.SetBlockProfileRate with n.
    See 'go doc runtime.SetBlockProfileRate'.
    The profiler aims to sample, on average, one blocking event every
    n nanoseconds the program spends blocked. By default,
    if -test.blockprofile is set without this flag, all blocking events
    are recorded, equivalent to -test.blockprofilerate=1.

-coverprofile cover.out
    Write a coverage profile to the file after all tests have passed.
    Sets -cover.

-cpuprofile cpu.out
    Write a CPU profile to the specified file before exiting.
    Writes test binary as -c would.

-memprofile mem.out
    Write an allocation profile to the file after all tests have passed.
    Writes test binary as -c would.

-memprofilerate n
    Enable more precise (and expensive) memory allocation profiles by
    setting runtime.MemProfileRate. See 'go doc runtime.MemProfileRate'.
    To profile all memory allocations, use -test.memprofilerate=1.

-mutexprofile mutex.out
    Write a mutex contention profile to the specified file
    when all tests are complete.
    Writes test binary as -c would.

-mutexprofilefraction n
    Sample 1 in n stack traces of goroutines holding a
    contended mutex.

-outputdir directory
    Place output files from profiling in the specified directory,
    by default the directory in which "go test" is running.

-trace trace.out
    Write an execution trace to the specified file before exiting.
Each of these flags is also recognized with an optional 'test.' prefix, as in -test.v. When invoking the generated test binary (the result of 'go test -c') directly, however, the prefix is mandatory.

The 'go test' command rewrites or removes recognized flags, as appropriate, both before and after the optional package list, before invoking the test binary.

For instance, the command

go test -v -myflag testdata -cpuprofile=prof.out -x
will compile the test binary and then run it as

pkg.test -test.v -myflag testdata -test.cpuprofile=prof.out
(The -x flag is removed because it applies only to the go command's execution, not to the test itself.)

The test flags that generate profiles (other than for coverage) also leave the test binary in pkg.test for use when analyzing the profiles.

When 'go test' runs a test binary, it does so from within the corresponding package's source code directory. Depending on the test, it may be necessary to do the same when invoking a generated test binary directly.

The command-line package list, if present, must appear before any flag not known to the go test command. Continuing the example above, the package list would have to appear before -myflag, but could appear on either side of -v.

When 'go test' runs in package list mode, 'go test' caches successful package test results to avoid unnecessary repeated running of tests. To disable test caching, use any test flag or argument other than the cacheable flags. The idiomatic way to disable test caching explicitly is to use -count=1.

To keep an argument for a test binary from being interpreted as a known flag or a package name, use -args (see 'go help test') which passes the remainder of the command line through to the test binary uninterpreted and unaltered.

For instance, the command

go test -v -args -x -v
will compile the test binary and then run it as

pkg.test -test.v -x -v
Similarly,

go test -args math
will compile the test binary and then run it as

pkg.test math
In the first example, the -x and the second -v are passed through to the test binary unchanged and with no effect on the go command itself. In the second example, the argument math is passed through to the test binary, instead of being interpreted as the package list.

Testing functions
The 'go test' command expects to find test, benchmark, and example functions in the "*_test.go" files corresponding to the package under test.

A test function is one named TestXxx (where Xxx does not start with a lower case letter) and should have the signature,

func TestXxx(t *testing.T) { ... }
A benchmark function is one named BenchmarkXxx and should have the signature,

func BenchmarkXxx(b *testing.B) { ... }
An example function is similar to a test function but, instead of using *testing.T to report success or failure, prints output to os.Stdout. If the last comment in the function starts with "Output:" then the output is compared exactly against the comment (see examples below). If the last comment begins with "Unordered output:" then the output is compared to the comment, however the order of the lines is ignored. An example with no such comment is compiled but not executed. An example with no text after "Output:" is compiled, executed, and expected to produce no output.

Godoc displays the body of ExampleXxx to demonstrate the use of the function, constant, or variable Xxx. An example of a method M with receiver type T or *T is named ExampleT_M. There may be multiple examples for a given function, constant, or variable, distinguished by a trailing _xxx, where xxx is a suffix not beginning with an upper case letter.

Here is an example of an example:

func ExamplePrintln() {
	Println("The output of\nthis example.")
	// Output: The output of
	// this example.
}
Here is another example where the ordering of the output is ignored:

func ExamplePerm() {
	for _, value := range Perm(4) {
		fmt.Println(value)
	}

	// Unordered output: 4
	// 2
	// 1
	// 3
	// 0
}
The entire test file is presented as the example when it contains a single example function, at least one other function, type, variable, or constant declaration, and no test or benchmark functions.

See the documentation of the testing package for more information.

Controlling version control with GOVCS
The 'go get' command can run version control commands like git to download imported code. This functionality is critical to the decentralized Go package ecosystem, in which code can be imported from any server, but it is also a potential security problem, if a malicious server finds a way to cause the invoked version control command to run unintended code.

To balance the functionality and security concerns, the 'go get' command by default will only use git and hg to download code from public servers. But it will use any known version control system (bzr, fossil, git, hg, svn) to download code from private servers, defined as those hosting packages matching the GOPRIVATE variable (see 'go help private'). The rationale behind allowing only Git and Mercurial is that these two systems have had the most attention to issues of being run as clients of untrusted servers. In contrast, Bazaar, Fossil, and Subversion have primarily been used in trusted, authenticated environments and are not as well scrutinized as attack surfaces.

The version control command restrictions only apply when using direct version control access to download code. When downloading modules from a proxy, 'go get' uses the proxy protocol instead, which is always permitted. By default, the 'go get' command uses the Go module mirror (proxy.golang.org) for public packages and only falls back to version control for private packages or when the mirror refuses to serve a public package (typically for legal reasons). Therefore, clients can still access public code served from Bazaar, Fossil, or Subversion repositories by default, because those downloads use the Go module mirror, which takes on the security risk of running the version control commands using a custom sandbox.

The GOVCS variable can be used to change the allowed version control systems for specific packages (identified by a module or import path). The GOVCS variable applies when building package in both module-aware mode and GOPATH mode. When using modules, the patterns match against the module path. When using GOPATH, the patterns match against the import path corresponding to the root of the version control repository.

The general form of the GOVCS setting is a comma-separated list of pattern:vcslist rules. The pattern is a glob pattern that must match one or more leading elements of the module or import path. The vcslist is a pipe-separated list of allowed version control commands, or "all" to allow use of any known command, or "off" to disallow all commands. Note that if a module matches a pattern with vcslist "off", it may still be downloaded if the origin server uses the "mod" scheme, which instructs the go command to download the module using the GOPROXY protocol. The earliest matching pattern in the list applies, even if later patterns might also match.

For example, consider:

GOVCS=github.com:git,evil.com:off,*:git|hg
With this setting, code with a module or import path beginning with github.com/ can only use git; paths on evil.com cannot use any version control command, and all other paths (* matches everything) can use only git or hg.

The special patterns "public" and "private" match public and private module or import paths. A path is private if it matches the GOPRIVATE variable; otherwise it is public.

If no rules in the GOVCS variable match a particular module or import path, the 'go get' command applies its default rule, which can now be summarized in GOVCS notation as 'public:git|hg,private:all'.

To allow unfettered use of any version control system for any package, use:

GOVCS=*:all
To disable all use of version control, use:

GOVCS=*:off
The 'go env -w' command (see 'go help env') can be used to set the GOVCS variable for future go command invocations.

The Go Gopher
Copyright
Terms of Service
Privacy Policy
Report a website issue
Supported by Google
