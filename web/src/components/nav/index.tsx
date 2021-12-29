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
import { Button } from 'reactstrap';
import LogoSpinner from '../logo-spinner';
import logo from '../../assets/kubelens-logo-inverse.png';
import { IGlobalState } from 'store';
import config from '../../config';
import { setSelectedCluster } from '../../actions/cluster';
import { getOverviews } from '../../actions/overviews';
import ClusterModal  from '../cluster-modal';
import _ from 'lodash';
import './styles.css';
import { AvailableCluster } from 'types';

export type NavBarState = {
  clusterSelectOpen: boolean,
  availableClusters: AvailableCluster[],
  filteredClusters: AvailableCluster[]
};

export interface NavBarProps extends
  RouteComponentProps {
  isLoading: boolean,
  selectedCluster: AvailableCluster,
  identityToken: string,
  setSelectedCluster(cluster: AvailableCluster): void,
  getGetOverviews(cluster: string, jwt: string): void
}

export class NavBar extends Component<NavBarProps, NavBarState> {
  public state: NavBarState = {
    clusterSelectOpen: false,
    availableClusters: [],
    filteredClusters: []
  };

  constructor(props) {
    super(props);

    this.toggleClusterSelect = this.toggleClusterSelect.bind(this);
    this.setSelectedCluster = this.setSelectedCluster.bind(this);
    this.onFilterChanged = this.onFilterChanged.bind(this);
  }

  async componentDidMount() {
    const cfg = await config();
    this.setState({ availableClusters: cfg.availableClusters });
    this.setState({ filteredClusters: cfg.availableClusters });

    if (!this.props.selectedCluster) {
      this.props.setSelectedCluster(cfg.availableClusters[0]);
    }
  }

  private toggleClusterSelect() {
    this.setState({ clusterSelectOpen: !this.state.clusterSelectOpen });
  }

  private async setSelectedCluster(cluster: AvailableCluster) {
    this.props.setSelectedCluster(cluster);
    this.props.getGetOverviews(cluster.url, this.props.identityToken);
    this.toggleClusterSelect();
  }

  private onFilterChanged(event) {
    const filtered = _.filter(this.state.availableClusters, (cluster: AvailableCluster) => {
      return cluster.cluster.includes(event.target.value)
        || cluster.name.includes(event.target.value)
        || cluster.url.includes(event.target.value);
    });
    this.setState({ filteredClusters: filtered });
  }

  public render() {
    return (
      <div className="kubelens-navbar" >
        <a id="logo" href="/">
          {this.props.isLoading
            ? <LogoSpinner />
            : <img height={50} src={logo} alt="Kubelens Logo" />
          }
        </a>
        <div id="navbar-right">
          <Button style={{marginTop: "10px"}} color="info" onClick={() => this.toggleClusterSelect()} block>{this.props.selectedCluster ? this.props.selectedCluster.name : 'Select Cluster'}</Button>
        </div>

        <ClusterModal
          handleClose={this.toggleClusterSelect}
          onSelect={this.setSelectedCluster}
          onFilterChanged={this.onFilterChanged}
          open={this.state.clusterSelectOpen}
          availableClusters={this.state.filteredClusters} />
      </div>
    );
  }
}

export const mapStateToProps = ({ loadingState, clustersState, authState }: IGlobalState) => {
  return {
    isLoading: loadingState.loading,
    selectedCluster: clustersState.cluster,
    identityToken: authState.identityToken
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    setSelectedCluster: (cluster: AvailableCluster) => dispatch(setSelectedCluster(cluster)),
    getGetOverviews: (cluster: string, jwt: string) => dispatch(getOverviews(cluster, jwt))
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(NavBar));

