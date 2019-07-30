jest.mock("@okta/okta-auth-js", () => {
  const getWithRedirectMock = jest.fn();
  const renewMock = jest.fn();
  const parseUrlMock = jest.fn();

  return (
    class {
      options = {
        redirectUri: "state=%2f"
      }
      token = {
        getWithRedirect: getWithRedirectMock,
        renew: renewMock,
        parseFromUrl: parseUrlMock,
      }
    }
  )
});

import * as OktaStrategy from "./okta-strategy";
import * as ClientExport from "../clients/auth-implicit-client";
import OktaAuth from "@okta/okta-auth-js";
import { AuthError } from '../AuthError';
import { ok } from "assert";
describe("OktaStrategy", () => {
  let testAuth0Config: ClientExport.oAuthConfig;
  let oauth: jest.Mock<OktaAuth>;
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
    const accessTokenObject = 'AwesomeTokenObject';
    const idTokenObject = 'IcyTokenObject';
    beforeEach(() => {
      jest.resetAllMocks();
      oauth = require("@okta/okta-auth-js");
    });

    it("removes invalid params from renew call", async () => {
      let mock = new oauth();

      const accessTokenObject = 'AwesomeTokenObject';
      const idTokenObject = 'IcyTokenObject';
      (<jest.Mock>mock.token.renew)
        .mockResolvedValueOnce((accessTokenObject) => { accessToken: testToken })
        .mockResolvedValueOnce((idTokenObject) => { identityToken: testToken })

      expect.hasAssertions();
      await OktaStrategy.reAuth(testAuth0Config).then(() => {
        expect(mock.token.renew).toHaveBeenCalledTimes(2);
      });
    });

    it("should return an error if getting a new token was unsuccessful", async () => {
      let mock = new oauth();
      const oktaError = new AuthError("INTERNAL", "No token");

      (<jest.Mock>mock.token.renew)
        .mockResolvedValueOnce({ accessToken: testToken })
        .mockRejectedValueOnce(new AuthError("INTERNAL", "No token"));

      expect.hasAssertions();
      await OktaStrategy.reAuth(testAuth0Config).catch(error => {
        expect({ ...error }).toEqual({ ...oktaError });
        expect(mock.token.getWithRedirect).toHaveBeenCalledTimes(1);
      });
    });

    it(`calls login if renew returns no token objects`, async () => {
      let mock = new oauth();

      const oktaError = new AuthError("INTERNAL", "Inactive Session");
      (<jest.Mock>mock.token.renew)
        .mockResolvedValueOnce({ accessToken: testToken })
        .mockRejectedValueOnce(oktaError)

      expect.hasAssertions();
      await OktaStrategy.reAuth(testAuth0Config).catch(error => {
        expect({ ...error }).toEqual({ ...oktaError });
        expect(mock.token.getWithRedirect).toHaveBeenCalledTimes(1);
      });
    });

    it("parses hash if no access token, rejecting if parseFromUrl throws", async () => {
      let mock = new oauth();
      const oktaError = new AuthError("INTERNAL", "Unable to parse url");
      expect.hasAssertions();
      (<jest.Mock>mock.token.parseFromUrl).mockRejectedValue(new AuthError("INTERNAL", "Unable to parse url"));
      await OktaStrategy.ensureAuthed(testAuth0Config)
        .then(() => {
          throw "this should never happen";
        })
        .catch(e => {
          expect({ ...e }).toEqual({ ...oktaError });
        });
    });

    it("parses hash if no access token, rejecting if parseFromUrl returns no accessToken", async () => {
      let mock = new oauth();
      const oktaError = new AuthError("INTERNAL", "Unable to parse url")
      expect.hasAssertions();
      (<jest.Mock>mock.token.parseFromUrl).mockRejectedValue(oktaError);
      await OktaStrategy.ensureAuthed(testAuth0Config)
        .then(() => {
          throw "this should never happen";
        })
        .catch(e => {
          expect(e).not.toBe("this should never happen");
          expect(mock.token.getWithRedirect).toHaveBeenCalledTimes(1);
        });
    });

    it("parses hash if no access token, preserve state and resolve if parseHash returns accessToken", async () => {
      let mock = new oauth();

      expect.hasAssertions();

      (<jest.Mock>mock.token.parseFromUrl).mockResolvedValue([{ accessToken: "a golden accessToken" }, { idToken: "a silver idtoken", claims: "I claim you" }]);
      await OktaStrategy.ensureAuthed(testAuth0Config).then(data => {
        expect(data.accessToken).toBe("a golden accessToken");
        expect(data.identityToken).toBe("a silver idtoken");
        expect(data.identity).toBe("I claim you");
        expect(window.location.hash).toBe("#/");
      });
    });

    it("parses state from okta redirectUri, returns state value if state property ends with an ampersand", async () => {

      expect.hasAssertions();
      let testRedirectUri = 'test.com/#/state=%2F&idtoken=awesome';
      var state = OktaStrategy.parseState(testRedirectUri);

      expect(state).toBe('%2F');
    });

    it("parses state from okta redirectUri, returns state value if state property is end of string", async () => {

      expect.hasAssertions();
      let testRedirectUri = 'test.com/#/idtoken=awesome&state=%2F';
      var state = OktaStrategy.parseState(testRedirectUri);

      expect(state).toBe('%2F');
    });
  });
});
