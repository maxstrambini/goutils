/*

main for testing

*/

package main

import (
	"fmt"
	"time"

	_ "github.com/maxstrambini/goutils"
)

func main() {
	//demo
	maxRotateWriter2 := NewMaxRotateWriter2("maxrotate2.log", 5000, true, 20)
	var buf string
	for i := 0; ; i++ {
		buf = fmt.Sprintf("line #%d ========================================\n", i)
		maxRotateWriter2.Write([]byte(buf))
		time.Sleep(10 * time.Millisecond)
	}
}
