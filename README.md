# goutils

Some golang utility functions I use


## History

[2018-05-10] added 'max_rotate_log_writer.go' defining 'MaxRotateWriter': a rotating file writer for logs [2018-05-10]

[2018-05-10] 'max_rotate_log_writer.go' can now rotate by number instaed of dates with a maximum number

[2019-03-25] added GetNowForFileName()

[2019-03-28] added PrettyPrintStruct() printing to stdout and PrettyFormatStruct() returning a string
They both can be used to print a struct one field per line with name and value

    for example: 
        goutils.PrettyPrintStruct(mystruct)
        or
        fmt.Printf("%s\n", goutils.PrettyFormatStruct(mystruct))

[2019-09-11] logs written to file and also to stdout

[2019-09-11] added PrettyPrintStructSimple, PrettyPrintInterface

## Example:

				package main

				import "log"
				import "github.com/maxstrambini/goutils"

				func main() {

					//print version:
					log.Printf(goutils.PackageVersion)

					log.Printf(goutils.GetNowString())

					// BEGIN test MaxRotateWriter: a rotating file writer for logs [2018-05-10]
					maxRotateWriter := NewMaxRotateWriter("maxrotate.log", 10000)
					var buf string
					for i := 0; i<10000; i++ {
						fmt.Println(i)
						buf = fmt.Sprintf("line #%d ========================================\n", i)
						maxRotateWriter.Write([]byte(buf))
						time.Sleep(10 * time.Millisecond)
					}
					// END test MaxRotateWriter

					/*
					//ExistsPath example:
					if !goutils.ExistsPath(zipFullName) {
						log.Printf("zipFilePathFromAssetID could not find '%s' returns default!", zipFullName)
						zipFullName = "./default.zip"
					}
					*/

					/*
					succ1 := goutils.WriteTextToFile("/Users/max/temp/test.txt", "Ciao da Max")
					log.Printf("goutils.WriteStub returned: '%v'", succ1)
					*/

					/*
					//test config:
					var conf goutils.ConfigReader

					log.Printf("This should *log* an error -> Total items found: '%v'", conf.CountItems())

					i, err := conf.Read("not-existing.ini")
					log.Printf("This should *return* error -> items found: '%v', err: '%v'", i, err)

					log.Println("==========")

					i, err = conf.Read("config.ini")
					log.Printf("items found: '%v', err: '%v'", i, err)

					conf.PrintItems(true)

					log.Printf("Total items found: '%v'", conf.CountItems())

					key := "name"
					v, found := conf.GetString(key)
					log.Printf("key: '%s' -> found: '%v', value: '%v'", key, found, v)

					key = "xxx"
					v, found = conf.GetString(key)
					log.Printf("key: '%s' -> found: '%v', value: '%v'", key, found, v)

					iv := 0
					key = "name"
					iv, found = conf.GetInt(key)
					log.Printf("this should fail -> key: '%s' -> found: '%v', value: '%v'", key, found, iv)

					key = "myint"
					iv, found = conf.GetInt(key)
					log.Printf("this should work -> key: '%s' -> found: '%v', value: '%v'", key, found, iv)

					bv := false
					key = "mybool_missing"
					bv, found = conf.GetBool(key)
					log.Printf("this should not be found -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)

					key = "mybool_1"
					bv, found = conf.GetBool(key)
					log.Printf("this should be true -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)

					key = "mybool_2"
					bv, found = conf.GetBool(key)
					log.Printf("this should be true -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)

					key = "mybool_3"
					bv, found = conf.GetBool(key)
					log.Printf("this should be true -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)

					key = "mybool_4"
					bv, found = conf.GetBool(key)
					log.Printf("this should be false -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)

					key = "mybool_5"
					bv, found = conf.GetBool(key)
					log.Printf("this should be true -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)

					key = "mybool_6"
					bv, found = conf.GetBool(key)
					log.Printf("this should be false -> key: '%s' -> found: '%v', value: '%v'", key, found, bv)
					*/					

				}



