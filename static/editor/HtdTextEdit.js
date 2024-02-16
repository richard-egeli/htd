
const supportedFonts = [
  "Menlo", "Monaco", "Consolas", "Courier New", "Monospace",
]

/**
 * @param {string} text
 * @returns {number, number} Cols and Rows 
 */
function calculateTextAreaSize(text) {
  const lines = text.split('\n');
  const rows = lines.length;
  const cols = lines.sort(function(a, b) { return b.length - a.length })[0].length;
  return { cols, rows };
}

export default class HtdTextEdit extends HTMLElement {
  /**
   * @param {HTMLElement} target
   */
  constructor() {
    super();
    /**
     * @type {HTMLElement | null} target
     */
    this.target = null;

    /**
     * @type {string} substring
     */
    this.substring = "";

    /**
     * @type {(string) => void} onselect
     *
     * Callback when you select a substring inside of the textarea
     */
    this.onTextSelect = null;

    this.setupTextArea();
    this.setupOptionsMenu();

    /**
     * @type {HTMLDivElement} container
     */
    this.container = document.createElement("div");
    this.container.appendChild(this.optionsMenu);
    this.container.appendChild(this.textarea);
  }

  setTarget(target) {
    this.target = target;
    const compStyle = window.getComputedStyle(this.target);
    const { cols, rows } = calculateTextAreaSize(this.target.innerHTML);

    this.textarea.value = this.target.innerHTML;
    this.textarea.style.fontFamily = compStyle.fontFamily;
    this.textarea.style.fontWeight = compStyle.fontWeight;
    this.textarea.style.fontSize = compStyle.fontSize;
    this.textarea.style.boxSizing = "border-box";
    this.textarea.style.display = "block";
    this.textarea.style.padding = 0;
    this.textarea.style.margin = 0;
    this.textarea.style.resize = "none";
    this.textarea.style.width = "100%";
    this.textarea.cols = cols;
    this.textarea.rows = rows;

    this.target.parentNode.insertBefore(this.container, this.target);
    this.textarea.focus();
  }

  setupOptionsMenu() {
    this.optionsMenu = document.createElement("div");
    this.optionsMenu.style.backgroundColor = "white";
    this.optionsMenu.style.borderRadius = "4px";
    this.optionsMenu.style.padding = "4px 8px";
    this.optionsMenu.style.boxShadow = "1px 1px 4px 4px #00000020";
    this.optionsMenu.style.position = "absolute";
    this.optionsMenu.style.transform = "translateY(-115%)";

    const form = document.createElement("form");
    const label = document.createElement("label");
    const select = document.createElement('select');

    label.innerHTML = "Font: ";
    select.name = "font";

    supportedFonts.forEach((font, index) => {
      const el = document.createElement('option');
      el.value = font.toLowerCase();
      el.innerHTML = font;
      select.appendChild(el);

      if (this.textarea?.style.fontFamily.toLowerCase() === font.toLowerCase()) {
        select.options[index].selected = true;
      }
    });

    form.oninput = function(event) {
      event.preventDefault();
      const formData = new FormData(form);
      const font = formData.get("font");

      if (font != null) this.textarea.style.fontFamily = font;

      this.textarea.focus();
    }.bind(this);

    form.appendChild(label);
    form.appendChild(select);
    this.optionsMenu.appendChild(form);
  }

  setupTextArea() {
    /**
     * @type {HTMLTextAreaElement} textarea
     */
    this.textarea = document.createElement("textarea");
    this.textarea.style = parent.style;
    this.textarea.style.width = "100%";
    this.textarea.style.height = "auto";
    this.textarea.style.resize = "none";
    this.textarea.onmouseup = function() {
      const start = this.textarea.selectionStart;
      const end = this.textarea.selectionEnd;
      const sub = this.textarea.value.substring(start, end);
      this.substring = sub;
      this.onTextSelect?.(sub);
    }.bind(this);

    this.textarea.oninput = function() {
      const { cols, rows } = calculateTextAreaSize(this.textarea.value);
      this.textarea.cols = cols;
      this.textarea.rows = rows;
    }.bind(this);
  }
}

customElements.define("htd-text-edit", HtdTextEdit);

