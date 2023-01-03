import { t } from 'i18next';

export const errorCodeType = (code: number | string): string => {
  let errMessage: string = t('requestError.default');
  switch (code) {
    case 400:
      errMessage = t('requestError.400');
      break;
    case 401:
      errMessage = t('requestError.401');
      break;
    case 403:
      errMessage = t('requestError.403');
      break;
    case 404:
      errMessage = t('requestError.404');
      break;
    case 405:
      errMessage = t('requestError.405');
      break;
    case 408:
      errMessage = t('requestError.408');
      break;
    case 500:
      errMessage = t('requestError.500');
      break;
    case 501:
      errMessage = t('requestError.501');
      break;
    case 502:
      errMessage = t('requestError.502');
      break;
    case 503:
      errMessage = t('requestError.503');
      break;
    case 504:
      errMessage = t('requestError.504');
      break;
    case 505:
      errMessage = t('requestError.505');
      break;
    default:
      errMessage = t('requestError.other', { code });
  }
  return errMessage;
};
