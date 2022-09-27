import stream from 'stream';
import {
  SplitStreamContext,
  getInitialDelimiter,
  SplitLineContext,
  splitQueryLine,
  finishSplitStream,
} from './splitQuery';
import { SplitterOptions } from './options';

export class SplitQueryStream extends stream.Transform {
  context: SplitStreamContext;
  lineBuffer: string;

  constructor(options: SplitterOptions) {
    super({ objectMode: true });
    this.context = {
      commandPart: '',
      commandStartLine: 0,
      commandStartColumn: 0,
      commandStartPosition: 0,
      streamPosition: 0,
      line: 0,
      column: 0,

      noWhiteLine: 0,
      noWhiteColumn: 0,
      noWhitePosition: 0,

      trimCommandStartPosition: 0,
      trimCommandStartLine: 0,
      trimCommandStartColumn: 0,

      wasDataInCommand: false,

      options,
      currentDelimiter: getInitialDelimiter(options),
      pushOutput: cmd => this.push(cmd),
    };
    this.lineBuffer = '';
  }

  flushBuffer() {
    const lineContext: SplitLineContext = {
      ...this.context,
      position: 0,
      currentCommandStart: 0,
      wasDataOnLine: false,
      source: this.lineBuffer,
      end: this.lineBuffer.length,

      streamPosition: this.context.streamPosition,
      line: this.context.line,
      column: this.context.column,

      commandStartPosition: this.context.commandStartPosition,
      commandStartLine: this.context.commandStartLine,
      commandStartColumn: this.context.commandStartColumn,

      noWhitePosition: this.context.noWhitePosition,
      noWhiteLine: this.context.noWhiteLine,
      noWhiteColumn: this.context.noWhiteColumn,

      trimCommandStartPosition: this.context.trimCommandStartPosition,
      trimCommandStartLine: this.context.trimCommandStartLine,
      trimCommandStartColumn: this.context.trimCommandStartColumn,

      wasDataInCommand: this.context.wasDataInCommand,
    };

    splitQueryLine(lineContext);

    this.context.commandPart = lineContext.commandPart;
    this.context.currentDelimiter = lineContext.currentDelimiter;

    this.context.streamPosition = lineContext.streamPosition;
    this.context.line = lineContext.line;
    this.context.column = lineContext.column;

    this.context.commandStartPosition = lineContext.commandStartPosition;
    this.context.commandStartLine = lineContext.commandStartLine;
    this.context.commandStartColumn = lineContext.commandStartColumn;

    this.context.noWhitePosition = lineContext.noWhitePosition;
    this.context.noWhiteLine = lineContext.noWhiteLine;
    this.context.noWhiteColumn = lineContext.noWhiteColumn;

    this.context.trimCommandStartPosition = lineContext.trimCommandStartPosition;
    this.context.trimCommandStartLine = lineContext.trimCommandStartLine;
    this.context.trimCommandStartColumn = lineContext.trimCommandStartColumn;

    this.context.wasDataInCommand = lineContext.wasDataInCommand;

    this.lineBuffer = '';
  }

  _transform(chunk, encoding, done) {
    for (let i = 0; i < chunk.length; i += 1) {
      const ch = chunk[i];
      this.lineBuffer += ch;
      if (ch == '\n') {
        this.flushBuffer();
      }
    }

    done();
  }
  _flush(done) {
    this.flushBuffer();
    finishSplitStream(this.context);
    done();
  }
}

export function splitQueryStream(sourceStream, options: SplitterOptions) {
  const splitter = new SplitQueryStream(options);
  sourceStream.pipe(splitter);
  return splitter;
}
