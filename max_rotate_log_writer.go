/*

max_rotate_log_writer.go

[2018-05-10] a custom file writer to rotate logs

## descrizione: quando il file supera 'maxBytes' viene chiuso, rinominato con la data e l'ora e ricreato

## esempio di utilizzo:

			func main() {
				maxRotateWriter := NewMaxRotateWriter(logName, 5*1024*1024, true, 30)
				log.SetOutput(maxRotateWriter)
				log.Printf("rotating log ...\n")
			}


[2019-09-11] logs written to file and also to stdout
*/

package goutils

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

// MaxRotateWriter defines a custom writer to rotate logs
type MaxRotateWriter struct {
	lock                    sync.Mutex
	filename                string // should be set to the actual filename
	writterBytes            int    // rotate > maxBytes
	maxBytes                int    // counter of written bytes
	fp                      *os.File
	rotateFilesByNumber     bool // when true rotated files are _1.log, _2.log, ecc
	maxRotatedFilesByNumber int  // default 9
}

// NewMaxRotateWriter Make a new MaxRotateWriter. Return nil if error occurs during setup.
func NewMaxRotateWriter(filename string, maxBytes int, rotateFilesByNumber bool, maxRotatedFilesByNumber int) *MaxRotateWriter {
	w := &MaxRotateWriter{filename: filename, maxBytes: maxBytes,
		rotateFilesByNumber: rotateFilesByNumber, maxRotatedFilesByNumber: maxRotatedFilesByNumber}

	err := w.Rotate()
	if err != nil {
		return nil
	}

	fmt.Printf("LOG: %+v\n", w)
	return w
}

// Write satisfies the io.Writer interface.
func (w *MaxRotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	os.Stdout.Write(output)
	n, err := w.fp.Write(output)
	w.writterBytes += n
	if w.maxBytes > 0 && w.writterBytes >= w.maxBytes {
		w.rotateWithoutLock()
	}
	return n, err
}

// Rotate performs the file rotation locked
func (w *MaxRotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.rotateWithoutLock()
	return
}

// rotateWithoutLock perform the actual act of rotating and reopening file.
func (w *MaxRotateWriter) rotateWithoutLock() (err error) {

	//fmt.Printf("Rotating logs ...\n")
	//time.Sleep(3 * time.Second)

	// Close existing file if open
	if w.fp != nil {
		//fmt.Printf("close file ...\n")
		err = w.fp.Sync()
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			fmt.Printf("rotating error on close current log: %v\n", err)
			return
		}
	} else {
		fmt.Printf("rotating log has no file to close ...\n")
	}

	if w.rotateFilesByNumber {

		var logBaseName = w.filename[0 : len(w.filename)-len(path.Ext(w.filename))] //removed extension from filename
		var logName string
		for i := w.maxRotatedFilesByNumber; i >= 0; i-- {

			//time.Sleep(2 * time.Second)
			//fmt.Printf("i: %d\n", i)

			logName = fmt.Sprintf("%s_%d%s", logBaseName, i, path.Ext(w.filename))

			if i == 0 {
				logNextName := fmt.Sprintf("%s_1%s", logBaseName, path.Ext(w.filename))
				//fmt.Printf("log rotate: renaming '%s' to '%s'\n", w.filename, logNextName)

				_, errs := os.Stat(w.filename)
				if errs != nil {
					fmt.Printf("log rotate: error stat '%s': %v\n", w.filename, errs)
				}

				err = os.Rename(w.filename, logNextName)
				if err != nil {
					//fmt.Printf("log rotate: error rotating '%s' to '%s': %v\n", w.filename, logNextName, err)
				}
			} else if i == w.maxRotatedFilesByNumber {
				//fmt.Printf("log rotate: removing '%s'\n", logName)
				err = os.Remove(logName)
				if err != nil {
					//fmt.Printf("log rotate: error deleting '%s': %v\n", logName, err)
				}
			} else {
				logNextName := fmt.Sprintf("%s_%d%s", logBaseName, i+1, path.Ext(w.filename))
				//fmt.Printf("log rotate: renaming '%s' to '%s'\n", logName, logNextName)
				err = os.Rename(logName, logNextName)
				if err != nil {
					//fmt.Printf("log rotate: error rotating '%s' to '%s': %v\n", logName, logNextName, err)
				}
			}

		}

	} else {
		// Rename dest file if it already exists
		_, err = os.Stat(w.filename)
		if err == nil {

			t := time.Now()

			//d, n := path.Split(w.filename)
			newName := w.filename[0 : len(w.filename)-len(path.Ext(w.filename))] //removed extension from filename

			newName = fmt.Sprintf("%s_%04d%02d%02dT%02d%02d%02d%s", newName, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), path.Ext(w.filename))

			fmt.Printf("rename '%s' to '%s' ...\n", w.filename, newName)
			err = os.Rename(w.filename, newName)
			if err != nil {
				fmt.Printf("rotating error on rename: %v\n", err)
				return
			}
		} else {
			fmt.Printf("stat %v failed: %v\n", w.filename, err)
		}
	}

	// Create a file.
	fmt.Printf("Creating '%s' ...\n", w.filename)
	w.fp, err = os.OpenFile(w.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) //os.Create(w.filename)
	if err != nil {
		fmt.Printf("rotating log error on create: %v\n", err)
		return
	}
	w.writterBytes = 0

	return
}

/*
func main() {
	maxRotateWriter := NewMaxRotateWriter("maxrotate.log", 100000)
	var buf string
	for i := 0; ; i++ {
		fmt.Println(i)
		buf = fmt.Sprintf("line #%d ========================================\n", i)
		maxRotateWriter.Write([]byte(buf))
		time.Sleep(10 * time.Millisecond)
	}
}
*/
