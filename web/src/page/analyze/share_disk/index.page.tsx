import { useEffect, useState } from 'react';
import { Alert, Button, Form, Input } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import ModalContent from '@/components/modal_content';
import { checkDiskPathApi, setDiskPathApi } from '@/api';

const ShareDiskPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const [form] = Form.useForm();
  const [needPath, setNeedPath] = useState(false);
  const [showServerError, setShowServerError] = useState(false);
  const path = Form.useWatch('path', form);

  useEffect(() => {
    setShowServerError(false);
  }, [path]);

  const handleBack = () => {
    navigate('/page/analyze/issue_map', {
      replace: true,
      state: location?.state,
    });
  };

  const handleNext = () => {
    navigate('/page/import_pack/init_password', { replace: true, state: location?.state });
  };

  const handleSubmit = (res) => {
    form.validateFields().then(() => {
      checkDiskPathApi(res?.path)
        .then(() => {
          setDiskPathApi(true, res?.path).then(() => {
            handleNext();
          });
        })
        .catch(() => {
          setShowServerError(true);
        });
    });
  };

  const handleUse = (e) => {
    e.preventDefault();
    setNeedPath(true);
  };

  const handleUnUse = () => {
    setDiskPathApi(false).then(() => {
      handleNext();
    });
  };

  return (
    <Form form={form} layout="vertical" onFinish={handleSubmit} autoComplete="off">
      <ModalContent
        title={t('shareDisk.title')}
        footer={
          <div className="flex flex-row-reverse">
            {needPath ? (
              <>
                <Button type="primary" className="mr-4" htmlType="submit">
                  {t('common.nextStep')}
                </Button>
                <Button className="mr-4" onClick={handleBack}>
                  {t('common.back')}
                </Button>
              </>
            ) : (
              <>
                <Button type="primary" className="mr-4" onClick={handleUse}>
                  {t('common.use')}
                </Button>
                <Button className="mr-4" onClick={handleUnUse}>
                  {t('common.unUse')}
                </Button>
              </>
            )}
          </div>
        }
      >
        <div className="p-6">
          <Alert
            className="mb-4"
            showIcon
            message={
              <div className="p-2">
                <div>{t('shareDisk.tip.message1')}</div>
                <div>{t('shareDisk.tip.message2')}</div>
                <div>{t('shareDisk.tip.message3')}</div>
              </div>
            }
            type="info"
          />
          {needPath ? (
            <Form.Item
              name="path"
              className="px-10"
              validateStatus={showServerError ? 'error' : undefined}
              help={showServerError ? t('shareDisk.form.serverError') : undefined}
              label={t('shareDisk.form.label')}
              rules={[{ required: true, message: t('shareDisk.form.emptyError') }]}
            >
              <Input autoFocus placeholder={t('common.placeholder')} />
            </Form.Item>
          ) : null}
        </div>
      </ModalContent>
    </Form>
  );
};

export default ShareDiskPage;
