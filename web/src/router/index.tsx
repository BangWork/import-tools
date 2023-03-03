import { FC, memo } from 'react';
import type { ComponentType } from 'react';
import { Route, Routes, Navigate } from 'react-router-dom';
import { map } from 'lodash-es';
import loadable from '@loadable/component';

import type { RedirectRoutesType, AutoRouterProps } from './type';

const pages = import.meta.glob('/src/**/*.page.*', {
  import: 'default',
});

/** get routes of auto analyze */
export const getRoutes = () => {

  return map(Object.keys(pages), (key) => {
    // route path transform ./src/page/**/*.page.xxx => /page/**
    const pathArr = key.slice(1).split('/').slice(0, -1);
    pathArr.shift();
    return {
      path: `/${pathArr.join('/')}`,
      loader: pages[key],
    };
  });
};

const getRouteList = (Loading?: ComponentType<any>) => {
  const routeList = getRoutes();

  const RouteList = map(routeList, (v) => {
    const Components = loadable<any>(v.loader);
    return (
      <Route
        key={v.path}
        path={v.path}
        element={<Components fallback={<Loading delay={500} />} />}
      />
    );
  });
  return RouteList;
};

function getRedirect(redirectRoutes: RedirectRoutesType[]) {
  return map(redirectRoutes, (v) => {
    return <Route key={v.from + v.to} path={v.from} element={<Navigate to={v.to} />} />;
  });
}

const AutoRouter: FC<AutoRouterProps> = ({ redirectRoutes, children, Loading, NoMatch }) => {
  return (
    <Routes>
      {children}
      {getRouteList(Loading)}
      {redirectRoutes && getRedirect(redirectRoutes)}
      <Route path="*" element={<NoMatch />} />
    </Routes>
  );
};

export default memo(AutoRouter);
