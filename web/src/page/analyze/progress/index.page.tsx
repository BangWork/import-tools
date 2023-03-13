import { useMemo } from 'react';
import { Alert } from '@ones-design/core';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import Big from 'big.js';
import { CheckmarkFilled, ErrorFilled } from '@ones-design/icons';
import type { BusinessProgressProps } from '@/components/business_progress';

import BusinessProgress from '@/components/business_progress';
import FrameworkContent from '@/components/framework_content';
import Footer from '@/components/footer';

import useNavigateBusiness from './use_navigate_business';

import styled from 'styled-components';

import ResultContent from './result_content';

const AnalyzeDescription = styled.div`
  padding-right: 20px;
  color: #909090;
`;
const AnalyzeDescriptionRight = styled.span`
  color: #606060;
  padding-left: 5px;
`;

const FailAnalyzeDescription = styled.div`
  color: #e52727;
`;
const ProgressPage = () => {
  const { t } = useTranslation();
  const { handleNext, handleBack, info, resultData, handleCancelMigrate } = useNavigateBusiness();

  const percent = useMemo(() => {
    if (info.spent_time && info.expected_time) {
      return new Big(info.spent_time).div(info.expected_time).times(100).toFixed(2);
    }
    return 0;
  }, [info?.spent_time, info?.expected_time]);

  const totalTime = info?.expected_time ? new Big(info?.expected_time).div(60).toFixed(0) : 0;
  const leftTime = info?.spent_time ? new Big(info?.spent_time).div(60).toFixed(0) : 0;

  const backupDataSource = [
    {
      date: '2020-04-12',
      key: 1,
      name: '组件库',
      number: 12,
      time: '14:07',
      user: 'htmlin',
    },
  ];
  const teamDataSource = [
    {
      date: '2020-04-12',
      key: 1,
      name: '组件库',
      number: 12,
      time: '14:07',
      user: 'htmlin',
    },
    {
      date: '2020-04-13',
      key: 2,
      name: '设计系统',
      number: 412,
      time: '14:08',
      user: 'lbg',
    },
    {
      date: '2020-04-12',
      key: 3,
      name: '组件库',
      number: 12,
      time: '14:07',
      user: 'htmlin',
    },
    {
      date: '2020-04-13',
      key: 4,
      name: '设计系统',
      number: 412,
      time: '14:08',
      user: 'lbg',
    },
  ];
  const backupName = window.localStorage.getItem('backupName');
  const progress: BusinessProgressProps = useMemo(() => {
    if (info?.status === 2) {
      return {
        status: 'success',
        statusText: 'analyzeProgress.backupMessage.status.success',
        percentTimeText: (
          <CheckmarkFilled style={{ color: '#24B47E', marginTop: '3px' }} fontSize="14" />
        ),
        bottomMessage: (
          <div className="oac-flex">
            <AnalyzeDescription>
              {t('analyzeProgress.backupMessage.analyzeBackupName')}
              <AnalyzeDescriptionRight>
                {t('analyzeProgress.tip.time', {
                  time: dayjs.unix(info?.start_time).format('YYYY-MM-DD HH:mm'),
                })}
              </AnalyzeDescriptionRight>
            </AnalyzeDescription>
            <AnalyzeDescription>
              {t('analyzeProgress.backupMessage.analyzeBackupTime')}
              <AnalyzeDescriptionRight>
                {t('analyzeProgress.tip.time', {
                  time: dayjs.unix(info?.start_time).format('YYYY-MM-DD HH:mm'),
                })}
              </AnalyzeDescriptionRight>
            </AnalyzeDescription>
            <AnalyzeDescription>
              {t('analyzeProgress.backupMessage.analyzeEnvironment')}
              <AnalyzeDescriptionRight>{t(backupName)}</AnalyzeDescriptionRight>
            </AnalyzeDescription>
          </div>
        ),
      };
    } else if (info?.status === 3) {
      return {
        status: 'fail',
        statusText: 'analyzeProgress.backupMessage.status.fail',
        percentTimeText: (
          <ErrorFilled style={{ color: '#24B47E', marginTop: '3px' }} fontSize="14" />
        ),
        bottomMessage: (
          <FailAnalyzeDescription>
            {t('analyzeProgress.backupMessage.analyzeFail')}
          </FailAnalyzeDescription>
        ),
      };
    } else {
      return {
        status: 'active',
        statusText: 'analyzeProgress.backupMessage.status.active',
        percentTimeText: t('analyzeProgress.timeMessage', {
          totalTime: Big(totalTime).lt(1) ? `1` : totalTime,
          leftTime: Big(leftTime).lt(1) ? `<1` : leftTime,
        }),
        bottomMessage: (
          <div className="oac-flex">
            <AnalyzeDescription>
              {t('analyzeProgress.backupMessage.analyzeTime')}
              <AnalyzeDescriptionRight>
                {t(dayjs.unix(info?.start_time).format('YYYY-MM-DD HH:mm'))}
              </AnalyzeDescriptionRight>
            </AnalyzeDescription>
            <AnalyzeDescription>
              {t('analyzeProgress.backupMessage.analyzeEnvironment')}
              <AnalyzeDescriptionRight>{'aaaa'}</AnalyzeDescriptionRight>
            </AnalyzeDescription>
          </div>
        ),
      };
    }
  }, [info?.status]);

  return (
    <FrameworkContent
      title={t('analyzeProgress.title')}
      footer={
        <Footer
          handleCancelMigrate={{ fun: handleCancelMigrate }}
          handleBack={{ fun: handleBack }}
          handleNext={{ fun: handleNext }}
        ></Footer>
      }
    >
      <div className=" oac-pr-4">
        <Alert type="info" className="oac-pb-4">
          {t('analyzeProgress.tip.environment', {
            name: info?.multi_team ? info?.org_name : info?.team_name,
          })}
        </Alert>
        <BusinessProgress
          title={t('analyzeProgress.backupMessage.title')}
          status={progress.status}
          statusText={t(progress.statusText)}
          percentDescription={t('analyzeProgress.backupMessage.analyzeProgress')}
          percentTimeText={progress.percentTimeText}
          bottomMessage={progress.bottomMessage}
          percent={percent}
        />
        {true ? (
          <ResultContent
            teamDataSource={teamDataSource}
            backupDataSource={backupDataSource}
            memory={'100G'}
          ></ResultContent>
        ) : null}
      </div>
    </FrameworkContent>
  );
};

export default ProgressPage;
