import { deepEqual, fail } from 'assert';

import { ErrorActionTypes, openErrorModal, closeErrorModal } from './error';

import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";

const middlewares = [ thunk ];
const mockStore = configureMockStore(middlewares);

describe('error actions should', () => {
  test('openErrorModal', async () => {
    const store = mockStore();
    const expectedActions = [ 
      ErrorActionTypes.OPEN_API_ERROR_MODAL
    ];
    return store.dispatch(openErrorModal())
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });

  test('closeErrorModal', async () => {
    const store = mockStore();
    const expectedActions = [ 
      ErrorActionTypes.CLOSE_API_ERROR_MODAL
    ];
    return store.dispatch(closeErrorModal())
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});