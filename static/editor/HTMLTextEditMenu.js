
import HTMLTextEditElement from "./HTMLTextEditElement.js";
import executeCommand from "./command.js";

import { createFontTypeForm as createTextTypeForm } from "./components/forms.js";

const configFonts = ["ui-monospace",
  "SFMono-Regular",
  "Menlo",
  "Monaco",
  "Consolas",
  "\"Liberation Mono\"",
  "\"Courier New\"",
  "monospace"
];

/**
 * @returns {number} Font size
 */
function getFontSize() {
  const selection = window.getSelection();
  if (!selection.rangeCount) return 0;

  const range = selection.getRangeAt(0);
  let element = range.startContainer;

  if (element.nodeType === Node.TEXT_NODE) {
    element = element.parentNode;
  }

  const computedStyle = window.getComputedStyle(element);
  return computedStyle.fontSize;
}

/**
 * @param {string} name
 * @returns {HTMLButtonElement}
 */
function createMenuButton(name) {
  const button = document.createElement('button');

  button.name = name;
  button.innerHTML = name;

  return button;
}

/**
 * @returns {HTMLFormElement}
 */
function createFontSizeForm() {
  const form = document.createElement('form');
  const label = document.createElement('label');
  const select = document.createElement('select');

  label.innerHTML = "Font Size: ";

  [1, 2, 3, 4, 5, 6, 7].forEach(size => {
    const option = document.createElement('option');

    option.value = size;
    option.innerHTML = size;

    select.appendChild(option);
  });

  form.appendChild(label);
  form.appendChild(select);
  form.onchange = (e) => executeCommand('fontSize', e.target.value);
  return form;
}

export default class HTMLTextEditMenu extends HTMLElement {

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

  updateTextType(event) {
    event.preventDefault();

    const type = event.target.value;
    const element = document.createElement(type);

    Array.from(this.target.attributes).forEach(attr => element.setAttribute(attr.name, attr.value));
    while (this.target.firstChild) element.appendChild(this.target.firstChild);

    this.target.parentNode.replaceChild(element, this.target);
    this.parent.setTarget(element);
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

  /**
   * @param {Event} event
   */
  onBoldButtonClicked(event) {
    event.preventDefault();
    executeCommand('bold');
  }

  setSelectionFont(event) {
    event.preventDefault();
    executeCommand('fontName', font);
  }

  /**
  * @param {HTMLTextEditElement} parent
  */
  initialize(parent) {
    this.parent = parent;
    this.fontSizeForm = createFontSizeForm();
    this.textTypeForm = createTextTypeForm();
    this.formMenu = document.createElement('form');
    this.boldButton = document.createElement('button');
    this.selectLabel = document.createElement('label');
    this.selectMenu = document.createElement('select');
    this.formMenu.onchange = this.updateFont.bind(this);
    this.italicButton = createMenuButton('Italic');

    this.boldButton.innerHTML = "Bold";
    this.boldButton.onclick = this.onBoldButtonClicked

    this.textTypeForm.onchange = this.updateTextType.bind(this);

    this.italicButton.onclick = () => executeCommand('italic');

    this.style.display = "flex";
    this.style.gap = "8px";
    this.style.position = "absolute";
    this.style.transform = "translateY(-120%)";
    this.style.boxShadow = "0px 0px 2px 2px #00000011";
    this.style.padding = "4px 8px";
    this.style.border = "4px";

    this.selectLabel.innerHTML = "Font: ";

    this.appendChild(this.boldButton);
    this.appendChild(this.italicButton);
    this.appendChild(this.formMenu);
    this.appendChild(this.fontSizeForm);
    this.formMenu.appendChild(this.selectLabel);
    this.formMenu.appendChild(this.selectMenu);
    this.appendChild(this.textTypeForm);
  }

  update() {
    this.target = this.nextElementSibling;
    const styles = window.getComputedStyle(this.target);
    const defaultFonts = styles.fontFamily.split(",").map(style => style.trim())
    const fonts = [...new Set([...configFonts, ...defaultFonts])];
    const activeFont = styles.fontFamily.split(',')[0];

    this.updateFontOptions(fonts);
    this.selectMenu.value = activeFont;
  }
}

customElements.define("text-edit-menu", HTMLTextEditMenu);
