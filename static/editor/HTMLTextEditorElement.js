const configFonts = ["ui-monospace",
  "SFMono-Regular",
  "Menlo",
  "Monaco",
  "Consolas",
  "\"Liberation Mono\"",
  "\"Courier New\"",
  "monospace"
];

export class HTMLTextEditMenu extends HTMLElement {

  /**
   * @param {Event} event
   */
  updateFont(event) {
    this.target.style.fontFamily = event.target.value;
    this.selectMenu.value = event.target.value;
  }

  removeFontOptions() {
    while (this.selectMenu.hasChildNodes()) {
      this.selectMenu.removeChild(this.selectMenu.firstChild);
    }
  }

  updateFontOptions(fonts) {
    this.removeFontOptions();

    fonts.forEach(font => {
      const option = document.createElement('option');
      font = font.replace(/^"|"$/g, "");

      option.value = font;
      option.innerHTML = font;

      this.selectMenu.appendChild(option);
    });
  }

  initialize() {
    this.formMenu = document.createElement('form');
    this.selectMenu = document.createElement('select');
    this.formMenu.onchange = this.updateFont.bind(this);

    this.style.position = "absolute";
    this.style.transform = "translateY(-120%)";
    this.style.boxShadow = "0px 0px 2px 2px #00000011";
    this.style.padding = "4px 8px";
    this.style.border = "4px";

    this.appendChild(this.formMenu);
    this.formMenu.appendChild(this.selectMenu);
  }

  update() {
    this.target = this.nextElementSibling;
    const styles = window.getComputedStyle(this.target);
    const defaultFonts = styles.fontFamily.split(",").map(style => style.trim())
    const fonts = [...new Set([...configFonts, ...defaultFonts])];

    this.updateFontOptions(fonts);

    const activeFont = styles.fontFamily.split(',')[0];

    this.selectMenu.value = activeFont;
  }
}

export default class HTMLTextEditElement extends HTMLElement {
  /**
   * @type {HTMLElement | null}
   */
  target = null

  /**
   * @type {HTMLTextEditMenu} menu
   */
  menu = null

  /**
   * @param {HTMLElement} target
   */
  addClickEvent(target) {
    const callback = function(event) {
      this.setTarget(event.target);
    }.bind(this);

    target.removeEventListener('click', callback);
    target.addEventListener('click', callback);
  }

  /**
   * @param {number} index The index from where to place the cursor
   */
  select(index) {
    if (this.target.firstChild !== null) {
      const range = document.createRange();
      const selection = window.getSelection();

      range.setStart(this.target.firstChild, index);
      range.collapse(true);

      selection.removeAllRanges();
      selection.addRange(range);
    }
  }

  focusOut() {
    this.target.removeAttribute('contentEditable');
    const clone = this.target.cloneNode(true);
    this.target.parentNode.replaceChild(clone, this.target);
    this.addClickEvent(clone);
    this.target = clone;
  }

  createLine() {
    console.log("Creating line")
    const selection = window.getSelection();
    const cursorIndex = selection.getRangeAt(0).startOffset;
    const content = this.target.innerHTML.slice(cursorIndex);
    const nextSibling = this.target.cloneNode(true);

    nextSibling.innerHTML = content;

    this.target.innerHTML = this.target.innerHTML.slice(0, cursorIndex);
    this.target.insertAdjacentElement('afterend', nextSibling);
    this.focusOut();
    this.setTarget(nextSibling);
    this.select(0);
    this.addClickEvent(this.target);
  }

  /**
   * returns bool true if it deleted a line, otherwise false
   */
  deleteLine() {
    const range = window.getSelection().getRangeAt(0);
    if (range.endOffset !== 0 || range.startOffset !== 0) return false;

    const targetAbove = this.target.previousSibling;
    const parent = this.target.parentNode;
    const text = this.target.innerHTML.trimEnd();

    if (targetAbove != null) {
      const length = targetAbove.innerHTML.length;

      parent.removeChild(this.target);
      this.setTarget(targetAbove);
      this.target.innerHTML += text;
      this.select(length);
      return true;
    }

    return false;
  }

  insertLine() {
    const selection = window.getSelection();
    const range = selection.getRangeAt(0);
    const br = document.createElement('br');
    const textNode = document.createTextNode('\u00a0');

    range.deleteContents();
    range.insertNode(br);
    range.collapse(false);
    range.insertNode(textNode);
    range.selectNodeContents(textNode);

    selection.removeAllRanges();
    selection.addRange(range);
    document.execCommand('delete');
  }

  /**
   * @param {Event} event
   */
  handleKeys(event) {
    switch (event.key) {
      case 'Enter':
        event.preventDefault();
        this.insertLine();
        break;
      // case 'Enter':
      //   if (!event.repeat) {
      //     event.preventDefault();
      //     this.createLine();
      //   }
      //
      //   break;
      // case 'Backspace':
      //   if (!event.repeat && this.deleteLine()) {
      //     event.preventDefault();
      //   }
      //   break;
      // case 'ArrowUp': {
      //   event.preventDefault();
      //   const cursorIndex = window.getSelection().getRangeAt(0).startOffset;
      //   const contentLength = this.target?.previousSibling?.innerHTML.length || null;
      //   const index = Math.min(cursorIndex, contentLength);
      //   this.setTarget(this.target.previousSibling);
      //   if (contentLength !== null) this.select(index);
      //   break;
      // }
      // case 'ArrowDown': {
      //   event.preventDefault();
      //   const cursorIndex = window.getSelection().getRangeAt(0).startOffset;
      //   const contentLength = this.target?.nextSibling?.innerHTML.length || null;
      //   const index = Math.min(cursorIndex, contentLength);
      //   this.setTarget(this.target.nextSibling);
      //   if (contentLength !== null) this.select(index);
      //   break;
      // }
      // case 'ArrowLeft': {
      //   const cursorIndex = window.getSelection().getRangeAt(0).startOffset;
      //   if (cursorIndex === 0 && this.target?.previousSibling) {
      //     event.preventDefault();
      //     this.setTarget(this.target.previousSibling);
      //     this.select(this.target.innerHTML.length || 0);
      //   }
      //   break;
      // }
      // case 'ArrowRight': {
      //   const cursorIndex = window.getSelection().getRangeAt(0).startOffset;
      //   if (cursorIndex === this.target.innerHTML.length && this.target?.nextSibling) {
      //     event.preventDefault();
      //     this.setTarget(this.target.nextSibling);
      //     this.select(0);
      //   }
      // }
    }
  }

  /**
   * @param {HTMLElement} target
   */
  setTarget(target) {
    if (!target || this.target === target) return;

    if (this.target) {
      this.target.removeEventListener('keydown', this.boundHandleKeys);
      this.target.removeAttribute('contentEditable');
    }

    if (!this.menu) {
      this.menu = document.createElement('text-edit-menu');
      this.menu.initialize();
    }

    if (this.menu.parentNode) {
      this.menu.parentNode.removeChild(this.menu);
    }

    this.target = target;
    this.target.contentEditable = "true";
    this.target.insertAdjacentElement('beforebegin', this.menu);
    this.target.focus();
    this.menu.update();

    if (!this.boundHandleKeys) {
      this.boundHandleKeys = this.handleKeys.bind(this);
    }

    this.target.addEventListener('keydown', this.boundHandleKeys);
  }
}

customElements.define("text-edit", HTMLTextEditElement);
customElements.define("text-edit-menu", HTMLTextEditMenu);
