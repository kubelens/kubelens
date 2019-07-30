jest.mock("auth0-js", () => {
  const authorizeMock = jest.fn();
  const checkSessionMock = jest.fn();
  const parseHashMock = jest.fn();

  return {
    WebAuth: class {
      authorize = authorizeMock;
      checkSession = checkSessionMock;
      parseHash = parseHashMock;
    }
  };
});

import * as Auth0Strategy from "./auth0-strategy";
import * as ClientExport from "../clients/auth-implicit-client";
import { WebAuth } from "auth0-js";
import { AuthError } from '../AuthError';

describe("Auth0Strategy", () => {
  let testAuth0Config: ClientExport.oAuthConfig;
  let oauth: { WebAuth: jest.Mock<WebAuth> };

  beforeEach(() => {
    testAuth0Config = {
      audience: "testAudience",
      clientId: "testClientId",
      connection: "testConnection",
      domain: "testDomain",
      redirectUri: "http://testRedirectUri",
      responseType: "token",
      scope: "profile"
    };
  });

  describe("methods and properties", () => {
    const testToken = "testToken";

    beforeEach(() => {
      jest.resetAllMocks();
      oauth = require("auth0-js");
    });

    it("removes invalid params from checkSession call", async () => {
      const mock = new oauth.WebAuth();

      (<jest.Mock>mock.checkSession).mockImplementation((configuration, cb) =>
        cb(null, { accessToken: testToken, identityToken: testToken })
      );
      expect.hasAssertions();
      await Auth0Strategy.reAuth(testAuth0Config).then(() => {
        const configuration = { ...testAuth0Config };
        expect(mock.checkSession).toHaveBeenLastCalledWith(
          configuration,
          expect.any(Function)
        );
      });
    });

    it("should return an error if getting a new token was unsuccessful", async () => {
      const mock = new oauth.WebAuth();
      const auth0Error = new AuthError("invalid_request", "Test Error Message");
      (<jest.Mock>mock.checkSession).mockImplementation((configuration, cb) =>
        cb({ error: "invalid_request", errorDescription: "Test Error Message" }, null)
      );
      expect.hasAssertions();
      await Auth0Strategy.reAuth(testAuth0Config).catch((error: AuthError) => {
        expect({ ...error }).toEqual({ ...auth0Error });
      });
    });

    it(`calls login if checkSession returns 'login_required'`, async () => {
      const mock = new oauth.WebAuth();

      const auth0Error = new AuthError("login_required", "invalid session");
      (<jest.Mock>mock.checkSession).mockImplementation((configuration, cb) =>
        cb({ error: "login_required", errorDescription: "invalid session" }, null)
      );

      expect.hasAssertions();
      await Auth0Strategy.reAuth(testAuth0Config).catch((error: AuthError) => {

        expect({ ...error }).toEqual({ ...auth0Error });
        expect(mock.authorize).toHaveBeenCalledTimes(1);
      });
    });

    it("parses hash if no access token, rejecting if parseHash throws", async () => {
      const mock = new oauth.WebAuth();

      const auth0Error = new AuthError("login_required", "invalid session");

      expect.hasAssertions();
      (<jest.Mock>mock.parseHash).mockImplementation(cb => {
        cb({ error: "login_required", errorDescription: "invalid session" });
      });
      await Auth0Strategy.ensureAuthed(testAuth0Config)
        .then(() => {
          throw "this should never happen";
        })
        .catch((e: AuthError) => {
          expect({ ...e }).toEqual({ ...auth0Error });
        });
    });

    it("parses hash if no access token, rejecting if parseHash returns no accessToken", async () => {
      const mock = new oauth.WebAuth();

      expect.hasAssertions();
      (<jest.Mock>mock.parseHash).mockImplementation(cb => {
        cb(null, null);
      });
      await Auth0Strategy.ensureAuthed(testAuth0Config)
        .then(() => {
          throw "this should never happen";
        })
        .catch(e => {
          expect(e).not.toBe("this should never happen");
          expect(mock.authorize).toHaveBeenCalledTimes(1);
        });
    });

    it("parses hash if no access token, resolving if parseHash returns accessToken", async () => {
      const mock = new oauth.WebAuth();

      expect.hasAssertions();

      (<jest.Mock>mock.parseHash).mockImplementation(cb => {
        cb(null, { accessToken: "shiny new token" });
      });
      await Auth0Strategy.ensureAuthed(testAuth0Config).then(data => {
        expect(data.accessToken).toBe("shiny new token");
      });
    });

    it("redirects user back to previous url upon redirecting back to app", async () => {
      const mock = new oauth.WebAuth();

      expect.hasAssertions();

      (<jest.Mock>mock.parseHash).mockImplementation(cb => {
        cb(null, { accessToken: "shiny new token", state: btoa("/somewhere") });
      });

      await Auth0Strategy.ensureAuthed(testAuth0Config).then(data => {
        expect(data.accessToken).toBe("shiny new token");
        expect(window.location.hash).toBe("#/somewhere");
      });
    });
  });
});
