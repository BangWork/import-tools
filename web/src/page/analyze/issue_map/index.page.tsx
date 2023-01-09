import { useRef, useState, useEffect } from 'react';
import { Button, Alert, Table } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSize } from 'ahooks';
import { map } from 'lodash-es';
import { getResultApi } from '@/api';

import useTableBusiness from './use_table_business';

const IssueMapPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const boxRef = useRef(null);
  const boxSize = useSize(boxRef);
  const location = useLocation();
  const [needDisk, setNeedDisk] = useState(false);
  const { loading, columns, select, jiraList } = useTableBusiness();

  useEffect(() => {
    getResultApi().then((res) => {
      setNeedDisk(res.body.resolve_result.disk_set_pop_ups);
    });
  }, []);

  const handleBack = () => {
    navigate('/page/analyze/import_project', { replace: true, state: location.state });
  };

  useEffect(() => {
    if (!location?.state) {
      handleBack();
    }
  }, [location]);

  const handleNext = () => {
    const finishSelect = map(jiraList, (item) => ({
      id: item.third_issue_type_id,
      type: select[item.third_issue_type_id] || item.ones_detail_type,
    }));

    const targetUrl = needDisk ? '/page/analyze/share_disk' : '/page/import_pack/init_password';
    navigate(targetUrl, {
      replace: true,
      state: {
        ...(location?.state || {}),
        issue_type_map: finishSelect,
      },
    });
  };

  return (
    <div className="h-full w-full">
      <div className="flex justify-between">
        <h2>{t('issueMap.title')}</h2>
        <div>
          <Button onClick={handleBack}>{t('common.cancel')}</Button>
          <Button className="ml-4" type="primary" onClick={handleNext}>
            {t('common.nextStep')}
          </Button>
        </div>
      </div>
      <Alert
        className="my-4"
        message={
          <div className="p-2">
            <div>{t('issueMap.tip.message1')}</div>
            <div>{t('issueMap.tip.message2')}</div>
            <div style={{ color: '#FF4D4F' }}>{t('issueMap.tip.message3')}</div>
          </div>
        }
        showIcon
        type="info"
      />
      <div className="flex h-4/6 justify-center">
        <div ref={boxRef} className="flex h-full w-2/3 justify-center">
          <Table
            columns={columns}
            loading={loading}
            dataSource={jiraList}
            pagination={false}
            scroll={{ y: boxSize?.height || 0 }}
          />
        </div>
      </div>
    </div>
  );
};

export default IssueMapPage;
