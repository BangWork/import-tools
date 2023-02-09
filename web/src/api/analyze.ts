import { pureRequest, Request } from '@/request';

import type {
  ProjectType,
  OnesIssueType,
  JiraIssueType,
  AnalyzeInfoType,
  ResultType,
} from './type';

export const submitEnvironmentApi = ({ localHome, backupName, url, email, password }) =>
  pureRequest.post('/resolve/start', {
    local_home: localHome,
    backup_name: backupName,
    url,
    email,
    password,
  });

export const getProjectsApi = (): Promise<{ body: ProjectType[]; code: number }> =>
  Request.get('/project_list');

export const getIssuesApi = (
  project_ids: string[]
): Promise<{
  body: { jira_list: JiraIssueType[]; ones_list: OnesIssueType[] };
  code: number;
}> => Request.post('/issue_type_list', { project_ids });

export const checkPathApi = (path) =>
  pureRequest.post('/check_jira_path_exist', {
    path,
  });

export const getBackUpApi = (path) =>
  pureRequest.post<any, { code: number; body: string[] }>('/jira_backup_list', { path });

export const getAnalyzeProgressInfoApi = (): Promise<{ body: AnalyzeInfoType; code: number }> =>
  Request.get('/resolve/progress');

export const cancelAnalyzeApi = () => pureRequest.post('/resolve/stop');

export const getResultApi = (): Promise<{ body: ResultType; code: number }> =>
  Request.get('/resolve/result');

export const chooseTeamApi = (uuid: string, name: string) =>
  pureRequest.post('/choose_team', { team_uuid: uuid, team_name: name });
