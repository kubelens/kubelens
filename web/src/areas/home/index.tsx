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
import { withRouter, RouteComponentProps } from 'react-router';
import NavBar from '../../components/nav';
import ErrorModal from '../../components/error-modal';
import View from './view';
import { App } from "../../types";
import _ from 'lodash';
import { IGlobalState } from '../../store';
import { getApps, getAppOverview, setSelectedAppName, filterApps, clearAppsErrors } from '../../actions/apps';

type initialState = {};

export type HomeProps = {
  cluster: string,
  identityToken?: string,
  isError: boolean,
  appsError: Error,
  appsRequested: boolean,
  selectedAppName: string,
  filteredApps: App[],
  apps: App[],
  getApps(cluster: string, jwt: string): void,
  getAppOverview(appname: string, queryString: string, cluster: string, jwt: string): void,
  setSelectedAppName(value: string): void,
  filterApps(value: string, apps: App[]): void,
  clearAppsErrors(): void
} | RouteComponentProps<{
  appName?: string
}>;

class Home extends Component<HomeProps, initialState> {
  constructor(props) {
    super(props);

    this.closeErrorModal = this.closeErrorModal.bind(this);
    this.onFilterChanged = this.onFilterChanged.bind(this);
    this.onViewApp = this.onViewApp.bind(this);
  }

  public componentDidMount() {
    const { match: { params } } = this.props;

    if (_.isEmpty(this.props.selectedAppName) && !_.isEmpty(params.appName)) {
      this.props.setSelectedAppName(params.appName);
    }

    if (!this.props.appsRequested) {
      this.props.getApps(this.props.cluster, this.props.identityToken);
    }
  }

  public shouldComponentUpdate(nextProps: HomeProps) {
    if (!this.props.selectedAppName
      || nextProps.selectedAppName !== this.props.selectedAppName
      || nextProps.appsRequested !== this.props.appsRequested
      || !_.isEqual(nextProps.filteredApps, this.props.filteredApps)) {
      return true;
    }
    return false;
  }

  private closeErrorModal() {
    this.props.clearAppsErrors();
  }

  private onFilterChanged(event) {
    this.props.filterApps(event.target.value, this.props.apps);
  }

  private onViewApp(appname: string, labelKey: string) {
    this.props.setSelectedAppName(appname);

    this.props.getAppOverview(appname, labelKey, this.props.cluster, this.props.identityToken);

    this.props.history.push(`/${appname}?labelKey=${labelKey}`);
  }

  public render() {
    return (
      <div>
        <NavBar {...this.props} />
        <View
          onFilterChanged={this.onFilterChanged}
          onViewApp={this.onViewApp}
          {...this.props} />
        <ErrorModal
          show={this.props.isError}
          handleClose={this.closeErrorModal}
          error={this.props.appsError} />

      </div>
    );
  }
}

export const mapStateToProps = ({ appsState, authState, clustersState }: IGlobalState) => {
  const isError = !_.isEmpty(appsState.appsError) ? true : false;

  return {
    cluster: clustersState.cluster,
    identityToken: authState.identityToken,
    apps: appsState.apps,
    filteredApps: appsState.filteredApps || appsState.apps,
    appsError: appsState.appsError,
    appsRequested: appsState.appsRequested,
    isError: isError,
    selectedAppName: appsState.selectedAppName
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getApps: (cluster: string, jwt: string) => dispatch(getApps(cluster, jwt)),
    getAppOverview: (appname: string, queryString: string, cluster: string, jwt: string) => dispatch(getAppOverview(appname, queryString, cluster, jwt)),
    setSelectedAppName: (value: string) => dispatch(setSelectedAppName(value)),
    filterApps: (value: string, apps: App[]) => dispatch(filterApps(value, apps)),
    clearAppsErrors: () => dispatch(clearAppsErrors())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Home));
