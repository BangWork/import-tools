import { map } from 'lodash-es';
import { useEffect, useState } from 'react';
import { Button, Card, Descriptions, Space, Modal, Skeleton } from 'antd';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';

import { getResultApi } from '@/api';
import type { ResultType } from '@/api';

import { getCurrentDescConfig } from './config';

const DescriptionsWrap = styled(Descriptions)`
  max-height: 280px;
  overflow-y: auto;
  width: 100%;
`;

import ModalContent from '@/components/modal_content';

const ResultPage = () => {
  const { t } = useTranslation();
  const [info, setInfo] = useState<Partial<ResultType>>({});
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    getResultApi().then((res) => {
      setInfo(res.body);
      setLoading(false);
    });

  }, []);

  const onBack = () => {
    navigate('/page/analyze/pack', { replace: true });
  };

  const handleBack = () => {
    Modal.confirm({
      title: t('analyzeResult.modal.back.title'),
      content: t('analyzeResult.modal.back.desc'),
      onOk: onBack,
    });
  };

  const handleNext = () => {
    navigate('/page/analyze/team', {
      state: {
        import_history: info?.import_history || [],
        resolve_result: info?.resolve_result || {},
      },
    });
  };

  const currentDescConfig = getCurrentDescConfig(info?.resolve_result);

  const teamList = map(info?.import_history, (item) => ({
    key: item.team_uuid,
    label: item.team_name,
    value: item?.import_list || [],
  }));

  return (
    <ModalContent
      title={t('analyzeResult.title')}
      footer={
        <div className="flex flex-row-reverse">
          <Button disabled={loading} type="primary" onClick={handleNext}>
            {t('common.nextStep')}
          </Button>
          <Button className="mr-4" onClick={handleBack}>
            {t('common.back')}
          </Button>
        </div>
      }
    >
      <div className="flex justify-center p-4">
        <div className="flex justify-between">
          <Card
            bordered={false}
            className="mr-6"
            title={t('analyzeResult.current.title')}
            style={{ width: 300 }}
          >
            <Skeleton loading={loading} active paragraph={{ rows: 6 }} title={false}>
              <DescriptionsWrap column={1}>
                {map(currentDescConfig, (item) => (
                  <Descriptions.Item key={item.label + item.value} label={item.label}>
                    <Space>
                      {item?.value}
                      {item?.unit}
                    </Space>
                  </Descriptions.Item>
                ))}
              </DescriptionsWrap>
            </Skeleton>
          </Card>
          <Card
            bordered={false}
            title={t('analyzeResult.environment.title')}
            style={{ width: 300 }}
          >
            <Skeleton loading={loading} active>
              <DescriptionsWrap
                colon={false}
                column={1}
                contentStyle={{ color: '#606060' }}
                labelStyle={{ fontWeight: 500 }}
                layout="vertical"
              >
                {map(teamList, (item) => (
                  <Descriptions.Item key={item.key + item.label} label={item.label}>
                    <div>
                      {item.value.length
                        ? map(item.value, (sub) => (
                            <div>
                              {t('analyzeResult.environment.history', {
                                time: dayjs.unix(sub.import_time).format('YYYY-MM-DD HH:mm'),
                                version: sub.jira_version,
                              })}
                            </div>
                          ))
                        : t('analyzeResult.environment.empty')}
                    </div>
                  </Descriptions.Item>
                ))}
              </DescriptionsWrap>
            </Skeleton>
          </Card>
        </div>
      </div>
    </ModalContent>
  );
};

export default ResultPage;
