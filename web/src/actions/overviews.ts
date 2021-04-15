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
import { IOverviewsState } from '../reducers/overviews';
import { Overview } from "../types";
import adapter from './adapter';
import _ from 'lodash';
import { ErrorActionTypes } from './error';
import { LoadingActionTypes } from './loading';

/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type OverviewActions = IGetOverview | IGetOverviews | ISetSelectedAppName | IFilterOverviews;

// Create Action Constants
export enum OverviewActionTypes {
  GET_OVERVIEW = 'GET_OVERVIEW',
  GET_OVERVIEWS = 'GET_OVERVIEWS',
  SET_SELECTED_APP_NAME = 'SET_SELECTED_APP_NAME',
  FILTER_OVERVIEWS = 'FILTER_APPS'
}

// IGetOverview interface .
export interface IGetOverview {
  type: OverviewActionTypes.GET_OVERVIEW,
  overview?: Overview
}

/* Get Overviews
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getOverview: ActionCreator<
  ThunkAction<Promise<any>, IOverviewsState, null, IGetOverview>
> = (linkedName: string, namespace: string, cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: true,
      });

      const response = await adapter.get(`overviews/${linkedName}?namespace=${namespace}`, cluster, jwt);

      dispatch({
        type: OverviewActionTypes.GET_OVERVIEW,
        overview: response.data
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

// IGetOverviews interface .
export interface IGetOverviews {
  type: OverviewActionTypes.GET_OVERVIEWS,
  overviews?: Overview[],
  filteredOverviews?: Overview[]
}

/* Get Overviews
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getOverviews: ActionCreator<
  ThunkAction<Promise<any>, IOverviewsState, null, IGetOverviews>
> = (cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: true,
      });

      const response = await adapter.get('overviews', cluster, jwt);
      const data = response.data as Overview[];
      const overviews = _.orderBy(data, 'linkedName');

      dispatch({
        type: OverviewActionTypes.GET_OVERVIEWS,
        overviews: overviews,
        filteredOverviews: overviews
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
  type: OverviewActionTypes.SET_SELECTED_APP_NAME,
  selectedAppName: string
}

/* <Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const setSelectedAppName: ActionCreator<
  ThunkAction<Promise<any>, IOverviewsState, null, ISetSelectedAppName>
> = (selectedAppName: string) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: OverviewActionTypes.SET_SELECTED_APP_NAME,
      selectedAppName: selectedAppName
    });
  };
};

// IFilterOverviews interface .
export interface IFilterOverviews {
  type: OverviewActionTypes.FILTER_OVERVIEWS,
  filteredOverviews: Overview[]
}

/* <Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const filterOverviews: ActionCreator<
  ThunkAction<Promise<any>, IOverviewsState, null, IFilterOverviews>
> = (value: string, overviews: Overview[]) => {
  return async (dispatch: Dispatch) => {

    const filtered = _.filter(overviews, (ov: Overview) => {
      return (ov.linkedName.includes(value)) 
        || ov.namespace.includes(value)
    });

    dispatch({
      type: OverviewActionTypes.FILTER_OVERVIEWS,
      filteredOverviews: filtered
    });
  };
};
