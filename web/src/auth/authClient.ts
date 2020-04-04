import OktaAuth from '@okta/okta-auth-js';
import { AuthError } from './AuthError';
import { getStateSession, setStateSession, removeStateSession, parseState, getRedirectUri } from './utils';
import decode from 'jwt-decode';
import { Config } from 'types';

const accessTokenKey = "accessToken";
const idTokenKey = "idToken";

export interface OAuthConfig {
  clientId: string,
  domain: string,
  responseType: string,
  redirectUri: string,
  audience: string,
  scope: string,
  history: any
}

export interface QueryParams {
  post_logout_redirect_uri: string
}

export interface RenewTokenResult {
  accessToken: string,
  identityToken?: string,
  identity?: any,
}
export interface CustomDecodedJwt {
  sub:string,
  name:string,
  email:string,
  ver:number,
  iss:string,
  aud:string,
  iat:number,
  exp: number,
  jti:string,
  amr:Array<string>,
  idp:string,
  nonce:string
  preferred_username: string,
  auth_time: number,
  at_hash: string
  samAccountName: string,
}

export interface AuthClient<T = CustomDecodedJwt> {
  logout: (queryParams: QueryParams) => Promise<void>,
  login: () => Promise<void>,
  getIdentity: () => Promise<T>,
  getAccessToken: () => Promise<string>,
  getIdentityToken: () => Promise<string>,
  ensureAuthed: () => Promise<RenewTokenResult>
}

export interface CommonAuthConfig {
  oAuthAudience: string
  oAuthClientId: string
  oAuthDomain: string
}

export default function createClient<T>(commonConfig: Config, overrides?: Partial<OAuthConfig>): AuthClient<T> {
  const badConfigOptions: string[] = [];

  if (!commonConfig) {
    throw new AuthError('missing_configuration', 'Missing required oAuth client configuration');
  }

  const config = {
    ...{
      redirectUri: getRedirectUri(),
      connection: commonConfig.oAuthConnection,
      scope: 'openid profile email',
      responseType: 'id_token token',
      audience: commonConfig.oAuthAudience,
      clientId: commonConfig.oAuthClientId,
      domain: commonConfig.oAuthJwtIssuer
    },
    ...overrides
  }

  // do validation of config blob
  if (!config.clientId) {
    badConfigOptions.push('No oAuth client id specified in configuration');
  }

  if (!config.domain) {
    badConfigOptions.push('No oAuth domain specified in configuration');
  }

  if (!config.responseType) {
    badConfigOptions.push('No oAuth response type specified in configuration');
  }

  if (!config.redirectUri) {
    badConfigOptions.push('No oAuth redirect uri specified in configuration');
  }

  if (!config.audience) {
    badConfigOptions.push('No oAuth audience specified in configuration');
  }

  if (!config.scope) {
    badConfigOptions.push('No oAuth scope specified in configuration');
  }

  if (badConfigOptions.length > 0) {
    throw new AuthError('missing_configuration', badConfigOptions.join(', '))
  }


  //new up the thing we will use internally
  const oktaAuthClient = new OktaAuth({
    issuer: config.domain,
    clientId: config.clientId,
    redirectUri: config.redirectUri,
    tokenManager: {
      storage: "memory",
      expireEarlySeconds: 120,
      autoRenew: true
    }
  });

  // declare anything that needs to persist w/in scope (any stateful stuff)
  oktaAuthClient.tokenManager.on('expired', async function (key: string) {
    try {
      await oktaAuthClient.tokenManager.renew(key);
    }
    catch (err) {
      //If Session is bad
      client.login();
      throw new AuthError(err.errorCode, err.message);
    }
  });

  const client = {
    logout: async (queryParams: QueryParams) => {
      if (!queryParams.post_logout_redirect_uri) {
        throw new Error('Please specify queryParams.post_logout_redirect_uri so we can send the user somewhere after successfully logging out');
      }
      await oktaAuthClient.signOut();
      window.location.assign(queryParams.post_logout_redirect_uri);
    },
    login: async () => {
      const url = window.location.href.replace(new RegExp(document.baseURI, 'gi'), '') || document.baseURI;
      const state = setStateSession(url);

      await oktaAuthClient.token.getWithRedirect({
        scopes: config.scope.split(' '),
        responseType: ['token', 'id_token'],
        state: state
      });
    },
    getIdentity: async (): Promise<T> => {
      try {
        let token = await oktaAuthClient.tokenManager.get(idTokenKey);
        let idToken = token.idToken
        return decode(idToken);
      }
      catch { return null as any; }
    },
    getAccessToken: async (): Promise<string> => {
      try {
        let token = await oktaAuthClient.tokenManager.get(accessTokenKey);
        return token.accessToken
      }
      catch (err) {
        throw new AuthError(err.errorCode, err.message);
      }
    },
    getIdentityToken: async (): Promise<string> => {
      try {
        let token = await oktaAuthClient.tokenManager.get(idTokenKey);
        return token.idToken
      }
      catch (err) {
        throw new AuthError(err.errorCode, err.message);
      }
    },
    ensureAuthed: async (): Promise<RenewTokenResult> => {
      const state = parseState(window.location.href);
      const exists = await oktaAuthClient.session.exists();

      //if session exists AND we are not being redirected after logging in
      if (exists && window.location.href.indexOf("access_token") === -1) {
        const session = await oktaAuthClient.session.get()
        if (session.status !== "ACTIVE") {
          client.login();
          throw new AuthError('inactive_session', `session was ${session.status}`);
        }
        else {
          try {
            let accessToken = await oktaAuthClient.tokenManager.get(accessTokenKey);
            let idToken = await oktaAuthClient.tokenManager.get(idTokenKey);
            if (accessToken && idToken) {
              return {
                accessToken: accessToken.accessToken,
                identityToken: idToken.idToken,
                identity: idToken.claims
              };
            }
            //We're missing a token!
            else {
              return await refreshTokensWithSession(oktaAuthClient, config, session, "missing");
            }
          }
          catch (err) {
            client.login();
            throw new AuthError(err.errorCode, err.message);
          }
        }
      }
      else {
        try {
          return await parseFromUrl(oktaAuthClient, state);
        } catch (err) {
          // If there are more error messages to check for, we should put the list
          // in a separate function.

          // If the error message is that the JWT was issued in the future, this means that the users
          // clock is off by more than five minutes. We don't want to login in this case, so the UI
          // can respond with a message to the user.
          // If it's any other error message, then let it login.
          if (err.message.indexOf("The JWT was issued in the future") === -1) {
            client.login();
          }
          throw new AuthError(err.errorCode, err.message);
        };
      }
    }
  }

  return client;
}

async function parseFromUrl(oktaAuthClient: any, state: string) {
  const tokens = await oktaAuthClient.token.parseFromUrl();
  oktaAuthClient.tokenManager.add(accessTokenKey, tokens[0]);
  oktaAuthClient.tokenManager.add(idTokenKey, tokens[1]);
  if (state) {
    const decodedState = decodeURIComponent(state);
    const sessionState = getStateSession(decodedState);
    removeStateSession(decodedState);
    if (sessionState) {
      window.history.replaceState(null, null as any, sessionState);
    }
  }
  return {
    accessToken: tokens[0].accessToken,
    identityToken: tokens[1].idToken,
    identity: tokens[1].claims
  };
}

async function refreshTokensWithSession(oktaAuthClient: any, config: any, session: any, state?: string) {
  const tokens = await oktaAuthClient.token.getWithoutPrompt({
    scopes: config.scope.split(' '),
    responseType: ['token', 'id_token'],
    state: state,
    sessionToken: session.id,
  });
  oktaAuthClient.tokenManager.add(accessTokenKey, tokens[0]);
  oktaAuthClient.tokenManager.add(idTokenKey, tokens[1]);
  return {
    accessToken: tokens[0].accessToken,
    identityToken: tokens[1].idToken,
    identity: tokens[1].claims
  };
}
