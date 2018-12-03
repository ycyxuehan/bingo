package bingdb

import (

)
//DBMInterface interface for database model
type DBMInterface interface {
	Columns()map[string]interface{}
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
	Exists(string, ...interface{})(bool, error)
	UnMarshal(interface{}, ...interface{})error
	UnMarshalI(interface{}, interface{}, interface{})([]interface{},error)
}