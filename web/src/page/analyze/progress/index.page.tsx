import { useMemo } from 'react';
import { Alert, Button } from 'antd';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import Big from 'big.js';

import BusinessProgress from '@/components/business_progress';
import ModalContent from '@/components/modal_content';
import Loading from '@/components/loading';

import useNavigateBusiness from './use_navigate_business';

const ProgressPage = () => {
  const { t } = useTranslation();
  const { handleCancel, contextHolder, info, loading } = useNavigateBusiness();

  const percent = useMemo(() => {
    if (info.spent_time && info.expected_time) {
      return new Big(info.spent_time).div(info.expected_time).times(100).toFixed(2);
    }
    return 0;
  }, [loading, info?.spent_time, info?.expected_time]);

  const totalTime = info?.expected_time ? new Big(info?.expected_time).div(60).toFixed(0) : 0;
  const leftTime = info?.spent_time ? new Big(info?.spent_time).div(60).toFixed(0) : 0;

  return (
    <ModalContent
      title={t('analyzeProgress.title')}
      footer={
        <div className='flex flex-row-reverse'>
          <Button disabled={loading} className='mr-4' onClick={handleCancel}>
            {t('analyzeProgress.cancel.text')}
          </Button>
        </div>
      }
    >
      {contextHolder}
      <div className='flex flex-col justify-center p-6'>
        {loading ? (
          <Loading />
        ) : (
          <>
            <Alert
              className='mb-12'
              showIcon
              message={
                <div className='p-2'>
                  <div>
                    {t('analyzeProgress.tip.environment', {
                      name: info?.multi_team ? info?.org_name : info?.team_name,
                    })}
                  </div>
                  <div>
                    {t('analyzeProgress.tip.time', {
                      time: dayjs.unix(info?.start_time).format('YYYY-MM-DD HH:mm'),
                    })}
                  </div>
                </div>
              }
              type='info'
            />
            <BusinessProgress
              message={t('analyzeProgress.status.doing')}
              bottomMessage={t('analyzeProgress.timeMessage', {
                totalTime: Big(totalTime).lt(1) ? `<1` : totalTime,
                leftTime: Big(leftTime).lt(1) ? `<1` : leftTime,
              })}
              percent={percent}
            />
          </>
        )}
      </div>
    </ModalContent>
  );
};

export default ProgressPage;
