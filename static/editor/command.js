
/**
 * @type {('backColor'|'bold'|'contentReadonly'|'copy'|'createLink'|'cut'|'decreaseFontSize'|'defaultParagraphSeparator'|'delete'|'enableAbsolutePositionEditor'|'enableInlineTableEditing'|'enableObjectResizing'|'fontName'|'fontSize'|'foreColor'|'formatBlock'|'forwardDelete'|'heading'|'hiliteColor'|'increaseFontSize'|'indent'|'insertBrOnReturn'|'insertHorizontalRule'|'insertHTML'|'insertImage'|'insertOrderedList'|'insertUnorderedList'|'insertParagraph'|'insertText'|'italic'|'justifyCenter'|'justifyFull'|'justifyRight'|'outdent'|'paste'|'redo'|'removeFormat'|'selectAll'|'strikeThrough'|'subscript'|'superscript'|'underline'|'undo'|'unlink'|'useCSS'|'styleWithCSS'|'autoUrlDetect')[]}
 */
const commandIds = [
  'backColor',
  'bold',
  'contentReadonly',
  'copy',
  'createLink',
  'cut',
  'decreaseFontSize',
  'defaultParagraphSeparator',
  'delete',
  'enableAbsolutePositionEditor',
  'enableInlineTableEditing',
  'enableObjectResizing',
  'fontName',
  'fontSize',
  'foreColor',
  'formatBlock',
  'forwardDelete',
  'heading',
  'hiliteColor',
  'increaseFontSize',
  'indent',
  'insertBrOnReturn',
  'insertHorizontalRule',
  'insertHTML',
  'insertImage',
  'insertOrderedList',
  'insertUnorderedList',
  'insertParagraph',
  'insertText',
  'italic',
  'justifyCenter',
  'justifyFull',
  'justifyRight',
  'outdent',
  'paste',
  'redo',
  'removeFormat',
  'selectAll',
  'strikeThrough',
  'subscript',
  'superscript',
  'underline',
  'undo',
  'unlink',
  'useCSS',
  'styleWithCSS',
  'autoUrlDetect'
];

/**
 * @param {('backColor'|'bold'|'contentReadonly'|'copy'|'createLink'|'cut'|'decreaseFontSize'|'defaultParagraphSeparator'|'delete'|'enableAbsolutePositionEditor'|'enableInlineTableEditing'|'enableObjectResizing'|'fontName'|'fontSize'|'foreColor'|'formatBlock'|'forwardDelete'|'heading'|'hiliteColor'|'increaseFontSize'|'indent'|'insertBrOnReturn'|'insertHorizontalRule'|'insertHTML'|'insertImage'|'insertOrderedList'|'insertUnorderedList'|'insertParagraph'|'insertText'|'italic'|'justifyCenter'|'justifyFull'|'justifyRight'|'outdent'|'paste'|'redo'|'removeFormat'|'selectAll'|'strikeThrough'|'subscript'|'superscript'|'underline'|'undo'|'unlink'|'useCSS'|'styleWithCSS'|'autoUrlDetect')[]} id
 * @param {...string} params 
 */
export default function command(id, ...params) {
  if (commandIds.includes(id)) document.execCommand(id, false, params);
}
