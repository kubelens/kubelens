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
import React, { Component } from 'react';
import { connect } from 'react-redux';
import ErrorModal from '../../components/error-modal';
import ServiceOverviewPage from './service-view';
import PodOverviewPage from './pod-view';
import { Service, PodOverview } from "../../types";
import { RouteComponentProps, withRouter } from 'react-router';
import _ from 'lodash';
import { IGlobalState } from '../../store';
import { getAppOverview, setSelectedAppName, clearAppsErrors } from '../../actions/apps';

type initialState = {
  specModalOpen: boolean,
  statusModalOpen: boolean
};

export type OverviewProps = {
  identityToken?: string,
  serviceOverviews: Service[],
  podOverview: PodOverview,
  appOverviewError: Error,
  isError: boolean,
  selectedAppName: string,
  appOverviewRequested: boolean,
  getAppOverview(appname: string, queryString: string): void,
  setSelectedAppName(value: string): void,
  clearAppsErrors(): void
} | RouteComponentProps<{
  appName?: string
}>;

class Overview extends Component<OverviewProps, initialState> {
  state: initialState = {
    specModalOpen: false,
    statusModalOpen: false
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

  closeErrorModal = () => {
    this.props.appActions &&
      this.props.clearAppsErrors();
  }

  toggleModalType = (type) => {
    switch (type) {
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

  render() {
    return (
      <div className="overview-container">
        <ServiceOverviewPage 
          serviceOverviews={this.props.serviceOverviews} 
          toggleModalType={this.toggleModalType}
          specModalOpen={this.state.specModalOpen}
          statusModalOpen={this.state.statusModalOpen} />
        <PodOverviewPage podOverview={this.props.podOverview} />
        <ErrorModal
          show={this.props.isError}
          handleClose={this.closeErrorModal}
          error={this.props.appOverviewError} />

      </div>
    );
  }
}

export const mapStateToProps = ({ appsState, authState, clustersState }: IGlobalState) => {
  const isError = !_.isEmpty(appsState.appOverviewError) ? true : false;

  let serviceOverviews,
    podOverview;

  if (appsState.appOverview && appsState.appOverview.serviceOverviews) {
    serviceOverviews = appsState.appOverview.serviceOverviews;
  }

  if (appsState.appOverview && appsState.appOverview.podOverviews) {
    podOverview = appsState.appOverview.podOverviews;
  }

  return {
    cluster: clustersState.cluster,
    identityToken: authState.identityToken,
    serviceOverviews: serviceOverviews,
    podOverview: podOverview,
    appOverviewError: appsState.appOverviewError,
    isError: isError,
    selectedAppName: appsState.selectedAppName,
    appOverviewRequested: appsState.appOverviewRequested
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getAppOverview: (appname: string, queryString: string, cluster: string, jwt: string) => dispatch(getAppOverview(appname, queryString, cluster, jwt)),
    setSelectedAppName: (value: string) => dispatch(setSelectedAppName(value)),
    clearAppsErrors: () => dispatch(clearAppsErrors())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Overview));