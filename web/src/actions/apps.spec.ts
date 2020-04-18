import moxios from 'moxios';
import sinon from 'sinon';
import { deepEqual, fail } from 'assert';

import adapter from './adapter';
import { AppsActionTypes, getApps, getAppOverview, setSelectedAppName, filterApps } from './apps';

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

describe('getApps should', () => {
  test('succeed', () => {
    const res = [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}];
    const store = mockStore(res);
    const expectedActions = [ 
        LoadingActionTypes.LOADING,
        AppsActionTypes.GET_APPS
    ];
    sinon.stub(adapter, 'get').resolves({data: res});
  
    return store.dispatch(getApps("appname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
        deepEqual(store.getState(), res) 
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
    sinon.stub(adapter, 'get').throws(errResponse);
  
    return store.dispatch(getApps("appname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        fail("should not be successful.");
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
    sinon.stub(adapter, 'get').throws('unexpected');
  
    return store.dispatch(getApps("appname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        fail("should not be successful.");
      })
      .catch(err => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      });
  });
});

describe('getAppOverviews should', () => {
  test('succeed', () => {
    const store = mockStore();
    const expectedActions = [ 
        LoadingActionTypes.LOADING,
        AppsActionTypes.GET_APP_OVERVIEW
    ];
    sinon.stub(adapter, 'get').resolves(successResponse);
  
    return store.dispatch(getAppOverview("appname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
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
  
    return store.dispatch(getAppOverview("appname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
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
  
    return store.dispatch(getAppOverview("appname", encodeURIComponent('{"key":"value"}'), "minikube", ""))
      .then(() => {
        fail("should not succeed.");
      })
      .catch(err => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      });
  });
});

describe('non-api calls actions should', () => {
  test('setSelectedAppName', async () => {
    const store = mockStore();
    const expectedActions = [ 
      AppsActionTypes.SET_SELECTED_APP_NAME
    ];
    return store.dispatch(setSelectedAppName('testapp'))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });

  test('filterApps', async () => {
    const store = mockStore();
    const expectedActions = [ 
      AppsActionTypes.FILTER_APPS
    ];
    return store.dispatch(filterApps('appname', [{name: 'appname'}]))
      .then(() => {
        deepEqual(store.getActions().map(action => action.type), expectedActions);
      })
      .catch(err => {
        fail(err);
      })
  });
});