import { generateRandomHash } from "./randhash.js";
import HtdTextEdit from "./HtdTextEdit.js"

/**
 * @param {string} text
 * @returns {number, number} Cols and Rows 
 */
function calculateTextarea(text) {
  const lines = text.split('\n');
  const rows = lines.length;
  const cols = lines.sort(function(a, b) { return b.length - a.length })[0];
  return { cols, rows };
}

/**
 * @param {HTMLTextAreaElement} textarea
 * @returns {HTMLFormElement} selection;
 */
function fontSelectionMenu(textarea) {
  const options = [
    "Menlo", "Monaco", "Consolas", "Courier New", "Monospace",
  ];

  const form = document.createElement("form");
  const label = document.createElement("label");
  const select = document.createElement('select');

  label.innerHTML = "Select a font:";
  select.name = "font";

  options.forEach((opt, index) => {
    const el = document.createElement('option');
    el.value = opt.toLowerCase();
    el.innerHTML = opt;
    select.appendChild(el);

    if (textarea.style.fontFamily.toLowerCase() === opt.toLowerCase()) {
      select.options[index].selected = true;
    }
  });

  form.oninput = function(event) {
    event.preventDefault();
    const formData = new FormData(form);
    const font = formData.get("font");

    if (font != null) textarea.style.fontFamily = font;

    textarea.focus();
  }

  form.appendChild(label);
  form.appendChild(select);
  return form;
}

/**
 * @param {HTMLTextAreaElement} textarea
 * @returns {HTMLDivElement} div
 */
function optionsMenu(textarea) {
  const div = document.createElement("div");
  const fontSelect = fontSelectionMenu(textarea, "Fantasy");

  div.style.position = "absolute";
  div.style.width = "100%";
  div.style.transform = "translateY(-100%)";
  div.appendChild(fontSelect);
  return div;
}


/**
 * @param {HTMLElement} el
 */
function initializeDragButton(el) {
  const button = document.createElement("button");

  /**
   * @type {HTMLDivElement} body
   */
  const body = document.querySelector('#main-content-area');

  button.style.position = "absolute";
  button.style.transform = "translateX(-100%)"
  button.innerHTML = "B";
  button.draggable = true;
  button.id = "test-id";

  button.onclick = function(event) {
    event.preventDefault();

    console.log('clicked');
  }

  button.ondragstart = function(event) {
    event.dataTransfer.setData("my-secret-thing", el.id);
  }

  body.ondragover = function(evt) {
    evt.preventDefault();
  }

  body.ondrop = function(event) {
    event.preventDefault();
    /**
     * @type {HTMLDivElement}
     */
    const target = event.currentTarget;
    const id = event.dataTransfer.getData("my-secret-thing");
    const element = document.querySelector(`#${id}`);
    const posY = event.clientY;
    const children = [...target.children];
    const array = [...children.map(child => child.getBoundingClientRect().y), posY]
    const sorted = array.sort((a, b) => b - a).reverse();
    const index = sorted.indexOf(posY);

    if (index >= children.length) {
      target.appendChild(element);
    } else {
      target.insertBefore(element, target.children[index]);
    }
  }

  button.ondragover = function(event) { event.preventDefault(); }
  el.appendChild(button);
}

/**
 * @param {HTMLTextAreaElement} textarea
 * @returns {HTMLDivElement} div
 */
function constructMenu(textarea) {
  const div = document.createElement("div");
  div.id = generateRandomHash(16);
  const button = initializeDragButton(div);
  const options = optionsMenu(textarea);

  div.style.width = "100%";
  div.style.height = "fit-content";
  div.style.boxSizing = "border-box";
  div.style.padding = "0px";
  div.style.margin = "0px";
  div.style.height = textarea.style.height;
  div.appendChild(options);
  div.appendChild(textarea);

  return div;
}

/**
 * @param {HTMLElement} el
 * @param {number} cols
 * @param {number} rows
 * @returns {HTMLTextAreaElement} 
 */
function initializeTextArea(el, cols, rows) {
  const textarea = document.createElement("textarea");
  const computedStyle = window.getComputedStyle(el);

  textarea.style = computedStyle;
  textarea.rows = rows;
  textarea.cols = cols;
  textarea.value = el.innerHTML;
  textarea.style.width = "100%";
  textarea.style.resize = "none";
  textarea.onmouseup = function() {
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const sub = textarea.value.substring(start, end);
    console.log(sub);
  }

  return textarea;
}


export function registerCallback() {
  document.querySelector("#hello").addEventListener('click', function(event) {
    /**
     * @type {HTMLParagraphElement}  
     */
    const target = event.target;
    // const { cols, rows } = calculateTextarea(p.innerHTML);
    // const menu = constructMenu(textarea);
    const text = document.createElement('htd-text-edit');

    text.setTarget(target);
    target.style.display = "none";

    text.focus();
    text.addEventListener('focusout', function(event) {
      if (event.currentTarget.contains(event.relatedTarget)) {
        event.preventDefault();
        return;
      }

      const textarea = text.textarea;
      target.innerhtml = textarea.value == "" ? "&nbsp" : textarea.value;
      target.style.fontWeight = textarea.style.fontWeight;
      target.style.fontFamily = textarea.style.fontFamily;
      target.style.fontSize = textarea.style.fontSize;
      target.parentNode.removeChild(text);
      target.style.display = "";
    })
    //
    // text.textarea.oninput = function() {
    //   const { cols, rows } = calculateTextarea(textarea.value);
    //   textarea.rows = rows;
    //   textarea.cols = cols;
    // }
  })
}
