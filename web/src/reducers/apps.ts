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
import { Reducer } from 'redux';
import { AppsActions, AppsActionTypes } from '../actions/apps';
import { AppOverview, App } from '../types';

export interface IAppsState {
  readonly appOverviewRequested: boolean,
  readonly appOverviewError?: Error,
  readonly appOverview?: AppOverview,
  readonly appsRequested: boolean,
  readonly appsError?: Error,
  readonly apps?: App[],
  readonly selectedAppName?: string,
  readonly filteredApps?: App[]
}

const INITIAL_STATE: IAppsState = {
  appOverviewRequested: false,
  appOverviewError: undefined,
  appOverview: undefined,
  appsRequested: false,
  appsError: undefined,
  apps: undefined,
  selectedAppName: undefined,
  filteredApps: undefined
};

export const appsReducer: Reducer<IAppsState, AppsActions> = (
  state = INITIAL_STATE,
  action
) => {
  switch (action.type) {
    case AppsActionTypes.GET_APP_OVERVIEW: {
      return {
        ...state,
        appOverview: action.appOverview,
        appOverviewError: action.appOverviewError,
        appOverviewRequested: action.appOverviewRequested
      };
    }

    case AppsActionTypes.GET_APPS: {
      return {
        ...state,
        apps: action.apps,
        filteredApps: action.apps,
        appsRequested: action.appsRequested,
        appsError: action.appsError
      }
    }

    case AppsActionTypes.SET_SELECTED_APP_NAME: {
      return {
        ...state,
        selectedAppName: action.selectedAppName
      }
    }

    case AppsActionTypes.FILTER_APPS: {
      return {
        ...state,
        filteredApps: action.filteredApps
      }
    }

    case AppsActionTypes.CLEAR_ERRORS: {
      return {
        ...state,
        appsError: undefined
      }
    }

    default:
      return state;
  }
};
