// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package restclient

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func prettyPrint(v interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	lineNo := strconv.Itoa(line)
	file = filepath.Base(file)
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Panic(err)
	}
	s := file + ":" + lineNo + ": \n" + string(b) + "\n"
	os.Stderr.WriteString(s)
}
