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
import { WebAuth } from "auth0-js";
import {
  oAuthConfig,
  RenewTokenResult
} from "../clients/auth-implicit-client";
import { AuthError } from '../AuthError';

let oauth: WebAuth;

function createClient(config: oAuthConfig) {
  return new WebAuth({
    domain: config.domain,
    clientID: config.clientId
  });
}

export function login(oauth: WebAuth, config: oAuthConfig): void {
  const url = `${window.location.pathname}${window.location.search}`;

  oauth.authorize({
    connection: config.connection,
    responseType: config.responseType,
    redirectUri: config.redirectUri,
    audience: config.audience,
    scope: config.scope,
    state: btoa(url)
  })
}

export function ensureAuthed(config: oAuthConfig): Promise<RenewTokenResult> {

  if (!oauth) {
    oauth = createClient(config);
  }

  return new Promise((resolve, reject) => {
    // try to find it upon redirecting back to app
    oauth.parseHash((err: any, result: any) => {
      if (err) {
        return reject(new AuthError(err.error, err.errorDescription));
      }

      if (!result) {
        login(oauth, config);
        return reject(new AuthError(err.error, err.errorDescription));
      }

      if (result.state) {
        config.history.replace(atob(result.state));
      }

      return resolve({
        accessToken: result.accessToken,
        identityToken: result.idToken,
        identity: result.idTokenPayload
      });
    });
  });
}

export function reAuth(config: oAuthConfig): Promise<RenewTokenResult> {

  if (!oauth) {
    oauth = createClient(config);
  }

  return new Promise((resolve, reject) => {
    const configuration = { ...config };

    oauth.checkSession(configuration, (error: any, result: any) => {
      if (error) {
        if (error.error === "login_required") {
          login(oauth, configuration);
        }
        reject(new AuthError(error.error, error.errorDescription));
      } else {
        resolve({
          accessToken: result.accessToken,
          identityToken: result.idToken
        });
      }
    });
  });
};