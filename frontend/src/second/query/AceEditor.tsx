import {defineComponent, onMounted, onBeforeUnmount, PropType, ref, toRefs} from 'vue'
import {useBootstrapStore} from "/@/store/modules/bootstrap";
import queryParserWorkerFallback from './queryParserWorkerFallback'
import * as ace from 'ace-builds/src-noconflict/ace'
import 'ace-builds/src-noconflict/mode-sql';
import 'ace-builds/src-noconflict/mode-mysql';
import 'ace-builds/src-noconflict/mode-pgsql';
import 'ace-builds/src-noconflict/mode-sqlserver';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/mode-javascript';
import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/mode-markdown';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/ext-language_tools';

// import 'ace-builds/src-noconflict/theme-github';
// import 'ace-builds/src-noconflict/theme-sqlserver';
// import 'ace-builds/src-noconflict/theme-twilight';
// import 'ace-builds/src-noconflict/theme-monokai';
// import 'ace-builds/src-noconflict/theme-chaos';
// import 'ace-builds/src-noconflict/theme-merbivore';
import 'ace-builds/src-noconflict/theme-ambiance';
import 'ace-builds/src-noconflict/theme-chaos';
import 'ace-builds/src-noconflict/theme-chrome';
import 'ace-builds/src-noconflict/theme-clouds';
import 'ace-builds/src-noconflict/theme-clouds_midnight';
import 'ace-builds/src-noconflict/theme-cobalt';
import 'ace-builds/src-noconflict/theme-crimson_editor';
import 'ace-builds/src-noconflict/theme-dawn';
import 'ace-builds/src-noconflict/theme-dracula';
import 'ace-builds/src-noconflict/theme-dreamweaver';
import 'ace-builds/src-noconflict/theme-eclipse';
import 'ace-builds/src-noconflict/theme-github';
import 'ace-builds/src-noconflict/theme-gob';
import 'ace-builds/src-noconflict/theme-gruvbox';
import 'ace-builds/src-noconflict/theme-idle_fingers';
import 'ace-builds/src-noconflict/theme-iplastic';
import 'ace-builds/src-noconflict/theme-katzenmilch';
import 'ace-builds/src-noconflict/theme-kr_theme';
import 'ace-builds/src-noconflict/theme-kuroir';
import 'ace-builds/src-noconflict/theme-merbivore';
import 'ace-builds/src-noconflict/theme-merbivore_soft';
import 'ace-builds/src-noconflict/theme-mono_industrial';
import 'ace-builds/src-noconflict/theme-monokai';
import 'ace-builds/src-noconflict/theme-nord_dark';
import 'ace-builds/src-noconflict/theme-pastel_on_dark';
import 'ace-builds/src-noconflict/theme-solarized_dark';
import 'ace-builds/src-noconflict/theme-solarized_light';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/theme-terminal';
import 'ace-builds/src-noconflict/theme-textmate';
import 'ace-builds/src-noconflict/theme-tomorrow';
import 'ace-builds/src-noconflict/theme-tomorrow_night_blue';
import 'ace-builds/src-noconflict/theme-tomorrow_night_bright';
import 'ace-builds/src-noconflict/theme-tomorrow_night_eighties';
import 'ace-builds/src-noconflict/theme-tomorrow_night';
import 'ace-builds/src-noconflict/theme-twilight';

interface IPart {
  start: {line: string; column: number; position: number}
  end: {line: string; column: number; position: number}
  trimStart: {line: string}
  text: string
}
const aceContainer = `{position: absolute; left: 0; top: 0; right: 0; bottom: 0;}`
export default defineComponent({
  name: "AceEditor",
  emits: ['init', 'focus', 'input'],
  props: {
    value: {
      type: String as PropType<string>,
      default: ''
    },
    mode: {
      type: String as PropType<string>,
      default: 'text'
    },
    options: {
      type: Object as PropType<{}>,
      default: {}
    },
    menu: {
      type: Array as unknown as PropType<[]>,
    },
    readOnly: {
      type: Boolean as PropType<Boolean>,
      default: false
    },
    splitterOptions: {
      type: [Array, Object] as unknown as PropType<[] | object>,
    },
    currentPart: {
      type: Object as PropType<IPart>,
    }
  },
  setup(props, {emit}) {

    const bootstrap = useBootstrapStore()

    const EDITOR_ID = `svelte-ace-editor-div:${Math.floor(Math.random() * 10000000000)}`
    const {value, mode, readOnly, splitterOptions, currentPart, options, menu} = toRefs(props)
    const editor = ref<Nullable<ace.Editor>>()
    const clientWidth = ref(0)
    const clientHeight = ref(0)
    const contentBackup = ref<string>('')
    const queryParts = ref<IPart[]>([])
    const queryParserWorker = ref<Nullable<{
      text: string;
      options: {returnRichInfo: boolean};
    } | string>>(null)

    const stdOptions = {
      showPrintMargin: false,
    }

    function getEditor(): ace.Editor {
      return editor.value
    }

    function processParserResult(data: []) {
      queryParts.value = data
      // editor.setHighlightActiveLine(queryParts.length <= 1);
      changedCurrentQueryPart();
      updateAnnotations();
    }

    function updateAnnotations() {}

    const handleContextMenu = e => {
      e.preventDefault()
      const left = e.pageX
      const top = e.pageY
      bootstrap.setCurrentDropDownMenu({ left, top, items: menu.value!, targetElement: e.target })
    }

    const handleKeyDown = (data, hash, keyString, keyCode, event) => {
      // if (event) handleCommandKeyDown(event)
    }

    function changedQueryParts() {
      const editor = getEditor()
      if (splitterOptions.value && editor && queryParserWorker.value) {
        const message = {
          text: editor.getValue(),
          options: {
            ...splitterOptions,
            returnRichInfo: true,
          },
        }

        if (queryParserWorker.value == 'fallback') {
          const res = queryParserWorkerFallback(message);
          processParserResult(res);
        } else {
          (queryParserWorker.value as unknown as Window).postMessage(message)
        }
      }
    }

    function changedCurrentQueryPart() {
      if (queryParts.value.length <= 1) {

      }

      const selectionRange = editor.value.getSelectionRange();

      // if (
      //   selectionRange.start.row != selectionRange.end.row ||
      //   selectionRange.start.column != selectionRange.end.column
      // ) {
      //   removeCurrentPartMarker();
      //   currentPart = null;
      //   return;
      // }

      const cursor = selectionRange.start;
      const part = queryParts.value.find(
        x =>
          ((cursor.row == x.start.line && cursor.column >= x.start.column) || cursor.row > x.start.line) &&
          ((cursor.row == x.end.line && cursor.column <= x.end.column) || cursor.row < x.end.line)
      );

      if (
        part?.text != currentPart.value?.text ||
        part?.start?.position != currentPart.value?.start?.position ||
        part?.end?.position != currentPart.value?.end?.position
      ) {
        removeCurrentPartMarker()
      }
    }




      function setEventCallBacks() {
        editor.value.on('focus', () => emit('focus'))

        editor.value.setReadOnly(readOnly.value)
        editor.value.on('change', () => {
          const content = editor.value.getValue()
          value.value = content
          emit('input', content)
          contentBackup.value = content
          changedQueryParts()
        })
      }

      function removeCurrentPartMarker() {

      }

      onMounted(() => {
        editor.value = ace.edit(EDITOR_ID)
        emit('init', editor.value)

        editor.value.$blockScrolling = Infinity
        editor.value.getSession().setMode('ace/mode/' + mode.value)
        editor.value.setValue(value.value, 1)
        editor.value.setHighlightActiveLine(false)
        contentBackup.value = value.value
        setEventCallBacks()
        if (options.value) {
          editor.value.setOptions({
            ...stdOptions,
            ...options.value
          })
        }

        editor.value.container.addEventListener('contextmenu', handleContextMenu)
        editor.value.keyBinding.addKeyboardHandler(handleKeyDown);
        changedQueryParts()

        editor.value.on('guttermousedown', e => {
          const row = e.getDocumentPosition().row
          const part = (queryParts.value || []).find(part => part.trimStart.line == row)
          // if (part && onExecuteFragment) {
          //   onExecuteFragment(part.text, part.trimStart.line);
          //   e.stop();
          //   editor.value.moveCursorTo(part.trimStart.line, 0);
          //   editor.value.selection.clearSelection();
          // }
        },
          true
        )
      })

    onBeforeUnmount(() => {
      if (editor.value) {
        editor.value.container.removeEventListener('contextmenu', handleContextMenu);
        editor.value.keyBinding.removeKeyboardHandler(handleKeyDown);
        editor.value.destroy();
        editor.value.container.remove();
      }
    })

      return () => (
        <div style={aceContainer}>
          <div id={EDITOR_ID}
               style={`width: ${clientWidth.value}px;height: ${clientHeight.value}px`}></div>
        </div>
      )
    }
  })

export const EDITOR_THEMES = [
  'ambiance',
  'chaos',
  'chrome',
  'clouds',
  'clouds_midnight',
  'cobalt',
  'crimson_editor',
  'dawn',
  'dracula',
  'dreamweaver',
  'eclipse',
  'github',
  'gob',
  'gruvbox',
  'idle_fingers',
  'iplastic',
  'katzenmilch',
  'kr_theme',
  'kuroir',
  'merbivore',
  'merbivore_soft',
  'mono_industrial',
  'monokai',
  'nord_dark',
  'pastel_on_dark',
  'solarized_dark',
  'solarized_light',
  'sqlserver',
  'terminal',
  'textmate',
  'tomorrow',
  'tomorrow_night_blue',
  'tomorrow_night_bright',
  'tomorrow_night_eighties',
  'tomorrow_night',
  'twilight',
];

export const FONT_SIZES = [
  {label: '8', value: '8'},
  {label: '9', value: '9'},
  {label: '10', value: '10'},
  {label: '11', value: '11'},
  {label: '12 - Normal', value: '12'},
  {label: '13', value: '13'},
  {label: '14', value: '14'},
  {label: '15', value: '15'},
  {label: '16', value: '16'},
  {label: '17', value: '17'},
  {label: 'Custom', value: 'custom'},
];
