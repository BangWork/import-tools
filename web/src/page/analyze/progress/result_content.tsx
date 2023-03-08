import { memo } from 'react';
import type { FC } from 'react';
import { Table } from '@ones-design/core';
import { CheckmarkFilled } from '@ones-design/icons';
import { t } from 'i18next';
interface columnsType {
  title: string;
  dataIndex?: string;
  render?: (text: any, record: any, index: number) => React.ReactNode;
  key: string;
}
interface ResultContentProgressProps {
  backupColumns: columnsType[];
  backupDataSource: any[];
  memory: string;
}
const ResultContent: FC<ResultContentProgressProps> = memo((props) => {
  const { backupColumns, backupDataSource, memory } = props;
  return (
    <div>
      <div>
        <div style={{ fontWeight: '500', fontSize: '16px' }} className="oac-pb-1">
          {t('title')}
        </div>
        <div className="oac-pb-2">{t('common.back')}</div>
        <Table
          columns={backupColumns}
          dataSource={backupDataSource}
          className="oac-overflow-auto"
          style={{ maxHeight: '100px' }}
        ></Table>
      </div>
      <div className="oac-pt-4">
        <div style={{ fontWeight: '500', fontSize: '16px' }}>{t('title')}</div>
        <div>
          {t('common.back') + t(memory)}
          <CheckmarkFilled />
        </div>
        <Table
          columns={backupColumns}
          dataSource={backupDataSource}
          className="oac-overflow-auto"
          style={{ maxHeight: '100px' }}
        ></Table>
      </div>
    </div>
  );
});

export default ResultContent;
