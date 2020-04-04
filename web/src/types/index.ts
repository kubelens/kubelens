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
export type Config = {
  oAuthJwtIssuer: string,
  oAuthAudience: string,
  oAuthClientId: string,
  oAuthRedirectUri: string,
  oAuthResponseType: string,
  oAuthRequestType: string,
  oAuthScope: string,
  oAuthConnection: string,
  oAuthEnabled: boolean,
  availableClusters: {}[],
  deployerLinkName: string
};

export type Name = {
  labelKey: string,
  value: string
};

export type Log = {
  pod: string,
  output: string
};

export type App = {
  name: string,
  labelKey: string,
  namespaces: string[],
  deployerLink?: string
};

export type PodDetail = {
  name: string,
  namespace: string,
  hostIP?: string,
  podIP?: string,
  startTime?: string,
  phase?: {},
  phaseMessage: string,
  containerStatus: [],
  status: {
    containerStatuses: [{
      image: string,
      restartCount: number,
      ready: boolean
    }],
    conditions: [],
    startTime: string,
    phase: string
  },
  spec: {
    containers: [{
      env: object
    }]
  },
  containerNames: string[]
};

export type PodOverview = {
  name: Name,
  namespace: string,
  clusterName?: string,
  deployerLink?: string,
  podDetails: PodDetail[]
}

export type Service = {
  appName: Name,
  name: string,
  namespace: string,
  type: {},
  selector: Map<string, string>,
  clusterIP?: string,
  loadBalancerIP?: string,
  externalIPs?: string[],
  ports?: [{}],
  spec?: {},
  status?: {},
  deployerLink?: string,
  configMaps?: [{}]
};

export type AppOverview = {
  serviceOverviews: Service[],
  podOverviews: PodOverview
}
