import { useState } from 'react';
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
    cancelAnalyzeApi().then(() => {
      navigate('/page/analyze/pack', { replace: true });
    });
  };

  const handleCancelMigrate = () => {
    cancelAnalyzeApi().then(() => {
      navigate('/page/home', { replace: true });
    });
  };

  return {
    handleBack,
    handleNext,
    info,
    resultData,
    handleCancelMigrate,
  };
};

export default useNavigateBusiness;
