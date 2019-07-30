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
import { IClustersState } from '../reducers/cluster';

/* 
Combine the action types with a union (we assume there are more)
example: export type CharacterActions = IGetAllAction | IGetOneAction ... 
*/
export type ClustersActions = ISetSelectedCluster;

// Create Action Constants
export enum ClustersActionTypes {
  SET_SELECTED_CLUSTER = 'SET_SELECTED_CLUSTER'
}

export interface ISetSelectedCluster {
  type: ClustersActionTypes.SET_SELECTED_CLUSTER,
  cluster: string
}

/* Clear errors
<Promise<Return Type>, State Interface, Type of Param, Type of Action> */
export const setSelectedCluster: ActionCreator<
  ThunkAction<Promise<any>, IClustersState, null, ISetSelectedCluster>
> = (cluster: string) => {
  return async (dispatch: Dispatch) => {
    dispatch({
      type: ClustersActionTypes.SET_SELECTED_CLUSTER,
      cluster: cluster
    });
  };
};
