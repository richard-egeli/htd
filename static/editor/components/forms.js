
const fontTypeOpts = {
  "Paragraph": "p",
  "Heading 1": "h1",
  "Heading 2": "h2",
  "Heading 3": "h3",
}

export function createFontTypeForm() {
  const form = document.createElement('form');
  const label = document.createElement('label');
  const select = document.createElement('select');

  Object.entries(fontTypeOpts).forEach(([key, value]) => {
    const option = document.createElement('option');

    option.value = value;
    option.innerHTML = key;

    select.appendChild(option);
  })

  form.appendChild(label);
  form.appendChild(select);

  return form;
}
