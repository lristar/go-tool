package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBContext struct {
	DBCommand             *MongoCommand
	DefaultCollectionName string
}

// NewMongoDBContext 获取mongo数据库操作ctx
func NewMongoDBContext(collectionName string) *MongoDBContext {
	db := new(MongoDBContext)
	db.DBCommand = new(MongoCommand)
	db.DefaultCollectionName = collectionName
	return db
}

// Init 初始化
func (ctx *MongoDBContext) Init(collectionName string) *MongoDBContext {
	ctx.DefaultCollectionName = collectionName
	return ctx
}

// InsertJSON 插入Json数据* 如果失败，则返回具体的error，成功则返回nil
func (ctx *MongoDBContext) InsertJSON(jsonData string) error {
	return ctx.DBCommand.InsertJSON(ctx.DefaultCollectionName, jsonData)
}

// Insert 新增实体数据插入指定Collection* 如果失败，则返回具体的error，成功则返回nil
func (ctx *MongoDBContext) Insert(data interface{}) error {
	return ctx.DBCommand.Insert(ctx.DefaultCollectionName, data)
}

// Update 更新指定的实体数据* 如果失败，则返回具体的error，成功则返回nil
func (ctx *MongoDBContext) Update(selector interface{}, data interface{}) error {
	return ctx.DBCommand.Update(ctx.DefaultCollectionName, selector, data)
}

// UpdateAll 更新指定字段* 如果失败，则返回具体的error，成功则返回nil
func (ctx *MongoDBContext) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return ctx.DBCommand.UpdateAll(ctx.DefaultCollectionName, selector, update)
}

// Remove 移除一条指定查询条件的记录* 如果失败，返回具体error，成功则返回nil
func (ctx *MongoDBContext) Remove(selector interface{}) error {
	return ctx.DBCommand.Remove(ctx.DefaultCollectionName, selector)
}

// RemoveAll 移除指定查询条件的记录 如果失败，返回具体error，成功则返回nil
func (ctx *MongoDBContext) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	return ctx.DBCommand.RemoveAll(ctx.DefaultCollectionName, selector)
}

// FindOne 查询指定查询条件的第一条单条数据* 如果失败，result为具体的document并返回具体error，成功则返回nil
func (ctx *MongoDBContext) FindOne(selector interface{}, result interface{}, sort ...string) (err error) {
	return ctx.DBCommand.FindOne(ctx.DefaultCollectionName, selector, result, sort...)
}

// FindByID 查询id指定查询条件的第一条单条数据* 如果失败，result为具体的document并返回具体error，成功则返回nil
func (ctx *MongoDBContext) FindByID(id string, result interface{}) (err error) {
	return ctx.DBCommand.FindByID(ctx.DefaultCollectionName, id, result)
}

// PipeOne 查询指定查询条件的第一条单条数据 如果失败，result为具体的document并返回具体error，成功则返回nil
func (ctx *MongoDBContext) PipeOne(pipe []bson.M, result interface{}) (err error) {
	return ctx.DBCommand.PipeOne(ctx.DefaultCollectionName, pipe, result)
}

// PipeAll 查询指定查询条件的数据* 如果失败，result为具体的document并返回具体error，成功则返回nil
func (ctx *MongoDBContext) PipeAll(pipe []bson.M, result interface{}) (err error) {
	return ctx.DBCommand.PipeAll(ctx.DefaultCollectionName, pipe, result)
}

// FindAll 查询指定查询条件的批量数据 支持分页，如果不需要分页，skip、limit请传入0 如果失败，则返回具体的document list，成功则返回nil
func (ctx *MongoDBContext) FindAll(selector interface{}, skip, limit int, result interface{}, sort ...string) (err error) {
	return ctx.DBCommand.FindList(ctx.DefaultCollectionName, selector, skip, limit, result, sort...)
}

// Distinct Distinct
func (ctx *MongoDBContext) Distinct(key string, selector interface{}, skip, limit int, result interface{}, sort ...string) (err error) {
	return ctx.DBCommand.Distinct(ctx.DefaultCollectionName, key, selector, skip, limit, result, sort...)
}

// Count 获取指定条件的记录行数 如果失败，则返回具体的error，成功则返回记录数
func (ctx *MongoDBContext) Count(selector interface{}) (count int, err error) {
	return ctx.DBCommand.Count(ctx.DefaultCollectionName, selector)
}

// Upsert 更新指定条件的数据，如果没有找到则直接插入数据 如果失败，则返回具体的error，成功则返回nil
func (ctx *MongoDBContext) Upsert(selector interface{}, data interface{}) (*mgo.ChangeInfo, error) {
	return ctx.DBCommand.Upsert(ctx.DefaultCollectionName, selector, data)

}

// Bulk Bulk批量
func (ctx *MongoDBContext) Bulk(f func(b *mgo.Bulk)) (*mgo.BulkResult, error) {
	return ctx.DBCommand.Bulk(ctx.DefaultCollectionName, f)
}
