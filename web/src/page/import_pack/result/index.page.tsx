import { Button, Result } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate, useLocation } from 'react-router-dom';
import { downloadFile } from '@/utils/download';

import { resetImportStatusApi } from '@/api';
import { getStatusConfig, ResultStatusEnum } from './config';

const ResultPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { status = ResultStatusEnum.info, err_code = 'ServerError' } = location?.state || {};

  const handleDownLoad = () => {
    downloadFile(
      `${
        import.meta.env.VITE_PROXY_DOMAIN_REAL || window.location.origin
      }/import/log/download/current`
    );
  };

  const handleBack = () => {
    resetImportStatusApi().then(() => {
      navigate('/page/analyze/pack', { replace: true });
    });
  };

  const statusConfig = getStatusConfig(err_code)[status];

  return (
    <div className="flex flex-col items-center p-6">
      <Result
        status={status}
        title={`${statusConfig.title}!`}
        subTitle={statusConfig.desc}
        extra={[
          <Button key="downloadText" onClick={handleDownLoad}>
            {t('importResultPage.downloadText')}
          </Button>,
          <Button key="finish" type="primary" onClick={handleBack}>
            {t('common.finish')}
          </Button>,
        ]}
      />
    </div>
  );
};

export default ResultPage;
