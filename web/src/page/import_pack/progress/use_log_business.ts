import { some } from 'lodash-es';
import { useState, useMemo, useEffect } from 'react';
import { useToggle, useRafInterval } from 'ahooks';

import { getLogApi, getLogHistoryApi } from '@/api';

import { LAST_END_FLAG, INTERVAL_TIME } from './config';

// import history log
const useLogBusiness = () => {
  const [logModalState, { setLeft, setRight }] = useToggle(true);
  const [logLoading, setLogLoading] = useState(true);
  const [logList, setLogList] = useState<string[]>([]);
  const hasEnd = useMemo(() => some(logList, (msg) => msg === LAST_END_FLAG), [logList]);

  useEffect(() => {
    setLogLoading(true);
    getLogHistoryApi().then((res) => {
      setLogList(res.body);
      setLogLoading(false);
    });
  }, []);

  const clearRafInterval = useRafInterval(() => {
    getLogApi(logList.length)
      .then((res) => {
        setLogList([...logList, ...res.body]);
      })
      .catch(() => {
        clearRafInterval();
      });
  }, INTERVAL_TIME);

  useEffect(() => {
    if (hasEnd) {
      clearRafInterval();
    }
  }, [hasEnd, logList]);

  return {
    logList,
    logModalState,
    logLoading,
    hasEnd,
    showLog: setLeft,
    hideLog: setRight,
  };
};

export default useLogBusiness;
