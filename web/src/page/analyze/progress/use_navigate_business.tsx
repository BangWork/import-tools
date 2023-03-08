import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { useRafInterval } from 'ahooks';

import { getAnalyzeProgressInfoApi, cancelAnalyzeApi, AnalyzeStatusEnum } from '@/api';
import type { AnalyzeInfoType } from '@/api';

const TIME = 5000;

/**
 * Abnormal Judgment Tips
 */
const useNavigateBusiness = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [info, setInfo] = useState<Partial<AnalyzeInfoType>>({});
  const [loading, setLoading] = useState(true);

  const handleBackPack = () => {
    navigate('/page/analyze/pack', { replace: true });
  };

  const onFail = (msg) => {
    return '11';
  };

  const cancelInterval = useRafInterval(
    () => {
      getAnalyzeProgressInfoApi()
        .then((res) => {
          setInfo(res.body);
          console.log(res.body);
          setLoading(false);

          if (res.body.status === AnalyzeStatusEnum.fail) {
            onFail(t('analyzeProgress.fail.normalDesc', { name: res.body.backup_name }));
          }
        })
        .catch((error) => {
          if (error.code === 404) {
            onFail(t('analyzeProgress.fail.onExistDesc', { name: error.body.backup_name }));
          }
          cancelInterval();
        });
    },
    TIME,
    { immediate: true }
  );

  const handleNext = () => {
    navigate('/page/analyze/result', { replace: true });
  };

  const handleBack = () => {
    navigate('/page/analyze/pack', { replace: true });
  };

  const handleModalOk = () => {
    cancelAnalyzeApi();
  };

  const handleCancel = () => {
    return 'aaa';
  };

  return {
    handleBack,
    handleNext,
    info,
    loading,
  };
};

export default useNavigateBusiness;
