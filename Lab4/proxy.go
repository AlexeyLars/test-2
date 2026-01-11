package main

import "fmt"

type Database interface {
	Query(sql string)
}

// RealDatabase - реальная база данных
type RealDatabase struct{}

func (db *RealDatabase) Query(sql string) {
	fmt.Printf("Executing query: %s\n", sql)
}

// DatabaseProxy - прокси для базы данных
type DatabaseProxy struct {
	realDatabase *RealDatabase
	hasAccess    bool
}

func NewDatabaseProxy(hasAccess bool) *DatabaseProxy {
	return &DatabaseProxy{
		realDatabase: &RealDatabase{},
		hasAccess:    hasAccess,
	}
}

func (proxy *DatabaseProxy) Query(sql string) {
	if proxy.hasAccess {
		proxy.realDatabase.Query(sql)
	} else {
		fmt.Println("Access denied. Query cannot be executed.")
	}
}

func main() {
	fmt.Println("=== Proxy Pattern ===")

	userDb := NewDatabaseProxy(false)
	adminDb := NewDatabaseProxy(true)

	fmt.Println("User trying to query:")
	userDb.Query("SELECT * FROM users")

	fmt.Println("Admin trying to query:")
	adminDb.Query("SELECT * FROM users")
	fmt.Println()
}
