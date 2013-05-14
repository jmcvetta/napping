// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package restclient

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func prettyPrint(v interface{}) {
	_, file, line, _ := runtime.Caller(1)
	lineNo := strconv.Itoa(line)
	file = filepath.Base(file)
	b, _ := json.MarshalIndent(v, "", "\t")
	s := file + ":" + lineNo + ": \n" + string(b) + "\n"
	os.Stderr.WriteString(s)
}
