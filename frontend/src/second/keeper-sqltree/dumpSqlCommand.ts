import {SqlDumper} from '/@/second/keeper-types'
import {Command, Select, Update, Delete, Insert} from './types'
import { dumpSqlExpression } from './dumpSqlExpression';
import { dumpSqlFromDefinition, dumpSqlSourceRef } from './dumpSqlSource';
import { dumpSqlCondition } from './dumpSqlCondition';


export function dumpSqlSelect(dmp: SqlDumper, cmd: Select) {
  dmp.put('^select ');
  if (cmd.distinct) {
    dmp.put('^distinct ');
  }
  if (cmd.topRecords) {
    if (!dmp.dialect.rangeSelect || dmp.dialect.offsetFetchRangeSyntax) dmp.put('^top %s ', cmd.topRecords);
  }
  if (cmd.selectAll) {
    dmp.put('* ');
  }
  if (cmd.columns) {
    if (cmd.selectAll) dmp.put('&n,');
    dmp.put('&>&n');
    dmp.putCollection(',&n', cmd.columns, fld => {
      dumpSqlExpression(dmp, fld);
      if (fld.alias) dmp.put(' ^as %i', fld.alias);
    });
    dmp.put('&n&<');
  }
  dmp.put('^from ');
  dumpSqlFromDefinition(dmp, cmd.from);
  if (cmd.where) {
    dmp.put('&n^where ');
    dumpSqlCondition(dmp, cmd.where);
    dmp.put('&n');
  }
  if (cmd.groupBy) {
    dmp.put('&n^group ^by ');
    dmp.putCollection(', ', cmd.groupBy, expr => dumpSqlExpression(dmp, expr));
    dmp.put('&n');
  }
  if (cmd.having) {
    dmp.put('&n^having ');
    dumpSqlCondition(dmp, cmd.having);
    dmp.put('&n');
  }
  if (cmd.orderBy) {
    dmp.put('&n^order ^by ');
    dmp.putCollection(', ', cmd.orderBy, expr => {
      dumpSqlExpression(dmp, expr);
      dmp.put(' %k', expr.direction);
    });
    dmp.put('&n');
  }
  if (cmd.range) {
    if (dmp.dialect.offsetFetchRangeSyntax) {
      dmp.put('^offset %s ^rows ^fetch ^next %s ^rows ^only', cmd.range.offset, cmd.range.limit);
    } else {
      dmp.put('^limit %s ^offset %s ', cmd.range.limit, cmd.range.offset);
    }
  }
  if (cmd.topRecords) {
    if (dmp.dialect.rangeSelect && !dmp.dialect.offsetFetchRangeSyntax) dmp.put('^limit %s ', cmd.topRecords);
  }
}

export function dumpSqlUpdate(dmp: SqlDumper, cmd: Update) {
  dmp.put('^update ');
  dumpSqlSourceRef(dmp, cmd.from);

  dmp.put('&n^set ');
  dmp.put('&>');
  dmp.putCollection(', ', cmd.fields, col => {
    dmp.put('%i=', col.targetColumn);
    dumpSqlExpression(dmp, col);
  });
  dmp.put('&<');

  if (cmd.where) {
    dmp.put('&n^where ');
    dumpSqlCondition(dmp, cmd.where);
    dmp.put('&n');
  }
}

export function dumpSqlDelete(dmp: SqlDumper, cmd: Delete) {
  dmp.put('^delete ^from ');
  dumpSqlSourceRef(dmp, cmd.from);

  if (cmd.where) {
    dmp.put('&n^where ');
    dumpSqlCondition(dmp, cmd.where);
    dmp.put('&n');
  }
}

export function dumpSqlInsert(dmp: SqlDumper, cmd: Insert) {
  dmp.put(
    '^insert ^into %f (%,i) ^values (',
    cmd.targetTable,
    cmd.fields.map(x => x.targetColumn)
  );
  dmp.putCollection(',', cmd.fields, x => dumpSqlExpression(dmp, x));
  dmp.put(')');
}

export function dumpSqlCommand(dmp: SqlDumper, cmd: Command) {
  switch (cmd.commandType) {
    case 'select':
      dumpSqlSelect(dmp, cmd);
      break;
    case 'update':
      dumpSqlUpdate(dmp, cmd);
      break;
    case 'delete':
      dumpSqlDelete(dmp, cmd);
      break;
    case 'insert':
      dumpSqlInsert(dmp, cmd);
      break;
    case 'allowIdentityInsert':
      dmp.allowIdentityInsert(cmd.targetTable, cmd.allow);
      break;
  }
}

