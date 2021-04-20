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
import { Overview, SelectedOverview } from "../../types";
import _ from 'lodash';
import { IGlobalState } from '../../store';
import { getOverviews, getOverview, setSelectedOverview, filterOverviews } from '../../actions/overviews';
import { closeErrorModal } from '../../actions/error';
import APIErrorModal from '../../components/error-modal';
import { IErrorState } from '../../reducers/error';

type initialState = {};

export type HomeProps = {
  cluster: string,
  identityToken?: string,
  selectedAppName: string,
  filteredApps: Overview[],
  overviews: Overview[],
  getOverviews(cluster: string, jwt: string): void,
  getOverview(appname: string, namespace: string, cluster: string, jwt: string): void,
  setSelectedOverview(value: SelectedOverview): void,
  filterOverviews(value: string, apps: Overview[]): void,
  error: IErrorState,
  isLoading: boolean
} | RouteComponentProps<{
  appName?: string
}>;

export class Home extends Component<HomeProps, initialState> {
  constructor(props) {
    super(props);

    this.onFilterChanged = this.onFilterChanged.bind(this);
    this.onViewOverview = this.onViewOverview.bind(this);
  }

  public componentDidMount() {
    const { match: { params }, location: { search } } = this.props;
    const query = new URLSearchParams(search);

    if (_.isEmpty(this.props.selectedAppName) && !_.isEmpty(params.linkedName)) {
      this.props.setSelectedOverview({linkName: params.linkedName, namespace: query.get('namespace')});
    }

    if (!this.props.isLoading && _.isEmpty(this.props.overviews)) {
      this.props.getOverviews(this.props.cluster, this.props.identityToken);
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
    this.props.filterOverviews(event.target.value, this.props.overviews);
  }

  private onViewOverview(linkedName: string, namespace: string) {
    this.props.setSelectedOverview({linkedName, namespace});

    this.props.getOverview(linkedName, namespace, this.props.cluster, this.props.identityToken);

    this.props.history.push(`/${linkedName}?namespace=${namespace}`);
  }

  public render() {
    return (
      <div>
        <NavBar {...this.props} />
        <View
          onFilterChanged={this.onFilterChanged}
          onViewOverview={this.onViewOverview}
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

export const mapStateToProps = ({ loadingState, overviewsState, authState, clustersState, errorState }: IGlobalState) => {
  return {
    cluster: clustersState.cluster && clustersState.cluster.url,
    identityToken: authState.identityToken,
    overviews: overviewsState.overviews,
    filteredOverviews: overviewsState.filteredOverviews || overviewsState.overviews,
    selectedOverview: overviewsState.selectedOverview,
    error: errorState,
    isLoading: loadingState.loading
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getOverviews: (cluster: string, jwt: string) => dispatch(getOverviews(cluster, jwt)),
    getOverview: (linkedName: string, namespace: string, cluster: string, jwt: string) => dispatch(getOverview(linkedName, namespace, cluster, jwt)),
    setSelectedOverview: (value: SelectedOverview) => dispatch(setSelectedOverview(value)),
    filterOverviews: (value: string, apps: Overview[]) => dispatch(filterOverviews(value, apps)),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Home));
