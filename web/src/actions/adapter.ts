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
import axios, { AxiosResponse } from 'axios';
import config from '../config';
import _ from 'lodash';

interface RequestConfig {
  headers: {}
}

const buildRequestConfig = (rqtype: string, jwt: string, headers: {}): RequestConfig => {
  let h = {
    "Authorization": `${rqtype} ${jwt}`
  };

  if (headers) {
    Object.assign(h, headers);
  }

  return {
    headers: h
  };
};

const get = async (path: string, cluster: string, jwt: string): Promise<AxiosResponse<any>> => {
  try {
    const cfg = await config();

    const rc = buildRequestConfig(cfg.oAuthRequestType, jwt, {});

    // add trailing slash if missing
    if (cluster.charAt(cluster.length - 1) != '/') {
      cluster += '/';
    }

    return await axios.get(`${cluster}${path}`, rc);
  } catch (err) {
    return Promise.reject(err);
  }
};

export default {
  get
}
