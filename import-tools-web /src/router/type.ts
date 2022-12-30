import type { ComponentType, ReactElement } from 'react';

export interface RedirectRoutesType {
  to: string;
  from: string;
}

export interface AutoRouterProps {
  redirectRoutes?: RedirectRoutesType[];
  NoMatch?: ComponentType<any>;
  Loading?: ComponentType<any>;
  children?: ReactElement;
}
