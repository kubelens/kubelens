jest.mock('jwt-decode', () => jest.fn())

jest.mock("../Strategies/Auth0Strategy", () => ({
  login: jest.fn(),
  ensureAuthed: jest.fn(),
  reAuth: jest.fn()
}));

jest.mock("../Strategies/OktaStrategy", () => ({
  login: jest.fn(),
  ensureAuthed: jest.fn(),
  reAuth: jest.fn()
}));

import * as Export from './auth-implicit-client';
import { Config } from '../../types';
import * as Auth0Strategy from '../Strategies/Auth0Strategy';
import * as OktaStrategy from '../Strategies/OktaStrategy';

describe('ChrAuth0ImplicitClient', () => {
  let testOAuthConfig: Export.oAuthConfig;
  const originalNow = Date.now;

  beforeEach(() => {
    const nowMs = Date.now();
    Date.now = () => nowMs;

    testOAuthConfig = {
      audience: 'testAudience',
      clientId: 'testClientId',
      connection: 'testConnection',
      domain: 'testDomain',
      redirectUri: 'http://testRedirectUri',
      responseType: 'token',
      scope: 'profile'
    };
  });

  afterEach(() => {
    Date.now = originalNow;
  });

  describe('constructing', () => {
    it('should instantiate without errors', () => {
      const testImplicitClient = new Export.AuthImplicitClient(testOAuthConfig);
      expect(testImplicitClient).toBeInstanceOf(Export.AuthImplicitClient);
    });

    it('should instantiate without errors if connection is an empty string', () => {
      const testAuth0ConfigMod = Object.assign(testOAuthConfig, { 'connection': '' });
      const testImplicitClient = new Export.AuthImplicitClient(testAuth0ConfigMod);
      expect(testImplicitClient).toBeInstanceOf(Export.AuthImplicitClient);
    });

    it('should throw error if specified config parameter is null', () => {
      expect(() => (new Export.AuthImplicitClient(null)))
        .toThrowError('Missing required oAuth client configuration');
    });

    it('should throw error if clientId is missing from the specified config', () => {
      delete testOAuthConfig.clientId;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth client id specified in configuration');
    });

    it('should throw error if domain is missing from the specified config', () => {
      delete testOAuthConfig.domain;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth domain specified in configuration');
    });

    it('should throw error if responseType is missing from the specified config', () => {
      delete testOAuthConfig.responseType;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth response type specified in configuration');
    });

    it('should throw error if redirectUri is missing from the specified config', () => {
      delete testOAuthConfig.redirectUri;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth redirect uri specified in configuration');
    });

    it('should throw error if audience is missing from the specified config', () => {
      delete testOAuthConfig.audience;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth audience specified in configuration');
    });

    it('should throw error if scope is missing from the specified config', () => {
      delete testOAuthConfig.scope;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth scope specified in configuration');
    });

    it('should throw error if connection is missing from the specified config', () => {
      delete testOAuthConfig.connection;
      expect(() => (new Export.AuthImplicitClient(testOAuthConfig)))
        .toThrowError('No oAuth connection specified in configuration');
    });

    it('should initialize oAuthStrategy with Auth0', () => {
      const testImplicitClient = new Export.AuthImplicitClient(testOAuthConfig);
      expect(testImplicitClient.oAuthStrategy)
        .toEqual(Auth0Strategy);
    });

    it('should initialize oAuthStrategy with Okta', () => {
      testOAuthConfig.domain = 'okta';
      const testImplicitClient = new Export.AuthImplicitClient(testOAuthConfig);
      expect(testImplicitClient.oAuthStrategy)
        .toEqual(OktaStrategy);
    });
  });

  describe('methods and properties', () => {
    let testImplicitClient: Export.AuthImplicitClient;

    beforeEach(() => {
      jest.resetAllMocks();
      jest.clearAllTimers();
      testImplicitClient = new Export.AuthImplicitClient(testOAuthConfig);
    });

    it('should remove tokens when logout is invoked and attempt to redirect', () => {
      window.location.assign = jest.fn();
      testImplicitClient.logout();
      expect(testImplicitClient.accessToken).toBeUndefined();
      expect(testImplicitClient.identityToken).toBeUndefined();
      expect(window.location.assign).toBeCalledWith('https://' + testImplicitClient.config.domain + '/logout?');
    });

    it('should remove tokens and redirect when logout is invoked with the federate parameter', () => {
      window.location.assign = jest.fn();
      testImplicitClient.logout({ 'federate': undefined });
      expect(testImplicitClient.accessToken).toBeUndefined();
      expect(testImplicitClient.identityToken).toBeUndefined();
      expect(window.location.assign).toBeCalledWith('https://' + testImplicitClient.config.domain + '/logout?federate');
    });

    it('should correctly append the redirect parameter when logout is called with all parameters', () => {
      window.location.assign = jest.fn();
      testImplicitClient.logout({ 'federate': undefined, 'returnTo': 'https://google.com' });
      expect(testImplicitClient.accessToken).toBeUndefined();
      expect(testImplicitClient.identityToken).toBeUndefined();
      expect(window.location.assign).toBeCalledWith('https://' + testImplicitClient.config.domain + '/logout?federate&returnTo=https%3A%2F%2Fgoogle.com');
    });

    it('ensures user is authed when access token is present and not expired', async () => {
      testImplicitClient.accessToken = 'a token';
      testImplicitClient.isTokenExpired = () => false;
      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      const result = await testImplicitClient.ensureAuthed();

      expect(result.accessToken).toBe('a token');
      expect((<any>testImplicitClient).postEnsureAuthed).toHaveBeenCalledTimes(1);
    });

    it('reauths auth0 when access token is present, but expired', async () => {
      const Auth0Strategy = require("../Strategies/Auth0Strategy");
      const Okta = require("../Strategies/OktaStrategy");

      testImplicitClient.accessToken = 'a token';
      testImplicitClient.isTokenExpired = () => true;
      (<jest.Mock>Auth0Strategy.reAuth).mockImplementation((testOAuthConfig) => Promise.resolve({ accessToken: 'new token', identityToken: "new id token" }));
      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      const result = await testImplicitClient.ensureAuthed();

      expect(result.accessToken).toBe('new token');
      expect(result.identityToken).toBe("new id token");
      expect(Okta.ensureAuthed).toHaveBeenCalledTimes(0);
      expect((<any>testImplicitClient).postEnsureAuthed).toHaveBeenCalledTimes(1);
    });

    it('reauths okta when access token is present, but expired', async () => {
      const Auth0Strategy = require("../Strategies/Auth0Strategy");
      testImplicitClient.oAuthStrategy = require("../Strategies/OktaStrategy");

      testImplicitClient.accessToken = 'a token';
      testImplicitClient.isTokenExpired = () => true;

      (<jest.Mock>testImplicitClient.oAuthStrategy.reAuth).mockImplementation((testOAuthConfig) => Promise.resolve({ accessToken: 'new token', identityToken: "new id token" }));
      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      const result = await testImplicitClient.ensureAuthed();

      expect(result.accessToken).toBe('new token');
      expect(result.identityToken).toBe("new id token");
      expect(Auth0Strategy.ensureAuthed).toHaveBeenCalledTimes(0);
      expect((<any>testImplicitClient).postEnsureAuthed).toHaveBeenCalledTimes(1);
    });

    it('ensures strategy auth0 if no access token, resolving if strategy ensureAuth returns promise with accessToken && identity token', async () => {
      const Auth0 = require("../Strategies/Auth0Strategy");
      const Okta = require("../Strategies/OktaStrategy");

      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      (<jest.Mock>Auth0.ensureAuthed).mockResolvedValue({ accessToken: 'shiny new token', identityToken: 'shiny new token' });

      expect.hasAssertions();

      await testImplicitClient.ensureAuthed().then((data) => {
        expect(data.accessToken).toBe('shiny new token');
        expect(data.identityToken).toBe('shiny new token');
        expect(Okta.ensureAuthed).toHaveBeenCalledTimes(0);
        expect((<any>testImplicitClient).postEnsureAuthed).toHaveBeenCalledTimes(1);
      });
    });

    it('ensures strategy okta if no access token, resolving if strategy ensureAuth returns promise with accessToken && identity token', async () => {
      const Auth0Strategy = require("../Strategies/Auth0Strategy");
      testImplicitClient.oAuthStrategy = require("../Strategies/OktaStrategy");

      testOAuthConfig.domain = 'https://chrobinson-poc.okta.com/oauth2/aus80c9fxZ5s16CFY356/';
      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      (<jest.Mock>testImplicitClient.oAuthStrategy.ensureAuthed).mockResolvedValue({ accessToken: 'shiny new token', identityToken: 'shiny new token' });

      expect.hasAssertions();

      await testImplicitClient.ensureAuthed().then((data) => {
        expect(data.accessToken).toBe('shiny new token');

        expect(data.identityToken).toBe('shiny new token');
        expect(Auth0Strategy.ensureAuthed).toHaveBeenCalledTimes(0);
        expect((<any>testImplicitClient).postEnsureAuthed).toHaveBeenCalledTimes(1);
      });
    });

    it('ensures strategy auth0 if no access token, rejecting if strategy ensureAuth returns rejected promise', async () => {
      const Auth0 = require("../Strategies/Auth0Strategy");

      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      (<jest.Mock>Auth0.ensureAuthed).mockRejectedValue("no access token");

      expect.hasAssertions();

      await testImplicitClient.ensureAuthed().then(() => {
        throw 'this should never happen'
      }).catch(e => {
        expect(e).not.toBe('this should never happen');
        expect((<any>testImplicitClient).postEnsureAuthed).not.toHaveBeenCalled();
      });
    });


    it('ensures strategy okta if no access token, rejecting if strategy ensureAuth returns rejected promise', async () => {
      testImplicitClient.oAuthStrategy = require("../Strategies/OktaStrategy");

      testOAuthConfig.domain = 'https://chrobinson-poc.okta.com/oauth2/aus80c9fxZ5s16CFY356/';
      (<any>testImplicitClient).postEnsureAuthed = jest.fn(value => Promise.resolve(value));

      (<jest.Mock>testImplicitClient.oAuthStrategy.ensureAuthed).mockRejectedValue("no access token");

      expect.hasAssertions();

      await testImplicitClient.ensureAuthed().then(() => {
        throw 'this should never happen'
      }).catch(e => {
        expect(e).not.toBe('this should never happen');
        expect((<any>testImplicitClient).postEnsureAuthed).not.toHaveBeenCalled();
      });
    });
    it('recursively sets a timer to reauth okta five minutes before current token expires', async () => {
      jest.useFakeTimers();
      testImplicitClient.oAuthStrategy = require("../Strategies/OktaStrategy");

      const fiveMinutes = 1000 * 60 * 5;
      const tenMinutes = fiveMinutes * 2;

      let latestPostEnsureAuthedPromise;
      let latestReAuthPromise;

      const postEnsureAuthed = (<any>testImplicitClient).postEnsureAuthed as
        (tokenInfo: Export.RenewTokenResult) => Promise<Export.RenewTokenResult>;

      jest.spyOn(testImplicitClient as any, 'postEnsureAuthed').mockImplementation((info) => {
        latestPostEnsureAuthedPromise = postEnsureAuthed.call(testImplicitClient, info);

        return latestPostEnsureAuthedPromise;
      });

      const whatReAuthReturnsFirst = Promise.resolve({ accessToken: 'my first new token' });
      const whatReAuthReturnsSecond = Promise.resolve({ accessToken: 'my second new token' });
      (<jest.Mock>testImplicitClient.oAuthStrategy.reAuth)
        .mockImplementationOnce(() => {
          latestReAuthPromise = whatReAuthReturnsFirst;
          return whatReAuthReturnsFirst;
        })
        .mockImplementationOnce(() => {
          latestReAuthPromise = whatReAuthReturnsSecond;
          return whatReAuthReturnsSecond;
        });

      testImplicitClient.getTokenExpirationDate = jest.fn()
        .mockImplementation(() => {
          const result = new Date(0);

          result.setUTCMilliseconds(Date.now() + tenMinutes);

          return result;
        });

      testImplicitClient.config.domain = 'https://chrobinson-poc.okta.com/oauth2/aus80c9fxZ5s16CFY356/';

      // kick it off
      (<any>testImplicitClient).postEnsureAuthed({ accessToken: 'my token', identityToken: "my id token" });

      let result = await latestPostEnsureAuthedPromise;

      expect(result).toEqual({ accessToken: 'my token', identityToken: "my id token" });

      // wait five minutes to trip the ten minute timeout
      jest.advanceTimersByTime(fiveMinutes - 1);

      expect(testImplicitClient.oAuthStrategy.reAuth).not.toHaveBeenCalled();

      jest.advanceTimersByTime(1);

      expect(testImplicitClient.oAuthStrategy.reAuth).toHaveBeenCalledTimes(1);

      // wait for reauth, which will call postEnsureAuthed
      await latestReAuthPromise;
      result = await latestPostEnsureAuthedPromise;

      expect(result).toEqual({ accessToken: 'my first new token' });

      // wait another five minutes to trip the twenty minute timeout
      jest.advanceTimersByTime(fiveMinutes - 1);

      expect(testImplicitClient.oAuthStrategy.reAuth).toHaveBeenCalledTimes(1);

      jest.advanceTimersByTime(1);

      expect(testImplicitClient.oAuthStrategy.reAuth).toHaveBeenCalledTimes(2);

      // wait for reauth, which will call postEnsureAuthed
      await latestReAuthPromise;
      result = await latestPostEnsureAuthedPromise;

      expect(result).toEqual({ accessToken: 'my second new token' });

      jest.useRealTimers();
    });

    it('recursively sets a timer to reauth auth0 five minutes before current token expires', async () => {
      jest.useFakeTimers();
      const Auth0 = require("../Strategies/Auth0Strategy");

      const fiveMinutes = 1000 * 60 * 5;
      const tenMinutes = fiveMinutes * 2;

      let latestPostEnsureAuthedPromise;
      let latestReAuthPromise;

      const postEnsureAuthed = (<any>testImplicitClient).postEnsureAuthed as
        (tokenInfo: Export.RenewTokenResult) => Promise<Export.RenewTokenResult>;

      jest.spyOn(testImplicitClient as any, 'postEnsureAuthed').mockImplementation((info) => {
        latestPostEnsureAuthedPromise = postEnsureAuthed.call(testImplicitClient, info);

        return latestPostEnsureAuthedPromise;
      });

      const whatReAuthReturnsFirst = Promise.resolve({ accessToken: 'my first new token', identityToken: `my first id token` });
      const whatReAuthReturnsSecond = Promise.resolve({ accessToken: 'my second new token', identityToken: `my second id token` });
      (<jest.Mock>Auth0.reAuth)
        .mockImplementationOnce(() => {
          latestReAuthPromise = whatReAuthReturnsFirst;
          return whatReAuthReturnsFirst;
        })
        .mockImplementationOnce(() => {
          latestReAuthPromise = whatReAuthReturnsSecond;
          return whatReAuthReturnsSecond;
        });

      testImplicitClient.getTokenExpirationDate = jest.fn()
        .mockImplementation(() => {
          const result = new Date(0);

          result.setUTCMilliseconds(Date.now() + tenMinutes);

          return result;
        });

      // kick it off
      (<any>testImplicitClient).postEnsureAuthed({ accessToken: 'my token', identityToken: `my id token` });

      let result = await latestPostEnsureAuthedPromise;

      expect(result).toEqual({ accessToken: 'my token', identityToken: `my id token` });

      // wait five minutes to trip the ten minute timeout
      jest.advanceTimersByTime(fiveMinutes - 1);

      expect(Auth0.reAuth).not.toHaveBeenCalled();

      jest.advanceTimersByTime(1);

      expect(Auth0.reAuth).toHaveBeenCalledTimes(1);

      // wait for reauth, which will call postEnsureAuthed
      await latestReAuthPromise;
      result = await latestPostEnsureAuthedPromise;

      expect(result).toEqual({ accessToken: 'my first new token', identityToken: `my first id token` });

      // wait another five minutes to trip the twenty minute timeout
      jest.advanceTimersByTime(fiveMinutes - 1);

      expect(Auth0.reAuth).toHaveBeenCalledTimes(1);

      jest.advanceTimersByTime(1);

      expect(Auth0.reAuth).toHaveBeenCalledTimes(2);

      // wait for reauth, which will call postEnsureAuthed
      await latestReAuthPromise;
      result = await latestPostEnsureAuthedPromise;

      expect(result).toEqual({ accessToken: 'my second new token', identityToken: `my second id token` });

      jest.useRealTimers();
    });

    it('reauths auth0 if current token expires in five minutes or less', async () => {
      const Auth0 = require("../Strategies/Auth0Strategy");
      jest.useFakeTimers();

      const fiveMinutes = 1000 * 60 * 5;

      const postEnsureAuthed = (<any>testImplicitClient).postEnsureAuthed as
        (tokenInfo: Export.RenewTokenResult) => Promise<Export.RenewTokenResult>;

      (<jest.Mock>Auth0.reAuth).mockImplementation(() => Promise.resolve({ accessToken: 'new token', identityToken: `my new id token` }));

      testImplicitClient.getTokenExpirationDate = jest.fn()
        .mockImplementationOnce(() => {
          const result = new Date(0);

          result.setUTCMilliseconds(Date.now() + fiveMinutes);

          return result;
        })
        .mockImplementationOnce(() => {
          const result = new Date(0);

          result.setUTCMilliseconds(Date.now() + fiveMinutes + fiveMinutes + fiveMinutes);

          return result;
        });

      const promise = postEnsureAuthed({ accessToken: 'my token', identityToken: `my id token` });

      expect(Auth0.reAuth).toHaveBeenCalledTimes(1);

      const result = await promise;

      expect(result).toEqual({ accessToken: 'new token', identityToken: `my new id token` });

      jest.useRealTimers();
    });

    it('reauths okta if current token expires in five minutes or less', async () => {
      testImplicitClient.oAuthStrategy = require("../Strategies/OktaStrategy");
      testOAuthConfig.domain = 'https://chrobinson-poc.okta.com/oauth2/aus80c9fxZ5s16CFY356/';
      jest.useFakeTimers();

      const fiveMinutes = 1000 * 60 * 5;

      const postEnsureAuthed = (<any>testImplicitClient).postEnsureAuthed as
        (tokenInfo: Export.RenewTokenResult) => Promise<Export.RenewTokenResult>;

      (<jest.Mock>testImplicitClient.oAuthStrategy.reAuth).mockImplementation(() => Promise.resolve({ accessToken: 'new token', identityToken: `new id token` }));

      testImplicitClient.getTokenExpirationDate = jest.fn()
        .mockImplementationOnce(() => {
          const result = new Date(0);

          result.setUTCMilliseconds(Date.now() + fiveMinutes);

          return result;
        })
        .mockImplementationOnce(() => {
          const result = new Date(0);

          result.setUTCMilliseconds(Date.now() + fiveMinutes + fiveMinutes + fiveMinutes);

          return result;
        });

      const promise = postEnsureAuthed({ accessToken: 'my token', identityToken: `my new id token` });

      expect(testImplicitClient.oAuthStrategy.reAuth).toHaveBeenCalledTimes(1);

      const result = await promise;

      expect(result).toEqual({ accessToken: 'new token', identityToken: `new id token` });

      jest.useRealTimers();
    });

    it('can be created with sensible defaults', () => {
      const config = {
        productUrl: 'product/v1',
        oAuthAudience: 'audience',
        oAuthClientId: 'client id',
        oAuthDomain: 'domain'
      };

      const instance = Export.AuthImplicitClient.createDefault(config as Config);

      expect(instance).toBeDefined();
      expect(instance.config.connection).toBe('chradfs');
      expect(instance.config.scope).toBe('openid profile email');
      expect(instance.config.responseType).toBe('id_token token');
      expect(instance.config.audience).toBe(config.oAuthAudience);
      expect(instance.config.clientId).toBe(config.oAuthClientId);
      expect(instance.config.domain).toBe(config.oAuthDomain);
      expect(instance.config.redirectUri).toBe(location.href.toLowerCase());
    })

    it('supports overriding sensible defaults', () => {
      const config = {
        productUrl: 'product/v1',
        oAuthAudience: 'audience',
        oAuthClientId: 'client id',
        oAuthDomain: 'domain'
      };

      const instance = Export.AuthImplicitClient.createDefault(config as Config, { domain: 'something different' });

      expect(instance.config.domain).toBe('something different');
    })

    it('exposes an identity property that returns a profile if idToken can be parsed', () => {
      const decode = require('jwt-decode');
      decode.mockImplementation(() => ({ profile: true }))

      const config = {
        productUrl: 'product/v1',
        oAuthAudience: 'audience',
        oAuthClientId: 'client id',
        oAuthDomain: 'domain'
      };

      const instance = Export.AuthImplicitClient.createDefault(config as Config, { domain: 'something different' });

      expect(instance.identity).toEqual({ profile: true });
    })

    it('exposes an identity property that returns null if idToken cannot be parsed', () => {
      const decode = require('jwt-decode');
      decode.mockImplementation(() => { throw 'bork!' })

      const config = {
        productUrl: 'product/v1',
        oAuthAudience: 'audience',
        oAuthClientId: 'client id',
        oAuthDomain: 'domain'
      };

      const instance = Export.AuthImplicitClient.createDefault(config as Config, { domain: 'something different' });

      expect(instance.identity).toEqual(null);
    });
  });
});
