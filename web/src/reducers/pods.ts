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
import { Reducer } from 'redux';
import { PodDetail } from '../types';
import { PodActions, PodActionTypes } from '../actions/pods';

export interface IPodState {
  readonly podRequested: boolean,
  readonly pod?: PodDetail,
  readonly podError?: Error,
  readonly selectedPodName?: string,
  readonly selectedContainerName?: string
}

const INITIAL_STATE: IPodState = {
  podRequested: false,
  podError: undefined,
  pod: undefined,
  selectedPodName: undefined,
  selectedContainerName: undefined
};

export const podsReducer: Reducer<IPodState, PodActions> = (
  state = INITIAL_STATE,
  action
) => {
  switch (action.type) {

    case PodActionTypes.GET_POD: {
      return {
        ...state,
        pod: action.pod,
        podRequested: action.podRequested
      }
    }

    case PodActionTypes.CLEAR_POD: {
      return {
        ...state,
        pod: action.pod
      }
    }

    case PodActionTypes.SET_SELECTED_POD_NAME: {
      return {
        ...state,
        selectedPodName: action.selectedPodName
      }
    }

    case PodActionTypes.SET_SELECTED_CONTAINER_NAME: {
      return {
        ...state,
        selectedContainerName: action.selectedContainerName
      }
    }

    case PodActionTypes.CLEAR_ERRORS: {
      return {
        ...state,
        podError: undefined
      }
    }

    default:
      return state;
  }
};