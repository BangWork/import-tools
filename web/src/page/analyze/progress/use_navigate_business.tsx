import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { useRafInterval } from 'ahooks';

import {
  getAnalyzeProgressInfoApi,
  cancelAnalyzeApi,
  AnalyzeStatusEnum,
  getResultApi,
} from '@/api';
import type { AnalyzeInfoType } from '@/api';

const TIME = 5000;

/**
 * Abnormal Judgment Tips
 */
const useNavigateBusiness = () => {
  const navigate = useNavigate();
  const [info, setInfo] = useState<Partial<AnalyzeInfoType>>({});
  const [resultData, setResultData] = useState({});
  const cancelInterval = useRafInterval(
    () => {
      getAnalyzeProgressInfoApi()
        .then((res) => {
          setInfo(res.body);
          if (res.body.status === AnalyzeStatusEnum.fail) {
            cancelInterval();
          }
          if (res.body.status === AnalyzeStatusEnum.done) {
            cancelInterval();
            getResultApi().then((res) => {
              setResultData(res.body);
            });
          }
        })
        .catch(() => {
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
    cancelAnalyze();
    navigate('/page/analyze/pack', { replace: true });
  };

  const cancelAnalyze = () => {
    cancelAnalyzeApi();
  };

  return {
    handleBack,
    handleNext,
    info,
    resultData,
    cancelAnalyze,
  };
};

export default useNavigateBusiness;
