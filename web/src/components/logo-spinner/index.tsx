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
import wheel from '../../assets/kubelens-wheel-alt.png';
import logoName from '../../assets/kubelens-logo-name-alt.png';
import './styles.css';

export type LogoState = {
  spinner: string,
  intervalRef?: NodeJS.Timeout
};

export type LogoProps = {}

export class Logo extends Component<LogoProps, LogoState> {
  private mounted: boolean = false;
  public state: LogoState = {
    spinner: ''
  };

  public componentDidMount() {
    this.mounted = true;
    setTimeout(() => {
      this.setState({ spinner: 'spinner-rotate-init' })
    }, 200);
    const itv = setInterval(() => {
      this.setState({ spinner: this.state.spinner === 'spinner-rotate-right' ? 'spinner-rotate-left' : 'spinner-rotate-right' });
    }, 1000);

    this.setState({ intervalRef: itv });
  }

  public componentWillUnmount() {
    if (this.mounted) {
      if (this.state.intervalRef) {
        clearInterval(this.state.intervalRef);
      }
    }
    this.mounted = false;
  }

  public render() {
    return (
      <span><img height={41} src={wheel} alt="KubeLens Wheel" className={`spinner ${this.state.spinner}`} />  <img height={50} src={logoName} alt="KubeLens" /></span>
    );
  }
}

export const mapStateToProps = () => {
  return {};
};

export default connect(
  mapStateToProps
)(Logo);

