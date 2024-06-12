package b

import (
	"fmt"
	"hundred_issues/internal/projectlayout/pkgdocumentation/a"
)

func DefaultPermissionUsage() {
	// We see only the comment above the const and its contents are hidden.
	copyOfDefaultPermission := a.DefaultPermission

	fmt.Println(copyOfDefaultPermission)
}
