package model

import (
	"github.com/ycyxuehan/bingo/bingdb"
	"strings"
	"fmt"
)

//BaseInterface base interface
type BaseInterface interface{
	Table()string
	Filter()map[string]interface{}
	Columns()map[string]interface{}
}

//Base base model
type Base struct {
	this BaseInterface
}

//SetThis set this
func (b *Base)SetThis(bi BaseInterface){
	b.this = bi
}

//Exists is this exists in database
func (b *Base)Exists(dbi bingdb.DBInterface)bool{
	if b.this == nil {
		return false
	}
	args := []interface{}{}
	exists, err := dbi.Exists(b.this.Table(), b.this)
	return exists && err == nil
}
//Table return database table name
func (b *Base)Table()string{
	return ""
}

//Filter return filter
func (b *Base)Filter()map[string]interface{}{
	return make(map[string]interface{})
}

//Save save this
func (b *Base)Save(dbi bingdb.DBInterface)(interface{}, error){
	if b.this == nil {
		return nil, fmt.Errorf("object is null")
	}
	if b.Exists(dbi) {
		//exists, update
		return dbi.Update(b.this.Table(),b.this.Filter(), b.this)
	}
	//not exists , insert
	return dbi.Insert(b.this.Table(), b.this)
}

//Columns return columns map
func (b *Base)Columns()map[string]interface{}{
	return make(map[string]interface{})
}

//columns return columns slice
func columns(columns map[string]interface{})([]string, []interface{}){
	cols := []string{}
	vals := []interface{}{}
	for col, val := range columns{
		cols = append(cols, col)
		vals = append(vals, val)
	}
	return cols, vals
}

//markChar return ? slice 
func markChar(cols []string)[]string{
	res := []string{}
	for range cols{
		res = append(res, "?")
	}
	return res
}

//Init init from database
func (b *Base)Init(dbi bingdb.DBInterface)error{
	_, vals := columns(b.this.Columns())
	row, err := dbi.SelectOne(b.this.Table(), b.this.Filter(), b.this)
	if err != nil {
		return err
	}
	return dbi.UnMarshal(row, vals)
}

//List list the interface
func List(dbi bingdb.DBInterface, filter map[string]interface{}, bi BaseInterface, limit int)([]interface{}, error){
	columns := bi.Columns()
	cols := []string{}
	for _, col := range columns{

	}
	rows, err := dbi.Select(bi.Table(), filter, bi, limit)
	if err != nil {
		return nil, err
	}
	res, err := dbi.UnMarshalI(rows, bi.Columns(), bi)
	return res, err
}

