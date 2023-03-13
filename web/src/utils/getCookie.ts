
const cookieName = 'IMPORTTOOLS'
export const isHasCookie = () => {
  const result = document.cookie.match('(^|[^;]+)\\s*' + cookieName + '\\s*=\\s*([^;]+)');
  return result ? true :false;
};
