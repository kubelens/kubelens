import uuid from 'uuid/v4';

export const queryStringify = (queryParams: any): string => {
  const encode = encodeURIComponent;

  return Object.keys(queryParams)
    .map(key => {
      if (queryParams[key] === undefined) {
        return encode(key);
      }
      if (Array.isArray(queryParams[key])) {
        return queryParams[key]
          .map((value: any) => `${encode(key)}=${encode(value)}`)
          .join('&');
      }
      return `${encode(key)}=${encode(queryParams[key])}`;
    })
    .join('&');
};

export const setStateSession = (path: string) => {
  const sessionKey = uuid();
  window.sessionStorage.setItem(sessionKey, path);

  return sessionKey;
}

export const getStateSession = (key: string) => {
  return window.sessionStorage.getItem(key);
}

export const removeStateSession = (key: string) => {
  return window.sessionStorage.removeItem(key);
}

export const parseState = (url: string)=> {
  let stateValue: string = "";
  let stateIndex = url.indexOf('state');
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

export const getRedirectUri=()=> {
  let baseUri = document.baseURI;
  if (!baseUri) {
    let baseNode = document.querySelector('base');
    if (baseNode) {
      baseUri = baseNode.getAttribute('href') || "";
    }
  }

  let result = baseUri || `${window.location.protocol}//${window.location.host}${window.location.pathname}`;

  const hashIndex = result.indexOf('#');
  if (hashIndex !== -1) {
    result = result.substring(0, hashIndex);
  }

  return result.toLowerCase();
}
