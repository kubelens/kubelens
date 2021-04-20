/*
MIT License

Copyright (c) 2020 The KubeLens Authors

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
  oAuthJwtIssuer:     string,
  oAuthAudience:      string,
  oAuthClientId:      string,
  oAuthRedirectUri:   string,
  oAuthResponseType:  string,
  oAuthRequestType:   string,
  oAuthScope:         string,
  oAuthConnection:    string,
  oAuthEnabled:       boolean,
  availableClusters:  AvailableCluster[]
};

export type AvailableCluster = {
  name:     string,
  cluster:  string,
  url:      string
}

export type SelectedOverview = {
  linkedName:  string,
	namespace:   string
}

export type Overview = {
  linkedName:  string,
	namespace:   string,
	daemonSets:  DaemonSet[],
	deployments: Deployment[],
	jobs:        Job[],
	pods:        Pod[],
	replicaSets: ReplicaSet[],
	services:    Service[],
  configMaps:  ConfigMap[]
}

export type DaemonSet = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	daemonSet:    any
}

export type Deployment = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	deployment:   any
}

export type Job = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	job:          any
}

export type Pod = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	pod:          any
}

export type ReplicaSet = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	replicaSet:   any
}

export type Service = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	service:      any
}

export type Log = {
  pod:    string,
  output: string
}

export type ConfigMap = {
  name:         string,
  linkedName:   string,
	namespace:    string,
	service:      any
}
