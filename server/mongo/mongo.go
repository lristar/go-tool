package mongo

import (
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"strings"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULTPOOLLIMIT = 4096
	DEFAULTTIMEOUT   = 30 * time.Second
)

var (
	session *mgo.Session
	db      *mgo.Database
	once    sync.Once
)

type Config struct {
	Host      string `json:"host"`
	Dbname    string `json:"dbname"`
	Rsname    string `json:"rsname"`
	Uname     string `json:"uname"`
	Pwd       string `json:"pwd"`
	PoolLimit int    `json:"pool_limit"`
	IsDirect  bool   `json:"is_direct"`
	TimeOut   int    `json:"time_out"`
}

type Server struct {
	session *mgo.Session
	db      *mgo.Database
	*mgo.Collection
}

func InitMongodb(c Config) {
	once.Do(func() {
		dialInfo := &mgo.DialInfo{
			Addrs:          strings.Split(c.Host, ","),
			Direct:         false,
			Timeout:        time.Second * time.Duration(c.TimeOut),
			Database:       c.Dbname,
			ReplicaSetName: c.Rsname,
			Username:       c.Uname,
			Password:       c.Pwd,
			PoolLimit:      c.PoolLimit, // Session.SetPoolLimit
		}
		if c.PoolLimit == 0 {
			dialInfo.PoolLimit = DEFAULTPOOLLIMIT
		}
		if c.TimeOut == 0 {
			dialInfo.Timeout = DEFAULTTIMEOUT
		}

		var err error
		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Strong, true)
		logger.Infof("mongo:host->%s,dbname->%s,username->%s\n", c.Host, c.Dbname, c.Uname)
	})
}

// GetServerCopy GetCon 返回到mongo collections的连接
func GetServerCopy(table string) *Server {
	s := session.Copy()
	d := s.DB("")
	return &Server{
		session:    s,
		Collection: d.C(table),
		db:         d,
	}
}

func GetDefaultServer() *Server {
	return &Server{
		session:    session,
		db:         db,
		Collection: nil,
	}
}

func (s *Server) Close() error {
	s.session.Close()
	return nil
}

// GetErrNotFound 获取未查询到数据的mongo错误
func GetErrNotFound() error {
	return mgo.ErrNotFound
}

// Populate 主键关联查询
type Populate struct {
	Query           bson.M //查询条件
	From            string // 关联表名
	LocalField      string // 关联字段
	Match           bson.M // 关联表过滤
	ForeignsSelects string // 关联表字段过滤
	Sort            string // 排序
	Skip            int
	Limit           int
	Count           string
}

// GetPopulatePramas 拼装连表查询条件
func (populate *Populate) GetPopulatePramas() []bson.M {
	query := populate.Query
	from := populate.From
	localField := populate.LocalField
	foreignsSelects := populate.ForeignsSelects
	sortBy := populate.Sort
	skip := populate.Skip
	limit := populate.Limit
	foreignsSelects = strings.Trim(foreignsSelects, " ")
	addFields := bson.M{}
	if foreignsSelects != "" {
		ss := strings.Split(foreignsSelects, " ")
		for _, s := range ss {
			addFields[from+"."+s] = "$__foreign_info." + s
		}
	}
	sortByStr := strings.Trim(sortBy, " ")
	sortByArr := strings.Split(sortByStr, " ")
	sort := bson.M{}
	for _, sortEle := range sortByArr {
		if strings.HasPrefix(sortEle, "-") {
			field := strings.Replace(sortEle, "-", "", 1)
			sort[field] = -1
		} else {
			sort[sortEle] = 1
		}
	}
	querys := []bson.M{}
	if len(query) > 0 {
		querys = append(querys, bson.M{"$match": query})
	}
	querys = append(querys, bson.M{"$lookup": bson.M{
		"from":         from,
		"localField":   localField,
		"foreignField": "_id",
		"as":           "__foreign"},
	})
	querys = append(querys, bson.M{
		"$addFields": bson.M{
			"__foreign_info": bson.M{
				"$arrayElemAt": []interface{}{"$__foreign", 0},
			},
		}})
	if len(populate.Match) > 0 {
		querys = append(querys, bson.M{"$match": getMatch(populate.Match)})
	}
	if len(addFields) > 0 {
		querys = append(querys, bson.M{
			"$addFields": addFields,
		})
	}
	querys = append(querys, bson.M{
		"$project": bson.M{
			"__foreign":      0,
			"__foreign_info": 0,
		},
	})
	if populate.Count != "" {
		querys = append(querys, bson.M{"$count": populate.Count})
	} else {
		if len(sort) != 0 {
			querys = append(querys, bson.M{"$sort": sort})
		}
		if skip > 0 {
			querys = append(querys, bson.M{"$skip": skip})
		}
		if limit > 0 {
			querys = append(querys, bson.M{"$limit": limit})
		}
	}
	return querys
}

func getMatch(pMatch bson.M) bson.M {
	match := bson.M{}
	for k, v := range pMatch {
		switch k {
		case "$or":
			array := make([]bson.M, 0)
			arr := v.([]bson.M)
			for _, ev := range arr {
				ele := bson.M{}
				for elek, elev := range ev {
					ele["__foreign_info."+elek] = elev
				}
				array = append(array, ele)
			}
			match[k] = array
		}
	}
	return match
}
