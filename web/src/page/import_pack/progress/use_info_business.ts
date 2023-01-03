import { isUndefined } from 'lodash-es';
import { useState } from 'react';
import { useRafInterval } from 'ahooks';
import { useNavigate } from 'react-router-dom';

import { getImportInfoApi, ImportStatusEnum, resetImportStatusApi } from '@/api';
import { ResultStatusEnum } from '@/page/import_pack/result/config';
import type { ImportInfoType } from '@/api';

import { INTERVAL_TIME, checkConfig } from './config';

// progress info
const useInfoBusiness = () => {
  const navigate = useNavigate();
  const [info, setInfo] = useState<Partial<ImportInfoType>>({});

  const handleFetchInfo = () => {
    return getImportInfoApi()
      .then((res) => {
        setInfo(res.body);

        if (res.body.status === ImportStatusEnum.none) {
          navigate('/page/analyze/pack', {
            replace: true,
          });
        }

        const targetStatus = checkConfig[res.body.status];
        if (!isUndefined(targetStatus)) {
          navigate('/page/import_pack/result', {
            replace: true,
            state: { status: targetStatus },
          });
        }
      })
      .catch((e) => {
        if (e.err_code) {
          navigate('/page/import_pack/result', {
            replace: true,
            state: { status: ResultStatusEnum.error, err_code: e.err_code },
          });
        }

        throw e;
      });
  };

  const cancelInterval = useRafInterval(
    () => {
      handleFetchInfo().catch(() => {
        cancelInterval();
      });
    },
    INTERVAL_TIME,
    { immediate: true }
  );

  return {
    info,
    handleFetchInfo,
  };
};

export default useInfoBusiness;
