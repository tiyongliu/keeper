import { mysqlSplitterOptions, mssqlSplitterOptions, postgreSplitterOptions, noSplitSplitterOptions } from './options';
import stream from 'stream';
import { splitQueryStream } from './splitQueryStream';
import fs, { fchownSync } from 'fs';

function createInputStream(...lines) {
  const pass = new stream.PassThrough({
    objectMode: true,
  });
  lines.forEach(line => pass.write(line));
  pass.end();
  return pass;
}

function streamToArray(streamSource) {
  return new Promise((resolve, reject) => {
    const res = [];
    streamSource.on('data', x => res.push(x));
    streamSource.on('end', () => resolve(res));
  });
}

test('stream: simple query', async () => {
  const output = await streamToArray(splitQueryStream(createInputStream('select * from A'), mysqlSplitterOptions));
  expect(output).toEqual(['select * from A']);
});

test('stream: query on 2 buffers', async () => {
  const output = await streamToArray(splitQueryStream(createInputStream('select * ', 'from A'), mysqlSplitterOptions));
  expect(output).toEqual(['select * from A']);
});

test('stream: 2 queries on more buffers', async () => {
  const output = await streamToArray(
    splitQueryStream(
      createInputStream('SELECT * ', 'FROM `table1`;', 'SELECT *', ' FROM `table2`'),
      mysqlSplitterOptions
    )
  );
  expect(output).toEqual(['SELECT * FROM `table1`', 'SELECT * FROM `table2`']);
});

test('file stream', async () => {
  fs.writeFileSync('testdata.sql', 'select *\n  from A;\nselect * from B');
  const fileStream = fs.createReadStream('testdata.sql', 'utf-8');
  const output = await streamToArray(splitQueryStream(fileStream, mysqlSplitterOptions));
  expect(output).toEqual(['select *\n  from A', 'select * from B']);
});

test('delimiter stream test', async () => {
  const output = await streamToArray(
    splitQueryStream(
      createInputStream('SELECT\n1;', '\n DELIMITER $$\n', ' SELECT\n2; SELECT\n3;', ' \n DELIMITER ;'),
      mysqlSplitterOptions
    )
  );
  expect(output).toEqual(['SELECT\n1', 'SELECT\n2; SELECT\n3;']);
});

test('splitted delimiter stream test', async () => {
  const output = await streamToArray(
    splitQueryStream(
      createInputStream('SELECT\n1;\n', 'DELI', 'MITER $', '$\n', ' SELECT\n2; SELECT\n3;', ' \n DELIMITER ;'),
      mysqlSplitterOptions
    )
  );
  expect(output).toEqual(['SELECT\n1', 'SELECT\n2; SELECT\n3;']);
});

test('northwind test', async () => {
  console.log(process.cwd());
  const fileStream = fs.createReadStream('sql/northwind.sql', 'utf-8');
  const output = await streamToArray(splitQueryStream(fileStream, mysqlSplitterOptions));
  expect(output['length']).toEqual(29);
});
