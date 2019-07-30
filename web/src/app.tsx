/*
MIT License

Copyright (c) 2019 The KubeLens Authors

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
import { Router } from 'react-router-dom';
import { ConnectedRouter } from 'connected-react-router';
import AsyncLoadingContainer from './components/async-container';
import AuthenticationWrapper from './components/authentication-wrapper';
import { AuthImplicitClient } from 'auth';
import ClusterIndex from './areas';

export interface IAppProps {
  authClient: AuthImplicitClient
  history: any
}

export default class App extends React.Component<IAppProps, any> {
  render() {
    return (
      <Router history={this.props.history}>
        <AsyncLoadingContainer >
          <ConnectedRouter history={this.props.history}>
            <AuthenticationWrapper authClient={this.props.authClient}>
              <ClusterIndex {...this.props} />
            </AuthenticationWrapper>
          </ConnectedRouter>
        </AsyncLoadingContainer>
      </Router>
    );
  }
}
