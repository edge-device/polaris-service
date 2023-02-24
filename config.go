package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type config struct {
	myInt  int
	myStr  string
	myBool bool
}

// getConfig() is used to get config values from environment variables.
func getConfig() config {
	var c config
	var err error
	v, ok := os.LookupEnv("MYINTVAR")
	if ok {
		c.myInt, err = strconv.Atoi(v)
		if err != nil {
			log.Println("error converting MYINTVAR to int, using default value")
			c.myInt = 1
		}
	} else {
		log.Println("environment variable MYINTVAR not found, using default value")
		c.myInt = 1
	}

	v, ok = os.LookupEnv("MYSTRVAR")
	if ok {
		c.myStr = v
	} else {
		log.Println("environment variable MYSTRVAR not found, using default value")
		c.myStr = "My String"
	}

	v, ok = os.LookupEnv("MYBOOLVAR")
	if ok {
		if strings.ToLower(v) == "true" {
			c.myBool = true
		} else if strings.ToLower(v) == "false" {
			c.myBool = false
		} else {
			log.Println("invalid MYBOOLVAR value, using default value")
			c.myBool = true
		}
	} else {
		log.Println("environment variable MYBOOLVAR not found, using default value")
		c.myBool = true
	}

	return c
}
