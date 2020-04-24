import moxios from 'moxios';
import sinon from 'sinon';
import { deepEqual, fail } from 'assert';

import adapter from './adapter';
import { LogsActionTypes, getLogs, toggleLogStream, clearLogsErrors } from './logs';

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
        LogsActionTypes.GET_LOGS,
        LoadingActionTypes.LOADING
    ];
    sinon.stub(adapter, 'get').resolves(successResponse);
  
    return store.dispatch(getLogs("podname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
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
      LogsActionTypes.GET_LOGS,
      ErrorActionTypes.OPEN_API_ERROR_MODAL
    ];
    sinon.stub(adapter, 'get').rejects(errResponse);
  
    return store.dispatch(getLogs("podname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
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
      LogsActionTypes.GET_LOGS,
      ErrorActionTypes.OPEN_API_ERROR_MODAL
    ];
    sinon.stub(adapter, 'get').rejects('unexpected');
  
    return store.dispatch(getLogs("podname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        fail("should not succeed.");
      })
      .catch(err => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      });
  });
});

describe('non-api calls actions should', () => {
  test('toggleLogStream', async () => {
    const store = mockStore();
    const expectedActions = [ 
      LogsActionTypes.TOGGLE_LOG_STREAM,
      LoadingActionTypes.LOADING
    ];
    return store.dispatch(toggleLogStream('testapp'))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });

  test('clearLogsErrors', async () => {
    const store = mockStore();
    const expectedActions = [ 
      LogsActionTypes.CLEAR_ERRORS
    ];
    return store.dispatch(clearLogsErrors())
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});