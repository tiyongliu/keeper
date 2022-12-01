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
	if sql == "" {
		return &modules.MysqlRowsResult{}, nil
	}

	rows, err := s.sqlDB.Raw(sql).Rows()
	if err != nil {
		logger.Errorf("sql err failed %v", err)
		return nil, err
	}

	/*
		col := rows.Columns()
		vals := make([]interface{}, len(cols))
		rows.Scan(&vals)
	*/

	columns, _ := rows.Columns()            //获取列的信息
	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i, _ := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}
	ret := make([]map[string]interface{}, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		//err = rows.Scan(values...)        //开始读行，Scan函数只接受指针变量
		//m := make(map[string]interface{}) //用于存放1列的 [键/值] 对
		//if err != nil {
		//	logger.Errorf("scan failed %v", err)
		//	return nil, err
		//}
		//
		//for i, colName := range columns {
		//	m[colName] = values[i] //colName是键，v是值
		//}
		//ret = append(ret, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	logger.Infof("ret !!!!!!!!!!!!!!! %s", ret)
	return nil, nil
}

//https://blog.csdn.net/rockage/article/details/103776251
