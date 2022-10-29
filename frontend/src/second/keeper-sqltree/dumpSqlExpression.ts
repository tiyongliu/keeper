import { SqlDumper } from '/@/second/keeper-types';
import { Expression } from './types';
import { dumpSqlSourceRef } from './dumpSqlSource';

export function dumpSqlExpression(dmp: SqlDumper, expr: Expression) {
  switch (expr.exprType) {
    case 'column':
    {
      if (expr.source) {
        if (dumpSqlSourceRef(dmp, expr.source)) {
          dmp.put('.');
        }
      }
      dmp.put('%i', expr.columnName);
    }
      break;

    case 'placeholder':
      dmp.putRaw('{PLACEHOLDER}');
      break;

    case 'value':
      dmp.put('%v', expr.value);
      break;

    case 'raw':
      dmp.put('%s', expr.sql);
      break;

    case 'call':
      dmp.put('%s(', expr.func);
      if (expr.argsPrefix) dmp.put('%s ', expr.argsPrefix);
      dmp.putCollection(',', expr.args, x => dumpSqlExpression(dmp, x));
      dmp.put(')');
      break;

    case 'methodCall':
      dumpSqlExpression(dmp, expr.thisObject)
      dmp.put('.%s(', expr.method);
      dmp.putCollection(',', expr.args, x => dumpSqlExpression(dmp, x));
      dmp.put(')');
      break;

    case 'transform':
      dmp.transform(expr.transform, () => dumpSqlExpression(dmp, expr.expr));
      break;

    case 'rowNumber':
      dmp.put(' ^row_number() ^over (^order ^by ');
      dmp.putCollection(', ', expr.orderBy, x => {
        dumpSqlExpression(dmp, x);
        dmp.put(' %k', x.direction);
      });
      dmp.put(')');
      break;
  }
}
