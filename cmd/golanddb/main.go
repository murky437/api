package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: golanddb <source_db> <dest_db>")
		os.Exit(1)
	}

	src := os.Args[1]
	dest := os.Args[2]

	// Copy the DB file
	if err := copyFile(src, dest); err != nil {
		log.Fatal("Failed to copy DB:", err)
	}
	fmt.Println("Copied DB to:", dest)

	// Open the copy and switch to DELETE journal mode
	db, err := sql.Open("sqlite", dest)
	if err != nil {
		log.Fatal("Failed to open copied DB:", err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA journal_mode=DELETE;")
	if err != nil {
		log.Fatal("Failed to set journal_mode=DELETE:", err)
	}

	fmt.Println("Set journal_mode=DELETE on copied DB")
}

// copyFile copies the contents of src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
