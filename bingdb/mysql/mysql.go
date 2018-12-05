package mysql

import (
	"strings"
	"github.com/ycyxuehan/bingo/bingdb"
	"reflect"
	"fmt"
	"database/sql"
	//for mysql
	_ "github.com/go-sql-driver/mysql"
)

//MySQL mysql 
type MySQL struct {
	URI string
	db *sql.DB
}

//Filter filter
type Filter map[string]interface{}


//New new a mysql
func New(uri string)*MySQL{
	return &MySQL{
		URI: uri,
		db: nil,
	}
}
//Connect connect to mysql
func (m *MySQL)Connect()error{
	if m.URI == "" {
		return fmt.Errorf("mysql connect uri error")
	}
	db, err := sql.Open("mysql", m.URI)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

//ConnectEx connect from a uri
func (m *MySQL)ConnectEx(uri string)error{
	if uri == "" {
		return fmt.Errorf("mysql connect uri error")
	}
	m.URI = uri
	return m.Connect()
}

//SelectOne select one row from database
func (m *MySQL)SelectOne(table string, filter interface{}, args ...interface{})(interface{}, error){
	if m.db == nil {
		return nil, fmt.Errorf("db not connect")
	}
	query := ""
	if cols,_ := m.columns(args...); len(cols)>0 {
		query = fmt.Sprintf("select %s from %s where 1=1 ", strings.Join(cols, ","), table)
	} else {
		query = fmt.Sprintf("select * from %s where 1=1 ", table)
	}
	filterStr, filterI := m.filter(filter)
	if len(filterStr) > 0 {
		query = fmt.Sprintf("%s and %s limit 1;", query, strings.Join(filterStr, " and"))
	}
	return m.db.QueryRow(query, filterI...), nil
}

//Select query select sql
func (m *MySQL)Select(table string, filter interface{}, args ...interface{})(interface{}, error){
	if m.db == nil {
		return nil, fmt.Errorf("db not connect")
	}
	query := ""
	cols, _ := m.columns(args...)
	if len(cols)>0 {
		query = fmt.Sprintf("select %s from %s where 1=1 ", strings.Join(cols, ","), table)
	} else {
		query = fmt.Sprintf("select * from %s where 1=1 ", table)
	}
	filterStr, filterI := m.filter(filter)
	if len(filterStr) > 0 {
		for _, i := range filterI {
			args = append(args, i)
		}
		query = fmt.Sprintf("%s and %s", query, strings.Join(filterStr, " and"))
	}
	query = fmt.Sprintf("%s %s;", query, m.limit(args...))
	return m.db.Query(query, filterI...)
}

//Marshal interface to sql
func (m *MySQL)Marshal(mi bingdb.DBMInterface)(string, error){
	return "", fmt.Errorf("not support")
}

//UnMarshal row to interface array
func (m *MySQL)UnMarshal(src interface{}, dest ...interface{})error{
	if row, ok := src.(*sql.Row); ok {
		return row.Scan(dest...)
	}
	return fmt.Errorf("i is null or is not sql.rows point")
}

//UnMarshalI row to interface
func (m *MySQL)UnMarshalI(src interface{}, columns []string, dest interface{})([]interface{}, error){
	if row, ok := src.(*sql.Row); ok {
		if destI, ok := dest.(bingdb.DBMInterface); ok {
			destCols := []interface{}{}
			if columns == nil || len(columns) == 0 {
				_, destCols = destI.Columns()

			}else {
				for _, col := range columns {
					
					destCols = append(destCols, m.columnsMap(destI.Columns())[col])
				}
			}
			return nil, row.Scan(destCols...)
		}
	}else if rows, ok := src.(*sql.Rows); ok {
		res := []interface{}{}
		t := reflect.TypeOf(dest)
		for rows.Next(){
			nv := reflect.New(t.Elem()).Interface()
			colIs := []interface{}{}
			if mbi, ok := nv.(bingdb.DBMInterface); ok {
				_, colIs = mbi.Columns()
				rows.Scan(colIs...)
				res = append(res, nv)
			}
		}
		return res, nil	
	}
	return nil, fmt.Errorf("src type error")
}

//Update update
func (m *MySQL)Update(table string,filter interface{}, args ...interface{})(interface{},error){
	if m.db == nil {
		return nil, fmt.Errorf("db not connect")
	}
	query := ""
	cols, vals := m.columns(args...)
	if  len(cols)>0 {
		query = fmt.Sprintf("update %s set %s =? where 1=1 ", table, strings.Join(cols, "=?,"))
	} else {
		return nil, fmt.Errorf("no column need to update")
	}
	filterStr, filterI := m.filter(filter)
	if len(filterStr) > 0 {
		for _, i := range filterI {
			vals = append(vals, i)
		}
		query = fmt.Sprintf("%s and %s;", query, strings.Join(filterStr, " and"))
	}
	return m.exec(query, vals...)
}

//Insert insert
func (m *MySQL)Insert(table string, args ...interface{})(interface{}, error){
	if m.db == nil {
		return nil, fmt.Errorf("db not connect")
	}
	query := ""
	cols, vals := m.columns(args...);
	if len(cols)>0 {
		tmp := []string{}
		for range cols {
			tmp = append(tmp, "?")
		}
		query = fmt.Sprintf("insert into  %s (%s) values( %s);", table, strings.Join(cols, ","), strings.Join(tmp, ","))
	} 
	return m.exec(query, vals...)
}
//InsertBatch batch insert
func (m *MySQL)InsertBatch(table string, args ...interface{})(interface{}, error){
	if m.db ==  nil {
		return nil, fmt.Errorf("not connect to a db")
	}
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	res := []sql.Result{}
	cols, _ := m.columns(args...)
	if len(cols) <1 {
		return nil, fmt.Errorf("no column to insert")
	}
	tmp := []string{}
	for range cols {
		tmp = append(tmp, "?")
	}
	query := fmt.Sprintf("insert into  %s (%s) values( %s);", table, strings.Join(cols, ","), strings.Join(tmp, ","))
	for _, arg := range args {
		_, vals := m.columns(arg)
		if len(vals) > 0 {
			r, err := tx.Exec(query, vals...)
			res = append(res, r)
			if err != nil {
				e := tx.Commit()
				if e != nil {
					return res, fmt.Errorf("query error %s and commit previous query error %s", err, e)
				}
				return res, err
			}
		} else {
			e := tx.Commit()
			if e != nil {
				return res, fmt.Errorf("no column found and commit previous query error %s", e)
			}
			return res, fmt.Errorf("no column found")
		}
	}
	return res, tx.Commit()
}
//Delete delete
func (m *MySQL)Delete(table string, args ...interface{})(interface{}, error){
	if m.db == nil{
		return nil, fmt.Errorf("db not connect")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("no recode to delete")
	}
	filterS, filterI := m.filter(args[0])
	query := ""
	if len(filterS) > 0 {
		query = fmt.Sprintf("delete from %s where %s;", table, strings.Join(filterS, " and"))
	}
	return m.exec(query, filterI...)
}

//exec exec sql
func (m *MySQL)exec(query string, args ...interface{})(sql.Result, error){
	if m.db ==  nil {
		return nil, fmt.Errorf("not connect to a db")
	}
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return res, tx.Commit()
}

//Count is the record in database
func (m *MySQL)Count(table string, filter interface{})(int, error){
	count := -1
	if m.db == nil {
		return count, fmt.Errorf("db not connect")
	}
	query := fmt.Sprintf("select count(1) from %s where 1=1 ", table)
	filterStr, filterI := m.filter(filter)
	if len(filterStr) > 0 {
		query = fmt.Sprintf("%s and %s limit 1;", query, strings.Join(filterStr, " and"))
	}
	row := m.db.QueryRow(query, filterI...)
	err := row.Scan(&count)
	return count, err
}

//
func (m *MySQL)filter(f interface{})([]string, []interface{}){
	resStr := []string{}
	resI := []interface{}{}
	if filter, ok :=f.(Filter); ok {
		for key, val := range filter {
			resStr = append(resStr, fmt.Sprintf("%s=?", key))
			resI = append(resI, val)
		}
	}else if filter, ok :=f.(map[string]interface{}); ok {
		for key, val := range filter {
			resStr = append(resStr, fmt.Sprintf("%s=?", key))
			resI = append(resI, val)
		}
	}
	return resStr, resI
}

//
func (m *MySQL)columns(args ...interface{})([]string, []interface{}){
	cols := []string{}
	vals := []interface{}{}
	if len(args) > 0 {	
		if res, ok := args[0].(bingdb.DBMInterface); ok {
			return res.Columns()
		}
	}
	return cols, vals
}
//
func (m *MySQL)limit(args ...interface{})string{
	if len(args) >1 {
		if arg, ok := args[1].(int); ok {
			if arg == 0 {
				return ""
			}
			return fmt.Sprintf("limit %d", arg)
		}
	}
	return ""
}

func (m *MySQL)removeFilterCol(cols []string, val []interface{}, filter string){
	for i := range cols {
		if cols[i] == filter{
			cols = append(cols[:i], cols[i+1:]...)
			val = append(val[:i], val[i+1:]...)
		}
	}
}

func (m *MySQL)columnsMap(cols []string, vals []interface{})map[string]interface{}{
	colsMap := make(map[string]interface{})
	for i := range cols {
		if i < len(vals){
			colsMap[cols[i]] = vals[i]
		}
	}
	return colsMap
}