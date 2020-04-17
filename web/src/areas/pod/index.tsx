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
import React, { Component } from 'react';
import { connect } from 'react-redux';
import LogSocket, { ILogSocket } from '../../io/logSocket';
import PodPage from './view';
import { Log, PodDetail } from '../../types';
import _ from 'lodash';
import { RouteComponentProps, withRouter } from 'react-router';
import { setSelectedAppName } from '../../actions/apps';
import { getPod, setSelectedContainerName, clearPod } from '../../actions/pods';
import { getLogs, toggleLogStream } from '../../actions/logs';
import config from '../../config';
import { AuthClient } from '../../auth/authClient';
import { closeErrorModal } from '../../actions/error';
import APIErrorModal from '../../components/error-modal';
import { IErrorState } from '../../reducers/error';

type initialState = {
  envModalOpen: boolean,
  specModalOpen: boolean,
  statusModalOpen: boolean,
  socket?: ILogSocket,
  logStream?: string,
  streamEnabled: boolean,
  containerNameSelectOpen: boolean
};

export type PodProps = {
  cluster: string,
  identityToken: string,
  pod: PodDetail,
  logs: Log,
  envBody: {},
  selectedAppName: string,
  hasLogAccess: boolean,
  authClient(): AuthClient,
  setSelectedAppName(value: string): void,
  getPod(podName: string, queryString: string): void,
  getLogs(podName: string, queryString: string): void,
  toggleLogStream(enabled: boolean): void,
  setSelectedContainerName(value: string): void,
  selectedContainerName?: string,
  error: IErrorState
} | RouteComponentProps<{
  appName: string,
  podName: string
}>;

class Pod extends Component<any, initialState> {
  // flag to determine if state can be updated in
  // an async function
  _isMounted = false;
  state: initialState = {
    envModalOpen: false,
    specModalOpen: false,
    statusModalOpen: false,
    socket: undefined,
    logStream: undefined,
    streamEnabled: false,
    containerNameSelectOpen: false
  }

  async componentDidMount() {
    this._isMounted = true;
    const { match: { params } } = this.props;
    const podname = params.podName.substring(0, params.podName.indexOf("?"));
    const search = params.podName.substring(params.podName.indexOf("?"));

    if (_.isEmpty(this.props.selectedAppName)) {
      this.props.setSelectedAppName(params.appName);
    }

    if (_.isEmpty(this.props.selectedContainerName) && !_.isEmpty(this.props.pod)) {
      this.props.setSelectedContainerName(this.props.pod.containerNames[0]);
    }

    if (_.isEmpty(this.props.pod)) {
      this.props.getPod(podname, search, this.props.cluster, this.props.identityToken);
    }
  }

  componentDidUpdate(prevProps) {
    const { match: { params } } = this.props;
    const search = params.podName.substring(params.podName.indexOf("?"))
    const char = _.isEmpty(search) ? "" : "&"

    // will be on first load, set and return to allow for update of props.
    if (_.isEmpty(this.props.selectedContainerName) && !_.isEmpty(this.props.pod) && (!_.isEmpty(this.props.pod.containerNames))) {
      this.props.setSelectedContainerName(this.props.pod.containerNames[0]);
      return;
    }

    // check if props updated and that it's not a fresh load
    if (prevProps.selectedContainerName !== this.props.selectedContainerName && !_.isEmpty(this.props.selectedContainerName)) {
      this.props.getLogs(this.props.pod.name, `${search}${char}containerName=${this.props.selectedContainerName}`, this.props.cluster, this.props.identityToken);
    }
  }

  componentWillUnmount() {
    this._isMounted = false;
    this.state.socket && this.state.socket.Close();
    this.props.toggleLogStream(false);
    this.props.setSelectedContainerName("");
    this.props.clearPod();
  }

  logStreamHandler = (stream: MessageEvent) => {
    this.setState({ logStream: `${this.state.logStream}\r${stream.data}` });
  }

  openLogStream = async () => {
    const cfg = await config();

    let endpoint = '';
    _.forEach(cfg.availableClusters, value => {
      if (!_.isEmpty(value[this.props.cluster])) {
        endpoint = `${value[this.props.cluster].replace('http', 'ws')}/io`;
      }
    })

    const socket = new LogSocket({
      cluster: this.props.cluster,
      podname: this.props.pod.name,
      containerName: this.props.selectedContainerName,
      namespace: this.props.pod.namespace,
      handler: this.logStreamHandler,
      wsBase: endpoint,
      accessToken: this.props.identityToken
    });

    const now = new Date().toLocaleString();

    this.props.toggleLogStream(true);

    this.setState({
      socket: socket,
      logStream: `\n\nStream Started ${now}\n\n`,
      streamEnabled: true
    });
  };

  closeLogStream = () => {
    this.state.socket && this.state.socket.Close();

    const now = new Date().toLocaleString();

    this.props.toggleLogStream(false);

    this.setState({
      socket: undefined,
      logStream: `\n\n${this.state.logStream}\n\nStream Ended ${now}\n\n`,
      streamEnabled: false
    });
  }

  toggleModalType = (type) => {
    switch (type) {
      case 'env':
        this.setState({ envModalOpen: !this.state.envModalOpen });
        break;
      case 'spec':
        this.setState({ specModalOpen: !this.state.specModalOpen });
        break;
      case 'status':
        this.setState({ statusModalOpen: !this.state.statusModalOpen });
        break;

      default:
        break;
    }
  }

  toggleContainerNameSelect = () => {
    this.setState({ containerNameSelectOpen: !this.state.containerNameSelectOpen });
  }

  render() {
    return (
      <div>
        <PodPage
          podDetail={this.props.pod}
          envBody={this.props.envBody}
          logs={this.props.logs}
          showEnvModal={this.state.envModalOpen}
          showSpecModal={this.state.specModalOpen}
          showStatusModal={this.state.statusModalOpen}
          toggleModalType={this.toggleModalType}
          openLogStream={this.openLogStream}
          closeLogStream={this.closeLogStream}
          logStream={this.state.logStream}
          streamEnabled={this.state.streamEnabled}
          hasLogAccess={this.props.hasLogAccess}
          toggleContainerNameSelect={this.toggleContainerNameSelect}
          containerNameSelectOpen={this.state.containerNameSelectOpen}
          setSelectedContainerName={this.props.setSelectedContainerName}
          selectedContainerName={this.props.selectedContainerName} />

        <APIErrorModal
          open={this.props.error.apiOpen}
          handleClose={this.props.closeErrorModal}
          status={this.props.error.status}
          statusText={this.props.error.statusText}
          message={this.props.error.message} />
      </div>
    )
  }
}

export const mapStateToProps = ({ appsState, podsState, logsState, authState, clustersState, errorState }) => {
  // api returns forbidden if the user doesn't have access to view logs for the pod. 
  // check and set the message.
  const hasLogAccess = logsState.logsError && logsState.logsError.status === 403 ? false : true;

  let envBody;
  if (!_.isEmpty(podsState.pod)) {
    const c = _.find(podsState.pod && podsState.pod.spec && podsState.pod.spec.containers, 'env');
    envBody = !_.isEmpty(c) ? c && c.env : '';
  }

  let identityToken = '';
  if (authState.identityToken) {
    identityToken = authState.identityToken;
  }

  return {
    cluster: clustersState.cluster,
    identityToken: identityToken,
    pod: podsState.pod,
    logs: logsState.logs,
    envBody: envBody,
    selectedAppName: appsState.selectedAppName,
    hasLogAccess: hasLogAccess,
    selectedContainerName: podsState.selectedContainerName,
    error: errorState
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    setSelectedAppName: (value: string) => dispatch(setSelectedAppName(value)),
    getPod: (podName: string, queryString: string, cluster: string, jwt: string) => dispatch(getPod(podName, queryString, cluster, jwt)),
    getLogs: (podName: string, queryString: string, cluster: string, jwt: string) => dispatch(getLogs(podName, queryString, cluster, jwt)),
    toggleLogStream: (enabled: boolean) => dispatch(toggleLogStream(enabled)),
    setSelectedContainerName: (selectedContainerName: string) => dispatch(setSelectedContainerName(selectedContainerName)),
    clearPod: () => dispatch(clearPod()),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Pod));
