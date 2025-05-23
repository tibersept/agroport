package main

import (
	"apps/migrations/sql_scripts"
	"fmt"
)

func Hello(name string) string {
	result := "Hello " + name
	return result
}

func main() {
	fmt.Println(Hello("migrations"))

	content, err := sql_scripts.Asset("../../schema/database/000001_create_initial_tables.up.sql")
	if err != nil {
		fmt.Println("An error occurred ...")
	}

	fmt.Println(string(content))
}
