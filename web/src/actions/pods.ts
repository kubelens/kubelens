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
import { PodDetail } from '../types';
import { IPodState } from '../reducers/pods';
import adapter from './adapter';
import { ErrorActionTypes } from './error';
import { LoadingActionTypes } from './loading';

/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type PodActions = IGetPod | ISetSelectedPodName | ISetSelectedContainerName | IClearPod;

// Create Action Constants
export enum PodActionTypes {
  GET_POD = 'GET_POD',
  CLEAR_POD = 'CLEAR_POD',
  SET_SELECTED_POD_NAME = 'SET_SELECTED_POD_NAME',
  SET_SELECTED_CONTAINER_NAME = 'SET_SELECTED_CONTAINER_NAME'
}

export interface IGetPod {
  type: PodActionTypes.GET_POD,
  pod?: PodDetail
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const getPod: ActionCreator<
  ThunkAction<Promise<any>, IPodState, null, IGetPod>
> = (podname: string, queryString: string, cluster: string, jwt: string) => {
  return async (dispatch: Dispatch) => {
    try {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: true,
      });

      const response = await adapter.get(`/pods/${podname}${queryString}`, cluster, jwt);

      dispatch({
        type: PodActionTypes.GET_POD,
        pod: response.data
      });
    } catch (err) {
      dispatch({
        type: LoadingActionTypes.LOADING,
        loading: true,
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

export interface ISetSelectedPodName {
  type: PodActionTypes.SET_SELECTED_POD_NAME,
  selectedPodName: string
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const setSelectedPodNameRequest: ActionCreator<
  ThunkAction<Promise<any>, IPodState, null, ISetSelectedPodName>
> = (selectedPodName: string) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: PodActionTypes.SET_SELECTED_POD_NAME,
      selectedPodName: selectedPodName
    });
  };
};

export interface ISetSelectedContainerName {
  type: PodActionTypes.SET_SELECTED_CONTAINER_NAME,
  selectedContainerName: string
}

/* Set container name
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const setSelectedContainerName: ActionCreator<
  ThunkAction<Promise<any>, IPodState, null, ISetSelectedContainerName>
> = (selectedContainerName: string) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: PodActionTypes.SET_SELECTED_CONTAINER_NAME,
      selectedContainerName: selectedContainerName
    });
  };
};

export interface IClearPod {
  type: PodActionTypes.CLEAR_POD,
  pod?: PodDetail
}

/* Set container name
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const clearPod: ActionCreator<
  ThunkAction<Promise<any>, IPodState, null, IClearPod>
> = () => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: PodActionTypes.CLEAR_POD,
      pod: null
    });
  };
};