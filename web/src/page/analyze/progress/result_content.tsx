import { memo } from 'react';
import type { FC } from 'react';
import { Table } from '@ones-design/core';
import { CheckmarkFilled, Warning, Launch } from '@ones-design/icons';
import { t } from 'i18next';

interface ResultContentProgressProps {
  teamDataSource: any[];
  backupDataSource: any[];
  memory: string;
}

const backupColumns = [
  {
    dataIndex: 'version',
    key: 'version',
    title: t('analyzeProgress.tableTitle.version'),
  },
  {
    dataIndex: 'projects',
    key: 'projects',
    title: t('analyzeProgress.tableTitle.projects'),
  },
  {
    dataIndex: 'date',
    key: 'date',
    title: '创建日期',
  },
  {
    dataIndex: 'works',
    key: 'works',
    title: t('analyzeProgress.tableTitle.works'),
  },
  {
    dataIndex: 'members',
    key: 'members',
    title: t('analyzeProgress.tableTitle.members'),
  },
  {
    dataIndex: 'fileSize',
    key: 'fileSize',
    title: t('analyzeProgress.tableTitle.fileSize'),
  },
  {
    dataIndex: 'files',
    key: 'files',
    title: t('analyzeProgress.tableTitle.files'),
  },
  {
    dataIndex: 'id',
    key: 'id',
    title: t('analyzeProgress.tableTitle.id'),
  },
];
const teamColumns = [
  {
    render: (text, record) => {
      return (
        <div className="oac-flex oac-items-center">
          <CheckmarkFilled fontSize="16" style={{ marginRight: '5px' }}></CheckmarkFilled>
          <div>{record.user}</div>
        </div>
      );
    },
    key: 'team',
    title: t('analyzeProgress.tableTitle.team'),
  },
  {
    dataIndex: 'status',
    key: 'status',
    title: t('analyzeProgress.tableTitle.status'),
  },
  {
    dataIndex: 'time',
    key: 'time',
    title: t('analyzeProgress.tableTitle.time'),
  },
  {
    dataIndex: 'id',
    key: 'id',
    title: t('analyzeProgress.tableTitle.id'),
  },
];
const ResultContent: FC<ResultContentProgressProps> = memo((props) => {
  const { backupDataSource, memory, teamDataSource } = props;
  return (
    <div>
      <div>
        <div style={{ fontWeight: '500', fontSize: '16px' }} className="oac-pt-4 oac-pb-1">
          {t('analyzeProgress.analyzeResult.title')}
        </div>
        <div className="oac-pb-2">{t('analyzeProgress.analyzeResult.jiraBackupResult')}</div>
        <Table
          columns={backupColumns}
          dataSource={backupDataSource}
          className="oac-overflow-auto"
          style={{ height: '72px', border: '1px solid #E5E5E5' }}
          locale={{ emptyText: t('common.noData') }}
        ></Table>
      </div>
      <div className="oac-pt-4">
        <div style={{ lineHeight: '22px' }}>
          {t('analyzeProgress.analyzeResult.onesTeamResult')}
        </div>
        {false ? (
          <div className="oac-pb-2 oac-pl-2">
            {t('analyzeProgress.analyzeResult.localStorage', { memory: memory })}
            <CheckmarkFilled style={{ marginLeft: '10px' }} />{' '}
            {t('analyzeProgress.analyzeResult.localStorageSupport')}
          </div>
        ) : (
          <div className="oac-pb-2 oac-pl-2">
            {t('analyzeProgress.analyzeResult.localStorage', { memory: memory })}
            <Warning style={{ color: '#F0A100', marginLeft: '10px' }} />
            <span style={{ color: '#F0A100' }}>
              {t('analyzeProgress.analyzeResult.localStorageNotSupport')}
            </span>
            <a target="_blank" rel="noopener noreferrer">
              {t('analyzeProgress.analyzeResult.localStorageRule')}
              <Launch style={{ marginLeft: '5px' }} />
            </a>
          </div>
        )}
        <Table
          columns={teamColumns}
          dataSource={teamDataSource}
          className="oac-overflow-auto"
          style={{ border: '1px solid #E5E5E5', maxHeight: '141px' }}
          locale={{ emptyText: t('common.emptyData') }}
        ></Table>
      </div>
    </div>
  );
});

export default ResultContent;
