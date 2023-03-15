import { useState, useEffect } from 'react';
import { Button, Transfer, Tooltip } from 'antd';
import { useTranslation } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import { map } from 'lodash-es';

import { getProjectsApi, saveProjectsApi } from '@/api';
import { LOCAL_CONFIG, listStyle } from './config';

const ImportProjectPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const [targetKeys, setTargetKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [projects, setProjects] = useState<{ id: string; title: string }[]>([]);

  const handleBack = () => {
    saveData().then(() => {
      navigate('/page/analyze/team', { replace: true, state: location?.state });
    });
  };

  useEffect(() => {
    if (!location?.state) {
      handleBack();
    }
  }, [location]);

  useEffect(() => {
    getProjectsApi().then((res) => {
      setProjects(map(res.body.projects, (item) => ({ id: item.id, title: item.name })));
      setTargetKeys(map(res.body.cache));
    });
  }, []);

  const handleSubmit = () => {
    saveData().then(() => {
      navigate('/page/analyze/issue_map', {
        replace: true,
        state: {
          ...(location?.state || {}),
          projects: targetKeys,
        },
      });
    });
  };

  const saveData = () => {
    return saveProjectsApi(targetKeys);
  };
  const handleChange = (newTargetKeys: string[]) => {
    setTargetKeys(newTargetKeys);
  };

  const handleSelectChange = (sourceSelectedKeys: string[], targetSelectedKeys: string[]) => {
    setSelectedKeys([...sourceSelectedKeys, ...targetSelectedKeys]);
  };

  const renderButton = () => (
    <Button disabled={!targetKeys.length} className="ml-4" type="primary" onClick={handleSubmit}>
      {t('common.nextStep')}
    </Button>
  );

  return (
    <div className="h-full w-full">
      <div className="flex justify-between px-60">
        <h2>{t('importProject.title')}</h2>
        <div>
          <Button onClick={handleBack}>{t('common.back')}</Button>
          {targetKeys.length ? (
            renderButton()
          ) : (
            <Tooltip title={t('importProject.buttonTip')}>{renderButton()}</Tooltip>
          )}
        </div>
      </div>
      <Transfer
        dataSource={projects}
        rowKey={(record) => record.id}
        titles={[t('importProject.sourceTitle'), t('importProject.targetTitle')]}
        targetKeys={targetKeys}
        showSearch
        className="mt-8 h-5/6 px-60 py-0"
        locale={LOCAL_CONFIG}
        selectedKeys={selectedKeys}
        onChange={handleChange}
        onSelectChange={handleSelectChange}
        listStyle={listStyle}
        render={(item) => item.title}
        oneWay
        style={{ marginBottom: 16 }}
      />
    </div>
  );
};

export default ImportProjectPage;
