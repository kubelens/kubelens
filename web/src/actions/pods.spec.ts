import moxios from 'moxios';
import sinon from 'sinon';
import { deepEqual, fail } from 'assert';

import adapter from './adapter';
import { PodActionTypes, getPod, setSelectedPodNameRequest, setSelectedContainerName, clearPod } from './pods';

import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";
import { LoadingActionTypes } from './loading';
import { ErrorActionTypes } from './error';

const middlewares = [ thunk ];
const mockStore = configureMockStore(middlewares);
const errResponse = {response: {status: 500}};
const successResponse = [{name:"app1"}];

beforeEach(function () {
  // import and pass your custom axios instance to this method
  moxios.install()
})

afterEach(function () { 
  // import and pass your custom axios instance to this method
  moxios.uninstall()
  sinon.restore()
})

describe('getLogs should', () => {
  test('succeed', () => {
    const store = mockStore();
    const expectedActions = [ 
        LoadingActionTypes.LOADING,
        PodActionTypes.GET_POD,
        LoadingActionTypes.LOADING
    ];
    sinon.stub(adapter, 'get').resolves(successResponse);
  
    return store.dispatch(getPod("podname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      });
  });

  test('catch error', () => {
    const store = mockStore();
    const expectedActions = [ 
      LoadingActionTypes.LOADING,
      LoadingActionTypes.LOADING,
      ErrorActionTypes.OPEN_API_ERROR_MODAL
    ];
    sinon.stub(adapter, 'get').rejects(errResponse);
  
    return store.dispatch(getPod("podname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        fail("should not succeed.");
      })
      .catch(err => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      });
  });

  test('catch error unexpected format', () => {
    const store = mockStore();
    const expectedActions = [ 
      LoadingActionTypes.LOADING,
      LoadingActionTypes.LOADING,
      ErrorActionTypes.OPEN_API_ERROR_MODAL
    ];
    sinon.stub(adapter, 'get').rejects('unexpected');
  
    return store.dispatch(getPod("podname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        fail("should not succeed.");
      })
      .catch(err => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      });
  });
});

describe('non-api calls actions should', () => {
  test('setSelectedPodNameRequest', async () => {
    const store = mockStore();
    const expectedActions = [ 
      PodActionTypes.SET_SELECTED_POD_NAME
    ];
    return store.dispatch(setSelectedPodNameRequest('podname'))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });

  test('setSelectedContainerName', async () => {
    const store = mockStore();
    const expectedActions = [ 
      PodActionTypes.SET_SELECTED_CONTAINER_NAME
    ];
    return store.dispatch(setSelectedContainerName('containername'))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });

  test('clearPod', async () => {
    const store = mockStore();
    const expectedActions = [ 
      PodActionTypes.CLEAR_POD
    ];
    return store.dispatch(clearPod())
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});