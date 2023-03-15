import {
  Alert,
  Button,
  Modal,
  Space,
  Divider,
  Descriptions,
  List,
  Typography,
  Skeleton,
  Empty,
} from 'antd';
import { useTranslation } from 'react-i18next';
import { map, last } from 'lodash-es';
import { useMemo, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import VirtualList from 'rc-virtual-list';
import dayjs from 'dayjs';
import Big from 'big.js';

import BusinessProgress from '@/components/business_progress';
import ModalContent from '@/components/modal_content';

import {
  cancelImportApi,
  pauseImportApi,
  continueImportApi,
  ImportStatusEnum,
  resetImportStatusApi,
} from '@/api';
import { ResultStatusEnum } from '@/page/import_pack/result/config';

import { getCurrentDescConfig, CONTAINER_HEIGHT, ITEM_HEIGHT } from './config';

import useScopeBusiness from './use_scope_business';
import useInfoBusiness from './use_info_business';
import useLogBusiness from './use_log_business';

const ProgressPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const logListRef = useRef(null);
  const { scope, scopeLoading, infoModalState, hideModal, showModal } = useScopeBusiness();
  const { info, handleFetchInfo } = useInfoBusiness();
  const { logList, logModalState, logLoading, hasEnd, showLog, hideLog } = useLogBusiness();

  useEffect(() => {
    if (logList.length && logListRef.current) {
      logListRef.current.scrollTo({
        index: logList.length - 2,
        align: 'bottom',
      });
    }
  }, [logList]);

  const handleStop = () => {
    pauseImportApi().then(() => {
      handleFetchInfo();
    });
  };

  const handleContinue = () => {
    continueImportApi()
      .then(() => {
        handleFetchInfo();
      })
      .catch((e) => {
        if (e.err_code === 'JiraFileNotFoundError') {
          Modal.error({
            title: t('progressPage.errorModal.miss.title'),
            okText: t('common.ok'),
            content: <div>{t('progressPage.errorModal.miss.desc')}</div>,
            onOk: () => {
              resetImportStatusApi().then(() => {
                navigate('/page/analyze/pack', {
                  replace: true,
                });
              });
            },
          });
        }
      });
  };

  const handleCancelImport = () => {
    cancelImportApi().then(() => {
      navigate('/page/import_pack/result', {
        replace: true,
        state: { status: ResultStatusEnum.info },
      });
    });
  };

  const handleCancel = () => {
    Modal.confirm({
      title: t('progressPage.cancelModal.title'),
      okText: t('common.ok'),
      okType: 'danger',
      content: <div>{t('progressPage.cancelModal.desc')}</div>,
      onOk: handleCancelImport,
    });
  };

  const currentDescConfig = getCurrentDescConfig(scope);

  const percent = useMemo(() => {
    if (info.spent_time && info.expected_time) {
      return new Big(info.spent_time).div(info.expected_time).times(100).toFixed(2);
    }
    return 0;
  }, [info?.spent_time, info?.expected_time]);

  const isPause = info.status === ImportStatusEnum.pause;
  const totalTime = info?.expected_time ? new Big(info?.expected_time).div(60).toFixed(0) : 0;
  const leftTime = info?.spent_time ? new Big(info?.spent_time).div(60).toFixed(0) : 0;

  return (
    <>
      {/* Scope Info modal */}
      <Modal
        title={t('progressPage.viewModal.title')}
        open={infoModalState}
        onCancel={hideModal}
        footer={
          <Button type="primary" onClick={hideModal}>
            {t('common.close')}
          </Button>
        }
      >
        <div>
          <Divider />
          <Skeleton className="py-4" paragraph={{ rows: 6 }} loading={scopeLoading}>
            <Descriptions title={t('progressPage.viewModal.currentInfo.title')} column={1}>
              {map(currentDescConfig, (item) => (
                <Descriptions.Item key={item.label + item.value} label={item.label}>
                  <Space>
                    {item.value == -1 ? (
                      t('progressPage.viewModal.loading')
                    ) : (
                      <>
                        {item.value}
                        {item?.unit}
                      </>
                    )}
                  </Space>
                </Descriptions.Item>
              ))}
            </Descriptions>
          </Skeleton>
        </div>
      </Modal>

      {/* Left Box */}
      <ModalContent
        title={t('progressPage.title')}
        width="572px"
        footer={
          <div className="flex flex-row-reverse">
            <Space>
              <Button onClick={handleCancel}>{t('progressPage.action.cancel')}</Button>
              {isPause ? (
                <Button onClick={handleContinue}>{t('progressPage.action.continue')}</Button>
              ) : (
                <Button onClick={handleStop}>{t('progressPage.action.stop')}</Button>
              )}
              <Button onClick={showModal}>{t('progressPage.action.view')}</Button>
              {logModalState ? null : (
                <Button onClick={showLog}>{t('progressPage.action.openLog')}</Button>
              )}
            </Space>
          </div>
        }
      >
        <div className="flex flex-col items-center p-6">
          <Alert
            showIcon
            className="mb-4 w-4/5"
            message={
              <div className="pl-2 pr-2">
                <div>{t('progressPage.tip.environment', { name: info.team_name || '' })}</div>
                <div>
                  {t('progressPage.tip.startTime', {
                    time: dayjs.unix(info?.start_time).format('YYYY-MM-DD HH:mm'),
                  })}
                </div>
                <div>{t('progressPage.tip.version', { name: info?.backup_name })}</div>
                <div>
                  {t('progressPage.tip.backUpTime', {
                    time: dayjs.unix(info?.backup_time).format('YYYY-MM-DD HH:mm'),
                  })}
                </div>
              </div>
            }
            type="info"
          />
          <BusinessProgress
            message={isPause ? t('progressPage.action.stop') : last(logList)}
            bottomMessage={t('progressPage.timeMessage', {
              totalTime: Big(totalTime).lt(1) ? `<1` : totalTime,
              leftTime: Big(leftTime).lt(1) ? `<1` : leftTime,
            })}
            percent={percent}
          />
        </div>
      </ModalContent>

      {/* Right Box */}
      {logModalState ? (
        <ModalContent
          title={t('progressPage.logTitle')}
          width="572px"
          className="ml-16"
          footer={
            <div className="flex flex-row-reverse">
              <Button onClick={hideLog}>{t('progressPage.action.closeLog')}</Button>
            </div>
          }
        >
          {/* The first loading shows the bottom skeleton */}
          {logLoading ? (
            <Skeleton className="px-4 py-8" paragraph={{ rows: 5 }} title={false} active />
          ) : (
            <List>
              {logList.length ? (
                <VirtualList
                  ref={logListRef}
                  data={hasEnd ? logList.slice(logList.length - 1) : logList}
                  height={CONTAINER_HEIGHT}
                  itemHeight={ITEM_HEIGHT}
                  itemKey="id"
                >
                  {(msg, index) => (
                    <List.Item style={{ height: `${ITEM_HEIGHT}px` }} key={index + msg}>
                      <Typography.Text ellipsis={{ tooltip: msg }}>{msg}</Typography.Text>
                    </List.Item>
                  )}
                </VirtualList>
              ) : (
                <Empty className="p-14" />
              )}
            </List>
          )}
        </ModalContent>
      ) : null}
    </>
  );
};

export default ProgressPage;
