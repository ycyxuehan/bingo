package bingdb

import (

)
//DBMInterface interface for database model
type DBMInterface interface {
	Columns()([]string, []interface{})
}



//DBInterface interface for database
type DBInterface interface {
	Connect()error
	ConnectEx(string)error
	SelectOne(string, interface{}, ...interface{})(interface{}, error)
	Select(string, interface{}, ...interface{})(interface{}, error)
	Update(string, interface{}, ...interface{})(interface{}, error)
	Insert(string, ...interface{})(interface{}, error)
	Delete(string, ...interface{})(interface{}, error)
	InsertBatch(string, ...interface{})(interface{}, error)
	Count(string, interface{})(int, error)
	UnMarshal(interface{}, ...interface{})error
	UnMarshalI(interface{}, []string, interface{})([]interface{},error)
}
