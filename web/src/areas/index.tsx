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
import React, { Component, lazy } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router';
import NavBar from '../components/nav';
import _ from 'lodash';
import { IGlobalState } from '../store';
import './styles.css';

const Home = lazy(() => import('./home'));

type initialState = {};

export class Cluster extends Component<any, initialState> {
  public render() {
    // make sure we have the selected cluster before rendering the home page
    return (
      <div>
        <NavBar {...this.props} />
        <div className="background-logo"></div>
        {!_.isEmpty(this.props.clustersState.cluster) ? <Home {...this.props} /> : null}
      </div>
    );
  }
}

export const mapStateToProps = (state: IGlobalState) => {
  return state;
};

export const mapActionsToProps = () => {
  return {};
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(Cluster));
