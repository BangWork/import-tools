export enum AnalyzeStatusEnum {
  none, // Jira package was never analyze
  doing,
  done,
  fail,
}

export enum ImportStatusEnum {
  none,
  doing,
  done,
  pause,
  cancel,
}

export interface ProjectType {
  id: string;
  name: string;
  key: string;
  orignial_key: string;
  assign: string;
  category: string;
  type: string;
  issue_count:number;
}

export interface JiraIssueType {
  third_issue_type_id: string;
  third_issue_type_name: string;
  ones_detail_type: number; // 绑定的ones detail type
}

export interface OnesIssueType {
  uuid: string;
  detail_type: number;
  name: string;
}

export interface AnalyzeInfoType {
  multi_team: boolean;
  team_name: string;
  backup_name: string;
  org_name: string;
  start_time: number;
  status: AnalyzeStatusEnum;
  expected_time: number; // unit second
  spent_time: number; // unit second
}

export interface ResultType {
  resolve_result: {
    jira_version: string;
    project_count: number;
    issue_count: number;
    member_count: number;
    attachment_size: number; // byte
    attachment_count: number;

    jira_server_id: string; // Judging whether it is the same source
    disk_set_pop_ups: boolean; // Whether to display the "Shared Disk" pop-up window
  };
  import_history: {
    team_uuid: string;
    team_name: string;
    import_list: {
      import_time: number;
      jira_version: string;
      jira_server_id: string;
    }[];
  }[];
}

export interface ImportInfoType {
  team_name: string;
  backup_name: string;
  start_time: number;
  backup_time: number;
  status: ImportStatusEnum;
  expected_time: number; // second
  spent_time: number; // second
}

// Dynamic data, when the value of count is -1, it means: calculation
export interface ImportScopeType {
  jira_version: string;
  project_count: number;
  issue_count: number;
  member_count: number;
  attachment_size: number; // byte
  attachment_count: number;
}
