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
import ServiceOverviewPage from './service-view';
import PodOverviewPage from './pod-view';
import DaemonSetOverviewPage from './daemonset-view';
import { Service, PodOverview } from "../../types";
import { RouteComponentProps, withRouter } from 'react-router';
import _ from 'lodash';
import { IGlobalState } from '../../store';
import { getAppOverview, setSelectedAppName } from '../../actions/apps';
import APIErrorModal from '../../components/error-modal';
import { closeErrorModal } from '../../actions/error';
import { IErrorState } from '../../reducers/error';

type initialState = {
  specModalOpen: boolean,
  statusModalOpen: boolean,
  configMapModalOpen: boolean,
  deploymentModalOpen: boolean,
  dsConfigMapModalOpen: boolean,
  dsDeploymentModalOpen: boolean
};

export type OverviewProps = {
  identityToken?: string,
  serviceOverviews: Service[],
  podOverview: PodOverview,
  selectedAppName: string,
  appOverviewRequested: boolean,
  getAppOverview(appname: string, queryString: string): void,
  setSelectedAppName(value: string): void,
  error: IErrorState
} | RouteComponentProps<{
  appName?: string
}>;

class Overview extends Component<OverviewProps, initialState> {
  state: initialState = {
    specModalOpen: false,
    statusModalOpen: false,
    configMapModalOpen: false,
    deploymentModalOpen: false,
    dsConfigMapModalOpen: false,
    dsDeploymentModalOpen: false
  }

  async componentDidMount() {
    const { match: { params }, location: { search } } = this.props;
    const query = new URLSearchParams(search);

    let queryString = query.get('labelKey');
    if (!queryString) {
      queryString = 'app';
    }

    let appname = '';
    if (params.appName) {
      appname = params.appName;
    }

    if (_.isEmpty(this.props.selectedAppName)) {
      this.props.setSelectedAppName(appname);
    }

    if (!this.props.appOverviewRequested
      && (_.isEmpty(this.props.podOverview)
        || _.isEmpty(this.props.serviceOverviews
          || params.appName !== this.props.selectedAppName))) {
      this.props.appActions &&
        this.props.getAppOverview(
          appname,
          queryString,
          this.props.cluster,
          this.props.identityToken
        );
    }
  }

  toggleModalType = (type) => {
    switch (type) {
      case 'spec':
        this.setState({ specModalOpen: !this.state.specModalOpen });
        break;
      case 'status':
        this.setState({ statusModalOpen: !this.state.statusModalOpen });
        break;
      case 'configMap':
        this.setState({ configMapModalOpen: !this.state.configMapModalOpen});
        break;
      case 'deployment':
        this.setState({ deploymentModalOpen: !this.state.deploymentModalOpen});
        break;
      default:
        break;
    }
  }

  render() {
    return (
      <div className="overview-container">
        <ServiceOverviewPage 
          serviceOverviews={this.props.serviceOverviews} 
          toggleModalType={this.toggleModalType}
          specModalOpen={this.state.specModalOpen}
          statusModalOpen={this.state.statusModalOpen}
          configMapModalOpen={this.state.configMapModalOpen}
          deploymentModalOpen={this.state.deploymentModalOpen} />
        <DaemonSetOverviewPage 
          daemonSetOverviews={this.props.daemonSetOverviews} 
          toggleModalType={this.toggleModalType}
          conditionsModalOpen={this.state.specModalOpen}
          configMapModalOpen={this.state.configMapModalOpen}
          deploymentModalOpen={this.state.deploymentModalOpen} />
        <PodOverviewPage podOverview={this.props.podOverview} />
        <APIErrorModal
          open={this.props.error.apiOpen}
          handleClose={this.props.closeErrorModal}
          status={this.props.error.status}
          statusText={this.props.error.statusText}
          message={this.props.error.message} />

      </div>
    );
  }
}

export const mapStateToProps = ({ appsState, authState, clustersState, errorState }: IGlobalState) => {
  let serviceOverviews,
    daemonSetOverviews,
    podOverview;

  if (appsState.appOverview && appsState.appOverview.serviceOverviews) {
    serviceOverviews = appsState.appOverview.serviceOverviews;
  }

  if (appsState.appOverview && appsState.appOverview.daemonSetOverviews) {
    daemonSetOverviews = appsState.appOverview.daemonSetOverviews;
  }

  if (appsState.appOverview && appsState.appOverview.podOverviews) {
    podOverview = appsState.appOverview.podOverviews;
  }

  return {
    cluster: clustersState.cluster,
    identityToken: authState.identityToken,
    serviceOverviews: serviceOverviews,
    daemonSetOverviews: daemonSetOverviews,
    podOverview: podOverview,
    selectedAppName: appsState.selectedAppName,
    appOverviewRequested: appsState.appOverviewRequested,
    error: errorState
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getAppOverview: (appname: string, queryString: string, cluster: string, jwt: string) => dispatch(getAppOverview(appname, queryString, cluster, jwt)),
    setSelectedAppName: (value: string) => dispatch(setSelectedAppName(value)),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Overview));