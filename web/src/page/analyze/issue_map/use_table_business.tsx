import { useState, useEffect } from 'react';
import { Select, Typography, Tooltip } from 'antd';
import { useTranslation } from 'react-i18next';
import { useLocation } from 'react-router-dom';
import type { ColumnsType } from 'antd/es/table';
import { map } from 'lodash-es';

import { getIssuesApi } from '@/api';
import type { OnesIssueType, JiraIssueType } from '@/api';

const selectedSet = new Set();

const useTableBusiness = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const [loading, setLoading] = useState(true);
  const [select, setSelect] = useState({});
  const [data, setData] = useState<{ jira_list: JiraIssueType[]; ones_list: OnesIssueType[] }>({
    jira_list: [],
    ones_list: [],
  });
  const negativeOneAndZero = new Set([0, -1]);
  const negativeOne = -1;
  useEffect(() => {
    const projectIds = location?.state?.projects || [];
    getIssuesApi(projectIds).then((res) => {
      setData(res.body.issue_types);
      if (res.body.issue_types.jira_list) {
        res.body.issue_types.jira_list.forEach((item) => {
          if (!negativeOneAndZero.has(item.ones_detail_type)) {
            selectedSet.add(item.ones_detail_type);
          }
        });
      }

      if (res.body.issue_type_map) {
        const temporaryObj = {};
        res.body.issue_type_map.forEach((item) => {
          if (item.type !== negativeOne) {
            temporaryObj[item.id] = item.type;
            item.type && selectedSet.add(item.type);
          }
        });
        setSelect(temporaryObj);
      }
      setLoading(false);
    });

    return () => {
      selectedSet.clear();
    };
  }, []);

  const handleSelect = (record) => (v) => {
    const { third_issue_type_id } = record;
    const preValue = select[third_issue_type_id];
    selectedSet.delete(preValue);

    // when option is diy,don’t add
    if (!negativeOneAndZero.has(v)) {
      selectedSet.add(v);
    }

    setSelect({
      ...select,
      [third_issue_type_id]: v,
    });
  };

  // The selected options need disabled, excluding what has a value of 0
  const options = map(data.ones_list, (item) => ({
    label: item.name,
    value: item.detail_type,
    disabled: selectedSet.has(item.detail_type),
  }));

  const columns: ColumnsType<JiraIssueType> = [
    {
      title: t('issueMap.table.columns.jira'),
      width: 100,
      dataIndex: 'third_issue_type_name',
      fixed: 'left',
      render: (text) => <Typography.Text ellipsis={{ tooltip: text }}>{text}</Typography.Text>,
    },
    {
      title: t('issueMap.table.columns.issueID'),
      width: 100,
      dataIndex: 'third_issue_type_id',
      fixed: 'left',
    },
    {
      title: t('issueMap.table.columns.onesIssue'),
      fixed: 'right',
      width: 100,
      render: (text, record) => (
        <Tooltip
          title={record.ones_detail_type !== negativeOne ? t('issueMap.table.disabledTip') : ''}
        >
          <Select
            value={
              record.ones_detail_type !== negativeOne
                ? record.ones_detail_type
                : select[record.third_issue_type_id]
            }
            disabled={record.ones_detail_type !== negativeOne}
            placeholder={t('issueMap.table.placeholder')}
            className="w-full"
            onSelect={handleSelect(record)}
            options={options}
          />
        </Tooltip>
      ),
    },
  ];

  return {
    loading,
    columns,
    select,
    jiraList: data.jira_list,
  };
};

export default useTableBusiness;
