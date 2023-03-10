import Big from 'big.js';
import { t } from 'i18next';

// Use byte base unit, automatically round up when 1024 is reached
export const getFileUnit = (fileSize: number, needNumber = true, startUnit = 0) => {
  if (fileSize === 0) {
    return { size: fileSize, unit: t(`common.fileUnit.1`) };
  }

  const size = new Big(fileSize);

  if (size.lt(1024)) {
    const unit = t(`common.fileUnit.${startUnit}`);
    return needNumber ? { size: size.toFixed(2), unit } : unit;
  }

  return getFileUnit(size.div(1024), needNumber, startUnit + 1);
};

export  const getCookieValue = (name: string) => {
  const result = document.cookie.match('(^|[^;]+)\\s*' + name + '\\s*=\\s*([^;]+)');
  return result ? result.pop() : '';
};
