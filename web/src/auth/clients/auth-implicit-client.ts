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
import { ClientBase } from './clientbase';
import decode from 'jwt-decode';
import { Config } from '../../types/index';
import * as OktaStrategy from '../strategies/okta-strategy';
import * as Auth0Strategy from '../strategies/auth0-strategy';
import { AuthError } from '../AuthError';

export interface oAuthConfig {
  clientId: string,
  domain: string,
  responseType: string,
  redirectUri: string,
  audience: string,
  scope: string,
  connection: string,
  history: any
}

export interface RenewTokenResult {
  accessToken: string,
  identityToken?: string,
  identity?: any,
}

export interface UserProfile {
  username: string,
  name: string,
  email: string
}

export class AuthImplicitClient extends ClientBase {
  static createDefault(config: Config, overrides?: Partial<oAuthConfig>): AuthImplicitClient {

    const result = {
      ...{
        redirectUri: window.location.href.toLowerCase(),
        connection: config.oAuthConnection,
        scope: config.oAuthScope,
        responseType: config.oAuthResponseType,
        audience: config.oAuthAudience,
        clientId: config.oAuthClientId,
        domain: config.oAuthJwtIssuer,
        history: undefined
      },
      ...overrides
    }
    return new AuthImplicitClient(result)
  }

  constructor(config: oAuthConfig) {
    super();

    if (!config) {
      throw new AuthError('missing_configuration', 'Missing required oAuth client configuration');
    }

    if (!config.clientId) {
      throw new AuthError('missing_configuration', 'No oAuth client id specified in configuration');
    }

    if (!config.domain) {
      throw new AuthError('missing_configuration', 'No oAuth domain specified in configuration');
    }

    if (!config.responseType) {
      throw new AuthError('missing_configuration', 'No oAuth response type specified in configuration');
    }

    if (!config.redirectUri) {
      throw new AuthError('missing_configuration', 'No oAuth redirect uri specified in configuration');
    }

    if (!config.audience) {
      throw new AuthError('missing_configuration', 'No oAuth audience specified in configuration');
    }

    if (!config.scope) {
      throw new AuthError('missing_configuration', 'No oAuth scope specified in configuration');
    }

    if (!config.connection && config.connection !== '') {
      throw new AuthError('missing_configuration', 'No oAuth connection specified in configuration');
    }

    this.config = config;

    if (this.config.domain.indexOf('okta') > -1) {
      this.oAuthStrategy = OktaStrategy;
    } else {
      this.oAuthStrategy = Auth0Strategy;
    }
  }

  config: oAuthConfig;

  accessToken?: string;

  identityToken?: string;

  oAuthStrategy: any;

  ensureAuthed(): Promise<RenewTokenResult> {
    return new Promise((resolve, reject) => {
      // got what we need
      if (this.isLoggedIn()) {
        let identity;
        try {
          identity = decode(this.identityToken);
        } catch (e) {
          // missing id token... is there a better way to handle this?
          this.ensureAuthed();
        }

        return resolve({
          accessToken: this.accessToken,
          identityToken: this.identityToken,
          identity
        })
      }

      // we have an access token, it's just expired
      if (this.accessToken) {
        return this.oAuthStrategy.reAuth(this.config).then((result: RenewTokenResult) => {
          this.accessToken = result.accessToken;
          this.identityToken = result.identityToken;
          resolve(result);
        });
      }
      // try to find it upon redirecting back to app
      this.oAuthStrategy.ensureAuthed(this.config)
        .then((result: RenewTokenResult) => {
          this.accessToken = result.accessToken;
          this.identityToken = result.identityToken;
          resolve(result);
        })
        .catch((e: AuthError) => {
          reject(e);
        });
      // @ts-ignore
    }).then(this.postEnsureAuthed);
  }

  private postEnsureAuthed: (tokenInfo: RenewTokenResult) => Promise<RenewTokenResult> = (tokenInfo: RenewTokenResult) => {
    let tokenExpiration = 0;

    const txp = this.getTokenExpirationDate(tokenInfo.accessToken);

    if (txp) {
      tokenExpiration = txp.valueOf();
    }

    const expiresIn = tokenExpiration - Date.now();

    const fiveMinutes = 1000 * 60 * 5;

    const renewIn = expiresIn - fiveMinutes;

    if (renewIn <= 0) {
      return this.oAuthStrategy.reAuth(this.config).then(this.postEnsureAuthed);
    }

    setTimeout(() => {
      this.oAuthStrategy.reAuth(this.config).then(this.postEnsureAuthed);
    }, renewIn);

    return Promise.resolve(tokenInfo);
  };

  logout(queryParams?: any): void {
    this.accessToken = undefined;
    this.identityToken = undefined;

    const encode = encodeURIComponent;
    const queryString = Object.keys(queryParams || {})
      .map(key => {
        if (queryParams[key] === undefined) {
          return encode(key);
        }
        if (Array.isArray(queryParams[key])) {
          return queryParams[key].map((value: any) => `${encode(key)}=${encode(value)}`).join('&');
        }
        return `${encode(key)}=${encode(queryParams[key])}`
      })
      .join('&');
    window.location.assign(`https://${this.config.domain}/logout?${queryString}`);
  }

  get identity(): UserProfile | null {
    try { return decode(this.identityToken); }
    catch { return null; }
  }

  isLoggedIn(): boolean {
    return !!this.accessToken && !this.isTokenExpired(this.accessToken);
  }
}
