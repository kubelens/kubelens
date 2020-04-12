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
import React, { lazy } from 'react';
import { Route } from 'react-router-dom';
import { Input, Card, CardHeader, CardBody, Row, Col } from 'reactstrap';
import { App } from '../../types';
import { RouteChildrenProps } from 'react-router';
import RightArrowInverse from '../../assets/right-arrow-yellow-inverse.png';
import _ from 'lodash';
import './home.css';

const Overview = lazy(() => import('../overview'));
const Pod = lazy(() => import('../pod'));

export type HomeViewProps =
  Partial<RouteChildrenProps<{
    appName?: string
  }>> & {
    onFilterChanged: Function,
    onViewApp(appname: string, labelSelector: string),
    filteredApps: App[],
    appsRequested: boolean,
    selectedAppName: string
  }

const HomePage = (props: HomeViewProps) => {
  const {
    match,
    onFilterChanged,
    onViewApp,
    filteredApps,
    selectedAppName
  } = props;

  return (
    <div className="container">
      <div className="inner-container">
        <div className="stuck">
          {/* search */}
          <Input className="search" title="Applications" type="text" placeholder="Search" onChange={onFilterChanged} />
          <p className="text-center search-title"><strong>Applications</strong></p>
          <hr className="separator" />

          {/* applications list */}
          {/* I don't understand css enough to get this to work without inline style, moving on. */}
          <div dir="rtl" style={{ maxHeight: '76vh', overflow: 'scroll', marginRight: -40, paddingRight: 40 }}>
            {filteredApps && filteredApps.map((value: any, index: number) => {
              const svc = value as App;
              const viewApp = () => {
                return onViewApp(svc.name, svc.labelSelector);
              };
              // if from a link, grab the name of the app so we can mark which one is being viewed.
              const selected =
                (_.isEmpty(selectedAppName) && match)
                  ? match.params.appName
                  : selectedAppName;

              return (
                <div key={`${svc.name}-${index}`} id="anti-shadow-div">
                  <div id="shadow-div" >
                    <Card dir="ltr" style={{ marginRight: (svc.labelSelector === selected) ? -40 : 0, marginBottom: '10px', border: '3px solid #4D5061' }}>
                      <CardHeader className="text-center" style={{ backgroundColor: 'white' }}>
                        <strong>
                          {svc.name}
                        </strong>
                      </CardHeader>
                      <CardBody>
                        <Row>
                          <Col sm={10}>
                            <div>
                              <div className="app-list-text-root">
                                <small>Namespace: <strong>{svc.namespace}</strong></small>
                              </div>
                              <div className="app-list-text-root">
                                <small>Kind: <strong>{svc.kind}</strong></small>
                              </div>
                              <div className="app-list-text-secondary">
                                <small>
                                  {
                                    svc.deployerLink
                                      ? <a href={svc.deployerLink} target="_blank" rel="noopener noreferrer"><strong>Deployer Link</strong></a>
                                      : <em>No Deployer Link</em>
                                  }
                                </small>
                              </div>
                            </div>
                          </Col>
                          <Col sm={2} className="action-right-container" >
                            <div onClick={viewApp}>
                              <span aria-hidden><img height={30} src={RightArrowInverse} alt="View" /></span>
                            </div>
                          </Col>
                        </Row>
                      </CardBody>
                    </Card>
                  </div>
                </div>
              )
            })}
          </div>
        </div>
        <div className="not-stuck">
          <Route exact path='/:appName' render={p => <Overview {...p} {...props} />} />
          <Route exact path='/:appName/pods/:podName' render={p => <Pod {...props} {...p} />} />
          <Route path='/' render={() => { return null }} />
        </div>
      </div>
    </div >
  );
};

export default HomePage;
