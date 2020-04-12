
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
import { IErrorState } from '../reducers/error';

/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type ErrorActions = IOpenAPIErrorModal | ICloseAPIErrorModal;

// Create Action Constants
export enum ErrorActionTypes {
  OPEN_API_ERROR_MODAL = 'OPEN_API_ERROR_MODAL',
  CLOSE_API_ERROR_MODAL = 'CLOSE_API_ERROR_MODAL'
}

export interface IOpenAPIErrorModal {
  type: ErrorActionTypes.OPEN_API_ERROR_MODAL,
  apiOpen: boolean,
  status: number,
  statusText: string,
  message: string
}

/* <Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const openErrorModal: ActionCreator<
  ThunkAction<Promise<any>, IErrorState, null, IOpenAPIErrorModal>
> = (status?: number, statusText?: string, message?: string) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: ErrorActionTypes.OPEN_API_ERROR_MODAL,
      apiOpen: true,
      status: status,
      statusText: statusText,
      message: message
    });
  };
};

export interface ICloseAPIErrorModal {
  type: ErrorActionTypes.CLOSE_API_ERROR_MODAL,
  apiOpen: boolean
}

/* <Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const closeErrorModal: ActionCreator<
  ThunkAction<Promise<any>, IErrorState, null, ICloseAPIErrorModal>
> = (apiOpen: boolean) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: ErrorActionTypes.CLOSE_API_ERROR_MODAL,
      apiOpen: false
    });
  };
};
