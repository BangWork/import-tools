import { Alert, Form, Input, toast, Modal } from '@ones-design/core';
import { useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';
import { trim } from 'lodash-es';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import FrameworkContent from '@/components/framework_content';
import Footer from '@/components/footer';
import { loginApi } from '@/api';

import { ERROR_MAP, COOKIENAME } from './config';

const EnvironmentPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [showServerError, setShowServerError] = useState(false);
  const [showAccountError, setShowAccountError] = useState(false);
  const url = Form.useWatch('url', form);
  const email = Form.useWatch('email', form);
  const password = Form.useWatch('password', form);
  const FormItemStyled = styled(Form.Item)`
    max-width: 570px;
  `;

  const ContentBox = styled.div`
    display: flex;
    padding-bottom: 20px;
  `;
  const ContentLeftBox = styled.div`
    width: 160px;
    padding-right: 20px;
  `;
  const ProfileStyled = styled.div`
    width: 100px;
    height: 100px;
    border-radius: 50px;
    overflow: hidden;
  `;
  const handleBack = () => {
    navigate('/page/home', {
      replace: true,
    });
  };

  const profile = window.localStorage.getItem('profile');
  const [isLogin, setIsLogin] = useState(false);
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

  useEffect(() => {
    setShowAccountError(false);
  }, [email, password]);

  useEffect(() => {
    if (getCookieValue(COOKIENAME)) {
      setIsLogin(true);
    }
  }, []);

  const handleNext = () => {
    navigate('/page/analyze/pack', {
      replace: true,
    });
  };

  const getCookieValue = (name: string) => {
    const result = document.cookie.match('(^|[^;]+)\\s*' + name + '\\s*=\\s*([^;]+)');
    return result ? result.pop() : '';
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
      window.localStorage.setItem('environmentUrl', url);
      loginApi({
        url,
        email,
        password,
      })
        .then((res) => {
          window.localStorage.setItem('profile', res?.body?.profile);
          handleNext();
        })
        .catch((e) => {
          const { err_code, body } = e;
          const retryCount = body?.retry_count || 0;
          const msg = ERROR_MAP[err_code];

          if (err_code === 'NetworkError') {
            setShowServerError(true);
            form.scrollToField('url');
            return;
          }

          if (err_code === 'AccountError' && retryCount <= 2) {
            setShowAccountError(true);
            form.scrollToField('email');
            return;
          }

          if (err_code === 'AccountError' && retryCount > 2) {
            Modal.warning({
              title: t('environment.serverError.count.title'),
              content: t('environment.serverError.count.desc'),
              okText: t('common.ok'),
              onOk: handleBack,
            });
          }
          if (msg) {
            toast.warning(t(msg));
          }
        });
    });
  };

  const canSubmit = !!(url && email && password);

  return !isLogin ? (
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
          <Footer
            handleCancelMigrate={{ fun: handleBack }}
            handleNext={{
              htmlType: 'submit',
              isDisabled: !canSubmit,
              text: t('environment.startButton'),
            }}
          ></Footer>
        }
      >
        <div className="oac-pr-4">
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
            <Input autoFocus placeholder={t('environment.url.placeholder')} />
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
          {showAccountError ? <Alert type="error" style={{ maxWidth: '570px' }}></Alert> : null}
        </div>
      </FrameworkContent>
    </Form>
  ) : (
    <FrameworkContent
      title={t('environment.title')}
      footer={
        <Footer handleCancelMigrate={{ fun: handleBack }} handleNext={{ fun: handleNext }}></Footer>
      }
    >
      <div className="oac-pt-2">
        <ContentBox>
          <ContentLeftBox>{t('common.back')}</ContentLeftBox>
          <ProfileStyled>
            <img src={profile} style={{ width: '100%' }} />
          </ProfileStyled>
        </ContentBox>
        <ContentBox>
          <ContentLeftBox>{t('common.back')}</ContentLeftBox>
          <div>{t('common.back')}</div>
        </ContentBox>
        <ContentBox>
          <ContentLeftBox>{t('common.back')}</ContentLeftBox>
          <div>{t('common.back')}</div>
        </ContentBox>
      </div>
    </FrameworkContent>
  );
};

export default EnvironmentPage;
