package bingo

import (
	"encoding/json"
	"github.com/ycyxuehan/bingo/bingdb"
	"fmt"
)

//ModelInterface base interface
type ModelInterface interface{
	Table()string
	Filter()map[string]interface{}
	Columns()([]string, []interface{})
	Show()string
}

//Base base model
type Model struct {
	this ModelInterface
}

//SetThis set this
func (b *Model)SetThis(bi ModelInterface){
	b.this = bi
}

//Exists is this exists in database
func (b *Model)Exists(dbi bingdb.DBInterface)bool{
	if b.this == nil {
		return false
	}
	count, err := dbi.Count(b.this.Table(), b.this.Filter())
	return count >0 && err == nil
}
//Table return database table name
func (b *Model)Table()string{
	return ""
}

//Filter return filter
func (b *Model)Filter()map[string]interface{}{
	return make(map[string]interface{})
}

//Save save this
func (b *Model)Save(dbi bingdb.DBInterface)(interface{}, error){
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
func (b *Model)Columns()([]string, []interface{}){
	return nil, nil
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
func (b *Model)Init(dbi bingdb.DBInterface)error{
	_, vals := b.this.Columns()
	row, err := dbi.SelectOne(b.this.Table(), b.this.Filter(), b.this)
	if err != nil {
		return err
	}
	return dbi.UnMarshal(row, vals...)
}

//List list the interface
func List(dbi bingdb.DBInterface, filter map[string]interface{}, bi ModelInterface, limit int)([]interface{}, error){
	rows, err := dbi.Select(bi.Table(), filter, bi, limit)
	if err != nil {
		return nil, err
	}
	cols, _ := bi.Columns()
	res, err := dbi.UnMarshalI(rows, cols, bi)
	return res, err
}

func (m *Model)Show()string{
	var data []byte
	err := json.Unmarshal(data, m.this)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

