import { pureRequest, Request } from '@/request';

import type {
  ProjectType,
  OnesIssueType,
  JiraIssueType,
  AnalyzeInfoType,
  ResultType,
} from './type';

export const submitEnvironmentApi = ({ localHome, backupName, url, email, password }): Promise<{
  body: {key:string}
  code: number;
}> =>
  pureRequest.post('/resolve/start', {
    local_home: localHome,
    backup_name: backupName,
    url,
    email,
    password,
  });

export const getProjectsApi = (): Promise<{ body: {projects:ProjectType[],cache:any[] }; code: number; }> =>
  Request.get('/project_list');


export const saveProjectsApi = ( key:string, project_ids:string[]) =>
  pureRequest.post('/project_list/save', {
    key,
    project_ids,
  })

export const getIssuesApi = (
  project_ids: string[]
): Promise<{
  body: { issue_types: { jira_list: JiraIssueType[]; ones_list: OnesIssueType[] } ,issue_type_map:any[]};
  code: number;
}> => Request.post('/issue_type_list', { project_ids });

export const saveIssuesApi = (key: string,issue_type_map:any[]) =>
  pureRequest.post('/issue_type_list/save',{key:key,issue_type_map:issue_type_map})

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
