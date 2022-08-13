import { DatabaseInfo, TableInfo } from 'keeper-types'
import {flatten} from 'lodash-es';

export function safeJsonParse(json, defaultValue?, logError = false) {
    try {
        return JSON.parse(json);
    } catch (err) {
        if (logError) {
            console.error(`Error parsing JSON value "${json}"`, err);
        }
        return defaultValue;
    }
}

export function addTableDependencies(db: DatabaseInfo): DatabaseInfo {
    const allForeignKeys = flatten(db.tables.map(x => x.foreignKeys || []));
    return {
        ...db,
        tables: db.tables.map(table => ({
            ...table,
            dependencies: allForeignKeys.filter(x => x.refSchemaName == table.schemaName && x.refTableName == table.pureName),
        })),
    };
}

export function extendTableInfo(table: TableInfo): TableInfo {
    return {
        ...table,
        objectTypeField: 'tables',
        columns: (table.columns || []).map(column => ({
            pureName: table.pureName,
            schemaName: table.schemaName,
            ...column,
        })),
        primaryKey: table.primaryKey
            ? {
                ...table.primaryKey,
                pureName: table.pureName,
                schemaName: table.schemaName,
                constraintType: 'primaryKey',
            }
            : undefined,
        foreignKeys: (table.foreignKeys || []).map(cnt => ({
            ...cnt,
            pureName: table.pureName,
            schemaName: table.schemaName,
            constraintType: 'foreignKey',
        })),
        indexes: (table.indexes || []).map(cnt => ({
            ...cnt,
            pureName: table.pureName,
            schemaName: table.schemaName,
            constraintType: 'index',
        })),
        checks: (table.checks || []).map(cnt => ({
            ...cnt,
            pureName: table.pureName,
            schemaName: table.schemaName,
            constraintType: 'check',
        })),
        uniques: (table.uniques || []).map(cnt => ({
            ...cnt,
            pureName: table.pureName,
            schemaName: table.schemaName,
            constraintType: 'unique',
        })),
    };
}

function fillDatabaseExtendedInfo(db: DatabaseInfo): DatabaseInfo {
    return {
        ...db,
        tables: (db.tables || []).map(extendTableInfo),
        collections: (db.collections || []).map(obj => ({
            ...obj,
            objectTypeField: 'collections',
        })),
        views: (db.views || []).map(obj => ({
            ...obj,
            objectTypeField: 'views',
        })),
        matviews: (db.matviews || []).map(obj => ({
            ...obj,
            objectTypeField: 'matviews',
        })),
        procedures: (db.procedures || []).map(obj => ({
            ...obj,
            objectTypeField: 'procedures',
        })),
        functions: (db.functions || []).map(obj => ({
            ...obj,
            objectTypeField: 'functions',
        })),
        triggers: (db.triggers || []).map(obj => ({
            ...obj,
            objectTypeField: 'triggers',
        })),
    };
}

export function extendDatabaseInfo(db: DatabaseInfo): DatabaseInfo {
    console.log(`extendDatabaseInfo`, db)
    return fillDatabaseExtendedInfo(addTableDependencies(db));
}