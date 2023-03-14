import { useState, useEffect } from 'react';
import { Alert, Table, Checkbox } from '@ones-design/core';
import FrameworkContent from '@/components/framework_content';
import Footer from '@/components/footer';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { containsSubstring } from '@/utils/containsSubstring';
import { filter, map } from 'lodash-es';

import { getProjectsApi } from '@/api';

let comparisonProjectData = [];
const ImportProjectPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [selectedSet, setSelectedSet] = useState<Set<string>>(new Set());
  const [projectData, setProjectData] = useState([]);

  const columns = [
    {
      title: '',
      dataIndex: 'selected',
      key: 'selected',
      width: '3%',
      render: (_, record) => (
        <Checkbox onChange={() => handleSelect(record)} checked={selectedSet.has(record.id)} />
      ),
    },
    {
      title: t('importProject.table.projectName'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('importProject.table.projectKey'),
      dataIndex: 'projectKey',
      key: 'projectKey',
    },
    {
      title: t('importProject.table.leader'),
      dataIndex: 'leader',
      key: 'leader',
    },
    {
      title: t('importProject.table.projectClassification'),
      dataIndex: 'category',
      key: 'category',
    },
    {
      title: t('importProject.table.issueCount'),
      dataIndex: 'issueCount',
      key: 'issueCount',
    },
  ];
  const handleBack = () => {
    navigate('/page/analyze/team', { replace: true });
  };

  const handleSelect = (record) => {
    if (selectedSet.has(record.id)) {
      selectedSet.delete(record.id);
    } else {
      selectedSet.add(record.id);
    }
    setSelectedSet(new Set(selectedSet));
  };

  useEffect(() => {
    getProjectsApi().then((res) => {
      comparisonProjectData = map(res.body, (item) => ({
        name: item.name,
        projectKey: item.key,
        leader: item.assign,
        category: item.category,
        issueCount: item.issue_count,
        type: item.type,
      }));
      setProjectData([...comparisonProjectData]);
    });
  }, []);

  const handleSubmit = () => {
    navigate('/page/analyze/issue_map', {
      replace: true,
      state: {},
    });
  };

  const handleSearch = (e) => {
    const projectDataValue = filter(comparisonProjectData, (item) => {
      return (
        containsSubstring(item.name, e.target.value) ||
        containsSubstring(item.projectKey, e.target.value)
      );
    });
    setProjectData(projectDataValue);
  };

  const handleDownload = () => {
    console.log('download');
  };

  const handleConfig = () => {
    window.open('www.baidu.com', '_blank');
  };
  return (
    <FrameworkContent
      title={t('importProject.title')}
      search={{ fun: handleSearch, text: t('importProject.search') }}
      download={{ fun: handleDownload }}
      config={{ fun: handleConfig, text: t('importProject.noSupportMigrateProject') }}
      footer={
        <Footer
          handleBack={{ fun: handleBack }}
          handleNext={{ fun: handleSubmit }}
          handleCancelMigrate={{}}
        ></Footer>
      }
    >
      <Alert>
        <div>{t('importProject.desc1')}</div>
        <div>
          {t('importProject.desc2')}
          <a target="_blank" rel="noopener noreferrer">
            {t('importProject.link')}
          </a>
        </div>
      </Alert>
      <div className="oac-pt-4">
        <Table dataSource={projectData} columns={columns} rowKey="id" bordered={true} />
      </div>
    </FrameworkContent>
  );
};

export default ImportProjectPage;
