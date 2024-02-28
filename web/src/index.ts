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

interface HtmxEvent extends Event {
  detail: {
    headers: {
      [key: string]: string;
    }
  }
}

const registerAddServerSentEvent = () => {
  try {
    const evtSource = new EventSource("/server/sent/event/browser/reload");
    evtSource.onmessage = function(_) {
      window.location.reload();
    };
  } catch (_) {
    window.location.reload();
  }
}

const registerAddCSRFTokenEvent = () => {
  document.body.addEventListener('htmx:configRequest', (event) => {
    const meta = document.querySelector('meta[name="csrf"]');

    if (meta) {
      const token = meta.getAttribute('content');
      const evt = event as HtmxEvent;

      if (evt.detail && evt.detail.headers && token) {
        evt.detail.headers['X-CSRF-Token'] = token;
      }
    }
  });
}

document.addEventListener('htmx:load', (_) => createCharts());

window.addEventListener('load', (_) => {
  sidebarButtonHighlight();
  registerAddCSRFTokenEvent();
  registerAddServerSentEvent();
})
