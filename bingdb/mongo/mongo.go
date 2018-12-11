package mongo

import (
	"github.com/ycyxuehan/bingo/bingdb"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"fmt"
	"gopkg.in/mgo.v2"
)

//Mongo mongo db 
type Mongo struct {
	URI string
	DBName string
	Session *mgo.Session
	Database *mgo.Database
}


//New create a Database object
func New(connString string)(*Mongo, error){
	var m Mongo
	if connString == "" {
		return nil, fmt.Errorf("connection string is empty")
	}
	m.URI = connString
	return &m, nil
}

//Connect connect to mysql
func (m *Mongo)Connect()error{
	session, err := mgo.Dial(m.URI)
	if err != nil {
		return err
	}
	m.Session = session
	m.Database = session.DB(getDBName(m.URI))
	return nil
}

//ConnectEx connect from a uri
func (m *Mongo)ConnectEx(uri string)error{
	if uri == "" {
		return fmt.Errorf("uri is empty")
	}
	m.URI = uri
	return m.Connect()
}


//ConnectDB connect a database
func (m *Mongo)ConnectDB(DB string){
	m.Database = m.Session.DB(DB)
}

//SelectOne select data
func (m *Mongo)SelectOne(table string, filter interface{}, args ...interface{})(interface{}, error){
	if m.Database == nil {
		return nil, fmt.Errorf("can not connect database %s", m.DBName)
	}
	if table == "" {
		return nil, fmt.Errorf("table name is empty")
	}

	collection := m.Database.C(table)
	query := collection.Find(filter)
	err := query.One(args[0])
	return args[0], err
}

//Select select data
func (m *Mongo)Select(table string, filter interface{}, args ...interface{})(interface{}, error){
	if m.Database == nil {
		return nil, fmt.Errorf("can not connect database %s", m.DBName)
	}
	if table == "" {
		return nil, fmt.Errorf("table name is empty")
	}

	collection := m.Database.C(table)
	query := collection.Find(filter)
	if limit := getLimit(args...); limit > 0 {
		query = query.Limit(limit)
	}
	err := query.All(args[1])
	return args[1], err
}


//Insert insert data
func (m *Mongo)Insert(table string, args ...interface{})(interface{},error){
	if m.Database == nil {
		return  nil, fmt.Errorf("can not connect database %s", m.DBName)
	}
	if table == "" {
		return nil, fmt.Errorf("table name is empty")
	}
	collection := m.Database.C(table)
	datas := []bson.M{}
	for _, arg := range args {
		data := bson.M{}
		if a, ok := arg.(bingdb.DBMInterface); ok {
			cols, vals := a.Columns()
			for i, col := range cols{
				data[col] = vals[i]
			}
			datas = append(datas, bson.M{"set":data})
		}
	}
	err := collection.Insert(datas)
	return nil, err
}

//InsertBatch insert data
func (m *Mongo)InsertBatch(table string, args ...interface{})(interface{},error){
	return m.Insert(table, args...)
}


//Update update data
func (m *Mongo)Update(table string, filter interface{}, args ...interface{})(interface{}, error){
	if m.Database == nil {
		return nil, fmt.Errorf("can not connect database %s", m.DBName)
	}
	if table == "" {
		return nil, fmt.Errorf("table name is empty")
	}
	collection := m.Database.C(table)
	multi := false
	upsert := false
	if len(args) > 1 {
		if m, ok := args[0].(bool); ok {
			multi = m
		}
		args = args[1:]
	}
	if len(args) > 1 {
		if u, ok := args[0].(bool); ok {
			upsert = u
		}
		args = args[1:]	
	}
	data := bson.M{}
	if dbm , ok := args[0].(bingdb.DBMInterface); ok {
		cols, vals := dbm.Columns()
		for i, col := range cols{
			data[col] = vals[i]
		}
	}
	if upsert {
		return collection.Upsert(filter, bson.M{"set":data})
	}
	if multi{
		return collection.UpdateAll(filter, bson.M{"set":data})
	}
	err := collection.Update(filter, bson.M{"set":data})
	return nil, err
}

//Delete delete data
func (m *Mongo)Delete(table string, args ...interface{})(interface{}, error){
	if m.Database == nil {
		return nil, fmt.Errorf("can not connect database %s", m.DBName)
	}
	if table == "" {
		return nil, fmt.Errorf("table name is empty")
	}
	collection := m.Database.C(table)
	multi := false
	if len(args) > 1 {
		if m, ok := args[1].(bool); ok {
			multi = m
		}
	}
	if multi {
		info, err := collection.RemoveAll(args[0])
		return info, err
	}
	err := collection.Remove(args[0])
	return nil, err
}
//Marshal interface to sql
func (m *Mongo)Marshal(mi bingdb.DBMInterface)(string, error){
	return "", fmt.Errorf("not support")
}

//UnMarshal row to interface array
func (m *Mongo)UnMarshal(src interface{}, dest ...interface{})error{
	return fmt.Errorf("not support")
}

//UnMarshalI row to interface
func (m *Mongo)UnMarshalI(src interface{}, columns []string, dest interface{})([]interface{}, error){
	return nil, fmt.Errorf("not support")
}

//Count return number of records in database
func (m *Mongo)Count(table string, filter interface{})(int, error){
	if m.Database == nil {
		return -1, fmt.Errorf("can not connect to database %s", m.DBName)
	}
	if table == "" {
		return -1, fmt.Errorf("table name is empty")
	}
	collection := m.Database.C(table)
	query := collection.Find(filter)
	return query.Count()
}

//Close close
func (m *Mongo)Close(){
	m.Session.Close()
}

func (m *Mongo)Query(string, ...interface{})(interface{}, error){
	return nil, nil
}

func getDBName(uri string)string{
	if uri == "" {
		return ""
	}
	uris := strings.Split(uri, "/")
	l := len(uri)
	dbname := uris[l-1]
	return strings.Split(dbname, "?")[0]
}

func getLimit(args ...interface{})int{
	if i, ok := args[0].(int); ok {
		return i
	}
	return 0
}