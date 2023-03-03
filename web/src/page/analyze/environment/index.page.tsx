import { Input, message, Modal } from 'antd';
import { Alert, Form } from '@ones-design/core';
import { useTranslation } from 'react-i18next';
import { useEffect, useState, useRef } from 'react';
import { trim } from 'lodash-es';
import { useNavigate, useLocation } from 'react-router-dom';
import type { InputRef } from 'antd';
import styled from 'styled-components';
import FrameworkContent from '@/components/framework_content';
import Foooter from '@/components/footer';
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

  const FormItemStyled = styled(Form.Item)`
    max-width: 570px;
  `;

  const handleBack = () => {
    navigate('/page/home', {
      replace: true,
    });
  };

  useEffect(() => {
    const initUrl = window.localStorage.getItem('environmentUrl');
    form.setFieldsValue({ url: initUrl });
  }, []);

  useEffect(() => {
    form.scrollToField('url');
  }, []);

  useEffect(() => {
    // localHome change need clear error tip
    setShowServerError(false);
  }, [url]);

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
          navigate('/page/analyze/pack', {
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
    <Form
      form={form}
      layout="vertical"
      onFinish={onFinish}
      autoComplete="off"
      className={'oac-h-full oac-w-full'}
    >
      <FrameworkContent
        title={t('environment.title')}
        footer={
          <Foooter
            handleCancleMigrate={{ fun: handleBack }}
            handleNext={{
              htmlType: 'submit',
              isDisabled: !canSubmit,
              text: t('environment.startButton'),
            }}
          ></Foooter>
        }
      >
        <div>
          <Alert className="oac-mb-4" type="info">
            <div className="p-2">
              <div>{t('environment.tip.message1')}</div>
              <div>{t('environment.tip.message2')}</div>
            </div>
          </Alert>

          {/* form */}
          <FormItemStyled
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
          </FormItemStyled>
          <FormItemStyled
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
          </FormItemStyled>
          <FormItemStyled
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
          </FormItemStyled>
        </div>
      </FrameworkContent>
    </Form>
  );
};

export default EnvironmentPage;
