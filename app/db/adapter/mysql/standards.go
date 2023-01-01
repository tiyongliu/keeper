package mysql

import (
	"fmt"
	"keeper/app/db"
	"keeper/app/db/standard/modules"
	"keeper/app/pkg/logger"
	"regexp"
)

func (s *Source) Dialect() string {
	return Adapter
}

func (s *Source) Ping() error {
	if s.sqlDB != nil {
		database, err := s.sqlDB.DB()
		if err != nil {
			return err
		}
		return database.Ping()
	}
	return db.ErrNotConnected
}

func (s *Source) Version() (*modules.Version, error) {
	var rows []string
	err := s.sqlDB.Raw("select version()").Scan(&rows).Error
	if err != nil {
		logger.Errorf("get mysql version failed: %v", err)
		return nil, err
	}

	if len(rows) > 0 && rows[0] != "" {
		subMath := regexp.MustCompile("(.*)-MariaDB-").FindAllSubmatch([]byte(rows[0]), -1)
		if len(subMath) >= 1 {
			return &modules.Version{
				Version:     rows[0],
				VersionText: fmt.Sprintf("MariaDB %s", subMath[0]),
			}, nil
		}
	}

	return &modules.Version{
		Version:     rows[0],
		VersionText: fmt.Sprintf("MySQL %s", rows[0]),
	}, nil
}

func (s *Source) Close() error {
	defer func() {
		s.sqlDBMu.Lock()
		s.sqlDB = nil
		s.sqlDBMu.Unlock()
	}()
	if s.sqlDB == nil {
		return nil
	}
	database, err := s.sqlDB.DB()
	if err != nil {
		return nil
	}
	return database.Close()
}

func (s *Source) ListDatabases() (interface{}, error) {
	if s.sqlDB != nil {
		var rows []string
		err := s.sqlDB.Raw("SHOW DATABASES").Scan(&rows).Error
		if err != nil {
			logger.Errorf("get mysql lastDatabases failed: %v", err)
			return nil, err
		}
		return transformMysqlDatabases(rows), nil
	}

	return nil, db.ErrNotConnected
}

func (s *Source) Query(sql string) (interface{}, error) {
	rows := make([]map[string]interface{}, 0)
	return rows, s.sqlDB.Raw(sql).Scan(&rows).Error
}

/*func (s *Source) Query(sql string) (interface{}, error) {
	if sql == "" {
		return &modules.MysqlRowsResult{}, nil
	}

	rows, err := s.sqlDB.Raw(sql).Rows()
	if err != nil {
		logger.Errorf("sql select failed %v", err)
		return nil, err
	}

	columns, err := rows.Columns() //获取列的信息
	if err != nil {
		return nil, err
	}

	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i, _ := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}

	res := make([]map[string]interface{}, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		err = rows.Scan(values...)        //开始读行，Scan函数只接受指针变量
		m := make(map[string]interface{}) //用于存放1列的 [键/值] 对
		for i, colName := range columns {
			typeof := reflect.TypeOf(values[i])
			valof := reflect.ValueOf(values[i])
			switch typeof.Elem().Kind() {
			case reflect.String:
				m[colName] = valof.String()
			case reflect.Int:
				m[colName] = valof.Int()
			}

			//var raw_value = *(values[i].(*interface{})) //读出raw数据，类型为byte
			//logger.Infof("raw_value %s", reflect.ValueOf(raw_value).Kind().String())
			//
			//switch raw_value.(type) {
			//case int8:
			//	m[colName] = raw_value.(int8)
			//case int16:
			//	m[colName] = raw_value.(int16)
			//case int32:
			//	m[colName] = raw_value.(int32)
			//case int64:
			//	m[colName] = raw_value.(int64)
			//case int:
			//	m[colName] = raw_value.(int)
			//case uint8:
			//	m[colName] = raw_value.(uint8)
			//case uint16:
			//	m[colName] = raw_value.(uint16)
			//case uint32:
			//	m[colName] = raw_value.(uint32)
			//case uint64:
			//	m[colName] = raw_value.(uint64)
			//case uint:
			//	m[colName] = raw_value.(uint)
			//case string:
			//	m[colName] = raw_value.(string)
			//case bool:
			//	m[colName] = raw_value.(bool)
			//case float32:
			//	m[colName] = raw_value.(float32)
			//case float64:
			//	m[colName] = raw_value.(float64)
			//case time.Time:
			//	m[colName] = raw_value.(time.Time)
			//case map[string]interface{}:
			//	m[colName] = raw_value.(map[string]interface{})
			//}
		}
		res = append(res, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	return res, nil
}*/

//https://blog.csdn.net/rockage/article/details/103776251
