# goutils
some golang utility functions I use

import (
	"github.com/maxstrambini/goutils"
)

//ExistsPath example:

	if !goutils.ExistsPath(zipFullName) {
		log.Printf("zipFilePathFromAssetID could not find '%s' returns default!", zipFullName)
		zipFullName = "./default.zip"
	}

