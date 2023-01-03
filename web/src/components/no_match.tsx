import { memo } from 'react';
import { Result, Button } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

/**
 * route lost tip
 */
const NoMatch = memo(() => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const handleBack = () => {
    navigate('/page/analyze/pack', {
      replace: true,
    });
  };

  return (
    <Result
      status="404"
      title="404"
      subTitle={t('common.noMatch')}
      extra={
        <Button type="primary" onClick={handleBack}>
          {t('common.back')}
        </Button>
      }
    />
  );
});

export default NoMatch;
