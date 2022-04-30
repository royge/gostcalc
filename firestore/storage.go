package firestore

import (
	"context"
	"math/big"
	"strings"
	"time"
)

const (
	// MonthlyFreeStorage is the 1GB free monthly storage.
	MonthlyFreeStorage = 1073741824

	// MonthNumOfDays is the fix number of days in a month.
	MonthNumOfDays = 30

	// OneGB is 1GB in bytes.
	OneGB = 1000000000

	// PricePerGB is storage price per GB after free.
	PricePerGB = 0.18

	// DocumentNamePadding is the default additional bytes for document name.
	DocumentNamePadding = 16 // bytes

	// DocumentPadding is the default additional bytes for document fields and
	// values.
	DocumentPadding = 32 // bytes
)

type DailyStorageCalculator struct {
	// Document on disk.
	Document *Document
}

func (ds *DailyStorageCalculator) Calculate(_ context.Context, count *big.Int) (*big.Float, error) {
	// Get document size in bytes.
	size := big.NewInt(ds.Document.Size())

	daily := size.Mul(size, count)

	return new(big.Float).SetInt(daily), nil
}

type MonthlyStorageCalculator struct {
	D *DailyStorageCalculator

	// Unit Price.
	// Price per GB.
	Price float64
}

func (ms *MonthlyStorageCalculator) Calculate(ctx context.Context, count *big.Int) (*big.Float, error) {
	daily, err := ms.D.Calculate(ctx, count)
	if err != nil {
		return new(big.Float), err
	}
	days := new(big.Float).SetInt64(MonthNumOfDays)

	monthly := daily.Mul(daily, days)

	free := new(big.Float).SetInt64(MonthlyFreeStorage)

	monthly = monthly.Sub(monthly, free)
	monthly = monthly.Quo(monthly, new(big.Float).SetInt64(OneGB))

	cost := monthly.Mul(monthly, new(big.Float).SetFloat64(ms.Price))

	return cost, nil
}

// Document defines the stored Firestore document.
type Document struct {
	// ID is the document id.
	ID string

	// Collection is the collection name.
	Collection string

	// Data contains the document fields and values.
	Data map[string]interface{}

	// SingleFieldIndexes represents a single-field indexes.
	SingleFieldIndexes []map[string]interface{}

	// CompositeIndexes represents composite indexes.
	CompositeIndexes []map[string]interface{}
}

// Calculate the document name total size size.
func (d *Document) nameSize() int64 {
	idSize := len(d.ID) + 1
	cols := strings.Split(d.Collection, "/")

	collSize := len(cols)
	for _, c := range cols {
		c := c
		collSize += len(c)
	}

	size := int64(idSize + collSize + DocumentNamePadding)

	return size
}

// Calculate the parent document name total size size.
func (d *Document) parentNameSize() (size int64) {
	cols := strings.Split(d.Collection, "/")

	parent := cols[:len(cols)-1]

	size += int64(len(parent)) + DocumentNamePadding
	for _, c := range parent {
		c := c
		size += int64(len(c))
	}

	if len(cols) < 3 {
		return 0
	}

	return size
}

// Calculate the document fields and values total size.
func (d *Document) dataSize() int64 {
	return getSize(d.Data)
}

func getSize(data map[string]interface{}) int64 {
	size := DocumentPadding
	size += len(data) // +1s for every field

	for k, v := range data {
		k, v := k, v
		size += len(k)
		size += getValueSize(v)
	}

	return int64(size)
}

// Calculate the total size of composite index.
// It is only for collection scope.
func (d *Document) singleFieldIndexSize() (size int64) {
	for _, sfi := range d.SingleFieldIndexes {
		sfi := sfi // to avoid possible race
		size += d.nameSize() + d.parentNameSize() + DocumentPadding

		for _, v := range sfi {
			v := v // to avoid possible race
			s := getValueSize(v)

			size += int64(s)
		}
	}

	return size
}

// Calculate the total size of composite index.
// It is only for collection scope.
func (d *Document) compositeIndexSize() (size int64) {
	for _, ci := range d.CompositeIndexes {
		ci := ci // to avoid possible race
		size += d.nameSize() + d.parentNameSize() + DocumentPadding

		for _, v := range ci {
			v := v // to avoid possible race
			s := getValueSize(v)

			size += int64(s)
		}
	}

	return size
}

// Size calculate and return document size.
func (d *Document) Size() (size int64) {
	size = d.nameSize() +
		d.dataSize() +
		d.singleFieldIndexSize() +
		d.compositeIndexSize()

	return size
}

// Get value size.
// TODO: Create unit tests.
func getValueSize(val interface{}) int {
	switch v := val.(type) {
	case string:
		return len(v) + 1
	case bool, byte:
		return 1
	case int, float64, time.Time:
		return 8
	case map[string]interface{}:
		return int(getSize(v))
	default:
		return 0
	}
}
