import {t } from 'i18next';
import { getFileUnit } from '@/utils/unit';

import {
  ImportStatusEnum,
} from '@/api';
import { ResultStatusEnum } from '@/page/import_pack/result/config';

export const CONTAINER_HEIGHT = 200;
export const ITEM_HEIGHT = 47;
export const LAST_END_FLAG = 'END_OF_FILE';
export const INTERVAL_TIME = 5000;

export const checkConfig = {
  [ImportStatusEnum.done]: ResultStatusEnum.success,
  [ImportStatusEnum.cancel]: ResultStatusEnum.info,
};

export const getCurrentDescConfig = (scope) => {
  const fileSize = getFileUnit(scope?.attachment_size || 0, true);
  const currentDescConfig = [
    {
      label: t('progressPage.viewModal.currentInfo.version'),
      value: scope.jira_version,
    },
    {
      label: t('progressPage.viewModal.currentInfo.projects'),
      value: scope.project_count,
      unit: t('progressPage.viewModal.unit'),
    },
    {
      label: t('progressPage.viewModal.currentInfo.works'),
      value: scope.issue_count,
      unit: t('progressPage.viewModal.unit'),
    },
    {
      label: t('progressPage.viewModal.currentInfo.members'),
      value: scope.member_count,
      unit: t('progressPage.viewModal.unit'),
    },
    {
      label: t('progressPage.viewModal.currentInfo.fileSize'),
      value: fileSize.size,
      unit: fileSize.unit,
    },
    {
      label: t('progressPage.viewModal.currentInfo.files'),
      value: scope.attachment_count,
      unit: t('progressPage.viewModal.unit'),
    },
  ];

  return currentDescConfig;
}
