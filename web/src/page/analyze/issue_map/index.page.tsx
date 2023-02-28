import { useRef, useEffect } from 'react';
import { Button, Alert, Table } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSize } from 'ahooks';
import { map } from 'lodash-es';
import { saveIssuesApi } from '@/api';

import useTableBusiness from './use_table_business';

const IssueMapPage = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const boxRef = useRef(null);
  const boxSize = useSize(boxRef);
  const location = useLocation();
  const { loading, columns, select, jiraList} = useTableBusiness();

  const handleBack = () => {
    navigate('/page/analyze/import_project', { replace: true, state: location.state });
  };

  const saveData = () => {
    const selectedArr = []
    Object.keys(select).forEach((key) => {
      selectedArr.push({ id: key, type: select[key] })
    })
    saveIssuesApi(selectedArr)
  };

  useEffect(() => {
    if (!location?.state) {
      handleBack();
    }
  }, [location]);

  const handleNext = () => {
    saveData()
    const finishSelect = map(jiraList, (item) => ({
      id: item.third_issue_type_id,
      type: select[item.third_issue_type_id] || item.ones_detail_type,
    }));

    navigate('/page/import_pack/init_password', {
      replace: true,
      state: {
        ...(location?.state || {}),
        issue_type_map: finishSelect,
      },
    });
  };

  const showDataList = map(jiraList, (item) => ({
   ...item,key:item.third_issue_type_id
  }))

  return (
    <div className="h-full w-full flex flex-col items-center">
      <div className="flex justify-between w-2/3">
        <h2>{t('issueMap.title')}</h2>
        <div>
          <Button onClick={handleBack}>{t('common.cancel')}</Button>
          <Button className="ml-4" type="primary" onClick={handleNext}>
            {t('common.nextStep')}
          </Button>
        </div>
      </div>
      <Alert
        className="my-4 w-2/3"
        message={
          <ol className="px-8">
            <li>{t('issueMap.tip.message1')}</li>
            <li>{t('issueMap.tip.message2')}</li>
            <li>{t('issueMap.tip.message3')}</li>
            <li>{t('issueMap.tip.message4')}</li>
          </ol>
        }
        showIcon
        type="info"
      />
      <div className="h-4/6 w-2/3">
        <div ref={boxRef} className="h-full">
          <Table
            columns={columns}
            loading={loading}
            dataSource={showDataList}
            pagination={false}
            scroll={{ y: boxSize?.height || 0 }}
          />
        </div>
      </div>
    </div>
  );
};

export default IssueMapPage;
