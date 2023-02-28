import { Alert, Button, Form, Input, message, Modal } from 'antd';
import { useTranslation } from 'react-i18next';
import { useEffect, useState, useRef } from 'react';
import { isEmpty, trim } from 'lodash-es';
import { useNavigate, useLocation } from 'react-router-dom';
import type { InputRef } from 'antd';

import ModalContent from '@/components/modal_content';
import { submitEnvironmentApi } from '@/api';

import { ERROR_MAP } from './config';

const EnvironmentPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { state: locationState = {} } = useLocation();
  const [form] = Form.useForm();
  const [showServerError, setShowServerError] = useState(false);
  const url = Form.useWatch('url', form);
  const email = Form.useWatch('email', form);
  const password = Form.useWatch('password', form);
  const urlInputRef = useRef<InputRef>(null);

  const handleBack = () => {
    navigate('/page/analyze/pack', {
      replace: true,
    });
  };

  useEffect(() => {
    const initUrl = window.localStorage.getItem('environmentUrl');
    form.setFieldValue('url', initUrl);
  }, []);

  useEffect(() => {
    // localHome change need clear error tip
    setShowServerError(false);
  }, [url]);

  useEffect(() => {
    if (isEmpty(locationState)) {
      handleBack();
    }
  }, [locationState]);

  const onFinish = (res) => {
    const url = trim(res.url);
    const email = trim(res.email);
    const password = trim(res.password);
    form.setFieldsValue({
      url,
      email,
      password,
    });

    form.validateFields().then(() => {
      window.localStorage.setItem('environmentUrl', url);
      submitEnvironmentApi({
        ...locationState,
        url,
        email,
        password,
      })
        .then(() => {

          navigate('/page/analyze/progress', {
            replace: true,
          });
        })
        .catch((e) => {
          const { err_code, body } = e;
          const retryCount = body?.retry_count || 0;
          const msg = ERROR_MAP[err_code];

          if (err_code === 'NetworkError') {
            setShowServerError(true);
            urlInputRef.current?.focus();
            return;
          }

          if (err_code === 'AccountError' && retryCount <= 2) {
            Modal.error({
              title: t('environment.serverError.count.title'),
              content: t('environment.serverError.count.desc'),
              okText: t('common.ok'),
              onOk: handleBack,
            });
            return;
          }

          if (msg) {
            message.error(t(msg));
          }
        });
    });
  };

  const canSubmit = !!(url && email && password);

  return (
    <Form form={form} layout="vertical" onFinish={onFinish} autoComplete="off">
      <ModalContent
        title={t('environment.title')}
        footer={
          <Form.Item className="flex flex-row-reverse">
            <Button className="mr-4" onClick={handleBack}>
              {t('common.back')}
            </Button>
            <Button type="primary" disabled={!canSubmit} htmlType="submit">
              {t('environment.startButton')}
            </Button>
          </Form.Item>
        }
      >
        <div className="p-6">
          <Alert
            showIcon
            className="mb-4"
            message={
              <div className="p-2">
                <div>{t('environment.tip.message1')}</div>
                <div>{t('environment.tip.message2')}</div>
              </div>
            }
            type="info"
          />

          {/* form */}
          <Form.Item
            name="url"
            label={t('environment.url.label')}
            validateStatus={showServerError ? 'error' : undefined}
            help={showServerError ? t('environment.url.serverError') : undefined}
            rules={[
              {
                required: true,
                message: t('environment.url.emptyError'),
              },
            ]}
          >
            <Input ref={urlInputRef} autoFocus placeholder={t('environment.url.placeholder')} />
          </Form.Item>
          <Form.Item
            name="email"
            label={t('environment.email.label')}
            rules={[
              {
                required: true,
                message: t('environment.email.emptyError'),
              },
            ]}
          >
            <Input placeholder={t('common.placeholder')} />
          </Form.Item>
          <Form.Item
            name="password"
            label={t('environment.password.label')}
            rules={[
              {
                required: true,
                message: t('environment.password.emptyError'),
              },
            ]}
          >
            <Input.Password placeholder={t('common.placeholder')} />
          </Form.Item>
        </div>
      </ModalContent>
    </Form>
  );
};

export default EnvironmentPage;
