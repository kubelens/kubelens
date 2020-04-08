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
export type LogSocketProps = {
  cluster: string,
  podname: string,
  containerName?: string,
  namespace: string,
  handler(ev: MessageEvent),
  wsBase: string,
  accessToken: string,
  socket?: WebSocket
};

export interface ILogSocket {
  Close()
}

export default class LogSocket {
  private socket: WebSocket;

  constructor(props: LogSocketProps) {

    if ("WebSocket" in window) {
      const keyQuery = props.accessToken ? `&key=${props.accessToken}` : "";
      const containerNameQuery = props.containerName ? `&containerName=${props.containerName}` : "";
      this.socket = props.socket
        ? props.socket
        : new WebSocket(`${props.wsBase}/${props.podname}/logs?namespace=${props.namespace}${containerNameQuery}${keyQuery}`);

      this.socket.onclose = () => {
        console.log('socket closed');
      }

      this.socket.onmessage = props.handler;

    } else {
      const err = { message: "WebSocket NOT supported by your Browser." };
      throw err;
    }
  }

  public Close = () => {
    if (this.socket) {
      this.socket.close();
    }
  }
}
