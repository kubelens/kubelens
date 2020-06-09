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
import React from 'react';
import { Row, Col, Card, CardHeader, CardBody } from 'reactstrap';
import { App } from '../../types/index';
import RightArrowInverse from '../../assets/right-arrow-yellow-inverse.png';
import _ from 'lodash';
import './styles.css';

export type AppCardProps = {
  app: App,
  index: number,
  match: any,
  selectedAppName: string,
  onViewApp(appname: string, namespace: string, labelSelector: string)
};

const AppCard = (props: AppCardProps) => {
  const { app, index, match, selectedAppName, onViewApp } = props;

  const viewApp = () => {
    return onViewApp(app.name, app.namespace, app.labelSelector);
  };
  // if from a link, grab the name of the app so we can mark which one is being viewed.
  const selected =
    (_.isEmpty(selectedAppName) && match)
      ? match.params.appName
      : selectedAppName;

  return (
    <div key={`${app.name}-${index}`} id="anti-shadow-div">
      <div id="shadow-div" >
        <Card dir="ltr" style={{ marginRight: (app.labelSelector === selected) ? -40 : 0, marginBottom: '10px', border: '3px solid #4D5061' }}>
          <CardHeader className="text-center" style={{ backgroundColor: 'white' }}>
            <strong>
              {app.name}
            </strong>
          </CardHeader>
          <CardBody>
            <Row>
              <Col sm={10}>
                <div>
                  <div className="app-list-text-root">
                    <small>Namespace: <strong>{app.namespace}</strong></small>
                  </div>
                  <div className="app-list-text-root">
                    <small>Kind: <strong>{app.kind}</strong></small>
                  </div>
                  <div className="app-list-text-secondary">
                    <small>
                      {
                        app.deployerLink
                          ? <a href={app.deployerLink} target="_blank" rel="noopener noreferrer"><strong>Deployer Link</strong></a>
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
  );
}

export default AppCard; 
