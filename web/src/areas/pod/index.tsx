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
import { Log, Pod, SelectedOverview } from '../../types';
import _ from 'lodash';
import { RouteComponentProps, withRouter } from 'react-router';
import { getPod, setSelectedContainerName, clearPod } from '../../actions/pods';
import { getLogs, toggleLogStream } from '../../actions/logs';
import { AuthClient } from '../../auth/authClient';
import { closeErrorModal } from '../../actions/error';
import APIErrorModal from '../../components/error-modal';
import { IErrorState } from '../../reducers/error';
import qs from 'qs';

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
  podOverview: Pod,
  logs: Log,
  envBody: {},
  selectedOverview: string,
  hasLogAccess: boolean,
  authClient(): AuthClient,
  setSelectedOverview(value: SelectedOverview): void,
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

export class PodView extends Component<any, initialState> {
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
    const { match: { params }, location } = this.props;
    const podName = params.podName.substring(0, params.podName.indexOf("?"));
    const search = location.pathname.substring(location.pathname.indexOf("?")+1);
    const query = qs.parse(search);

    if (_.isEmpty(this.props.selectedOverview)) {
      this.props.setSelectedOverview({linkedName: params.linkedName, namespace: query.namespace});
    }

    if (_.isEmpty(this.props.selectedContainerName) && !_.isEmpty(this.props.podOverview)) {
      this.props.setSelectedContainerName(this.props.podOverview.pod.spec.containers[0].name);
    }

    if (_.isEmpty(this.props.podOverview)) {
      this.props.getPod(podName, query.namespace, this.props.cluster, this.props.identityToken);
    }
  }

  componentDidUpdate(prevProps) {
    const { match: { params }, location } = this.props;
    const podName = params.podName.substring(0, params.podName.indexOf("?"));
    const search = location.pathname.substring(location.pathname.indexOf("?")+1);
    const query = qs.parse(search);

    // will be on first load, set and return to allow for update of props.
    if (_.isEmpty(this.props.selectedContainerName) && !_.isEmpty(this.props.podOverview)) {
      this.props.setSelectedContainerName(this.props.podOverview.pod.spec.containers[0].name);
      return;
    }

    // check if props updated and that it's not a fresh load
    if (prevProps.selectedContainerName !== this.props.selectedContainerName && !_.isEmpty(this.props.selectedContainerName)) {
      this.props.getLogs(podName, query.namespace, this.props.selectedContainerName, query.tail, this.props.cluster, this.props.identityToken);
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
    const socket = new LogSocket({
      cluster: `${this.props.cluster.replace('http', 'ws')}/io`,
      podname: this.props.podOverview.pod.metadata.name,
      containerName: this.props.selectedContainerName,
      namespace: this.props.podOverview.pod.metadata.namespace,
      handler: this.logStreamHandler,
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
          podDetail={this.props.podOverview}
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

export const mapStateToProps = ({ overviewsState, podsState, logsState, authState, clustersState, errorState }) => {
  // api returns forbidden if the user doesn't have access to view logs for the pod. 
  // check and set the message.
  const hasLogAccess = logsState.logsError && logsState.logsError.status === 403 ? false : true;

  let envBody;
  if (!_.isEmpty(podsState.podOverview) && !_.isEmpty(podsState.podOverview.pod)) {
    envBody = podsState.podOverview.pod.spec.containers[0].env;
  }

  let identityToken = '';
  if (authState.identityToken) {
    identityToken = authState.identityToken;
  }

  return {
    cluster: clustersState.cluster && clustersState.cluster.url,
    identityToken: identityToken,
    podOverview: podsState.podOverview,
    logs: logsState.logs,
    envBody: envBody,
    selectedCluster: overviewsState.selectedCluster,
    hasLogAccess: hasLogAccess,
    selectedContainerName: podsState.selectedContainerName,
    error: errorState
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getPod: (podName: string, namespace: string, cluster: string, jwt: string) => dispatch(getPod(podName, namespace, cluster, jwt)),
    getLogs: (podName: string, namespace: string, containerName: string, tail: number, cluster: string, jwt: string) => dispatch(getLogs(podName, namespace, containerName, tail, cluster, jwt)),
    toggleLogStream: (enabled: boolean) => dispatch(toggleLogStream(enabled)),
    setSelectedContainerName: (selectedContainerName: string) => dispatch(setSelectedContainerName(selectedContainerName)),
    clearPod: () => dispatch(clearPod()),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(PodView));
