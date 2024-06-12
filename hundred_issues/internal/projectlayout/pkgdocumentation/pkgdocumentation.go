package pkgdocumentation

/*
	PACKAGE DOCUMENTATION
1. Each exported field of a package must be documented properly.
	1) It's required to start a comment with the name of a unit and express its purpose.
2. Each comment must be a full sentence followed with a dot.
3.When commenting a function/method, we must point what it does instead of how it does.
4. If an element is old and should not be in use, we must comment it in the following way:
	// Deprecated
	and follow it with a usual comment
5. Commenting a var/const.
	1) The purpose must be arranged above the definition (public)
	2) The contents are not needed to be public. It might be arranged both at the top and on the right of a var/const
6. Each package must be commented to be properly maintained.
	1) By convention the comment of the package starts with "// Package packageName"
	2) The first row must be short.
	3) We should arrange the documentation of a package in a file has the same name with the package doccumented or into the doc.go file.
	4) All the comments before the first empty line will be included within the package documentation and displayed to a client.
*/