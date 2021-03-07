package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
}

func main() {
	//Parse flags
	path, format := parseFlags()
	//Collect users information
	users := collectUsers()

	//Write to file
	var output io.Writer

	// If path is empty then write to stdout, if not, then write to file
	if path != "" {
		f, err := os.Create(path)
		handleError(err)
		defer f.Close()
		output = f
	} else {
		output = os.Stdout
	}

	//Handle format of output
	if format == "json" {
		data, err := json.MarshalIndent(users, "", "  ") //we indent by 2 spaces
		handleError(err)
		output.Write(data)
	} else if format == "csv" {
		output.Write([]byte("name,id,home,shell\n"))
		//Opens csv writer
		writer := csv.NewWriter(output)
		for _, user := range users {
			err := writer.Write([]string{user.Name, strconv.Itoa(user.Id), user.Home, user.Shell})
			handleError(err)
		}
		//Need to call flush after writing individually
		writer.Flush()
	}

}

//Parse flags funciton
func parseFlags() (path, format string) {
	flags.StringVar(&path, "path", "", "the path to export file.")
	flags.StringVar(&format, "format", "json", "the output format for the user information. Available options are 'csv' and 'json'. The default is json")
	flag.Parse()

	format = strings.ToLower(format)

	if format != "csv" && format != "json" {
		fmt.Println("Error: invalid format. Use 'json' or 'csv'")
		flag.Usage()
		os.Exit(1)
	}
}

//Function which handles errors
func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Function which return clice of users with struct User
func collectUsers() (users []User) {
	//Open a file
	f, err := os.Open("/etc/passwd")
	//Handle error if problems with opening
	handleError(err)
	defer f.Close()

	//Create csv reader
	reader := csv.NewReader(f)
	//Change csv delimiter to :
	reader.Comma = ':'

	//Now read all lines
	lines, err := reader.ReadAll()
	handleError(err)

	//nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
	//Because output will be string, we have to convert ID to int
	//So second element of line is UID
	for _, line := range lines {
		id, err := strconv.ParseInt(line[2], 10, 64)
		handleError(err)

		//Ignore system users (with UID <1000)
		if id < 1000 {
			continue
		}
		//make user
		user := User{
			Name:  line[0],
			Id:    int(id),
			Home:  line[5],
			Shell: line[6],
		}
		//append users to slice
		users = append(users, user)

	}

	return
}
