package chapter10

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.1 INTRODUCTION
1. The purpose of any package system is to make the design and maintenance of large programs practical by grouping related features together into units that can be easily understood and changed
independent of the other packages of the program. This "modularity" allows packages to be shared and reused by different projects, distributed within an orgsnization, or made available to the
wider world.

2. Each package defines a distinct name space that encloses its identifiers. Each name is associated with a practicular package, letting us choose short, clear, names for the types, functions,
and so on on that we use most often, without creating conflicts with other parts of the program.

3. Packages also provdide encapsulation by controlling which names are visible or exported outside the package. Restricting the visibility of package members hides helper functions and
types behind the package's API, allowing the package maintainer to change the implementation with confidence that no code outside the package will be affected.

4. Three reasons of Go's compiler speed:
	1) All imports must be explicitly listed at the beginning of each source file, so the compiler doen't have to read and process an entire file to determine its dependencies.
	2) The dependencies of a package for a directed acyclic graph, and because there are no cycles, packages can be compiled separately and perhaps in parallel.
	3) The object file for a compiled Go package records export information not just for the package itself, but for its dependencies too.
When compiling a package, the compiler must read one object file for each import but needn't look beyond these files.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.2 IMPORT PATHS
1) For packages you intend to share or publish, import paths should be globally unique.

2) To avoid conflicts, the import paths of all packages other than those from the standart library should start with the Internet domain name of the organization that owns or hosts
the package. This also makes it possible to find packages.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.3 THE PACKAGE DECLARAION
1. A package declaration is required at the start of every Go source file. Its main purpose is to determine the default identifier for that package (called the "package name") when it's
imported by another package.

2. Conventionally the package name is the last segment of the import path, and as a result, two packages may have the same name even though their import psths necessarily differ.

3. Three major exceptions to the "last segment" convention:
	1) A package defining a command (an executable program) always has the name "main", regardless of the package's import path.
	This is a signal to "go build" that it must invoke the linker to make an executable file.

	2) Some files in the directory may have the suffix "_test" on their package name if the file name ends with "_test.go". So a directory may define 2 packages:
		1) The usual one
		2) Another called "externall test package"
	The "_test" suffix signals to "go test" that it must build both packages, and it indicates which files belong to each package. External test packages are used to avoid cycles
	in the import graph arising from dependencies of the test;

	3) Some tool for dependency management append version numbere suffixes to package import paths, such as "gopkg.in/yaml.v2".

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.4 IMPORT DECLARATIONS
1. Parenthesized form is more common

2. Imported packages may be grouped by introducing blank lines, such groupings usually indicate different domains.

3. The order isn't significant, but by convention the lines of each group are sorted alphabetically. Both "gofmt" and "goimports" will group and sort it for us.

4. If we need to import two packages whose names are the same, like "math/rand" and "crypto/rand", into a third package, the import declaration must specify an alternative name for at least
one of them to avoid a conflict. This is called "renaming import"
	1) A renaming import may be useful even when there's no conflict: if the name of the imported package is unwieldy, as is sometimes the case for automatically generated code, an abbreviated
	name may be more convenieny.
	2) Choosing an alternative name can help avoid conflicts with common local variable names. For example: in a file with many local variables named path, we might import the standart "path"
	package as "pathpkg".

5. Each import declaration establishes a dependency from the current package to the imported package. The "go build" tool reports an error if these dependencies form a cycle.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.5 BLANK imports
1. It's an error to import a package into a file but not refer to the name it defines within that file. However, on occasion we must import a package merely for the side effects of doing so:
evaluation of the initializer expressions of its package-level variabled and execution of its "init()" functions.
	import _ "image/png" // register PNG decoder

2. To suppress the "unused import" error we would otherwise encounter, we must use a "renaming import" in which the alternative name is "_", the blank identifier. As usual, the blank
identifier can never be referenced.

3. This technique is known as "blank import". It's most often used to implement to implement a compile-time mechanism whereby the main program can enable optional features by blank importing
additional packages.

4. The standart library provides decoders for GIF, PNG and JPEG and users may provide others, but to keep executables small, decoders aren't included in an application unless explicitly requested. An entry of the table
is added to the table by calling image.RegisterFormat, typically from within the package initializer of the supporting package for each format, like this one in "image/png". The effect is that an application need only
blank-import the package for the format it needs to make the "image.Decode" function able to decode it.

5. The "database/sql" package uses a similar mechanism to let users install just the database drivers they need.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.6 PACKAGES AND NAMING

PACKAGES
1. While creating a package, keep its name short, but not so short to be cryptic.

2. Be descriptive and unambiguous where possible. For example: don't name a utility package "util" when a name such as "imageutil" or "ioutil" is specific yet still concise.

3. Avoid choosing package names that are commonly used for related local variables, or you may compel the package's clients to use renaming imports, as with the "path" package.

4. Package names usually take the singular form. The standart packages: "bytes", "errors", "strings" use the plural to avoid hiding the corresponding predeclared types and, in
the case of "go/types" to conflict with a keyword.

5. Avoid package names that already have other connotations. For example: we used the name "temp" for the temperatures. It was a terrible idea because "temp" is an almost universal
synonym for "temporary". In the end, it became tempconv, which is shorter and parallel with strconv.


MEMBERS OF PACKAGE
1. The burden of describing the package member is borne equally by the package name and the member name.

2. When desining a package consider how the two parts of a qualified identifier work together, not the member name alone. For example:
	bytes.Equal
	flag.Int
	http.Get
	json.Marshal

3. We can identify some common naming patterns.
3.1 The string package provides a number of independent functions for manipulating string:
	Index(...) int

	type Replacer struct {...}
	func NewReplacer(...) *Replacer

	type Reader struct {...}
	func NewReplacer(...) *Reader

	1) The word "strings" doesn't appear in any of their names. Clients refer to them as "strings.Index(...)", "strings.Replacer()"

3.2 Other packages that we might describe as "single-type packages", such as "html/template" and "math/rand" expose one principal data type plus its methods, and often a "New()"
function to create instances.
	package rand // math/rand

	type Rand struct {...}
	func New(...) *Rand

	1) This can lead to repetition as in "template.Template" or "rand.Rand" which is why the names of these kinds of packages are often especially short.

3.3 At the other extreme, there are packages like "net/http" that have a lot of names without a lot of structure, because they perform a complicated task. Despite having over twenty
types and many more fuctions the package's most important members have the simplest names: "Get", "Post", "Handle", "Error", "Client", "Server"
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
10.7 THE GO TOOL

1. The "go" tool is a package manager that:
	1) answers queries about inventory of packages
	2) computes package dependencies
	3) downloads them from remote-version control systems
It's a build system that computes file dependencies and invokes: compilers, assemblers and linkers. Although it's intentionally less complete than the standart Unix "make". And it's
a test driver, as we'll see.

2. To keep the need for configuration to a minimum, the "go" tool relies heavily on conventions.

10.7.1 WORKSPACE ORGANIZATION
1. The only configuration most users ever need is the "GOPATH" environment variable, which specifies the root of the workspace.


2. When switching to a different workspace users update the value of "GOPATH". GOPATH" has three subdirectories:
	1) "src" holds source code.
	2) "pkg" is where the build tools store compiled packages
	3) "bin" holds executable programs like "helloworld"

3. The second environment variable, "GOROOT" specifies the root directory of the Go distribution, which provides all the packages of the standart library. Its structure resembles that
of "GOPATH". Users never need to set "GOROOT" since, by default the "go" tool will use the location where it was installed.

4. "go env" prints the effective values of the environment variables relevant to the toolchain, including the default values for the missing ones.
	1) "GOOS" specifies the target operating system ("android, linux, darwin, windows")
	2) "GOARCH" specifies the target processor architecture

10.7.2 DOWNLOADING PACKAGES
1. "go get" command can download a single package or an entire subrtree or repository using the "..." notation. The tool also computes and downloads all the dependencies of the initial
packages, which is why the "golang.org/x/net/htnl" package appeared in the workspace in the previous example.

2. Once "go get" has downloaded the packages, it builds them and then installs the libraries and commands. For example:
	1) go get golang.org/x/lint/golint  

3. "go get" command has support for popular code-hosting sites like "GitHub", "Bitbucket", "Launchpad" and can make the appropriate requests to their version-control systems. For less
known sites we have to indicate which version-control protocol to use in the import path, such as "Git" or "Mercurial".
	1) Run "go help importpath" for details.

4. The directories that "go get" creates are true client of the remote repository, not just copies of the files, so we can use version-control commands to see a difference of local edits
we've made or to update to a different revision. The feature of the go tool lets packages use a custom domain in their import path while being hosted by a generic service such as 
"googlesource.com" or any else. HTML pages beneath "https://golang.org/x/net/html" include the metadata which redirects the "go" tool to the git repository at the actual hosting site.
	1) "go get -u" generally retrivies the latest version of each package, which is convenient when we're getting started but may be inappropriate for deployed projects where precise 
	control of dependencies is critical for release. The usual solution to this problem is to "vendor" the code, that is, to make a persistent local copy of all the necessary
	dependencies, and to update this copy carefully and deliberately.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
