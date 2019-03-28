package goutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"time"
)

//PackageDesc string description
var PackageDesc = "utility functions in goutils package"

//PackageVersion string version
var PackageVersion = "1.2.0"

//PrintMap prints a map to console/log
func PrintMap(title string, m map[string]interface{}, toLog bool) {

	fmt.Println(title)
	if toLog {
		log.Println(title)
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		//fmt.Println("> [", k, "]:", m[k])
		//if toLog {log.Println("> [", k, "]:", m[k])}
		fmt.Println("> [", k, "](", reflect.TypeOf(m[k]), "):", m[k])
		if toLog {
			log.Println("> [", k, "]:(", reflect.TypeOf(m[k]), "):", m[k])
		}
	}

	fmt.Println("********************")
	if toLog {
		log.Println("********************")
	}
}

//PrettyPrintMap pretty prints a map to log
func PrettyPrintMap(m map[string]interface{}) {
	fmt.Println("********************")
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("PrettyPrintMap error: '%v'", err)
	} else {
		log.Printf("PrettyPrintMap: \n%s", string(b))
	}
	fmt.Println("********************")
}

//EncodeXMLText change XML forbidden chars to their safe equivalents
func EncodeXMLText(text string) (encodedtext string) {
	encodedtext = strings.Replace(text, "&", "&amp;", -1)
	encodedtext = strings.Replace(encodedtext, "<", "&lt;", -1)
	encodedtext = strings.Replace(encodedtext, ">", "&gt;", -1)
	encodedtext = strings.Replace(encodedtext, "\"", "&quot;", -1)
	encodedtext = strings.Replace(encodedtext, "'", "&#x0027;", -1)
	return
}

//GetNowString returns a now string
func GetNowString() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//GetNowForFileName returns a now string with a compact format for file names
func GetNowForFileName() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//ExistsPath check if a path exists by doing a stat, no distinction between file and folder
func ExistsPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

//ExistsPathE returns whether the given file or directory exists or not with the Error
func ExistsPathE(path string) (bool, error) {
	_, err := os.Stat(path)

	/*
		if err == nil {
			//file was found
			return true, nil
		}
		if os.IsNotExist(err) {
			// error reports that file does not exist
			// return false, without any error
			return false, nil
		}
		// more serious errors, maybe stat itself has failed
		return false, err
	*/

	//better written:
	if err != nil {
		if os.IsNotExist(err) {
			// error reports that file does not exist
			// return false, without any error
			return false, nil
		} else {
			// more serious errors, maybe stat itself has failed
			return false, err
		}
	}
	//file was found
	return true, nil

}

//ExistsDirE returns whether the given directory exists or not and the Error
func ExistsDirE(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("'%s' is not a folder", path)
	}
	return false, fmt.Errorf("'%s' not found", path)
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

// ReadDirMax is ioutil.ReadDir copied to allow different sorting (or none)
// Example:
//		fileInfos, _ := goutils.ReadDirMax("D:\\Temp\\in1", "size")
//		for _, f := range fileInfos {
//			if !f.IsDir() {
//				log.Printf("FILE %+v", f.Name())
//			}
//		}
func ReadDirMax(dirname string, sortBy string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	if sortBy == "name" {
		sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	} else if sortBy == "modified" {
		sort.Slice(list, func(i, j int) bool { return list[i].ModTime().UnixNano() < list[j].ModTime().UnixNano() })
	} else if sortBy == "size" {
		sort.Slice(list, func(i, j int) bool { return list[i].Size() < list[j].Size() })
	}
	return list, nil
}

//WriteTextToFile write some text to a file
func WriteTextToFile(fullName, text string) (success bool) {
	f, err := os.Create(fullName) // create/truncate the file
	if err != nil {
		log.Printf("ERROR while creating/truncating file '%s': %s", fullName, err)
	} else {
		defer f.Close() // make sure it gets closed after

		_, errw := f.WriteString(text)
		if errw != nil {
			log.Printf("ERROR writing text to file '%s': %s", fullName, errw)
		}
		success = true
	}
	return
}

//ReadTextFromFile reads text from a file
func ReadTextFromFile(fullName string) (text string, err error) {
	var f *os.File
	f, err = os.Open(fullName)
	if err != nil {
		log.Printf("ERROR while opening file '%s': %s", fullName, err)
	} else {
		defer f.Close() // make sure it gets closed after

		var b []byte
		b, err = ioutil.ReadAll(f)
		if err != nil {
			log.Printf("ERROR reading from file '%s': %s", fullName, err)
		}
		text = string(b) // convert content to a 'string'
	}
	return
}

//WaitForever block waiting forever a goroutine
func WaitForever() {
	select {}
}

//WaitForeverMain block waiting forever a main func
func WaitForeverMain() {
	<-time.After(time.Duration(math.MaxInt64))
}

//GetFileAgeInMinutes returns a file age in minutes
func GetFileAgeInMinutes(fullName string) (success bool, age int64) {

	info, err := os.Stat(fullName)
	if err != nil {
		log.Printf("GetFileAgeInMinutes: ERROR in stat for '%s' -> %v", fullName, err)
		return
	}

	age = int64(time.Since(info.ModTime()).Minutes())
	success = true
	return
}

//GetFileAgeInSeconds returns a file age in minutes
func GetFileAgeInSeconds(fullName string) (success bool, age int64) {

	info, err := os.Stat(fullName)
	if err != nil {
		log.Printf("GetFileAgeInSeconds: ERROR in stat for '%s' -> %v", fullName, err)
		return
	}

	age = int64(time.Since(info.ModTime()).Seconds())
	success = true
	return
}

//GetFileAge returns a file age as duration than use for example duration.Minutes() or duration.Seconds()
func GetFileAge(fullName string) (success bool, duration time.Duration) {

	info, err := os.Stat(fullName)
	if err != nil {
		log.Printf("GetFileAge: ERROR in stat for '%s' -> %v", fullName, err)
		return
	}

	duration = time.Since(info.ModTime())
	success = true
	return
}

//CopyFile copy from source path to destination path and return error
func CopyFile(src string, dst string, overwriteExisting bool) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	if !overwriteExisting {
		_, err := os.Stat(dst)
		if err == nil {
			//file exists, return error
			return errors.New("Destination '" + dst + "' already exists: overwrite not set")
		}
	}
	to, err2 := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err2 != nil {
		return err2
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

//CopyFileEx copy from source path to destination path and return error
func CopyFileEx(src string, dst string, overwriteExisting bool) (overwritten bool, err error) {
	from, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer from.Close()

	_, erre := os.Stat(dst)
	if erre == nil {
		//file exists
		if overwriteExisting {
			overwritten = true
		} else {
			return false, errors.New("Destination '" + dst + "' already exists: overwrite not set")
		}
	} else {
		//file does not exist
		//overwritten = false
	}

	to, err2 := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err2 != nil {
		return overwritten, err2
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return overwritten, err
	}
	return overwritten, nil
}

//PrettyPrintStruct prints a struct using reflection, simple version
func PrettyPrintStruct(obj interface{}) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		if val.Field(i).CanInterface() {
			fieldValue := val.Field(i).Interface()
			//fmt.Println(fieldValue)
			fmt.Printf("%d: %s %s = %v\n", i,
				typ.Field(i).Name, val.Field(i).Type(), fieldValue)
		} else {
			fmt.Printf("%d: %s (private value)\n", i, typ.Field(i).Name)
		}
	}
}

//PrettyFormatStruct format a struct to a string for printing using reflection, simple version
func PrettyFormatStruct(obj interface{}) string {
	s := ""
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		if val.Field(i).CanInterface() {
			fieldValue := val.Field(i).Interface()
			//fmt.Println(fieldValue)
			s += fmt.Sprintf("%d: %s %s = %v\n", i,
				typ.Field(i).Name, val.Field(i).Type(), fieldValue)
		} else {
			s += fmt.Sprintf("%d: %s (private value)\n", i, typ.Field(i).Name)
		}
	}
	return s
}
