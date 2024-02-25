import HTMLTextEditMenu from "./HTMLTextEditMenu.js";
import executeCommand from "./command.js";


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
   * @param {HTMLElement} element
   * @returns bool if it's a valid text element 
   */
  validTextElement(element) {
    return element && ["P", "PRE", "SPAN", "H1", "H2", "H3", "H4", "H5", "H6", "H7"].includes(element.nodeName.toUpperCase());
  }

  /**
   * @param {HTMLElement} target
   */
  addClickEvent(target) {
    console.log(target);
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
    executeCommand('delete');
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
    if (!this.validTextElement(target) || this.target === target) return;

    if (this.target) {
      this.target.removeEventListener('keydown', this.boundHandleKeys);
      this.target.removeAttribute('contentEditable');
    }

    if (!this.menu) {
      this.menu = document.createElement('text-edit-menu');
      this.menu.initialize(this);
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
