import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleFilled,
} from '@ant-design/icons';
import { useEffect, useState } from 'react';
import { message, Modal } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate ,useLocation} from 'react-router-dom';
import { useRafInterval } from 'ahooks';

import { getAnalyzeProgressInfoApi, cancelAnalyzeApi, AnalyzeStatusEnum } from '@/api';
import type { AnalyzeInfoType } from '@/api';

const TIME = 5000;

/**
 * Abnormal Judgment Tips
 */
const useNavigateBusiness = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();
  const [info, setInfo] = useState<Partial<AnalyzeInfoType>>({});
  const [loading, setLoading] = useState(true);
  const location = useLocation()
  const handleBackPack = () => {
    navigate('/page/analyze/pack', { replace: true });
  };
  const onFail = (msg) => {
    Modal.error({
      title: t('analyzeProgress.fail.title'),
      content: msg,
      okText: t('common.ok'),
      onOk: handleBackPack,
    });
  };

  const cancelInterval = useRafInterval(
    () => {
      getAnalyzeProgressInfoApi()
        .then((res) => {
          setInfo(res.body);
          setLoading(false);

          if (res.body.status === AnalyzeStatusEnum.fail) {
            onFail(t('analyzeProgress.fail.normalDesc', { name: res.body.backup_name }));
          }
        })
        .catch((error) => {
          if (error.code === 404) {
            onFail(t('analyzeProgress.fail.onExistDesc', { name: error.body.backup_name }));
          }
          cancelInterval();
        });
    },
    TIME,
    { immediate: true }
  );

  useEffect(() => {
    if (info?.status === AnalyzeStatusEnum.done) {
      navigate('/page/analyze/result', { replace: true ,state:{key:location?.state?.key}});
    }

    // Jira package was never analyze, need to go back to first page
    if (info?.status === AnalyzeStatusEnum.none) {
      handleBackPack();
    }
  }, [info.status]);

  const handleModalOk = () => {
    messageApi.open({
      type: 'loading',
      content: t('analyzeProgress.cancel.loading'),
      duration: 0,
    });

    cancelAnalyzeApi()
      .then(() => {
        messageApi.destroy();
        Modal.success({
          title: t('analyzeProgress.cancel.success'),
          icon: <CheckCircleOutlined />,
          okText: t('common.back'),
          okType: 'primary',
          onOk: handleBackPack,
        });
      })
      .catch(() => {
        messageApi.destroy();
        Modal.error({
          title: t('analyzeProgress.cancel.fail'),
          icon: <CloseCircleOutlined />,
          okText: t('common.ok'),
          okType: 'primary',
        });
      });
  };

  const handleCancel = () => {
    Modal.confirm({
      title: t('analyzeProgress.cancel.text'),
      icon: <ExclamationCircleFilled />,
      content: t('analyzeProgress.cancel.desc'),
      okText: t('common.ok'),
      okType: 'danger',
      cancelText: t('common.cancel'),
      onOk: handleModalOk,
    });
  };

  return {
    handleCancel,
    contextHolder,
    info,
    loading,
  };
};

export default useNavigateBusiness;
