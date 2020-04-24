import { deepEqual, fail } from 'assert';

import { ClustersActionTypes, setSelectedCluster } from './cluster';

import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";

const middlewares = [ thunk ];
const mockStore = configureMockStore(middlewares);

describe('cluster actions should', () => {
  test('setIdentityToken', async () => {
    const store = mockStore();
    const expectedActions = [ 
      ClustersActionTypes.SET_SELECTED_CLUSTER
    ];
    return store.dispatch(setSelectedCluster('clustername'))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});