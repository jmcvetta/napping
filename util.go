// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package napping

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
)

func pretty(v interface{}) string {
	// Get source file and line
	_, file, line, _ := runtime.Caller(1)

	// Get relative filename
	filename := filepath.Base(file)

	// Convert to JSON
	b, _ := json.MarshalIndent(v, "", "\t")

	// Make that all pretty together
	return fmt.Sprintf("%s:%d: \n%s\n", filename, line, string(b))
}
