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
import OktaAuth from '@okta/okta-auth-js';
import {
    oAuthConfig,
    RenewTokenResult
} from "../clients/auth-implicit-client";
import { AuthError } from '../AuthError';

let oauth: any;
let accessTokenObject: any;
let idTokenObject: any;

function createClient(config: oAuthConfig) {
    return new OktaAuth({
        issuer: config.domain,
        clientId: config.clientId,
        redirectUri: config.redirectUri,
    });
}

export function parseState(url: string) {
    let stateValue: string = "";
    let stateIndex = url.indexOf("state");
    if (stateIndex > -1) {
        let subState = url.substring(stateIndex, url.length);
        let equalIndex = subState.indexOf('=');
        let ampIndex = subState.indexOf('&');
        if (ampIndex > -1) {
            stateValue = subState.substring(equalIndex + 1, ampIndex);
        } else {
            stateValue = subState.substring(equalIndex + 1, subState.length);
        }

    }
    return stateValue;
}

export function login(oauth: any, config: oAuthConfig) {
    const url = `${window.location.pathname}${window.location.search}`;
    oauth.options.redirectUri = window.location.protocol + '//' + window.location.host + window.location.pathname;

    oauth.token.getWithRedirect({
        scope: config.scope,
        responseType: ['token', 'id_token'],
        state: url,
    });
}

export function ensureAuthed(config: oAuthConfig): Promise<RenewTokenResult> {

    oauth = createClient(config);

    return new Promise((resolve, reject) => {
        let state: string;
        oauth.token.parseFromUrl()
            .then((token: any) => {
                accessTokenObject = token[0];
                idTokenObject = token[1];

                state = parseState(oauth.options.redirectUri);
                if (state) {
                    config.history.replace(atob(state));
                }

                return resolve({
                    accessToken: token[0].accessToken,
                    identityToken: token[1].idToken,
                    identity: token[1].claims
                });
            })
            .catch((err: any) => {
                login(oauth, config);
                return reject(new AuthError(err.errorCode, err.message));
            })
    });
}

export function reAuth(config: oAuthConfig): Promise<RenewTokenResult> {

    if (!oauth) {
        oauth = createClient(config);
    }

    return new Promise(async (resolve, reject) => {

        try {
            accessTokenObject = await oauth.token.renew(accessTokenObject);
            idTokenObject = await oauth.token.renew(idTokenObject);

            resolve({
                accessToken: accessTokenObject.accessToken,
                identityToken: idTokenObject.idToken
            });
        } catch (e) {
            login(oauth, config)
            reject(new AuthError(e.errorCode, e.message))
        }
    });
}