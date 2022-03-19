package sql

func TablesSQL() string {
	return `select 
	TABLE_NAME as pureName, 
	TABLE_ROWS as tableRowCount,
	case when ENGINE='InnoDB' then CREATE_TIME else coalesce(UPDATE_TIME, CREATE_TIME) end as modifyDate 
from information_schema.tables 
where TABLE_SCHEMA = '#DATABASE#' and TABLE_TYPE='BASE TABLE' and TABLE_NAME =OBJECT_ID_CONDITION`
}

func TablesSQL1() string {
	return `select
        TABLE_NAME as pureName,
        TABLE_ROWS as tableRowCount,
        case when ENGINE='InnoDB' then CREATE_TIME else coalesce(UPDATE_TIME, CREATE_TIME) end as modifyDate
from information_schema.tables
where TABLE_SCHEMA = 'yami_shops' and TABLE_TYPE='BASE TABLE' and TABLE_NAME  is not null`
}
