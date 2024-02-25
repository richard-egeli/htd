
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
 * @returns {HTMLElement} The active element currently focused
 */
function getActiveElement() {
  const selection = document.getSelection();
  const range = selection.getRangeAt(0);
  let parent = range.commonAncestorContainer.parentNode;
  while (!parent.hasAttribute || !parent.hasAttribute("contentEditable")) {
    parent = parent.parentNode;
  }

  return parent;
}

function mergeTextNodes(parent) {
  if (parent) {
    let text = '';

    for (let i = 0; i < parent.childNodes.length;) {
      const node = parent.childNodes[i];
      if (node.nodeType != Node.TEXT_NODE) {
        parent.insertBefore(document.createTextNode(text), parent.childNodes[i]);
        text = '';
        i += 2;
        continue;
      }

      parent.removeChild(parent.childNodes[i]);
      text += node.textContent;
    }

    if (text.length != '') parent.appendChild(document.createTextNode(text));
  }
}

function handleBoldCommand() {
  const parent = getActiveElement();
  const selection = document.getSelection();
  const focusNode = selection.focusNode;
  const range = selection.getRangeAt(0);

  if (range.startOffset == range.endOffset) {
    console.log(parent);
    range.setStart(parent.firstChild, 0);
    range.setEnd(parent.firstChild, parent.innerHTML.length);
  }

  selection.removeAllRanges();
  selection.addRange(range);
  document.execCommand('bold', false, null);
  mergeTextNodes(parent);

  document.querySelectorAll('b').forEach(el => {
    if (el.innerHTML == '') el.remove();
  });
}

/**
 * @param {('backColor'|'bold'|'contentReadonly'|'copy'|'createLink'|'cut'|'decreaseFontSize'|'defaultParagraphSeparator'|'delete'|'enableAbsolutePositionEditor'|'enableInlineTableEditing'|'enableObjectResizing'|'fontName'|'fontSize'|'foreColor'|'formatBlock'|'forwardDelete'|'heading'|'hiliteColor'|'increaseFontSize'|'indent'|'insertBrOnReturn'|'insertHorizontalRule'|'insertHTML'|'insertImage'|'insertOrderedList'|'insertUnorderedList'|'insertParagraph'|'insertText'|'italic'|'justifyCenter'|'justifyFull'|'justifyRight'|'outdent'|'paste'|'redo'|'removeFormat'|'selectAll'|'strikeThrough'|'subscript'|'superscript'|'underline'|'undo'|'unlink'|'useCSS'|'styleWithCSS'|'autoUrlDetect')[]} id
 * @param {...string} params 
 */
export default function executeCommand(id, ...params) {
  switch (id) {
    case 'bold':
      handleBoldCommand();
      break;
    default:
      document.execCommand(id, false, params);
      break;

  }
}
