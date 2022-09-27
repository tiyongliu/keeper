import { SplitterOptions, defaultSplitterOptions } from './options';

const SEMICOLON = ';';

export interface SplitStreamContext {
  options: SplitterOptions;
  currentDelimiter: string;
  pushOutput: (item: SplitResultItem) => void;
  commandPart: string;

  line: number;
  column: number;
  streamPosition: number;

  noWhiteLine: number;
  noWhiteColumn: number;
  noWhitePosition: number;

  commandStartPosition: number;
  commandStartLine: number;
  commandStartColumn: number;

  trimCommandStartPosition: number;
  trimCommandStartLine: number;
  trimCommandStartColumn: number;

  wasDataInCommand: boolean;
}

interface ScannerContext {
  readonly options: SplitterOptions;
  readonly source: string;
  readonly position: number;
  readonly currentDelimiter: string;
  readonly end: number;
  readonly wasDataOnLine: boolean;
}

export interface SplitLineContext extends SplitStreamContext {
  source: string;
  position: number;
  // output: string[];
  end: number;
  wasDataOnLine: boolean;
  currentCommandStart: number;

  //   unread: string;
  //   currentStatement: string;
  //   semicolonKeyTokenRegex: RegExp;
}

export interface SplitPositionDefinition {
  position: number;
  line: number;
  column: number;
}

export interface SplitResultItemRich {
  text: string;
  start: SplitPositionDefinition;
  end: SplitPositionDefinition;
  trimStart?: SplitPositionDefinition;
  trimEnd?: SplitPositionDefinition;
}

export type SplitResultItem = string | SplitResultItemRich;

function movePosition(context: SplitLineContext, count: number, isWhite: boolean) {
  let { source, position, line, column, streamPosition } = context;
  while (count > 0) {
    if (source[position] == '\n') {
      line += 1;
      column = 0;
    } else {
      column += 1;
    }
    position += 1;
    streamPosition += 1;
    count -= 1;
  }
  context.position = position;
  context.streamPosition = streamPosition;
  context.line = line;
  context.column = column;

  if (!context.wasDataInCommand) {
    if (isWhite) {
      context.trimCommandStartPosition = streamPosition;
      context.trimCommandStartLine = line;
      context.trimCommandStartColumn = column;
    } else {
      context.wasDataInCommand = true;
    }
  }

  if (!isWhite) {
    context.noWhitePosition = streamPosition;
    context.noWhiteLine = line;
    context.noWhiteColumn = column;
  }
}

function isStringEnd(s: string, pos: number, endch: string, escapech: string) {
  if (!escapech) {
    return s[pos] == endch;
  }
  if (endch == escapech) {
    return s[pos] == endch && s[pos + 1] != endch;
  } else {
    return s[pos] == endch && s[pos - 1] != escapech;
  }
}

interface Token {
  type:
    | 'string'
    | 'delimiter'
    | 'whitespace'
    | 'eoln'
    | 'data'
    | 'set_delimiter'
    | 'comment'
    | 'go_delimiter'
    | 'create_routine';
  length: number;
  value?: string;
}

const WHITESPACE_TOKEN: Token = {
  type: 'whitespace',
  length: 1,
};
const EOLN_TOKEN: Token = {
  type: 'eoln',
  length: 1,
};
const DATA_TOKEN: Token = {
  type: 'data',
  length: 1,
};

function scanDollarQuotedString(context: ScannerContext): Token {
  if (!context.options.allowDollarDollarString) return null;

  let pos = context.position;
  const s = context.source;

  const match = /^(\$[a-zA-Z0-9_]*\$)/.exec(s.slice(pos));
  if (!match) return null;
  const label = match[1];
  pos += label.length;

  while (pos < context.end) {
    if (s.slice(pos).startsWith(label)) {
      return {
        type: 'string',
        length: pos + label.length - context.position,
      };
    }
    pos++;
  }

  return null;
}

function scanToken(context: ScannerContext): Token {
  let pos = context.position;
  const s = context.source;
  const ch = s[pos];

  if (context.options.stringsBegins.includes(ch)) {
    pos++;
    const endch = context.options.stringsEnds[ch];
    const escapech = context.options.stringEscapes[ch];
    while (pos < context.end && !isStringEnd(s, pos, endch, escapech)) {
      if (endch == escapech && s[pos] == endch && s[pos + 1] == endch) {
        pos += 2;
      } else {
        pos++;
      }
    }
    return {
      type: 'string',
      length: pos - context.position + 1,
    };
  }

  if (context.currentDelimiter && s.slice(pos).startsWith(context.currentDelimiter)) {
    return {
      type: 'delimiter',
      length: context.currentDelimiter.length,
    };
  }

  if (ch == ' ' || ch == '\t' || ch == '\r') {
    return WHITESPACE_TOKEN;
  }

  if (ch == '\n') {
    return EOLN_TOKEN;
  }

  if (context.options.doubleDashComments && ch == '-' && s[pos + 1] == '-') {
    while (pos < context.end && s[pos] != '\n') pos++;
    return {
      type: 'comment',
      length: pos - context.position,
    };
  }

  if (context.options.multilineComments && ch == '/' && s[pos + 1] == '*') {
    pos += 2;
    while (pos < context.end) {
      if (s[pos] == '*' && s[pos + 1] == '/') break;
      pos++;
    }
    return {
      type: 'comment',
      length: pos - context.position + 2,
    };
  }

  if (context.options.allowCustomDelimiter && !context.wasDataOnLine) {
    const m = s.slice(pos).match(/^DELIMITER[ \t]+([^\n]+)/i);
    if (m) {
      return {
        type: 'set_delimiter',
        value: m[1].trim(),
        length: m[0].length,
      };
    }
  }

  if ((context.options.allowGoDelimiter || context.options.adaptiveGoSplit) && !context.wasDataOnLine) {
    const m = s.slice(pos).match(/^GO[\t\r ]*(\n|$)/i);
    if (m) {
      return {
        type: 'go_delimiter',
        length: m[0].length - 1,
      };
    }
  }

  if (context.options.adaptiveGoSplit) {
    const m = s.slice(pos).match(/^CREATE\s*(PROCEDURE|FUNCTION|TRIGGER)/i);
    if (m) {
      return {
        type: 'create_routine',
        length: m[0].length - 1,
      };
    }
  }

  const dollarString = scanDollarQuotedString(context);
  if (dollarString) return dollarString;

  return DATA_TOKEN;
}

function containsDataAfterDelimiterOnLine(context: ScannerContext, delimiter: Token): boolean {
  const cloned = {
    options: context.options,
    source: context.source,
    position: context.position,
    currentDelimiter: context.currentDelimiter,
    end: context.end,
    wasDataOnLine: context.wasDataOnLine,
  };

  cloned.position += delimiter.length;

  while (cloned.position < cloned.end) {
    const token = scanToken(cloned);

    if (!token) {
      cloned.position += 1;
      continue;
    }

    switch (token.type) {
      case 'whitespace':
        cloned.position += token.length;
        continue;
      case 'eoln':
        return false;
      case 'comment':
        if (token.value?.includes('\n')) return true;
        cloned.position += token.length;
        continue;
      default:
        return true;
    }
  }
}

function pushQuery(context: SplitLineContext) {
  context.commandPart += context.source.slice(context.currentCommandStart, context.position);
  pushCurrentQueryPart(context);
}

function pushCurrentQueryPart(context: SplitStreamContext) {
  const trimmed = context.commandPart.substring(
    context.trimCommandStartPosition - context.commandStartPosition,
    context.noWhitePosition - context.commandStartPosition
  );

  if (trimmed.trim()) {
    if (context.options.returnRichInfo) {
      context.pushOutput({
        text: trimmed,

        start: {
          position: context.commandStartPosition,
          line: context.commandStartLine,
          column: context.commandStartColumn,
        },

        end: {
          position: context.streamPosition,
          line: context.line,
          column: context.column,
        },

        trimStart: {
          position: context.trimCommandStartPosition,
          line: context.trimCommandStartLine,
          column: context.trimCommandStartColumn,
        },

        trimEnd: {
          position: context.noWhitePosition,
          line: context.noWhiteLine,
          column: context.noWhiteColumn,
        },
      });
    } else {
      context.pushOutput(trimmed);
    }
  }
}

function markStartCommand(context: SplitLineContext) {
  context.commandStartPosition = context.streamPosition;
  context.commandStartLine = context.line;
  context.commandStartColumn = context.column;

  context.trimCommandStartPosition = context.streamPosition;
  context.trimCommandStartLine = context.line;
  context.trimCommandStartColumn = context.column;

  context.wasDataInCommand = false;
}

function splitByLines(context: SplitLineContext) {
  while (context.position < context.end) {
    if (context.source[context.position] == '\n') {
      pushQuery(context);
      context.commandPart = '';
      movePosition(context, 1, true);
      context.currentCommandStart = context.position;
      markStartCommand(context);
    } else {
      movePosition(context, 1, /\s/.test(context.source[context.position]));
    }
  }

  if (context.end > context.currentCommandStart) {
    context.commandPart += context.source.slice(context.currentCommandStart, context.position);
  }
}

export function splitQueryLine(context: SplitLineContext) {
  if (context.options.splitByLines) {
    splitByLines(context);
    return;
  }

  while (context.position < context.end) {
    const token = scanToken(context);
    if (!token) {
      // nothing special, move forward
      movePosition(context, 1, false);
      continue;
    }
    switch (token.type) {
      case 'string':
        movePosition(context, token.length, false);
        context.wasDataOnLine = true;
        break;
      case 'comment':
        movePosition(context, token.length, !!context.options.ignoreComments);
        context.wasDataOnLine = true;
        break;
      case 'eoln':
        movePosition(context, token.length, true);
        context.wasDataOnLine = false;
        break;
      case 'data':
        movePosition(context, token.length, false);
        context.wasDataOnLine = true;
        break;
      case 'whitespace':
        movePosition(context, token.length, true);
        break;
      case 'set_delimiter':
        pushQuery(context);
        context.commandPart = '';
        context.currentDelimiter = token.value;
        movePosition(context, token.length, false);
        context.currentCommandStart = context.position;
        markStartCommand(context);
        break;
      case 'go_delimiter':
        pushQuery(context);
        context.commandPart = '';
        movePosition(context, token.length, false);
        context.currentCommandStart = context.position;
        markStartCommand(context);
        if (context.options.adaptiveGoSplit) {
          context.currentDelimiter = SEMICOLON;
        }
        break;
      case 'create_routine':
        movePosition(context, token.length, false);
        if (context.options.adaptiveGoSplit) {
          context.currentDelimiter = null;
        }
        break;
      case 'delimiter':
        if (context.options.preventSingleLineSplit && containsDataAfterDelimiterOnLine(context, token)) {
          movePosition(context, token.length, false);
          context.wasDataOnLine = true;
          break;
        }
        pushQuery(context);
        context.commandPart = '';
        movePosition(context, token.length, false);
        context.currentCommandStart = context.position;
        markStartCommand(context);
        break;
    }
  }

  if (context.end > context.currentCommandStart) {
    context.commandPart += context.source.slice(context.currentCommandStart, context.position);
  }
}

export function getInitialDelimiter(options: SplitterOptions) {
  if (options?.adaptiveGoSplit) return SEMICOLON;
  return options?.allowSemicolon === false ? null : SEMICOLON;
}

export function finishSplitStream(context: SplitStreamContext) {
  pushCurrentQueryPart(context);
}

export function splitQuery(sql: string, options: SplitterOptions = null): SplitResultItem[] {
  const usedOptions = {
    ...defaultSplitterOptions,
    ...options,
  };

  if (usedOptions.noSplit) {
    if (usedOptions.returnRichInfo) {
      const lines = sql.split('\n');
      return [
        {
          text: sql,
          start: {
            position: 0,
            line: 0,
            column: 0,
          },
          end: {
            position: sql.length,
            line: lines.length,
            column: lines[lines.length - 1]?.length || 0,
          },
        },
      ];
    }
    return [sql];
  }

  const output = [];
  const context: SplitLineContext = {
    source: sql,
    end: sql.length,
    currentDelimiter: getInitialDelimiter(options),
    position: 0,
    column: 0,
    line: 0,
    currentCommandStart: 0,
    commandStartLine: 0,
    commandStartColumn: 0,
    commandStartPosition: 0,
    streamPosition: 0,

    noWhiteLine: 0,
    noWhiteColumn: 0,
    noWhitePosition: 0,

    trimCommandStartPosition: 0,
    trimCommandStartLine: 0,
    trimCommandStartColumn: 0,

    wasDataInCommand: false,

    pushOutput: cmd => output.push(cmd),
    wasDataOnLine: false,
    options: usedOptions,
    commandPart: '',
  };

  splitQueryLine(context);
  finishSplitStream(context);

  return output;
}
