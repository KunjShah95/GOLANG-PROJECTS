package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Record represents a single row in the database
type Record struct {
	ID    int
	Name  string
	Value string
}

// VersionedRecord represents a record with a version for MVCC
type VersionedRecord struct {
	Record  Record
	Version int
}

// Table represents a database table
type Table struct {
	records map[int]VersionedRecord
	indexes map[string]map[string]int
	mu      sync.RWMutex
}

// Database represents an in-memory database
type Database struct {
	tables map[string]*Table
	cache  map[string]Record
	mu     sync.RWMutex
}

// NewDatabase creates a new database
func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
		cache:  make(map[string]Record),
	}
}

// CreateTable creates a new table in the database
func (db *Database) CreateTable(name string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.tables[name] = &Table{
		records: make(map[int]VersionedRecord),
		indexes: make(map[string]map[string]int),
	}
}

// Insert adds a new record to a table
func (db *Database) Insert(ctx context.Context, tableName string, record Record) error {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return errors.New("table does not exist")
	}

	table.mu.Lock()
	defer table.mu.Unlock()
	table.records[record.ID] = VersionedRecord{Record: record, Version: 1}
	for _, index := range table.indexes {
		index[record.Name] = record.ID
	}
	return nil
}

// Get retrieves a record by ID from a table, using cache if available
func (db *Database) Get(ctx context.Context, tableName string, id int) (Record, bool) {
	cacheKey := fmt.Sprintf("%s:%d", tableName, id)
	if record, found := db.cache[cacheKey]; found {
		return record, true
	}

	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return Record{}, false
	}

	table.mu.RLock()
	defer table.mu.RUnlock()
	vRecord, exists := table.records[id]
	if exists {
		db.cache[cacheKey] = vRecord.Record // Add to cache
	}
	return vRecord.Record, exists
}

// Update modifies an existing record in a table
func (db *Database) Update(ctx context.Context, tableName string, id int, newRecord Record) error {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return errors.New("table does not exist")
	}

	table.mu.Lock()
	defer table.mu.Unlock()
	vRecord, exists := table.records[id]
	if !exists {
		return errors.New("record does not exist")
	}
	newVersion := vRecord.Version + 1
	table.records[id] = VersionedRecord{Record: newRecord, Version: newVersion}
	for _, index := range table.indexes {
		index[newRecord.Name] = newRecord.ID
	}
	return nil
}

// Delete removes a record by ID from a table
func (db *Database) Delete(ctx context.Context, tableName string, id int) error {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return errors.New("table does not exist")
	}

	table.mu.Lock()
	defer table.mu.Unlock()
	vRecord, exists := table.records[id]
	if !exists {
		return errors.New("record does not exist")
	}
	delete(table.records, id)
	for _, index := range table.indexes {
		delete(index, vRecord.Record.Name)
	}
	return nil
}

// List returns all records in a table
func (db *Database) List(ctx context.Context, tableName string) ([]Record, error) {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return nil, errors.New("table does not exist")
	}

	table.mu.RLock()
	defer table.mu.RUnlock()
	records := make([]Record, 0, len(table.records))
	for _, vRecord := range table.records {
		records = append(records, vRecord.Record)
	}
	return records, nil
}

// CreateIndex creates an index on a column in a table
func (db *Database) CreateIndex(ctx context.Context, tableName string, columnName string) error {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return errors.New("table does not exist")
	}

	table.mu.Lock()
	defer table.mu.Unlock()
	index := make(map[string]int)
	for _, vRecord := range table.records {
		index[vRecord.Record.Name] = vRecord.Record.ID
	}
	table.indexes[columnName] = index
	return nil
}

// Query retrieves records from a table based on a column value
func (db *Database) Query(ctx context.Context, tableName string, columnName string, value string) ([]Record, error) {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return nil, errors.New("table does not exist")
	}

	table.mu.RLock()
	defer table.mu.RUnlock()
	index, exists := table.indexes[columnName]
	if !exists {
		return nil, errors.New("index does not exist")
	}

	recordID, exists := index[value]
	if !exists {
		return nil, errors.New("record not found")
	}

	vRecord, exists := table.records[recordID]
	if !exists {
		return nil, errors.New("record not found")
	}

	return []Record{vRecord.Record}, nil
}

// Save to persist the database to a file
func (db *Database) Save(filename string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(db.tables)
}

// Load to load the database from a file
func (db *Database) Load(filename string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&db.tables)
}

// Simple query parser for demonstration purposes
func (db *Database) QueryAdvanced(ctx context.Context, query string) ([]Record, error) {
	parts := strings.Fields(query)
	if len(parts) < 4 || strings.ToLower(parts[0]) != "select" || strings.ToLower(parts[2]) != "from" {
		return nil, errors.New("invalid query")
	}

	tableName := parts[3]
	columnName := parts[1] // Assuming simple select column from table
	if strings.ToLower(columnName) == "*" {
		return db.List(ctx, tableName)
	}

	value := ""
	if len(parts) > 5 && strings.ToLower(parts[4]) == "where" {
		value = parts[5] // Assuming where column = value
	}

	return db.Query(ctx, tableName, columnName, value)
}

// Transaction represents a database transaction
type Transaction struct {
	db       *Database
	table    *Table
	ctx      context.Context
	cancel   context.CancelFunc
	versions map[int]int // Track versions of records in the transaction
}

// BeginTransaction starts a new transaction
func (db *Database) BeginTransaction(ctx context.Context, tableName string) (*Transaction, error) {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()
	if !exists {
		return nil, errors.New("table does not exist")
	}

	ctx, cancel := context.WithCancel(ctx)
	return &Transaction{db: db, table: table, ctx: ctx, cancel: cancel, versions: make(map[int]int)}, nil
}

// Get retrieves a record by ID within a transaction
func (tx *Transaction) Get(id int) (Record, bool) {
	tx.table.mu.RLock()
	defer tx.table.mu.RUnlock()
	vRecord, exists := tx.table.records[id]
	if !exists {
		return Record{}, false
	}
	if tx.versions[id] > 0 && tx.versions[id] != vRecord.Version {
		return Record{}, false
	}
	return vRecord.Record, true
}

// Commit commits the transaction
func (tx *Transaction) Commit() {
	tx.cancel()
}

// Rollback rolls back the transaction
func (tx *Transaction) Rollback() {
	tx.cancel()
}

func main() {
	db := NewDatabase()

	// Create a table
	db.CreateTable("users")

	// Insert records
	ctx := context.Background()
	db.Insert(ctx, "users", Record{ID: 1, Name: "Alice", Value: "Value1"})
	db.Insert(ctx, "users", Record{ID: 2, Name: "Bob", Value: "Value2"})

	// Create an index on the Name column
	db.CreateIndex(ctx, "users", "Name")

	// Query a record by Name
	records, err := db.Query(ctx, "users", "Name", "Alice")
	if err != nil {
		fmt.Println("Query error:", err)
	} else {
		fmt.Println("Query result:", records)
	}

	// Get and print a record by ID
	record, exists := db.Get(ctx, "users", 1)
	if exists {
		fmt.Printf("Record: %+v\n", record)
	} else {
		fmt.Println("Record not found")
	}

	// Update a record
	err = db.Update(ctx, "users", 1, Record{ID: 1, Name: "AliceUpdated", Value: "UpdatedValue1"})
	if err != nil {
		fmt.Println("Update error:", err)
	} else {
		fmt.Println("Record updated")
	}

	// List all records
	records, err = db.List(ctx, "users")
	if err != nil {
		fmt.Println("List error:", err)
	} else {
		fmt.Println("All Records:")
		for _, rec := range records {
			fmt.Printf("%+v\n", rec)
		}
	}

	// Delete a record
	err = db.Delete(ctx, "users", 1)
	if err != nil {
		fmt.Println("Delete error:", err)
	} else {
		fmt.Println("Record deleted")
	}

	// List all records after deletion
	records, err = db.List(ctx, "users")
	if err != nil {
		fmt.Println("List error:", err)
	} else {
		fmt.Println("All Records after deletion:")
		for _, rec := range records {
			fmt.Printf("%+v\n", rec)
		}
	}

	// Begin a transaction
	tx, err := db.BeginTransaction(ctx, "users")
	if err != nil {
		fmt.Println("Transaction error:", err)
		return
	}

	// Insert a record in the transaction
	err = db.Insert(tx.ctx, "users", Record{ID: 3, Name: "Charlie", Value: "Value3"})
	if err != nil {
		fmt.Println("Insert error:", err)
		tx.Rollback()
		return
	}

	// Commit the transaction
	tx.Commit()
}