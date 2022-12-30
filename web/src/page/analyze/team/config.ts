import { t } from 'i18next';

export enum WarningEnum {
  import,
  importVersion,
  version,
}

// major version
export const CompatibleList = ['7'];

export const WARNING_CONFIG = {
  [WarningEnum.import]: {
    title: t('teamPage.error.importDiff.title'),
    renderDesc: (params) => t('teamPage.error.importDiff.desc', params),
  },
  [WarningEnum.importVersion]: {
    title: t('teamPage.error.importVersionDiff.title'),
    renderDesc: (params) => t('teamPage.error.importVersionDiff.desc', params),
    backPath: '/page/analyze/pack',
  },
  [WarningEnum.version]: {
    title: t('teamPage.error.versionDiff.title'),
    renderDesc: (params) => t('teamPage.error.versionDiff.desc', params),
    backPath: '/page/analyze/pack',
  },
};
