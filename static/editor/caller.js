

export function registerCallback() {
  const elements = document.querySelectorAll("#hello");

  /**
   * @type {HTMLTextEditElement} text
   */
  const text = document.createElement("text-edit");

  elements.forEach(el => {
    text.addClickEvent(el);
  })
}
