import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import {
  getAnalyzeProgressInfoApi,
  getImportInfoApi,
  ImportStatusEnum,
  AnalyzeStatusEnum,
} from '@/api';

const analyzeUrlMap = {
  [AnalyzeStatusEnum.none]: '/page/analyze/pack',
  [AnalyzeStatusEnum.doing]: '/page/analyze/progress',
  [AnalyzeStatusEnum.done]: '/page/analyze/result',
  [AnalyzeStatusEnum.fail]: '/page/analyze/progress',
};

const importProgressUrl = '/page/import_pack/progress';

const Page = () => {
  const navigate = useNavigate();

  useEffect(() => {
    getAnalyzeProgressInfoApi()
      .then((res) => {
        const analyzeStatus = res.body.status;
        const targetUrl = analyzeUrlMap[analyzeStatus];

        // After the analysis is complete, you need to confirm the status of the import
        if (analyzeStatus === AnalyzeStatusEnum.done) {
          return getImportInfoApi();
        }

        navigate(targetUrl || analyzeUrlMap[AnalyzeStatusEnum.none], { replace: true });
        return null;
      })
      .then((res) => {
        if (!res) return null;

        const importStatus = res.body.status;
        let targetUrl = importProgressUrl;

        // If not imported, you need to go to the analysis result page
        if (importStatus === ImportStatusEnum.none) {
          targetUrl = analyzeUrlMap[AnalyzeStatusEnum.done];
        }

        navigate(targetUrl, { replace: true });
      })
      .catch(() => {
        navigate(importProgressUrl, { replace: true });
      });
  }, []);

  return null;
};

export default Page;
