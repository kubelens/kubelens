import { deepEqual, fail } from 'assert';

import { AuthActionTypes, setIdentityToken } from './auth';

import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";

const middlewares = [ thunk ];
const mockStore = configureMockStore(middlewares);

describe('auth actions should', () => {
  test('setIdentityToken', async () => {
    const store = mockStore();
    const expectedActions = [ 
      AuthActionTypes.SET_IDENTITY_TOKEN
    ];
    return store.dispatch(setIdentityToken('token'))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});