import { queryStringify, setStateSession, getStateSession, removeStateSession, parseState, getRedirectUri } from './utils';
jest.mock('uuid/v4', () => () => '3f56b96a-46e9-4d3a-9b15-898b04350a29');

describe('utils', () => {
  beforeEach(() => {
    (window as any).sessionStorage = (() => {
      const cache = {};
      return {
        getItem: jest.fn().mockImplementation((key: string) => { return cache[key]; }),
        setItem: jest.fn().mockImplementation((key: string, value: string) => { cache[key] = value; }),
        removeItem: jest.fn().mockImplementation((key: string) => { delete cache[key]; })
      };
    })();
  });
  describe('queryStringify', () => {
    it('returns empty string when no params are given', () => {
      expect(queryStringify({})).toBe('');
    });

    it('stringifies a single parameter', () => {
      expect(queryStringify({ testing: 'hello!' })).toBe('testing=hello!');
    });

    it('uri encodes characters', () => {
      expect(queryStringify({ testing: '🚀' })).toBe('testing=%F0%9F%9A%80');
    });

    it('skips a parameter value when it is undefined', () => {
      expect(queryStringify({ testing: undefined })).toBe('testing');
    });

    it('stringifies arrays', () => {
      expect(
        queryStringify({
          testing: ['deeply', 'nested', 'and complex  parameters']
        })
      ).toBe(
        'testing=deeply&testing=nested&testing=and%20complex%20%20parameters'
      );
    });

    it('stringifies arrays, special characters, and multiple values', () => {
      expect(
        queryStringify({
          testing: ['deeply', 'nested', 'and complex  parameters'],
          emoji: '🐙',
          'some-undefined-param': undefined
        })
      ).toBe(
        'testing=deeply&testing=nested&testing=and%20complex%20%20parameters&emoji=%F0%9F%90%99&some-undefined-param'
      );
    });
  });
  describe('setStateSession', () => {
    it('sets session and returns key', () => {
      const uuid = require('uuid/v4');
      const path = '/path/12345';
      expect(setStateSession(path)).toBe(uuid());
    })
  });
  describe('getStateSession', () => {
    it('returns key when given a uuid', () => {
      const uuid = require('uuid/v4');
      const path = '/path/12345';
      const key = setStateSession(path);
      expect(sessionStorage.setItem).toHaveBeenCalledWith(key, path);
      expect(getStateSession(key)).toBe(path)
    })
  });
  describe('removeStateSession', () => {
    it('removes session', () => {
      const uuid = require('uuid/v4');
      const path = '/path/12345';
      const key = setStateSession(path);
      removeStateSession(key);
      expect(getStateSession(key)).toBe(undefined);
    })
  });
  describe('parseState', () => {
    it('parses state from okta redirectUri, returns state value if state property ends with an ampersand', async () => {
      expect.hasAssertions();
      const testRedirectUri = 'test.com/#/state=%2F&idtoken=awesome';
      const state = parseState(testRedirectUri);

      expect(state).toBe('%2F');
    });

    it('parses state from okta redirectUri, returns state value if state property is end of string', async () => {
      expect.hasAssertions();
      const testRedirectUri = 'test.com/#/idtoken=awesome&state=%2F';
      const state = parseState(testRedirectUri);

      expect(state).toBe('%2F');
    });
  });
  describe('getRedirectUri', () => {
    beforeEach(() => {
      jest.resetAllMocks();
    })
    it('Should get href', async () => {
      const spyFunc = jest.fn();
      const globalAny: any = global;
      Object.defineProperty(globalAny.document, 'baseURI', { value: "" });
      Object.defineProperty(globalAny.document, 'querySelector', {
        value: () => {
          return {
            getAttribute: () => "test"
          }
        },
        configurable:true
      });
      let uri = getRedirectUri();
      expect(uri).toBe("test")
    });
    it('Should default to window location', async () => {
      const spyFunc = jest.fn();
      const globalAny: any = global;
      Object.defineProperty(globalAny.document, 'baseURI', { value: "" });
      Object.defineProperty(globalAny.document, 'querySelector', {
        value: () => {
          null
        },
        configurable: true
      });
      let uri = getRedirectUri();
      expect(uri).toBe("http://localhost/")
    });
    it('Should remove after hash', async () => {
      const spyFunc = jest.fn();
      const globalAny: any = global;
      Object.defineProperty(globalAny.document, 'baseURI', { value: "" });
      Object.defineProperty(globalAny.document, 'querySelector', {
        value: () => {
          return {
            getAttribute: () => "test.com#helooooo"
          }
        },
        configurable: true
      });
      let uri = getRedirectUri();
      expect(uri).toBe("test.com")
    });
    it('parses state from okta redirectUri, returns state value if state property is end of string', async () => {
      expect.hasAssertions();
      const testRedirectUri = 'test.com/#/idtoken=awesome&state=%2F';
      const state = parseState(testRedirectUri);

      expect(state).toBe('%2F');
    });
  });
});
