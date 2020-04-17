jest.mock('jwt-decode', () => jest.fn())
jest.mock('@okta/okta-auth-js', () => {
  const getWithRedirectMock = jest.fn();
  const renewMock = jest.fn();
  const parseUrlMock = jest.fn();
  const signoutMock = jest.fn();
  const existsMock = jest.fn();
  const addMock = jest.fn();
  const getTokenMock = jest.fn();
  const getMock = jest.fn();
  const onMock = jest.fn();
  const getWithoutPromptMock = jest.fn();

  return class {
    options = {
      redirectUri: 'state=%2f'
    };
    session = {
      exists: existsMock, 
      get: getMock
    };
    token = {
      getWithoutPrompt: getWithoutPromptMock,
      getWithRedirect: getWithRedirectMock,
      parseFromUrl: parseUrlMock

    };
    tokenManager={
      add: addMock,
      renew: renewMock,
      get: getTokenMock,
      on: onMock
    }
    signOut = signoutMock;
  };
});

import * as utils from './utils';
import createClient, {OAuthConfig, AuthClient}from './authClient';
import OktaAuth from '@okta/okta-auth-js';
import { AuthError } from './AuthError';
import { Config } from 'types';

describe('OktaStrategy', () => {
  let testOktaConfig: OAuthConfig;
  let testCommonConfig: Config;
  let oauth: jest.Mock<OktaAuth>;
  let authClient: AuthClient;
  beforeEach(() => {
    oauth = require('@okta/okta-auth-js');
    testOktaConfig = {
      audience: 'testAudience',
      clientId: 'testClientId',
      domain: 'testDomain',
      redirectUri: 'http://testRedirectUri',
      responseType: 'token',
      scope: 'profile',
      history: null
    };

    testCommonConfig = {
      oAuthAudience: 'testAudience',
      oAuthClientId: 'testClientId',
      oAuthJwtIssuer: 'testdomain',
      oAuthRedirectUri: 'http://testRedirectUri',
      oAuthResponseType: 'token',
      oAuthRequestType: 'Bearer',
      oAuthScope: 'openid profile email',
      oAuthConnection: 'adfs',
      oAuthEnabled: false,
      availableClusters: [],
      deployerLinkName: ''
    };

    window = Object.create(window);
    Object.defineProperty(window, "location", {
      value: {
        hash: '',
        href: ' http://localhost',
        pathname: "no_token"
      },
      writable: true
    });
  });
  describe('constructing', () => {
    it('should instantiate without errors', () => {
      const testImplicitClient = createClient(testCommonConfig);
      expect(testImplicitClient).toHaveProperty("login");
    });

    it('should call the on function, and the on function should renew a token',async () => {
      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.renew).mockResolvedValue("done");

      expect(mock.tokenManager.on).toHaveBeenCalled();
      let fn = (mock.tokenManager.on).mock.calls[0][1];
      await fn("testkey")
      expect(mock.tokenManager.renew).toHaveBeenCalled()
    });

    it('should call the on function, and the on function should throw error if renew errors',async () => {
      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.renew).mockRejectedValue("done");

      expect(mock.tokenManager.on).toHaveBeenCalled();
      let fn = (mock.tokenManager.on).mock.calls[0][1];
      expect(fn("testkey")).rejects.toBeTruthy();
    });

    it('should throw error if specified config parameter is null', () => {
      expect(() => (createClient(null)))
        .toThrowError('Missing required oAuth client configuration');
    });

    it('should throw error if clientId is missing from the specified config', () => {
      const testOktaConfigMod = Object.assign(testOktaConfig, { 'clientId': '' });
      expect(() => (createClient(testCommonConfig, testOktaConfigMod)))
        .toThrowError('No oAuth client id specified in configuration');
    });

    it('should throw error if domain is missing from the specified config', () => {
      const testOktaConfigMod = Object.assign(testOktaConfig, { 'domain': '' });
      expect(() => (createClient(testCommonConfig, testOktaConfigMod)))
        .toThrowError('No oAuth domain specified in configuration');
    });

    it('should throw error if responseType is missing from the specified config', () => {
      const testOktaConfigMod = Object.assign(testOktaConfig, { 'responseType': '' });
      expect(() => (createClient(testCommonConfig, testOktaConfigMod)))
        .toThrowError('No oAuth response type specified in configuration');
    });

    it('should throw error if redirectUri is missing from the specified config', () => {
      const testOktaConfigMod = Object.assign(testOktaConfig, { 'redirectUri': '' });
      expect(() => (createClient(testCommonConfig, testOktaConfigMod)))
        .toThrowError('No oAuth redirect uri specified in configuration');
    });

    it('should throw multiple errors if scope and audience is missing from the specified config', () => {
      const testOktaConfigMod = Object.assign(testOktaConfig, { 'scope': '', 'audience': '' });
      expect(() => (createClient(testCommonConfig, testOktaConfigMod)))
        .toThrowError('No oAuth audience specified in configuration, No oAuth scope specified in configuration');
    });
  });

  describe('methods and properties', () => {
    const testToken = {accessToken:"test", idToken:"test"};
    beforeEach(() => {
      jest.resetAllMocks();

      authClient = createClient(testCommonConfig);
      (window as any).sessionStorage = {
        getItem: jest.fn(),
        setItem: jest.fn(),
        removeItem: jest.fn()
      };
    });

    it(`login calls getWithRedirect with scopes string array`, () => {
      jest.spyOn(utils, 'setStateSession').mockImplementationOnce(() => 'test-state-value');
      const mock = new oauth();

      authClient.login();

      expect(mock.token.getWithRedirect).toHaveBeenCalledWith({
        responseType: ['token', 'id_token'],
        scopes: ['openid', 'profile','email'],
        state: 'test-state-value'
      });
    });

    it(`getAccessToken returns its token`,async () => {
      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.get).mockResolvedValue(testToken);
      let token = await authClient.getAccessToken();

      expect(mock.tokenManager.get).toHaveBeenCalledWith("accessToken");
      expect(token).toEqual(testToken.accessToken)
    });

    it(`getAccessToken throws okta error`, async () => {
      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.get).mockRejectedValue(
        new AuthError('INTERNAL', 'Unable to parse url')
      );
      expect(authClient.getAccessToken()).rejects.toEqual(new AuthError('INTERNAL', 'Unable to parse url'))
    });

    it(`getIdentityToken returns its token`, async () => {
      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.get).mockResolvedValue(testToken);
      let token = await authClient.getIdentityToken();

      expect(mock.tokenManager.get).toHaveBeenCalledWith("idToken");
      expect(token).toEqual(testToken.idToken)
    });

    it(`getIdentityToken throws okta error`, async () => {
      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.get).mockRejectedValue(
        new AuthError('INTERNAL', 'Unable to parse url')
      );
      expect(authClient.getIdentityToken()).rejects.toEqual(new AuthError('INTERNAL', 'Unable to parse url'))
    });

    it(`getIdentity returns its token`, async () => {
      const decode = require('jwt-decode');
      decode.mockImplementation(() => { return "token!" })

      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.get).mockResolvedValue(testToken);
      let token = await authClient.getIdentity();

      expect(decode).toHaveBeenCalled();
      expect(token).toEqual("token!")
    });

    it(`getIdentity throws error if cannot decode`, async () => {
      const decode = require('jwt-decode');
      decode.mockImplementation(() => { throw "token!" })

      const mock = new oauth();
      (<jest.Mock>mock.tokenManager.get).mockResolvedValue(testToken);
      let token = await authClient.getIdentity();

      expect(token).toEqual(null)
    });
    it('parses hash if no access token, rejecting if parseFromUrl throws', async () => {
      const mock = new oauth();
      const oktaError = new AuthError('INTERNAL', 'Unable to parse url');
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(false);
      (<jest.Mock>mock.token.parseFromUrl).mockRejectedValue(
        new AuthError('INTERNAL', 'Unable to parse url')
      );
      await authClient.ensureAuthed()
        .then(() => {
          throw 'this should never happen';
        })
        .catch(e => {
          expect({ ...e }).toEqual({ ...oktaError });
        });
    });

    it('parses hash if no access token, rejecting if parseFromUrl returns no accessToken', async () => {
      const mock = new oauth();
      const oktaError = new AuthError('INTERNAL', 'Unable to parse url');
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(false);
      (<jest.Mock>mock.token.parseFromUrl).mockRejectedValue(oktaError);
      await authClient.ensureAuthed()
        .then(() => {
          throw 'this should never happen';
        })
        .catch(e => {
          expect(e).not.toBe('this should never happen');
          expect(mock.token.getWithRedirect).toHaveBeenCalledTimes(1);
        });
    });

    it('parses hash if no access token, rejecting if parseFromUrl returns error that JWT is in the future', async () => {
      const mock = new oauth();
      const oktaError = new AuthError("INTERNAL", "The JWT was issued in the future");
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(false);
      (<jest.Mock>mock.token.parseFromUrl).mockRejectedValue(oktaError);
      await authClient.ensureAuthed()
        .then(() => {
          throw 'this should never happen';
        })
        .catch(e => {
          expect(e.message).toBe('The JWT was issued in the future');
          expect(e.errorCode).toBe('INTERNAL');
          expect(mock.token.getWithRedirect).toHaveBeenCalledTimes(0);
        });
    });

    it('parses hash if no access token, preserve state and resolve if parseHash returns accessToken', async () => {
      const mock = new oauth();

      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(false);
      (<jest.Mock>mock.token.parseFromUrl).mockResolvedValue([
        { accessToken: 'a golden accessToken' },
        { idToken: 'a silver idtoken', claims: 'I claim you' }
      ]);
      await authClient.ensureAuthed().then(data => {
        expect(data.accessToken).toBe('a golden accessToken');
        expect(data.identityToken).toBe('a silver idtoken');
        expect(data.identity).toBe('I claim you');
        expect(window.location.hash).toBe('');
      });
    });

    it('checks for active session, if session is active, return access & id token', async () => {
      const mock = new oauth();
      jest.spyOn(window.history, 'replaceState');
      window.location.hash = '#state=%23%2Fthe-coolest-route';
      window.location.pathname = "no_token"
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(true);
      (<jest.Mock>mock.session.get).mockResolvedValue(
        {
          id: 'a session token',
          userId: 'some user id',
          status: 'ACTIVE'
        }
      );
      (<jest.Mock>mock.tokenManager.get).mockResolvedValueOnce(
        { accessToken: 'a golden accessToken' }).mockResolvedValueOnce(
          { idToken: 'a silver idtoken', claims: 'I claim you' }
        );

      await authClient.ensureAuthed().then(data => {
        expect(data.accessToken).toBe('a golden accessToken');
        expect(data.identityToken).toBe('a silver idtoken');
        expect(data.identity).toBe('I claim you');
      });
    });

    it('checks for active session, if session is active and getWithoutPrompt fails', async () => {
      const mock = new oauth();
      jest.spyOn(window.history, 'replaceState');
      const oktaError = new AuthError('INTERNAL', 'invalid session token');
      window.location.hash = '#state=%23%2Fthe-coolest-route';
      window.location.pathname = "no_token"
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(true);
      (<jest.Mock>mock.session.get).mockResolvedValue(
        {
          id: 'a session token',
          userId: 'some user id',
          status: 'ACTIVE'
        }
      );
      (<jest.Mock>mock.token.getWithoutPrompt).mockRejectedValue(oktaError);

      await authClient.ensureAuthed()
        .then(() => {
          throw 'this should never happen';
        })
        .catch(e => {
          expect(e).not.toBe('this should never happen');
        });
    });


    it('checks for active session, if session is active and getWithoutPrompt suceeds', async () => {
      const mock = new oauth();
      jest.spyOn(window.history, 'replaceState');
      const oktaError = new AuthError('INTERNAL', 'invalid session token');
      window.location.hash = '#state=%23%2Fthe-coolest-route';
      window.location.pathname = "no_token"
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(true);
      (<jest.Mock>mock.session.get).mockResolvedValue(
        {
          id: 'a session token',
          userId: 'some user id',
          status: 'ACTIVE'
        }
      );
      (<jest.Mock>mock.token.getWithoutPrompt).mockResolvedValue([{ accessToken: "test" }, { idToken: "test", claims:"test"}]);

      let result = await authClient.ensureAuthed()

      expect(result).toEqual({ accessToken: "test", identityToken: "test", identity: "test" })

    });


    it('checks for active session, if session is NOT ACTIVE, redirect login', async () => {
      const mock = new oauth();
      jest.spyOn(window.history, 'replaceState');
      window.location.hash = '#state=%23%2Fthe-coolest-route';
      window.location.pathname = "no_token_in_path"
      expect.hasAssertions();
      (<jest.Mock>mock.session.exists).mockResolvedValue(true);
      (<jest.Mock>mock.session.get).mockResolvedValue(
        {
          id: 'a session token',
          userId: 'some user id',
          status: 'MFA_REQUIRED'
        }
      );

      await authClient.ensureAuthed().catch(data => {
        expect(mock.token.getWithoutPrompt).not.toBeCalled();
        expect(mock.token.getWithRedirect).toHaveBeenCalledTimes(1);
      });
    });

    it('signs out the user and redirects to configured redirectUri', async () => {
      const mock = new oauth();
      const post_logout_redirect_uri = testOktaConfig.redirectUri;
      (<jest.Mock>mock.signOut).mockResolvedValue({});

      window.location.assign = jest.fn();
      await authClient.logout({post_logout_redirect_uri});

      expect(mock.signOut).toHaveBeenCalled()
      expect(window.location.assign).toHaveBeenCalledWith(
        testOktaConfig.redirectUri
      );
    });
  });
});
