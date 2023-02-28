package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultMySQLport    = 3306
	defaultMySQLhost    = "polaris"
	defaultMySQLname    = "polaris"
	defaultMySQLcharset = "utf8"
)

type config struct {
	myInt  int
	myStr  string
	myBool bool
	DB     dbConfig
}

type dbConfig struct {
	username string
	password string
	host     string
	port     int
	name     string
	charset  string
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

	v, ok = os.LookupEnv("DBUSER")
	if ok {
		c.DB.username = v
	} else {
		log.Println("environment variable DBUSER not found")
	}

	v, ok = os.LookupEnv("DBPASSWD")
	if ok {
		c.DB.password = v
	} else {
		log.Println("environment variable DBPASSWD not found")
	}

	v, ok = os.LookupEnv("DBHOST")
	if ok {
		c.DB.host = v
	} else {
		log.Println("environment variable DBHOST not foundt, using default value")
		c.DB.host = defaultMySQLhost
	}

	v, ok = os.LookupEnv("DBPORT")
	if ok {
		c.DB.port, err = strconv.Atoi(v)
		if err != nil {
			log.Println("error converting DBPORT to int, using default value")
			c.DB.port = defaultMySQLport
		}
	} else {
		log.Println("environment variable DBPORT not found, using default value")
		c.DB.port = defaultMySQLport
	}

	v, ok = os.LookupEnv("DBNAME")
	if ok {
		c.DB.name = v
	} else {
		log.Println("environment variable DBNAME not foundt, using default value")
		c.DB.name = defaultMySQLname
	}

	v, ok = os.LookupEnv("DBCHARSET")
	if ok {
		c.DB.charset = v
	} else {
		log.Println("environment variable DBCHARSET not foundt, using default value")
		c.DB.charset = defaultMySQLcharset
	}

	return c
}
