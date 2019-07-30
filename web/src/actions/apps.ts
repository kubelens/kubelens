/*
MIT License

Copyright (c) 2019 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
import { ActionCreator, Dispatch } from 'redux';
import { ThunkAction } from 'redux-thunk';
import { IAppsState } from '../reducers/apps';
import { AppOverview, App, Service } from "../types";
import adapter from './adapter';
import _ from 'lodash';

/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type AppsActions = IGetAppOverview | IGetApps | ISetSelectedAppName | IFilterApps | IClearAppsErrors;

// Create Action Constants
export enum AppsActionTypes {
  GET_APP_OVERVIEW = 'GET_APP_OVERVIEW',
  GET_APPS = 'GET_APPS',
  CLEAR_ERRORS = 'CLEAR_ERRORS',
  SET_SELECTED_APP_NAME = 'SET_SELECTED_APP_NAME',
  FILTER_APPS = 'FILTER_APPS'
}

// IGetAppOverview interface .
export interface IGetAppOverview {
  type: AppsActionTypes.GET_APP_OVERVIEW,
  appOverview?: AppOverview,
  appOverviewRequested: boolean,
  appOverviewError?: Error
}

/* Get App Overviews
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getAppOverview: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetAppOverview>
> = (appname: string, labelKey: string, cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: AppsActionTypes.GET_APP_OVERVIEW,
        appOverviewRequested: true,
      });

      const response = await adapter.get(`/apps/${appname}?labelKey=${labelKey}&detailed=true`, cluster, jwt);

      dispatch({
        type: AppsActionTypes.GET_APP_OVERVIEW,
        appOverview: response.data,
        appOverviewRequested: false,
        appOverviewError: null
      });
    } catch (err) {
      dispatch({
        type: AppsActionTypes.GET_APP_OVERVIEW,
        appOverview: null,
        appOverviewRequested: false,
        appOverviewError: err
      });
    }
  };
};

// IGetApps interface .
export interface IGetApps {
  type: AppsActionTypes.GET_APPS,
  apps?: App[],
  appsRequested: boolean,
  appsError?: Error,
  filteredApps?: App[]
}

/* Get Apps
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getApps: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetApps>
> = (cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: AppsActionTypes.GET_APPS,
        appsRequested: true
      });

      const response = await adapter.get('/services?detailed=false', cluster, jwt);
      const data = response.data as Service[];

      const grp = _.groupBy(data, 'appName.value');
      let apps = new Array<App>();

      _.forEach(grp, (svc, key) => {
        let overview: App = {
          name: key,
          labelKey: '',
          namespaces: [],
          deployerLink: undefined
        };

        let i = 0;
        _.forEach(svc, detail => {
          if (i === 0) {
            overview.name = detail.appName.value || '';
            overview.labelKey = detail.appName.labelKey || '';

            if (detail.deployerLink) {
              overview.deployerLink = detail.deployerLink;
            }
          }
          overview.namespaces.push(detail.namespace || '');
          i++;
        })
        overview.namespaces = _.uniq(overview.namespaces);
        apps.push(overview);
      });

      dispatch({
        type: AppsActionTypes.GET_APPS,
        apps: apps,
        filteredApps: apps,
        appsRequested: false,
        appsError: null
      });
    } catch (err) {
      dispatch({
        type: AppsActionTypes.GET_APPS,
        apps: null,
        appsRequested: false,
        appsError: err
      });
    }
  };
};

// IClearAppsErrors interface .
export interface IClearAppsErrors {
  type: AppsActionTypes.CLEAR_ERRORS,
  appsError?: Error,
  appOverviewError?: Error
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const clearAppsErrors: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetApps>
> = () => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: AppsActionTypes.GET_APP_OVERVIEW,
      appOverviewError: null,
      appsError: null
    });
  };
};



// ISetSelectedAppName interface .
export interface ISetSelectedAppName {
  type: AppsActionTypes.SET_SELECTED_APP_NAME,
  selectedAppName: string
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const setSelectedAppName: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetApps>
> = (selectedAppName: string) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: AppsActionTypes.SET_SELECTED_APP_NAME,
      selectedAppName: selectedAppName
    });
  };
};

// IFilterApps interface .
export interface IFilterApps {
  type: AppsActionTypes.FILTER_APPS,
  filteredApps: App[]
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const filterApps: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetApps>
> = (value: string, apps: App[]) => {
  return async (dispatch: Dispatch) => {

    const filtered = _.filter(apps, (svc: App) => {
      return (svc.name.includes(value))
        || _.join(svc.namespaces, ",").includes(value);
    });

    dispatch({
      type: AppsActionTypes.FILTER_APPS,
      filteredApps: filtered
    });
  };
};
