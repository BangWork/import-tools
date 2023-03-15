import { useEffect, useState } from 'react';
import { Button, Form, Input, Select, Steps, Divider, Tooltip, Space, Alert, Popover } from 'antd';
import { useTranslation, Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { map, debounce, trim } from 'lodash-es';
import styled from 'styled-components';
import { SettingOutlined, QuestionCircleOutlined } from '@ant-design/icons';

import Image from '@/components/image';
import { checkPathApi, getBackUpApi } from '@/api';
import Guide1Image from './images/guide-1.png';
import Guide2Image from './images/guide-2.png';

const QuestionCircleOutlinedStyled = styled(QuestionCircleOutlined)`
  font-size: 14px;
  margin: 0 4px;
  color: #909090;
`;

const INIT_LOCAL_HOME = '/var/atlassian/application-data/jira';

const AnalyzePage = () => {
  const { t } = useTranslation();

  const navigate = useNavigate();
  const [showServerError, setShowServerError] = useState(false);
  const [showPackItem, setShowPackItem] = useState(false);
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

    // when patch is empty,need hide next item
    if (showPackItem) {
      setShowPackItem(false);
    }

    //
    const openBackup = window.sessionStorage.getItem('openBackup');
    if (localHome && openBackup) {
      handleCheckPath();
      window.sessionStorage.setItem('openBackup', '');
    }
  }, [localHome]);

  useEffect(() => {
    if (!showPackItem) {
      form.setFieldValue('backupName', undefined);
    }
  }, [showPackItem]);

  const handleFinish = (res) => {
    form.validateFields().then(() => {
      const { localHome, backupName } = res;
      window.localStorage.setItem('backupName', backupName);
      window.sessionStorage.setItem('openBackup', '1');
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
          setShowPackItem(true);
        })
        .then(() => {
          setPack({ options: [], loading: true });
          getBackUpApi(path).then(({ body = [] }) => {
            const initValue = window.localStorage.getItem('backupName');
            const options = map(body, (key) => {
              if (key === initValue && !backupName) {
                form.setFieldValue('backupName', initValue);
              }

              return { label: key, value: key };
            });

            setPack({
              options,
              loading: false,
            });
          });
        })
        .catch(() => {
          setShowPackItem(false);
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
    <div className="flex h-full w-full justify-center">
      <Form form={form} layout="vertical" onFinish={handleFinish}>
        <Space direction="vertical" size="large">
          <Alert
            style={{ width: '1000px' }}
            message={
              <div>
                {t('backupPage.guide.alert.desc')}
                <a
                  href="https://confluence.atlassian.com/jirakb/find-the-location-of-the-jira-home-directory-313466063.html"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {t('backupPage.guide.alert.link')}
                </a>
              </div>
            }
            type="info"
            showIcon
          />

          <Steps
            direction="vertical"
            items={[
              {
                status: 'process',
                title: (
                  <div style={{ fontWeight: '500' }}>
                    {t('backupPage.guide.step1.title')}
                    <Popover
                      placement="right"
                      content={<Image style={{ width: '780px' }} src={Guide1Image} />}
                    >
                      <QuestionCircleOutlinedStyled />
                    </Popover>
                  </div>
                ),
                description: (
                  <Space direction="vertical">
                    <span>{t('backupPage.guide.step1.desc1')}</span>
                    <Trans i18nKey="backupPage.guide.step1.desc2">
                      1.2 In the top navigation bar of Jira, click <SettingOutlined /> Jira
                      administration {'>'} System；
                    </Trans>
                    <span>{t('backupPage.guide.step1.desc3')}</span>
                  </Space>
                ),
              },
              {
                status: 'process',
                title: (
                  <div style={{ fontWeight: '500' }}>
                    {t('backupPage.guide.step2.title')}
                    <Popover
                      placement="right"
                      content={<Image style={{ width: '780px' }} src={Guide2Image} />}
                    >
                      <QuestionCircleOutlinedStyled />
                    </Popover>
                  </div>
                ),
                description: (
                  <div>
                    <div>
                      {t('backupPage.guide.step2.desc')}
                      <a
                        href="https://confluence.atlassian.com/adminjiraserver/jira-application-home-directory-938847746.html#:~:text=If%20Jira%20was%20installed%20using,data%2FJIRA%20(on%20Linux"
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {t('backupPage.guide.step2.link')}
                      </a>
                    </div>

                    <Space direction="vertical">
                      <div className="flex">
                        <div>
                          <Form.Item
                            style={{ paddingTop: '8px' }}
                            name="localHome"
                            initialValue={INIT_LOCAL_HOME}
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
                          {showPackItem ? (
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
                          style={{ marginTop: '37px', marginLeft: '16px' }}
                          type="primary"
                          onClick={handleCheckPath}
                        >
                          {t('backupPage.form.localHome.get')}
                        </Button>
                      </div>
                    </Space>
                  </div>
                ),
              },
            ]}
          />
        </Space>

        {/* bottom button */}

        <div
          className="mb-12"
          style={{ borderTop: '1px solid #eaeaea', marginTop: '-10px', paddingTop: '24px' }}
        >
          {!backupName ? (
            <Tooltip title={t('backupPage.form.tip')}>{renderButton()}</Tooltip>
          ) : (
            renderButton()
          )}
        </div>
      </Form>
    </div>
  );
};

export default AnalyzePage;
