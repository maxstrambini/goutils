package goutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"time"
)

var PackageDesc string = "utility functions in goutils package"

func PrintMap(title string, m map[string]interface{}, to_log bool) {

	fmt.Println(title)
	if to_log {
		log.Println(title)
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		//fmt.Println("> [", k, "]:", m[k])
		//if to_log {log.Println("> [", k, "]:", m[k])}
		fmt.Println("> [", k, "](", reflect.TypeOf(m[k]), "):", m[k])
		if to_log {
			log.Println("> [", k, "]:(", reflect.TypeOf(m[k]), "):", m[k])
		}
	}

	fmt.Println("********************")
	if to_log {
		log.Println("********************")
	}
}

func PrettyPrintMap(m map[string]interface{}) {
	fmt.Println("********************")
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("PrettyPrintMap error:", err)
	} else {
		log.Printf("PrettyPrintMap: \n%s", string(b))
	}
	fmt.Println("********************")
}

func EncodeXMLText(text string) (encodedtext string) {
	encodedtext = strings.Replace(text, "&", "&amp;", -1)
	encodedtext = strings.Replace(encodedtext, "<", "&lt;", -1)
	encodedtext = strings.Replace(encodedtext, ">", "&gt;", -1)
	encodedtext = strings.Replace(encodedtext, "\"", "&quot;", -1)
	encodedtext = strings.Replace(encodedtext, "'", "&#x0027;", -1)
	return
}

func GetNowString() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//ExistsPath check if a path exists by doing a stat, no distinction between file and folder
func ExistsPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

//SafeCreateFolder create all path with logging, folder rights are '0777'
func SafeCreateFolder(path string) (success bool) {
	if !ExistsPath(path) {
		log.Printf("SafeCreateFolder creates: '%s'", path)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			log.Panicf("ERROR: Could not create folder '%s' Error: %v", path, err)
		}
	} else {
		log.Printf("SafeCreateFolder has found '%s' already existing", path)
	}
	success = true
	return
}

//SafeMove with logging
func SafeMove(source, dest string) (success bool) {
	log.Printf("SafeMove '%s' -> '%s'", source, dest)
	err := os.Rename(source, dest)
	if err != nil {
		log.Panicf("EXIT: Could not rename XML file to .done! Error: %v", err)
	}
	success = true
	return
}

//SafeMoveWithStub moves a file and writes a text stub
func SafeMoveWithStub(source, dest, stubText string) (success bool) {
	log.Printf("SafeMoveWithStub '%s' -> '%s'", source, dest)
	err := os.Rename(source, dest)
	if err != nil {
		log.Panicf("EXIT: Could not rename XML file to .done! Error: %v", err)
	}

	f, err := os.Create(dest + ".stub") // create/truncate the file
	if err != nil {
		log.Panicf("EXIT: Could not create stub '%s'! Error: %v", dest+".stub", err)
	} // panic if error
	defer f.Close() // make sure it gets closed after

	fmt.Fprintln(f, stubText)

	success = true
	return
}

//GetFileNameWithoutExtension get the name of a file without the extension
func GetFileNameWithoutExtension(fullName string) (nameWithoutExt string) {
	_, n := path.Split(fullName)
	nameWithoutExt = n[0 : len(n)-len(path.Ext(n))]
	return
}
