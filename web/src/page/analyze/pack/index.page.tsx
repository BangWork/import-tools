import { useEffect, useState } from 'react';
import { Select } from 'antd';
import { Alert, Button, Space, Input, Form, Tooltip, toast, Popover } from '@ones-design/core';
import { useTranslation, Trans } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { map, debounce, trim } from 'lodash-es';
import styled from 'styled-components';
import { SettingOutlined } from '@ant-design/icons';
import { Launch, Help } from '@ones-design/icons';

import Image from '@/components/image';
import { checkPathApi, getBackUpApi } from '@/api';
import Guide1Image from './images/guide-1.png';
import Guide2Image from './images/guide-2.png';
import FrameworkContent from '@/components/framework_content';
import Footer from '@/components/footer';

const HelpStyled = styled(Help)`
  font-size: 14px;
  margin: 0 4px;
  color: #909090;
`;

const TitleBoxStyled = styled.div`
  display: flex;
  font-size: 16px;
  font-weight: 500;
  line-height: 24px;
  align-items: center;
  padding-bottom: 5px;
`;
const TitleNumber = styled.div`
  width: 24px;
  height: 24px;
  border-radius: 12px;
  background: #e8e8e8;
  font-size: 14px;
  line-height: 22px;
  text-align: center;
  margin-right: 10px;
`;

const AnalyzePage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [showServerError, setShowServerError] = useState(false);
  const [showPackItem, setShowPackItem] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isPackItemDisabled, setIsPackItemDisabled] = useState(false);
  const [pack, setPack] = useState({
    loading: true,
    options: [],
  });
  const [isSetted, setIsSetted] = useState(false);
  const [form] = Form.useForm();
  const localHome = Form.useWatch('localHome', form);
  const backupName = Form.useWatch('backupName', form);

  useEffect(() => {
    // localHome change need clear error tip
    setShowServerError(false);

    // when patch is empty,need hide next item
    if (!isPackItemDisabled) {
      setIsPackItemDisabled(true);
    }
  }, [localHome]);

  useEffect(() => {
    if (!showPackItem) {
      form.setFieldsValue({ backupName: undefined });
    }
  }, [showPackItem]);

  useEffect(() => {
    setIsSetted(true);
  }, []);
  const handleFinish = (res) => {
    if (isPackItemDisabled) {
      toast.warning(t('common.back'));
    } else {
      form.validateFields().then(() => {
        const { localHome, backupName } = res;
        window.localStorage.setItem('backupName', backupName);
        window.localStorage.setItem('localHome', localHome);
        navigate('/page/analyze/progress', {
          state: {
            localHome,
            backupName,
          },
        });
      });
    }
  };

  const handleCheckPath = debounce(() => {
    const path = trim(localHome);
    form.setFieldsValue({ localHome: path });
    form.validateFields(['localHome']).then(() => {
      setIsLoading(true);
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
                form.setFieldsValue({ backupName: initValue });
              }

              return { label: key, value: key };
            });

            setPack({
              options,
              loading: false,
            });
            setIsPackItemDisabled(false);
          });
        })
        .catch(() => {
          setShowServerError(true);
        })
        .finally(() => {
          setIsLoading(false);
        });
    });
  }, 500);

  const handleBack = () => {
    navigate('/page/analyze/environment');
  };
  return (
    <Form form={form} layout="vertical" onFinish={handleFinish} className={'oac-h-full oac-w-full'}>
      <FrameworkContent
        className=" oac-h-full oac-w-full "
        title={t('common.nextStep')}
        footer={
          <Footer
            handleBack={{
              fun: handleBack,
              text: 'common.back',
            }}
            handleNext={{
              isDisabled: !backupName,
              text: 'common.nextStep',
              type: 'primary',
              htmlType: 'submit',
            }}
          ></Footer>
        }
      >
        {isSetted ? (
          <div>
            <Alert type="info">
              <div>
                {t('backupPage.guide.alert.desc')}
                <a
                  href="https://confluence.atlassian.com/jirakb/find-the-location-of-the-jira-home-directory-313466063.html"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {t('backupPage.guide.alert.link')}&nbsp;
                  <Launch />
                </a>
              </div>
            </Alert>
            <div className="oac-pt-4">
              <TitleBoxStyled>
                <TitleNumber>1</TitleNumber>
                {t('backupPage.guide.step1.title')}
                <Popover placement="bottom" content={<Image src={Guide1Image} />}>
                  <HelpStyled />
                </Popover>
              </TitleBoxStyled>

              <Space direction="vertical">
                <span>{t('backupPage.guide.step1.desc1')}</span>
                <Trans i18nKey="backupPage.guide.step1.desc2">
                  1.2 In the top navigation bar of Jira, click <SettingOutlined /> Jira
                  administration {'>'} System；
                </Trans>
                <span>{t('backupPage.guide.step1.desc3')}</span>
              </Space>
            </div>
            <div className="oac-pt-4">
              <TitleBoxStyled>
                <TitleNumber>2</TitleNumber>
                {t('backupPage.guide.step2.title')}
                <Popover placement="bottom" content={<Image src={Guide2Image} />}>
                  <HelpStyled />
                </Popover>
              </TitleBoxStyled>
              <div>
                <div>
                  {t('backupPage.guide.step2.desc')}
                  <a
                    href="https://confluence.atlassian.com/adminjiraserver/jira-application-home-directory-938847746.html#:~:text=If%20Jira%20was%20installed%20using,data%2FJIRA%20(on%20Linux"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {t('backupPage.guide.step2.link')}&nbsp;
                    <Launch />
                  </a>
                </div>
              </div>
              <div>
                <div>{t('backupPage.form.desc')}</div>
                <div className="oac-flex oac-pt-2 ">
                  <div>
                    <Form.Item
                      name="localHome"
                      initialValue="/var/atlassian/application-data/jira"
                      validateStatus={showServerError ? 'error' : undefined}
                      help={
                        showServerError ? t('backupPage.form.localHome.serverError') : undefined
                      }
                      label={t('backupPage.form.localHome.label')}
                      rules={[
                        {
                          required: true,
                          message: t('backupPage.form.localHome.emptyError'),
                        },
                      ]}
                    >
                      <Input style={{ width: '300px' }} placeholder={t('common.placeholder')} />
                    </Form.Item>
                    {showPackItem ? (
                      <Tooltip title={isPackItemDisabled ? t('issueMap.table.disabledTip') : ''}>
                        <Form.Item
                          name="backupName"
                          label={t('backupPage.form.backup.label')}
                          rules={[
                            { required: true, message: t('backupPage.form.backup.emptyError') },
                          ]}
                        >
                          <Select
                            style={{ width: '300px' }}
                            disabled={isPackItemDisabled}
                            loading={pack.loading}
                            options={pack.options}
                            placeholder={t('backupPage.form.backup.placeholder')}
                          />
                        </Form.Item>
                      </Tooltip>
                    ) : null}
                  </div>
                  <Button
                    className="oac-ml-2"
                    style={{ marginTop: '27px' }}
                    type="primary"
                    onClick={handleCheckPath}
                    loading={isLoading}
                  >
                    {t('backupPage.form.localHome.get')}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div>
            <Alert type="info">
              <div>
                {t('backupPage.guide.alert.desc')}
                <a
                  href="https://confluence.atlassian.com/jirakb/find-the-location-of-the-jira-home-directory-313466063.html"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {t('backupPage.guide.alert.link')}&nbsp;
                  <Launch />
                </a>
              </div>
            </Alert>
            <Form.Item
              name="backupName"
              label={t('backupPage.form.backup.label')}
              rules={[{ required: true, message: t('backupPage.form.backup.emptyError') }]}
            >
              <Select
                style={{ width: '300px' }}
                options={pack.options}
                placeholder={t('backupPage.form.backup.placeholder')}
              />
            </Form.Item>
          </div>
        )}
      </FrameworkContent>
    </Form>
  );
};

export default AnalyzePage;
