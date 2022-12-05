package db

import (
	"errors"
)

// Error messages
var (
	ErrMissingDriverName = errors.New(`keeper: missing database driver name`)
	ErrNotConnected      = errors.New(`keeper: not connected to a database`)
	ErrNilRecord         = errors.New(`keeper: invalid item (nil)`)
	ErrInvalidCollection = errors.New(`keeper: invalid collection`)

	ErrMissingAdapter           = errors.New(`keeper: missing adapter`)
	ErrAlreadyWithinTransaction = errors.New(`keeper: already within a transaction`)
	ErrCollectionDoesNotExist   = errors.New(`keeper: collection does not exist`)
	ErrExpectingNonNilModel     = errors.New(`keeper: expecting non nil model`)
	ErrExpectingPointerToStruct = errors.New(`keeper: expecting pointer to struct`)
	ErrGivingUpTryingToConnect  = errors.New(`keeper: giving up trying to connect: too many clients`)
	ErrMissingCollectionName    = errors.New(`keeper: missing collection name`)
	ErrMissingConditions        = errors.New(`keeper: missing selector conditions`)
	ErrMissingConnURL           = errors.New(`keeper: missing DSN`)
	ErrMissingDatabaseName      = errors.New(`keeper: missing database name`)
	ErrNoMoreRows               = errors.New(`keeper: no more rows in this result set`)
	ErrNotImplemented           = errors.New(`keeper: call not implemented`)
	ErrQueryIsPending           = errors.New(`keeper: can't execute this instruction while the result set is still open`)
	ErrQueryLimitParam          = errors.New(`keeper: a query can accept only one limit parameter`)
	ErrQueryOffsetParam         = errors.New(`keeper: a query can accept only one offset parameter`)
	ErrQuerySortParam           = errors.New(`keeper: a query can accept only one order-by parameter`)
	ErrSockerOrHost             = errors.New(`keeper: you may connect either to a UNIX socket or a TCP address, but not both`)
	ErrTooManyClients           = errors.New(`keeper: can't connect to database server: too many clients`)
	ErrUndefined                = errors.New(`keeper: value is undefined`)
	ErrUnknownConditionType     = errors.New(`keeper: arguments of type %T can't be used as constraints`)
	ErrUnsupported              = errors.New(`keeper: action is not supported by the DBMS`)
	ErrUnsupportedDestination   = errors.New(`keeper: unsupported destination type`)
	ErrUnsupportedType          = errors.New(`keeper: type does not support marshaling`)
	ErrUnsupportedValue         = errors.New(`keeper: value does not support unmarshaling`)
	ErrRecordIDIsZero           = errors.New(`keeper: item ID is not defined`)
	ErrMissingPrimaryKeys       = errors.New(`keeper: collection %q has no primary keys`)
	ErrWarnSlowQuery            = errors.New(`keeper: slow query`)
	ErrTransactionAborted       = errors.New(`keeper: transaction was aborted`)
	ErrNotWithinTransaction     = errors.New(`keeper: not within transaction`)
	ErrNotSupportedByAdapter    = errors.New(`keeper: not supported by adapter`)
)
