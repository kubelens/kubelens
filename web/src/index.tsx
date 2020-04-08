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
import React from 'react';
import ReactDOM from 'react-dom';
import * as serviceWorker from './serviceWorker';
import { Provider } from 'react-redux';
import { Store } from 'redux';
import configureStore, { IGlobalState, history } from './store';
import App from './app';
import { createClient } from './auth';
import { AuthClient } from './auth/authClient';
import config from './config';

import './index.css';

interface IProps {
  store: Store<IGlobalState>;
}

export const bootstrap = (async () => {
  const cfg = await config();
  let authClient: AuthClient;

  if (cfg.oAuthEnabled) {
    authClient = createClient(cfg, { history })
  }

  /* 
  Create a root component that receives the store via props
  and wraps the App component with Provider, giving props to containers
  */
  const Root: React.SFC<IProps> = props => {
    return (
      <Provider store={props.store}>
        {/* add authClient to props for easy access, e.g. web socket connections, user profile, etc. */}
        <App {...Object.assign({ authClient, history }, props)} />
      </Provider>
    );
  };

  // Generate the store
  const store = configureStore();

  // Render the App
  ReactDOM.render(<Root store={store} />, document.getElementById(
    'root'
  ) as HTMLElement);

  // If you want your app to work offline and load faster, you can change
  // unregister() to register() below. Note this comes with some pitfalls.
  // Learn more about service workers: https://bit.ly/CRA-PWA
  serviceWorker.unregister();
})();