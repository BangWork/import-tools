import { Alert, Button, Form, Input, Modal } from 'antd';
import { useTranslation } from 'react-i18next';
import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

import ModalContent from '@/components/modal_content';
import { importStartApi } from '@/api';

const passwordRegexp = /^(?=.*[a-zA-Z])(?=.*\d)[^]{8,32}$/;

const InitPasswordPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const [form] = Form.useForm();
  const password = Form.useWatch('password', form);
  const againPassword = Form.useWatch('againPassword', form);

  const onBack = () => {
    navigate('/page/analyze/result', { replace: true });
  };

  const handleBack = () => {
    Modal.confirm({
      title: t('initPassword.modal.back.title'),
      content: t('initPassword.modal.back.desc'),
      onOk: onBack,
    });
  };

  useEffect(() => {
    if (!location?.state) {
      onBack();
    }
  }, [location]);

  const onFinish = () => {
    form.validateFields().then(() => {
      importStartApi({
        password,
        projectIds: location?.state.projects,
        issueTypeMap: location?.state.issue_type_map,
      }).then(() => {
        navigate('/page/import_pack/progress', { replace: true });
      });
    });
  };

  const handleValidate = (value: string) => {
    return new Promise((resolve, reject) => {
      if (passwordRegexp.test(value)) {
        return resolve('success');
      }
      return reject('fail');
    });
  };

  return (
    <Form form={form} onFinish={onFinish} autoComplete='off'>
      <ModalContent
        title={t('initPassword.title')}
        width='572px'
        footer={
          <Form.Item className='flex flex-row-reverse'>
            <Button className='mr-4' onClick={handleBack}>
              {t('common.back')}
            </Button>
            <Button type='primary' htmlType='submit'>
              {t('initPassword.startButton')}
            </Button>
          </Form.Item>
        }
      >
        <div className='p-6'>
          <Alert
            showIcon
            className='mb-4'
            message={
              <div className='p-2'>
                <div>{t('initPassword.tip.message1')}</div>
                <div style={{ color: '#FF4D4F' }}>{t('initPassword.tip.message2')}</div>
              </div>
            }
            type='info'
          />

          {/* form */}
          <div className='mx-8'>
            <Form.Item
              name='password'
              label={t('initPassword.form.init.label')}
              rules={[
                {
                  required: true,
                  message: t(
                    password
                      ? 'initPassword.form.init.error.rule'
                      : 'initPassword.form.init.error.empty'
                  ),
                  validator: (rule, value) => handleValidate(value),
                },
              ]}
            >
              <Input.Password allowClear placeholder={t('initPassword.form.init.placeholder')} />
            </Form.Item>
            <Form.Item
              name='againPassword'
              label={t('initPassword.form.again.label')}
              rules={[
                {
                  required: true,
                  message: t(
                    againPassword
                      ? 'initPassword.form.again.error.diff'
                      : 'initPassword.form.again.error.empty'
                  ),
                  validator: (rule, value) =>
                    handleValidate(password).then(() => {
                      if (value === password) {
                        return Promise.resolve('success');
                      }
                      return Promise.reject('fail');
                    }),
                },
              ]}
            >
              <Input.Password allowClear placeholder={t('initPassword.form.again.placeholder')} />
            </Form.Item>
          </div>

          <div className='opacity-50'>{t('initPassword.form.tip')}</div>
        </div>
      </ModalContent>
    </Form>
  );
};

export default InitPasswordPage;
