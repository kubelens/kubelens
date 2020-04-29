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
import { withRouter, RouteComponentProps } from 'react-router';
import NavBar from '../../components/nav';
import View from './view';
import { App } from "../../types";
import _ from 'lodash';
import { IGlobalState } from '../../store';
import { getApps, getAppOverview, setSelectedAppName, filterApps } from '../../actions/apps';
import { closeErrorModal } from '../../actions/error';
import APIErrorModal from '../../components/error-modal';
import { IErrorState } from '../../reducers/error';

type initialState = {};

export type HomeProps = {
  cluster: string,
  identityToken?: string,
  selectedAppName: string,
  filteredApps: App[],
  apps: App[],
  getApps(cluster: string, jwt: string): void,
  getAppOverview(appname: string, namespace: string, queryString: string, cluster: string, jwt: string): void,
  setSelectedAppName(value: string): void,
  filterApps(value: string, apps: App[]): void,
  error: IErrorState,
  isLoading: boolean
} | RouteComponentProps<{
  appName?: string
}>;

export class Home extends Component<HomeProps, initialState> {
  constructor(props) {
    super(props);

    this.onFilterChanged = this.onFilterChanged.bind(this);
    this.onViewApp = this.onViewApp.bind(this);
  }

  public componentDidMount() {
    const { match: { params } } = this.props;

    if (_.isEmpty(this.props.selectedAppName) && !_.isEmpty(params.appName)) {
      this.props.setSelectedAppName(params.appName);
    }

    if (!this.props.isLoading && _.isEmpty(this.props.apps)) {
      this.props.getApps(this.props.cluster, this.props.identityToken);
    }
  }

  public shouldComponentUpdate(nextProps: HomeProps) {
    if (!this.props.selectedAppName
      || nextProps.selectedAppName !== this.props.selectedAppName
      || nextProps.isLoading !== this.props.isLoading
      || !_.isEqual(nextProps.filteredApps, this.props.filteredApps)) {
      return true;
    }
    return false;
  }

  private onFilterChanged(event) {
    this.props.filterApps(event.target.value, this.props.apps);
  }

  private onViewApp(appname: string, namespace: string, labelSelector: string) {
    this.props.setSelectedAppName(labelSelector);

    this.props.getAppOverview(appname, namespace, labelSelector, this.props.cluster, this.props.identityToken);

    this.props.history.push(`/${appname}?namespace=${namespace}&labelSelector=${encodeURIComponent(labelSelector)}`);
  }

  public render() {
    return (
      <div>
        <NavBar {...this.props} />
        <View
          onFilterChanged={this.onFilterChanged}
          onViewApp={this.onViewApp}
          {...this.props} />
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

export const mapStateToProps = ({ loadingState, appsState, authState, clustersState, errorState }: IGlobalState) => {
  return {
    cluster: clustersState.cluster,
    identityToken: authState.identityToken,
    apps: appsState.apps,
    filteredApps: appsState.filteredApps || appsState.apps,
    selectedAppName: appsState.selectedAppName,
    error: errorState,
    isLoading: loadingState.loading
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getApps: (cluster: string, jwt: string) => dispatch(getApps(cluster, jwt)),
    getAppOverview: (appname: string, namespace: string, queryString: string, cluster: string, jwt: string) => dispatch(getAppOverview(appname, namespace, queryString, cluster, jwt)),
    setSelectedAppName: (value: string) => dispatch(setSelectedAppName(value)),
    filterApps: (value: string, apps: App[]) => dispatch(filterApps(value, apps)),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Home));
