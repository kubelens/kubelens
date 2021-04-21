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
import ViewPage from './view';
import PodOverviewPage from './pod-view';
import { Service, Pod, DaemonSet, Deployment, Job, ReplicaSet, SelectedOverview } from "../../types";
import { RouteComponentProps, withRouter } from 'react-router';
import _ from 'lodash';
import { IGlobalState } from '../../store';
import { getOverview, setSelectedOverview } from '../../actions/overviews';
import APIErrorModal from '../../components/error-modal';
import { closeErrorModal } from '../../actions/error';
import { IErrorState } from '../../reducers/error';
import qs from 'qs';

type initialState = {
  daemonSetsModalOpen: boolean,
  deploymentsModalOpen: boolean,
  jobsModalOpen: boolean,
  podsModalOpen: boolean,
  replicaSetsModalOpen: boolean,
  servicesModalOpen: boolean,
  configMapsModalOpen: boolean
};

export type OverviewProps = {
  identityToken?: string,
  serviceOverviews?: Service[],
  daemonSetOverviews?: DaemonSet[],
  deploymentOverviews?: Deployment[],
  podOverviews: Pod[],
  jobOverviews?: Job[],
  replicaSetOverviews?: ReplicaSet[],
  selectedOverview: SelectedOverview,
  getAppOverview(appname: string, queryString: string): void,
  setSelectedOverview(value: SelectedOverview): void,
  error: IErrorState,
  overviewsEmpty: boolean,
} | RouteComponentProps<{
  appName?: string
}>;

export class Overview extends Component<OverviewProps, initialState> {
  state: initialState = {
    daemonSetsModalOpen: false,
    deploymentsModalOpen: false,
    jobsModalOpen: false,
    podsModalOpen: false,
    replicaSetsModalOpen: false,
    servicesModalOpen: false,
    configMapsModalOpen: false
  }

  async componentDidMount() {
    const { match: { params }, location: { search } } = this.props;
    const query = qs.parse(search.replace('?',''));

    if (params.linkedName) {
      let overview = {
        linkedName: params.linkedName,
        namespace: query.namespace
      };

      if (_.isEmpty(this.props.selectedOverview)) {
        this.props.setSelectedOverview({linkedName: overview.linkedName, namespace: overview.namespace});
      } else {
        overview = this.props.selectedOverview;
      }
      
      if (!this.props.isLoading && !_.isEmpty(overview.linkedName) && !_.isEmpty(overview.namespace) && this.props.overviewsEmpty) {
        this.props.getOverview(overview.linkedName, overview.namespace, this.props.cluster, this.props.identityToken);
      }
    }
  }

  toggleModalType = (type) => {
    switch (type) {
      case 'daemonSets':
        this.setState({ daemonSetsModalOpen: !this.state.daemonSetsModalOpen});
        break;
      case 'deployments':
        this.setState({ deploymentsModalOpen: !this.state.deploymentsModalOpen});
        break;
      case 'jobs':
          this.setState({ jobsModalOpen: !this.state.jobsModalOpen});
          break;
      case 'pods':
        this.setState({ podsModalOpen: !this.state.podsModalOpen});
        break;
      case 'replicaSets':
        this.setState({ replicaSetsModalOpen: !this.state.replicaSetsModalOpen});
        break;
      case 'services':
        this.setState({ servicesModalOpen: !this.state.servicesModalOpen});
        break;
      case 'configMaps':
        this.setState({ configMapsModalOpen: !this.state.configMapsModalOpen});
        break;
      default:
        break;
    }
  }

  render() {
    return (
      <div className="overview-container">
        <ViewPage 
          linkedName={this.props.linkedName}
          namespace={this.props.namespace}
          daemonSetOverviews={this.props.daemonSetOverviews}
          deploymentOverviews={this.props.deploymentOverviews}
          jobOverviews={this.props.jobsOverviews}
          podOverviews={this.props.podsOverviews}
          replicaSetOverviews={this.props.replicaSetOverviews}
          serviceOverviews={this.props.serviceOverviews}
          configMapOverviews={this.props.configMapOverviews}
          toggleModalType={this.toggleModalType}
          daemonSetsModalOpen={this.state.daemonSetsModalOpen}
          deploymentsModalOpen={this.state.deploymentsModalOpen}
          jobsModalOpen={this.state.jobsModalOpen}
          podsModalOpen={this.state.podsModalOpen}
          replicaSetsModalOpen={this.state.replicaSetsModalOpen}
          servicesModalOpen={this.state.servicesModalOpen}
          configMapsModalOpen={this.state.configMapsModalOpen} />

        <PodOverviewPage podOverviews={this.props.podOverviews} />

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

export const mapStateToProps = ({ overviewsState, authState, clustersState, errorState }: IGlobalState) => {
  let linkedName,
    namespace,
    serviceOverviews,
    daemonSetOverviews,
    deploymentOverviews,
    jobOverviews,
    podOverviews,
    replicaSetOverviews,
    configMapOverviews;

  if (overviewsState.overview && overviewsState.overview.linkedName) {
    linkedName = overviewsState.overview.linkedName;
  }

  if (overviewsState.overview && overviewsState.overview.namespace) {
    namespace = overviewsState.overview.namespace;
  }

  if (overviewsState.overview && overviewsState.overview.services) {
    serviceOverviews = overviewsState.overview.services;
  }

  if (overviewsState.overview && overviewsState.overview.daemonSets) {
    daemonSetOverviews = overviewsState.overview.daemonSets;
  }

  if (overviewsState.overview && overviewsState.overview.deployments) {
    deploymentOverviews = overviewsState.overview.deployments;
  }

  if (overviewsState.overview && overviewsState.overview.jobs) {
    jobOverviews = overviewsState.overview.jobs;
  }

  if (overviewsState.overview && overviewsState.overview.pods) {
    podOverviews = overviewsState.overview.pods;
  }

  if (overviewsState.overview && overviewsState.overview.replicaSets) {
    replicaSetOverviews = overviewsState.overview.replicaSets;
  }

  if (overviewsState.overview && overviewsState.overview.configMaps) {
    configMapOverviews = overviewsState.overview.configMaps;
  }

  const overviewsEmpty = _.isEmpty(podOverviews) && 
    (_.isEmpty(serviceOverviews) || 
    _.isEmpty(daemonSetOverviews) || 
    _.isEmpty(deploymentOverviews) ||
    _.isEmpty(jobOverviews) || 
    _.isEmpty(replicaSetOverviews) ||
    _.isEmpty(configMapOverviews));

  return {
    cluster: clustersState.cluster && clustersState.cluster.url,
    identityToken: authState.identityToken,
    linkedName: linkedName,
    namespace: namespace,
    serviceOverviews: serviceOverviews,
    daemonSetOverviews: daemonSetOverviews,
    deploymentOverviews: deploymentOverviews,
    jobOverviews: jobOverviews,
    podOverviews: podOverviews,
    replicaSetOverviews: replicaSetOverviews,
    configMapOverviews:configMapOverviews,
    selectedOverview: overviewsState.selectedOverview,
    error: errorState,
    overviewsEmpty: overviewsEmpty
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    getOverview: (linkedName: string, namespace: string, cluster: string, jwt: string) => dispatch(getOverview(linkedName, namespace, cluster, jwt)),
    setSelectedOverview: (value: SelectedOverview) => dispatch(setSelectedOverview(value)),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Overview));