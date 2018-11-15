package rdb

import "errors"

var (
	// ErrRecordNotFound is not found any record in database
	ErrRecordNotFound = errors.New("record not found")
	// ErrCanNotCreateRecord is can not create record
	ErrCanNotCreateRecord = errors.New("can not create record")
)
