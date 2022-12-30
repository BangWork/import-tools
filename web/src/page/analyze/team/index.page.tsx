import { Button, Form, Modal, Select, Tooltip } from 'antd';
import { useMemo, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { map, find, includes } from 'lodash-es';
import dayjs from 'dayjs';

import ModalContent from '@/components/modal_content';
import { chooseTeamApi } from '@/api';

import { WARNING_CONFIG, CompatibleList, WarningEnum } from './config';

const SelectStyled = styled(Select)`
  width: 210px;
`;

const TeamPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const [form] = Form.useForm();
  const teamUUID = Form.useWatch('teamUUID', form);

  const handleBack = () => {
    navigate('/page/analyze/result', {
      replace: true,
    });
  };

  useEffect(() => {
    if (!location?.state) {
      handleBack();
    }
  }, [location]);

  // reTry analyze process
  const handleBackPack = () => {
    navigate('/page/analyze/pack', {
      replace: true,
    });
  };

  const handleNext = () => {
    const selectTeam = find(location?.state?.import_history, { team_uuid: teamUUID });
    chooseTeamApi(teamUUID, selectTeam.team_name).then(() => {
      navigate('/page/analyze/import_project', {
        replace: true,
        state: location?.state,
      });
    });
  };

  const showFail = ({ time, version, packVersion }) => {
    Modal.error({
      title: t('teamPage.error.packDiff.title'),
      content: t('teamPage.error.packDiff.desc', {
        time,
        version,
        packVersion,
      }),
      okText: t('common.ok'),
      okType: 'primary',
      onOk: handleBackPack,
    });
  };

  const showWarning = (type, { time = '', version = '', packVersion = '' }) => {
    const config = WARNING_CONFIG[type];

    if (!config) return;

    Modal.confirm({
      title: config.title,
      content: config.renderDesc({ time, version, packVersion }),
      okText: t('common.ok'),
      okType: 'primary',
      cancelText: t(type === WarningEnum.import ? 'common.cancel' : 'teamPage.backButton'),
      onOk: handleNext,
      onCancel: () => {
        if (config.backPath) {
          navigate(config.backPath, {
            replace: true,
          });
        }
      },
    });
  };

  const onFinish = (res) => {
    form.validateFields().then(() => {
      const selectTeam = find(location?.state?.import_history, { team_uuid: res.teamUUID });
      const resolveResult = location?.state.resolve_result || {};
      const isCompatibleVersion = includes(CompatibleList, resolveResult.jira_version[0]);
      const isHistoryEmpty = !selectTeam.import_list.length;

      const prePack = selectTeam.import_list[0];
      // had import history and the Jira pack different
      if (!isHistoryEmpty && prePack.jira_server_id !== resolveResult?.jira_server_id) {
        showFail({
          time: dayjs.unix(prePack.import_time).format('YYYY-MM-DD HH:mm'),
          version: prePack.jira_version,
          packVersion: resolveResult.jira_version,
        });
        return;
      }

      // no import history and inCompatible
      if (isHistoryEmpty && !isCompatibleVersion) {
        showWarning(WarningEnum.version, { packVersion: resolveResult.jira_version });
        return;
      }

      // no import history and compatible version
      if (isHistoryEmpty && isCompatibleVersion) {
        handleNext();
        return;
      }

      // had import history and compatible version
      if (!isHistoryEmpty && isCompatibleVersion) {
        showWarning(WarningEnum.import, {
          time: dayjs.unix(prePack.import_time).format('YYYY-MM-DD HH:mm'),
          version: prePack.jira_version,
          packVersion: resolveResult.jira_version,
        });
        return;
      }

      // had import history and inCompatible version
      if (!isHistoryEmpty && !isCompatibleVersion) {
        showWarning(WarningEnum.importVersion, {
          time: dayjs.unix(prePack.import_time).format('YYYY-MM-DD HH:mm'),
          version: prePack.jira_version,
          packVersion: resolveResult.jira_version,
        });
        return;
      }
    });
  };

  const options = useMemo(() => {
    const data = location?.state?.import_history || [];
    return map(data, (item) => ({
      label: item.team_name,
      value: item.team_uuid,
    }));
  }, [location]);

  const renderButton = () => (
    <Button disabled={!teamUUID} type='primary' htmlType='submit'>
      {t('common.nextStep')}
    </Button>
  );

  return (
    <Form form={form} layout='vertical' onFinish={onFinish} autoComplete='off'>
      <ModalContent
        title={t('teamPage.title')}
        footer={
          <Form.Item className='flex flex-row-reverse'>
            <Button className='mr-4' onClick={handleBack}>
              {t('common.back')}
            </Button>
            {!teamUUID ? (
              <Tooltip title={t('teamPage.buttonTip')}>{renderButton()}</Tooltip>
            ) : (
              renderButton()
            )}
          </Form.Item>
        }
      >
        <div className='flex justify-center p-6'>
          {/* form */}
          <Form.Item
            name='teamUUID'
            style={{ width: '220px' }}
            label={t('teamPage.form.label')}
            rules={[{ required: true }]}
          >
            <SelectStyled options={options} placeholder={t('teamPage.form.placeholder')} />
          </Form.Item>
        </div>
      </ModalContent>
    </Form>
  );
};

export default TeamPage;
