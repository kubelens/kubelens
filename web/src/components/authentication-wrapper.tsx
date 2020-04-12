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
import { withRouter } from 'react-router';
import { connect } from 'react-redux';
import { setIdentityToken } from '../actions/auth';
import { IGlobalState } from 'store';
import config from '../config';
import { AuthClient } from '../auth/authClient';
import { closeErrorModal } from '../actions/error';
import { IErrorState } from '../reducers/error';
import APIErrorModal from './error-modal';
import unauthorized from '../assets/access-denied.png';
import logo from '../assets/kubelens-logo-inverse.png';
import { Col, Row } from 'reactstrap';

export interface AuthProps {
  authClient?: AuthClient,
  identityToken?: string,
  setIdentityToken(identityToken?: string): void
  closeErrorModal(): void,
  error: IErrorState
}

class AuthenticationWrapper extends Component<AuthProps, any> {
  state = {
    error: null,
    authEnabled: true
  }

  async componentDidMount() {
    const cfg = await config();

    if (cfg.oAuthEnabled) {
      if (this.props.authClient) {
        try {
          let token = await this.props.authClient.ensureAuthed();
          // using the accessToken since all we need is the email of the user for authorization.
          this.props.setIdentityToken(token.accessToken);
        } catch (e) {
          this.setState({ error: e }); 
        }
      }
    } else {
      this.setState({ authEnabled: false });
    }
  }

  render() {
    const { error } = this.state;

    if (!this.state.authEnabled) {
      return this.props.children;
    }

    if (!error && this.props.identityToken) {
      return this.props.children;
    }

    if (error) {
      return (
        <div>
          <div className="kubelens-navbar" >
            <img height={50} src={logo} alt="Kubelens Logo" />
          </div>
          <Row className="text-center">
            <Col sx={{size: 10, offset: 6}}>
                <img className="text-center" src={unauthorized} alt="Unauthorized"/>
            </Col>
          </Row>
          <APIErrorModal
            open={this.props.error.apiOpen}
            handleClose={this.props.closeErrorModal}
            status={this.props.error.status}
            statusText={this.props.error.statusText}
            message={this.props.error.message} />
        </div>
      )
    }

    return <div className='ns-icon ns-loading' />;
  }
}

export const mapStateToProps = ({ authState, errorState }: IGlobalState) => {
  return {
    identityToken: authState.identityToken,
    error: errorState
  };
};

export const mapActionsToProps = (dispatch) => {
  return {
    setIdentityToken: (identityToken?: string) => dispatch(setIdentityToken(identityToken)),
    closeErrorModal: () => dispatch(closeErrorModal())
  };
};

export default withRouter(connect(
  mapStateToProps,
  mapActionsToProps
)(AuthenticationWrapper));