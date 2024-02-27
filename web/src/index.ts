import { createCharts } from "./charts";

function sidebarButtonHighlight() {
  const buttonActiveColor = (el: HTMLButtonElement): boolean => {
    const attr = el.attributes.getNamedItem('hx-push-url');

    if (attr?.nodeValue === window.location.pathname) {
      const button = el as HTMLButtonElement;

      button.style.color = "#0000FF";
      return true;
    }

    return false;
  }

  let activeButton: HTMLButtonElement | null;
  document.querySelectorAll('.sidebar-button').forEach(el => {
    const button = el as HTMLButtonElement;

    button.onclick = (_) => {
      const color = activeButton?.style.color;

      if (activeButton) {
        activeButton.style.color = button.style.color;
      }

      activeButton = button;
      activeButton.style.color = color || button.style.color;
    }

    if (buttonActiveColor(button)) activeButton = button;
  })
}

document.addEventListener('htmx:load', (_) => createCharts());

window.addEventListener('load', (_) => {
  sidebarButtonHighlight();
})
