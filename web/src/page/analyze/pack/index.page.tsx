import { useEffect, useState } from 'react';
import { Button, Form, Input, Select, Steps, Divider, Tooltip, Space } from 'antd';
import { useTranslation, Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { map, debounce, trim } from 'lodash-es';
import { SettingOutlined } from '@ant-design/icons';

import Guide1Image from './images/guide-1.png';
import Guide2Image from './images/guide-2.png';
import Guide3Image from './images/guide-3.png';

import Image from '@/components/image';
import { checkPathApi, getBackUpApi } from '@/api';

const AnalyzePage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [showServerError, setShowServerError] = useState(false);
  const [showBackUp, setShowBackUp] = useState(false);
  const [pack, setPack] = useState({
    loading: true,
    options: [],
  });
  const [form] = Form.useForm();
  const localHome = Form.useWatch('localHome', form);
  const backupName = Form.useWatch('backupName', form);

  useEffect(() => {
    // localHome change need clear error tip
    setShowServerError(false);

    // when patch is empty,need lock submit
    if (showBackUp) {
      setShowBackUp(false);
    }
  }, [localHome]);

  useEffect(() => {
    if (!showBackUp) {
      form.setFieldValue('backupName', undefined);
    }
  }, [showBackUp]);

  const handleFinish = (res) => {
    form.validateFields().then(() => {
      const { localHome, backupName } = res;
      navigate('/page/analyze/environment', {
        state: {
          localHome,
          backupName,
        },
      });
    });
  };

  const handleCheckPath = debounce(() => {
    const path = trim(localHome);
    form.setFieldValue('localHome', path);
    form.validateFields(['localHome']).then(() => {
      checkPathApi(path)
        .then(() => {
          setShowBackUp(true);
        })
        .then(() => {
          setPack({ options: [], loading: true });
          getBackUpApi(path).then(({ body = [] }) => {
            const options = map(body, (key) => ({ label: key, value: key }));
            setPack({
              options,
              loading: false,
            });
          });
        })
        .catch(() => {
          setShowBackUp(false);
          setShowServerError(true);
        });
    });
  }, 500);

  const renderButton = () => (
    <Button disabled={!backupName} type="primary" className="mr-4" htmlType="submit">
      {t('common.nextStep')}
    </Button>
  );

  return (
    <div className="flex h-full w-full">
      <Form form={form} layout="vertical" onFinish={handleFinish}>
        <Space direction="vertical" size="large">
          <h3>{t('backupPage.guide.title')}</h3>
          <div>
            <div>{t('backupPage.guide.desc')}</div>

            <Steps
              direction="vertical"
              className="mt-8"
              items={[
                {
                  title: (
                    <Trans i18nKey="backupPage.guide.step1.title">
                      Click the <SettingOutlined /> icon in the top menu bar, then click System
                    </Trans>
                  ),
                  description: <Image src={Guide1Image} />,
                  status: 'process',
                },
                {
                  title: t('backupPage.guide.step2.title'),
                  description: <Image src={Guide2Image} />,
                  status: 'process',
                },
                {
                  title: t('backupPage.guide.step3.title'),
                  description: <Image src={Guide3Image} />,
                  status: 'process',
                },
                {
                  status: 'process',
                  title: t('backupPage.form.title'),
                  description: (
                    <Space direction="vertical">
                      <div>{t('backupPage.form.desc')}</div>
                      <div className="flex">
                        <div>
                          <Form.Item
                            name="localHome"
                            validateStatus={showServerError ? 'error' : undefined}
                            help={
                              showServerError
                                ? t('backupPage.form.localHome.serverError')
                                : undefined
                            }
                            label={t('backupPage.form.localHome.label')}
                            rules={[
                              {
                                required: true,
                                message: t('backupPage.form.localHome.emptyError'),
                              },
                            ]}
                          >
                            <Input
                              style={{ width: '300px' }}
                              placeholder={t('common.placeholder')}
                            />
                          </Form.Item>
                          {showBackUp ? (
                            <Form.Item
                              name="backupName"
                              label={t('backupPage.form.backup.label')}
                              rules={[
                                { required: true, message: t('backupPage.form.backup.emptyError') },
                              ]}
                            >
                              <Select
                                style={{ width: '300px' }}
                                loading={pack.loading}
                                options={pack.options}
                                placeholder={t('backupPage.form.backup.placeholder')}
                              />
                            </Form.Item>
                          ) : null}
                        </div>
                        <Button
                          className="ml-8"
                          style={{ marginTop: '30px' }}
                          type="primary"
                          onClick={handleCheckPath}
                        >
                          {t('backupPage.form.localHome.get')}
                        </Button>
                      </div>
                      {/* button groups */}
                      <Divider />
                      <div className="mb-12 flex flex-row-reverse">
                        {!backupName ? (
                          <Tooltip title={t('backupPage.form.tip')}>{renderButton()}</Tooltip>
                        ) : (
                          renderButton()
                        )}
                      </div>
                    </Space>
                  ),
                },
              ]}
            />
          </div>
        </Space>
      </Form>
    </div>
  );
};

export default AnalyzePage;
