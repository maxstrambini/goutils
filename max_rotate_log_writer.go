/*

max_rotate_log_writer.go

[2018-05-10] a custom file writer to rotate logs

## descrizione: quando il file supera 'maxBytes' viene chiuso, rinominato con la data e l'ora e ricreato

## esempio di utilizzo:

			func main() {
				maxRotateWriter := NewMaxRotateWriter("maxrotate.log", 100000) //log file name, max bytes per log file
				var buf string
				for i := 0; ; i++ {
					fmt.Println(i)
					buf = fmt.Sprintf("line #%d ========================================\n", i)
					maxRotateWriter.Write([]byte(buf))
					time.Sleep(10 * time.Millisecond)
				}
			}

## TODO:
- ruotare usando numeri invece che date
- cancellare vecchi log

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
	lock         sync.Mutex
	filename     string // should be set to the actual filename
	writterBytes int    // rotate > maxBytes
	maxBytes     int    // counter of written bytes
	fp           *os.File
}

// NewMaxRotateWriter Make a new MaxRotateWriter. Return nil if error occurs during setup.
func NewMaxRotateWriter(filename string, maxBytes int) *MaxRotateWriter {
	w := &MaxRotateWriter{filename: filename, maxBytes: maxBytes}
	err := w.Rotate()
	if err != nil {
		return nil
	}
	return w
}

// Write satisfies the io.Writer interface.
func (w *MaxRotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
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

	fmt.Printf("rotating ...\n")
	//time.Sleep(3 * time.Second)

	// Close existing file if open
	if w.fp != nil {
		fmt.Printf("close file ...\n")
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			fmt.Printf("rotating error on close (1/3) -> %v\n", err)
			return
		}
	} else {
		fmt.Printf("no file to close ...\n")
	}

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
			fmt.Printf("rotating error on rename (2/3) -> %v\n", err)
			return
		}
	} else {
		fmt.Printf("stat %v failed\n", w.filename)
	}

	// Create a file.
	fmt.Printf("create '%s' ...\n", w.filename)
	w.fp, err = os.Create(w.filename)
	if err != nil {
		fmt.Printf("rotating error on create (3/3) -> %v\n", err)
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
