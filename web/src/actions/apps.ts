/*
MIT License

Copyright (c) 2020 The KubeLens Authors

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
import { AppOverview, App, Apps } from "../types";
import adapter from './adapter';
import _ from 'lodash';
import { ErrorActionTypes } from './error';
import { LoadingActionTypes } from './loading';

/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type AppsActions = IGetAppOverview | IGetApps | ISetSelectedAppName | IFilterApps;

// Create Action Constants
export enum AppsActionTypes {
  GET_APP_OVERVIEW = 'GET_APP_OVERVIEW',
  GET_APPS = 'GET_APPS',
  SET_SELECTED_APP_NAME = 'SET_SELECTED_APP_NAME',
  FILTER_APPS = 'FILTER_APPS'
}

// IGetAppOverview interface .
export interface IGetAppOverview {
  type: AppsActionTypes.GET_APP_OVERVIEW,
  appOverview?: AppOverview
}

/* Get App Overviews
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getAppOverview: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetAppOverview>
> = (appname: string, labelSelector: string, cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: true,
      });

      const response = await adapter.get(`/apps/${appname}?labelSelector=${encodeURIComponent(labelSelector)}&detailed=true`, cluster, jwt);

      dispatch({
        type: AppsActionTypes.GET_APP_OVERVIEW,
        appOverview: response.data
      });

      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: false,
      });
    } catch (err) {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: false,
      });

      dispatch({
        type: ErrorActionTypes.OPEN_API_ERROR_MODAL,
        status: err.response ? err.response.status : 500,
        statusText: err.response ? err.response.statusText : 'Internal Server Error',
        message: err.response ? err.response.data : err
      });
    }
  };
};

// IGetApps interface .
export interface IGetApps {
  type: AppsActionTypes.GET_APPS,
  apps?: App[],
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
        type: LoadingActionTypes.LOADING,
        loading: true,
      });

      const response = await adapter.get('/apps', cluster, jwt);
      const data = response.data as Apps;
      const apps = _.orderBy(data, 'name');

      dispatch({
        type: AppsActionTypes.GET_APPS,
        apps: apps,
        filteredApps: apps
      });
      
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: false,
      });
    } catch (err) {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: false,
      });

      dispatch({
        type: ErrorActionTypes.OPEN_API_ERROR_MODAL,
        status: err.response ? err.response.status : 500,
        statusText: err.response ? err.response.statusText : 'Internal Server Error',
        message: err.response ? err.response.data : err
      });
    }
  };
};

// ISetSelectedAppName interface .
export interface ISetSelectedAppName {
  type: AppsActionTypes.SET_SELECTED_APP_NAME,
  selectedAppName: string
}

/* <Promise<Return Type>, State Interface, Type of Param, Type of Action> */
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

/* <Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const filterApps: ActionCreator<
  ThunkAction<Promise<any>, IAppsState, null, IGetApps>
> = (value: string, apps: App[]) => {
  return async (dispatch: Dispatch) => {

    const filtered = _.filter(apps, (svc: App) => {
      return (svc.name.includes(value)) 
        || svc.namespace.includes(value)
        || svc.labelSelector.includes(value)
        || svc.kind.includes(value);
    });

    dispatch({
      type: AppsActionTypes.FILTER_APPS,
      filteredApps: filtered
    });
  };
};
