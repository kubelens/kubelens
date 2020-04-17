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

/*  Imports from Redux:
 applyMiddleware: Applies middleware to the dispatch method of the Redux store
 combineReducers: Merges reducers into one
 createStore: Creates a Redux store that holds the state tree
 Store: The TS Type used for the store, or state tree
 */
import { applyMiddleware, combineReducers, createStore, Store } from 'redux';
import { routerMiddleware, connectRouter } from 'connected-react-router';
import { createBrowserHistory } from 'history';
/*  Thunk
Redux Thunk middleware allows you to write action creators that return a function instead of an action. The thunk can be used to delay the dispatch of an action, or to dispatch only if a certain condition is met. The inner function receives the store methods dispatch and getState as parameters.
*/
import thunk from 'redux-thunk';

// Import reducers and state type
import * as apps from './reducers/apps';
import * as logs from './reducers/logs';
import * as pods from './reducers/pods';
import * as auth from './reducers/auth';
import * as cluster from './reducers/cluster';
import * as errors from './reducers/error';
import * as loading from './reducers/loading';

export const history = createBrowserHistory();

// Create an interface for the application state
export interface IGlobalState {
  appsState: apps.IAppsState,
  logsState: logs.ILogsState,
  podsState: pods.IPodState,
  authState: auth.IAuthState,
  clustersState: cluster.IClustersState,
  errorState: errors.IErrorState,
  loadingState: loading.ILoadingState,
  router: any
}

// Create the root reducer
const rootReducer = combineReducers<IGlobalState>({
  appsState: apps.reducer,
  logsState: logs.reducer,
  podsState: pods.reducer,
  authState: auth.reducer,
  clustersState: cluster.reducer,
  errorState: errors.reducer,
  loadingState: loading.reducer,
  router: connectRouter(history),
});

const middleware = [
  thunk,
  routerMiddleware(history)
];

// Create a configure store function of type `IAppState`
export default function configureStore(): Store<IGlobalState, any> {
  return createStore(
    rootReducer,
    undefined,
    applyMiddleware(...middleware)
  );
}
