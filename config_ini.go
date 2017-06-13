package goutils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//ConfigReader struct representing a config file
type ConfigReader struct {
	configFilePath string
	nItems         int
	items          map[string]string
}

//Read and parse a configuration file
func (c *ConfigReader) Read(configPath string) (numItemsFound int, err error) {

	c.configFilePath = configPath
	c.items = make(map[string]string)

	var file *os.File
	file, err = os.Open(configPath)
	if err != nil {
		log.Printf("goutils.ConfigReader.Read(%s) open error: %s", configPath, err)
		return
	}
	defer file.Close()

	tempItems := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) > 0 && t[0] != '#' {

			if strings.Contains(t, "=") {

				a := strings.SplitN(t, "=", 2)

				if len(a) == 2 {
					//items += 1
					tempItems[a[0]] = a[1]
					log.Printf("from '%s' decoded key: '%v', value: '%v'", t, a[0], a[1])

				} else {
					log.Printf("goutils.ConfigReader.Read(%s) scan invalid line: '%s'", configPath, t)
				}

			} else {
				log.Printf("goutils.ConfigReader.Read(%s) scan invalid line: '%s'", configPath, t)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		log.Printf("goutils.ConfigReader.Read(%s) scan error: %s", configPath, err)
		return
	}

	numItemsFound = len(tempItems)
	c.items = tempItems
	log.Printf("goutils.ConfigReader.Read(%s) success, decoded %d items", configPath, numItemsFound)
	return
}

//GetString get the string value of a configured item
func (c *ConfigReader) GetString(itemName string) (itemValue string, found bool) {
	if c.items != nil {
		itemValue, found = c.items[itemName]
	} else {
		log.Printf("goutils.ConfigReader.GetString failed, there are no items read")
	}
	return
}

//GetInt get the integer value 32 bit of a configured item
func (c *ConfigReader) GetInt(itemName string) (itemValue int, found bool) {
	if c.items != nil {
		stemp := ""
		stemp, found = c.items[itemName]
		if found {
			i64, err := strconv.ParseInt(stemp, 10, 32)
			if err == nil {
				itemValue = int(i64)
				found = true
			} else {
				found = false
				log.Printf("goutils.ConfigReader.GetInt found item but failed conversion to integer: '%v'", err)
			}
		}
	} else {
		log.Printf("goutils.ConfigReader.GetInt failed, there are no items read")
	}
	return
}

//GetInt64 get the integer 64 value of a configured item
func (c *ConfigReader) GetInt64(itemName string) (itemValue int64, found bool) {
	if c.items != nil {
		stemp := ""
		stemp, found = c.items[itemName]
		if found {
			i64, err := strconv.ParseInt(stemp, 10, 64)
			if err == nil {
				itemValue = i64
				found = true
			} else {
				found = false
				log.Printf("goutils.ConfigReader.GetInt64 found item but failed conversion to integer: '%v'", err)
			}
		}
	} else {
		log.Printf("goutils.ConfigReader.GetInt64 failed, there are no items read")
	}
	return
}

//GetBool get the boolean value of a configured item
func (c *ConfigReader) GetBool(itemName string) (itemValue bool, found bool) {
	if c.items != nil {
		stemp := ""
		stemp, found = c.items[itemName]
		if found {
			b, err := strconv.ParseBool(stemp)
			if err == nil {
				itemValue = b
				found = true
			} else {
				found = false
				log.Printf("goutils.ConfigReader.GetBool found item but failed conversion to bool: '%v'", err)
			}
		}
	} else {
		log.Printf("goutils.ConfigReader.GetBool failed, there are no items read")
	}
	return
}

//CountItems get the total number of counfigured items
func (c *ConfigReader) CountItems() (numItemsFound int) {
	if c.items != nil {
		numItemsFound = len(c.items)
	} else {
		log.Printf("goutils.ConfigReader.CountItems failed, there are no items read")
	}
	return
}

// PrintItems print all items found while reading the config
// The parameter toLog directs the output to the log when it is true otherwise to the console
func (c *ConfigReader) PrintItems(toLog bool) {
	if c.items != nil {
		if toLog {
			log.Printf("Elements in [%s]:\n", c.configFilePath)
		} else {
			fmt.Printf("Elements in [%s]:\n", c.configFilePath)
		}
		i := 0
		for k, v := range c.items {
			i++
			if toLog {
				log.Printf("- [%d/%d] key[%s] value[%s]\n", i, len(c.items), k, v)
			} else {
				fmt.Printf("- [%d/%d] key[%s] value[%s]\n", i, len(c.items), k, v)
			}
		}
	} else {
		log.Printf("goutils.ConfigReader.GetString failed, there are no items read")
	}
	return
}
