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
import { Log } from '../types';
import { ILogsState } from '../reducers/logs';
import adapter from './adapter';
import { ErrorActionTypes } from './error';
import { LoadingActionTypes } from './loading';
/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type LogsActions = IGetLogs | IToggleLogStream | IClearLogsErrors;

// Create Action Constants
export enum LogsActionTypes {
  GET_LOGS = 'GET_LOGS',
  GET_LOGS_ERROR = 'GET_LOGS_ERROR',
  CLEAR_ERRORS = 'CLEAR_ERRORS',
  TOGGLE_LOG_STREAM = 'TOGGLE_LOG_STREAM'
}

export interface IGetLogs {
  type: LogsActionTypes.GET_LOGS,
  logs?: Log,
  logsError?: Error
}

/* Get Logs
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getLogs: ActionCreator<
  ThunkAction<Promise<any>, ILogsState, null, IGetLogs>
> = (podname: string, queryString: string, cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: true,
      });

      const response = await adapter.get(`/logs/${podname}${queryString}`, cluster, jwt);

      dispatch({
        type: LogsActionTypes.GET_LOGS,
        logs: response.data,
        logsError: null
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
        type: LogsActionTypes.GET_LOGS,
        logsError: err
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

export interface IToggleLogStream {
  type: LogsActionTypes.TOGGLE_LOG_STREAM,
  logStreamEnabled: boolean
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const toggleLogStream: ActionCreator<
  ThunkAction<Promise<any>, ILogsState, null, IToggleLogStream>
> = (logStreamEnabled: boolean) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: LogsActionTypes.TOGGLE_LOG_STREAM,
      logStreamEnabled: logStreamEnabled
    });
    dispatch({
      type: LoadingActionTypes.LOADING,
      loading: logStreamEnabled,
    });
  };
};


// IClearLogsErrors interface .
export interface IClearLogsErrors {
  type: LogsActionTypes.CLEAR_ERRORS,
  logsError?: Error
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const clearLogsErrors: ActionCreator<
  ThunkAction<Promise<any>, ILogsState, null, IClearLogsErrors>
> = () => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: LogsActionTypes.CLEAR_ERRORS,
      logsError: null
    });
  };
};