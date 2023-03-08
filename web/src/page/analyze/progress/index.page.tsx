import { useMemo } from 'react';
import { Alert } from '@ones-design/core';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';
import Big from 'big.js';
import { CheckmarkFilled, ErrorFilled } from '@ones-design/icons';

import BusinessProgress from '@/components/business_progress';
import FrameworkContent from '@/components/framework_content';
import Footer from '@/components/footer';

import useNavigateBusiness from './use_navigate_business';

import styled from 'styled-components';

import ResultContent from './result_content';
const ProgressPage = () => {
  const { t } = useTranslation();
  const { handleNext, handleBack, info, loading } = useNavigateBusiness();

  const percent = useMemo(() => {
    if (info.spent_time && info.expected_time) {
      return new Big(info.spent_time).div(info.expected_time).times(100).toFixed(2);
    }
    return 0;
  }, [loading, info?.spent_time, info?.expected_time]);

  const totalTime = info?.expected_time ? new Big(info?.expected_time).div(60).toFixed(0) : 0;
  const leftTime = info?.spent_time ? new Big(info?.spent_time).div(60).toFixed(0) : 0;

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

  const backupColumns = [
    {
      render: (text, record) => {
        return (
          <div className="oac-flex oac-items-center">
            <CheckmarkFilled fontSize="16" style={{ marginRight: '5px' }}></CheckmarkFilled>
            <div>{record.user}</div>
          </div>
        );
      },
      key: 'name',
      title: '产品名称',
    },
    {
      dataIndex: 'user',
      key: 'user',
      title: '产品负责人',
    },
    {
      dataIndex: 'date',
      key: 'date',
      title: '创建日期',
    },
    {
      dataIndex: 'number',
      key: 'number',
      title: '工作项数量',
    },
  ];
  const backupDataSource = [
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
  ];
  return (
    <FrameworkContent
      title={t('analyzeProgress.title')}
      footer={
        <Footer
          handleCancelMigrate={{}}
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
          <div>
            {t('analyzeProgress.tip.time', {
              time: dayjs.unix(info?.start_time).format('YYYY-MM-DD HH:mm'),
            })}
          </div>
        </Alert>

        <BusinessProgress
          title={t('analyzeProgress.status.doing')}
          status="active"
          statusText={'ddddd'}
          percentDescription={'addd'}
          percentTimeText={
            false ? (
              t('analyzeProgress.timeMessage', {
                totalTime: Big(totalTime).lt(1) ? `<1` : totalTime,
                leftTime: Big(leftTime).lt(1) ? `<1` : leftTime,
              })
            ) : true ? (
              <CheckmarkFilled style={{ color: '#24B47E', marginTop: '3px' }} fontSize="14" />
            ) : (
              <ErrorFilled style={{ color: '#24B47E', marginTop: '3px' }} fontSize="14" />
            )
          }
          bottomMessage={
            true ? (
              <div className="oac-flex">
                <AnalyzeDescription>
                  {'aaa:'}
                  <AnalyzeDescriptionRight>{'aaaa'}</AnalyzeDescriptionRight>
                </AnalyzeDescription>
                <AnalyzeDescription>
                  {'nbb:'}
                  <AnalyzeDescriptionRight>{'aaaa'}</AnalyzeDescriptionRight>
                </AnalyzeDescription>
                <AnalyzeDescription>
                  {'ddd:'}
                  <AnalyzeDescriptionRight>{'aaaa'}</AnalyzeDescriptionRight>
                </AnalyzeDescription>
              </div>
            ) : (
              <FailAnalyzeDescription>{'adaa'}</FailAnalyzeDescription>
            )
          }
          percent={percent}
        />
        <ResultContent
          backupColumns={backupColumns}
          backupDataSource={backupDataSource}
          memory={'100G'}
        ></ResultContent>
      </div>
    </FrameworkContent>
  );
};

export default ProgressPage;
