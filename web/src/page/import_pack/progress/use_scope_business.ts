import { useState, useEffect } from 'react';
import { useToggle, useRafInterval } from 'ahooks';

import { getImportScopeApi } from '@/api';
import type { ImportScopeType } from '@/api';

import { INTERVAL_TIME } from './config';

// scope of import
const useScopeBusiness = () => {
  const [infoModalState, { setLeft, setRight }] = useToggle(false);
  const [scope, setScope] = useState<Partial<ImportScopeType>>({});
  const [scopeLoading, setScopeLoading] = useState(true);

  const cancelInterval = useRafInterval(
    () => {
      if (infoModalState) {
        getImportScopeApi()
          .then((res) => {
            setScope(res.body);
          })
          .catch(() => {
            cancelInterval();
          });
      }
    },
    INTERVAL_TIME,
    { immediate: true }
  );

  useEffect(() => {
    if (infoModalState) {
      setScopeLoading(true);
      getImportScopeApi().then((res) => {
        setScope(res.body);
        setScopeLoading(false);
      });
    }
  }, [infoModalState]);

  return {
    scope,
    scopeLoading,
    infoModalState,
    hideModal: setLeft,
    showModal: setRight,
  };
};

export default useScopeBusiness;
