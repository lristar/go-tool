package mongo

import (
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// var mgoSessionPool *sync.Map

// func init() {
// 	mgoSessionPool = new(sync.Map)
// }

type MongoCommand struct{}

type Selector map[string]interface{}

func CreateUpdateSet(selector Selector) Selector {
	return Selector{"$set": selector}
}

/*  新增Json数据插入指定的Collection
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) InsertJSON(collectionName, jsonData string) error {
	// logTitle := getLogTitle("InsertJSON", collectionName)
	_, db := GetDB()

	defer db.Session.Close()
	c := db.C(collectionName)

	var f interface{}
	data := []byte(jsonData)
	//特殊处理，需要先json序列化，成为标准json，再进一步转为bosn，才能成功
	errJsonunmar := json.Unmarshal(data, &f)
	if errJsonunmar != nil {
		// cmd.Error(errJsonunmar, logTitle+"json.Unmarshal["+jsonData+"]error - "+err_jsonunmar.Error())
		return errJsonunmar
	}
	bjsondata, errJsonunmar := json.Marshal(f)
	if errJsonunmar != nil {
		// cmd.Error(err_jsonmar, logTitle+"json.Marshal["+jsonData+"]error - "+err_jsonmar.Error())
		return errJsonunmar
	}

	var bf interface{}
	bsonerr := bson.UnmarshalJSON(bjsondata, &bf)
	if bsonerr != nil {
		// cmd.Error(bsonerr, logTitle+"bson.UnmarshalJSON["+jsonData+"]error - "+bsonerr.Error())
		return bsonerr
	}
	err := c.Insert(&bf)
	if err != nil {
		// cmd.Error(err, logTitle+"["+jsonData+"]error - "+err.Error())
		return err
	} else {
		// cmd.Debug(logTitle + "[" + jsonData + "]Success - " + jsonData)
	}
	return nil
}

/*新增实体数据插入指定Collection
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) Insert(collectionName string, data interface{}) error {
	// logTitle := getLogTitle("insertBlob", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	err := c.Insert(data)
	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(data)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success")
	// }
	return err
}

/*更新一条指定的实体数据
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) Update(collectionName string, selector interface{}, data interface{}) error {
	// logTitle := getLogTitle("updateBlob", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	err := c.Update(selector, data)
	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(data)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success")
	// }
	return err
}

/*更新指定的实体数据
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) UpdateAll(collectionName string, selector interface{}, data interface{}) (info *mgo.ChangeInfo, err error) {
	// logTitle := getLogTitle("updateBlob", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	info, err = c.UpdateAll(selector, data)
	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(data)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success")
	// }
	return info, err
}

/*移除一条指定查询条件的记录
* 如果失败，返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) Remove(collectionName string, selector interface{}) error {
	// logTitle := getLogTitle("remove", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	err := c.Remove(selector)
	// if err != nil {
	// 	cmd.Error(err, logTitle+"error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "Success")
	// }
	return err
}

/*移除指定查询条件的记录
* 如果失败，返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) RemoveAll(collectionName string, selector interface{}) (info *mgo.ChangeInfo, err error) {
	// logTitle := getLogTitle("remove", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	info, err = c.RemoveAll(selector)
	// if err != nil {
	// 	cmd.Error(err, logTitle+"error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "Success")
	// }
	return info, err
}

/*查询指定查询条件的第一条单条数据
* 如果失败，result为具体的document并返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) FindOne(collectionName string, selector interface{}, result interface{}, sort ...string) (err error) {
	// logTitle := getLogTitle("findOne", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	if len(sort) > 0 {
		err = c.Find(selector).Sort(sort...).One(result)
	} else {
		err = c.Find(selector).One(result)
	}

	if err != nil && err.Error() == GetErrNotFound().Error() {
		err = nil
	}
	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "]Success")
	// }
	return err
}

/*根据id查询指定查询条件的第一条单条数据
* 如果失败，result为具体的document并返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) FindByID(collectionName string, id string, result interface{}) (err error) {
	// logTitle := getLogTitle("findOne", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)

	err = c.FindId(bson.ObjectIdHex(id)).One(result)

	if err != nil && err.Error() == GetErrNotFound().Error() {
		err = nil
	}
	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "]Success")
	// }
	return err
}

/*查询指定查询条件的批量数据
* 支持分页，如果不需要分页，skip、limit请传入0
* 如果失败，则返回具体的document list，成功则返回nil
 */
func (cmd *MongoCommand) FindList(collectionName string, selector interface{}, skip, limit int, result interface{}, sort ...string) (err error) {
	// logTitle := getLogTitle("findList", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)

	query := c.Find(selector)
	if len(sort) > 0 {
		query = query.Sort(sort...)
	}

	if limit > 0 {
		err = query.Skip(skip).Limit(limit).All(result)
	} else {
		err = query.Skip(skip).All(result)
	}

	if err != nil && err.Error() == GetErrNotFound().Error() {
		err = nil
	}
	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]["+fmt.Sprint(skip, limit)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "][" + fmt.Sprint(skip, limit) + "]Success")
	// }
	return err
}

// Distinct collectionName 集合名称 key Distinct字段名 selector 查询条件 result 结果 sort 排序字段
func (cmd *MongoCommand) Distinct(collectionName, key string, selector interface{}, skip, limit int, result interface{}, sort ...string) (err error) {
	// logTitle := getLogTitle("findList", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)

	query := c.Find(selector)
	if len(sort) > 0 {
		query = query.Sort(sort...)
	}

	if limit > 0 {
		err = query.Skip(skip).Limit(limit).Distinct(key, result)
	} else {
		err = query.Skip(skip).Distinct(key, result)
	}

	// if err != nil {
	// 	cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]["+fmt.Sprint(skip, limit)+"]error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "][" + fmt.Sprint(skip, limit) + "]Success")
	// }
	return err
}

/*更新指定条件的数据，如果没有找到则直接插入数据
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) Upsert(collectionName string, selector interface{}, data interface{}) (*mgo.ChangeInfo, error) {
	// logTitle := getLogTitle("upsertBlob", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	change, errUpsert := c.Upsert(selector, data)
	// if errUpsert != nil {
	// 	cmd.Error(errUpsert, logTitle+"["+fmt.Sprint(data)+"]error - "+errUpsert.Error())
	// } else {
	// 	cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success -> " + fmt.Sprint(change))
	// }
	return change, errUpsert
}

/*获取指定条件的记录行数
* 如果失败，则返回具体的error，成功则返回记录数
 */
func (cmd *MongoCommand) Count(collectionName string, selector interface{}) (count int, err error) {
	// logTitle := getLogTitle("Count", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	count, err = c.Find(selector).Count()
	// if err != nil {
	// 	cmd.Error(err, logTitle+"error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "Success")
	// }
	return count, err
}

// Bulk Bulk批量
func (cmd *MongoCommand) Bulk(collectionName string, f func(b *mgo.Bulk)) (*mgo.BulkResult, error) {
	// logTitle := getLogTitle("Count", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	bulk := c.Bulk()
	f(bulk)
	res, err := bulk.Run()
	// if err != nil {
	// 	cmd.Error(err, logTitle+"error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "Success")
	// }
	return res, err
}

// PipeOne Pipe
func (cmd *MongoCommand) PipeOne(collectionName string, pipe []bson.M, result interface{}) error {
	// logTitle := getLogTitle("Count", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	err := c.Pipe(pipe).One(result)
	if err != nil && err.Error() == GetErrNotFound().Error() {
		err = nil
	}
	// if err != nil {
	// 	cmd.Error(err, logTitle+"error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "Success")
	// }
	return err
}

// PipeAll PipeAll
func (cmd *MongoCommand) PipeAll(collectionName string, pipe []bson.M, result interface{}) error {
	// logTitle := getLogTitle("Count", collectionName)
	_, db := GetDB()
	defer db.Session.Close()

	c := db.C(collectionName)
	err := c.Pipe(pipe).All(result)
	if err != nil && err.Error() == GetErrNotFound().Error() {
		err = nil
	}
	// if err != nil {
	// 	cmd.Error(err, logTitle+"error - "+err.Error())
	// } else {
	// 	cmd.Debug(logTitle + "Success")
	// }
	return err
}

func (cmd *MongoCommand) GetStat() (result map[string]interface{}, err error) {
	_, db := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Session.Close()
	result = make(map[string]interface{})
	err = db.Run(bson.M{"serverStatus": 1}, result)
	return result, err
}

// getLogTitle return log title
func getLogTitle(commandName string, collectionName string) string {
	return "database.MongoCommand:" + commandName + "[" + collectionName + "]:"
}

// // getSessionCopy get seesion copy with conn from pool
// func getSessionCopy(conn string) (*mgo.Session, error) {
// 	data, isOk := mgoSessionPool.Load(conn)
// 	if isOk {
// 		session, isSuccess := data.(*mgo.Session)
// 		if isSuccess {
// 			return session.Clone(), nil
// 		}
// 	}
// 	session, err := mgo.Dial(conn)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		mgoSessionPool.Store(conn, session)
// 		return session.Clone(), nil
// 	}
// }
