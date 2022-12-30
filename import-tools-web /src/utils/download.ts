export const downloadFile = (url, filename = '', newTab = false) => {
  const elem = window.document.createElement('a');
  elem.download = filename;
  elem.href = url;
  if (newTab) {
    elem.target = '_blank';
  }
  document.body.appendChild(elem);
  elem.click();
  document.body.removeChild(elem);
};
