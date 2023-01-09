import { Alert, Button, Form, Input, message, Modal } from 'antd';
import { useTranslation } from 'react-i18next';
import { useEffect } from 'react';
import { isEmpty, trim } from 'lodash-es';
import { useNavigate, useLocation } from 'react-router-dom';

import ModalContent from '@/components/modal_content';
import { submitEnvironmentApi } from '@/api';

import { ERROR_MAP } from './config';

const EnvironmentPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { state: locationState = {} } = useLocation();
  const [form] = Form.useForm();

  const handleBack = () => {
    navigate('/page/analyze/pack', {
      replace: true,
    });
  };

  useEffect(() => {
    if (isEmpty(locationState)) {
      handleBack();
    }
  }, [locationState]);

  const handleClear = () => {
    form.resetFields();
  };

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
            Modal.error({
              title: t('environment.serverError.network.title'),
              content: t('environment.serverError.network.desc'),
              okText: t('common.ok'),
              onOk: handleClear,
            });
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

  return (
    <Form form={form} layout="vertical" onFinish={onFinish} autoComplete="off">
      <ModalContent
        title={t('environment.title')}
        footer={
          <Form.Item className="flex flex-row-reverse">
            <Button className="mr-4" onClick={handleBack}>
              {t('common.back')}
            </Button>
            <Button type="primary" htmlType="submit">
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
            rules={[
              {
                required: true,
                message: t('environment.url.emptyError'),
              },
            ]}
          >
            <Input placeholder={t('common.placeholder')} />
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
