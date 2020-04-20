import { deepEqual, fail } from 'assert';

import { LoadingActionTypes, setLoading } from './loading';

import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";

const middlewares = [ thunk ];
const mockStore = configureMockStore(middlewares);

describe('loading actions should', () => {
  test('setLoading', async () => {
    const store = mockStore();
    const expectedActions = [ 
      LoadingActionTypes.LOADING
    ];
    return store.dispatch(setLoading())
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});