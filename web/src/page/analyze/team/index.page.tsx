import { Alert, Table, Radio, Highlight } from '@ones-design/core';
import FrameworkContent from '@/components/framework_content';
import Footer from '@/components/footer';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { CheckmarkFilled } from '@ones-design/icons';
import { useNavigate } from 'react-router-dom';
import { filter, map } from 'lodash-es';
import dayjs from 'dayjs';
import { containsSubstring } from '@/utils/containsSubstring';
import { getTeamListApi } from '@/api';
import { WARNING_CONFIG, CompatibleList, WarningEnum } from './config';

const statusMap = {
  'Not migrated': ' #338FE5',
  Migrated: '#24B47E',
};
let comparisonTeamData = [];
const TeamPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [keyword, setKeyword] = useState('');
  const [selectedId, setSelectedId] = useState();
  const [selectTeam, setSelectTeam] = useState({});
  const [showErrorTips, setShowErrorTips] = useState(false);
  const [teamData, setTeamData] = useState([]);

  useEffect(() => {
    getTeamListApi()
      .then((res) => {
        comparisonTeamData = map(res.body, (item) => ({
          id: item.team_uuid,
          name: item.team_name,
          status: item.import_list.length === 0 ? 'Not migrated' : 'Migrated',
          backupName: item.import_list.length === 0 ? '-' : item.import_list[0].backup_name || '-',
          version: item.import_list.length === 0 ? '-' : item.import_list[0].jira_version || '-',
          jiraId: item.import_list.length === 0 ? '-' : item.import_list[0].jira_server_id || '-',
          time:
            item.import_list.length === 0
              ? '-'
              : dayjs(item.import_list[0].import_time).format('YYYY-MM-DD HH:mm:ss'),
        }));

        setTeamData([...comparisonTeamData]);
        if (comparisonTeamData.length === 1) {
          setSelectedId(comparisonTeamData[0].jiraId);
          setSelectTeam({ ...comparisonTeamData[0] });
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }, []);

  const columns = [
    {
      title: '',
      dataIndex: 'selected',
      key: 'selected',
      width: '3%',
      render: (_, record) => (
        <Radio onChange={() => handleSelect(record)} checked={selectedId === record.id} />
      ),
    },
    {
      render: (text, record) => {
        return (
          <div className="oac-flex oac-items-center">
            <CheckmarkFilled fontSize="16" style={{ marginRight: '5px' }}></CheckmarkFilled>
            <Highlight keyword={keyword}>{record.name}</Highlight>
          </div>
        );
      },
      title: t('teamPage.table.teamName'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('teamPage.table.migrateStatus'),
      dataIndex: 'status',
      key: 'status',
      render: (_, record) => (
        <div style={{ marginLeft: '10px' }}>
          <span
            style={{
              border: '1px solid',
              borderColor: statusMap[record.status],
              color: statusMap[record.status],
              borderRadius: '10px',
              padding: '0 5px',
            }}
          >
            {record.status}
          </span>
        </div>
      ),
    },
    {
      title: t('teamPage.table.jiraBackupName'),
      dataIndex: 'backupName',
      key: 'backupName',
    },
    {
      title: t('teamPage.table.jiraVersion'),
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: t('teamPage.table.jiraId'),
      dataIndex: 'jiraId',
      key: 'jiraId',
    },
    {
      title: t('teamPage.table.migrateTime'),
      dataIndex: 'time',
      key: 'time',
    },
  ];
  const handleBack = () => {
    navigate('/page/analyze/progress', {
      replace: true,
    });
  };

  const handleSelect = (record) => {
    setSelectedId(record.id);
    setSelectTeam({ ...record });
    setShowErrorTips(false);
  };

  const handleSearch = (e) => {
    const dataValue = filter(comparisonTeamData, (item) => {
      return containsSubstring(item.name, e.target.value);
    });
    setTeamData(dataValue);
    setKeyword(e.target.value);
  };
  const handleNext = () => {
    if (selectedId === undefined) {
      setShowErrorTips(true);
    } else {
      navigate('/page/analyze/import_project', {
        replace: true,
      });
    }
  };

  return (
    <FrameworkContent
      title={t('teamPage.title')}
      search={{ fun: handleSearch, text: t('teamPage.search') }}
      footer={
        <Footer
          handleBack={{ fun: handleBack }}
          handleNext={{ fun: handleNext }}
          handleCancelMigrate={{}}
        >
          {showErrorTips ? (
            <div style={{ color: '#E52727' }}>{t('teamPage.toSelectTeam')}</div>
          ) : selectedId === undefined ? (
            <div>{t('teamPage.selectZero')}</div>
          ) : (
            <div>{t('teamPage.selectTeam', { teamName: selectTeam.name })}</div>
          )}
        </Footer>
      }
    >
      <div>
        <Alert className="oac-pb-4">{t('teamPage.desc')}</Alert>

        <Table dataSource={teamData} columns={columns} rowKey="id" bordered={true} />
      </div>
    </FrameworkContent>
  );
};

export default TeamPage;
