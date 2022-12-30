import { Request, pureRequest } from '@/request';
import type { ImportInfoType, ImportScopeType } from './type';

export const importStartApi = ({ password, projectIds, issueTypeMap }) =>
  pureRequest.post('/import/start', {
    password,
    project_ids: projectIds,
    issue_type_map: issueTypeMap,
  });

export const getImportInfoApi = (): Promise<{ body?: ImportInfoType; code: number }> =>
  Request.get('/import/progress');

export const getLogHistoryApi = (): Promise<{ body?: string[]; code: number }> =>
  Request.get(`/import/log`);

export const getLogApi = (lineNumber: number): Promise<{ body?: string[]; code: number }> =>
  Request.get(`/import/log/start_line/${lineNumber}`);

export const cancelImportApi = () => Request.post('/import/stop');

export const pauseImportApi = () => Request.post('/import/pause');

export const continueImportApi = () => Request.post('/import/continue');

export const getImportScopeApi = (): Promise<{ body: ImportScopeType; code: number }> => Request.get('/import/scope');

export const resetImportStatusApi = () => Request.get('/import/reset');
