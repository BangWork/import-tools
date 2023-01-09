import { t } from 'i18next';
import { getFileUnit } from '@/utils/unit';

export const getCurrentDescConfig = (resolve_result) => {
  const fileSize = getFileUnit(resolve_result?.attachment_size || 0, true);
  console.log(fileSize);

  const currentDescConfig = [
    {
      label: t('analyzeResult.current.version'),
      value: `Jira ${resolve_result?.jira_version}`,
    },
    {
      label: t('analyzeResult.current.projects'),
      value: resolve_result?.project_count,
      unit: t('analyzeResult.current.unit'),
    },
    {
      label: t('analyzeResult.current.works'),
      value: resolve_result?.issue_count,
      unit: t('analyzeResult.current.unit'),
    },
    {
      label: t('analyzeResult.current.members'),
      value: resolve_result?.member_count,
      unit: t('analyzeResult.current.unit'),
    },
    {
      label: t('analyzeResult.current.fileSize'),
      value: fileSize?.size,
      unit: fileSize?.unit || '',
    },
    {
      label: t('analyzeResult.current.files'),
      value: resolve_result?.attachment_count,
      unit: t('analyzeResult.current.unit'),
    },
  ];

  return currentDescConfig;
};
