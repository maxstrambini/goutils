# goutils

Some golang utility functions I use

##Example:

				package main

				import "log"
				import "github.com/maxstrambini/goutils"

				func main() {

					//print version:
					log.Printf(goutils.PackageVersion)

					/*
					//ExistsPath example:
					if !goutils.ExistsPath(zipFullName) {
						log.Printf("zipFilePathFromAssetID could not find '%s' returns default!", zipFullName)
						zipFullName = "./default.zip"
					}
					*/
				}



