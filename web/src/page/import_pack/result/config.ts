import { t } from 'i18next';

export enum ResultStatusEnum {
  info = 'info',
  success = 'success',
  error = 'error',
}

export const getStatusConfig = (err_code) => ({
  [ResultStatusEnum.info]: {
    title: t('importResultPage.info.title'),
    desc: t('importResultPage.info.desc'),
  },
  [ResultStatusEnum.success]: {
    title: t('importResultPage.info.title'),
    desc: t('importResultPage.info.desc'),
  },
  [ResultStatusEnum.error]: {
    title: t('importResultPage.error.title'),
    desc: t(
      err_code === 'UnknownError' ? 'importResultPage.error.desc2' : 'importResultPage.error.desc1'
    ),
  },
});
