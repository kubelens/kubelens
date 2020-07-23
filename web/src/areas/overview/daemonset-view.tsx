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
import { Row, Col, Card, CardBody, CardFooter, Button } from 'reactstrap';
import CardText from '../../components/text';
import JsonViewModal from '../../components/json-view-modal';
import { DaemonSetOverview, ReplicaSetOverview } from '../../types';
import DeploymentOverviews from '../../components/deployment-overviews';
import _ from 'lodash';

export type DaemonSetOverviewProps = {
  daemonSetOverviews: DaemonSetOverview[],
  toggleModalType: (type: string) => void,
  conditionsModalOpen: boolean,
  configMapModalOpen: boolean,
  deploymentModalOpen: boolean,
  replicaSetOverviews: ReplicaSetOverview[]
};

const DaemonSetView = ({
  daemonSetOverviews,
  toggleModalType,
  conditionsModalOpen,
  configMapModalOpen,
  deploymentModalOpen,
  replicaSetOverviews
}: DaemonSetOverviewProps) => {
  return (
    <div>
      {!_.isEmpty(daemonSetOverviews) &&
        <span>
          <h4>DaemonSet</h4>
          <hr />
        </span>
      }
      {!_.isEmpty(daemonSetOverviews) && daemonSetOverviews.map((overview: DaemonSetOverview) => {
      return (
      <div key={overview.name}>
        <Card className="mb-4">
          <CardBody>
            <small>
              <Row>
                <Col sm={!_.isEmpty(overview.conditions) || !_.isEmpty(overview.deploymentOverviews) ||!_.isEmpty(overview.configMaps) ? 7 : 12}>
                  <Row>
                    <Col sm={3}>
                      <CardText label="Ready" value={overview.numberReady} />
                    </Col>
                    <Col sm={3}>
                      <CardText label="Available" value={overview.numberAvailable} />
                    </Col>
                    <Col sm={3}>
                      <CardText label="Unavailable" value={overview.numberUnavailable} />
                    </Col>
                    <Col sm={3}>
                      <CardText label="Misscheduled" value={overview.numberMisscheduled} />
                    </Col>
                  </Row>
                  <hr />
                  <Row>
                    <Col sm={4}>
                      <CardText label="CurrentNumberScheduled" value={overview.currentNumberScheduled} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="UpdatedNumberScheduled" value={overview.updatedNumberScheduled} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="DesiredNumberScheduled" value={overview.desiredNumberScheduled} />
                    </Col>
                  </Row>
                </Col>
                {!_.isEmpty(overview.conditions) || !_.isEmpty(overview.deploymentOverviews) ||!_.isEmpty(overview.configMaps)
                ? <Col sm={5}>
                  <Row>
                    <Col sm={6}>
                      {!_.isEmpty(overview.conditions)
                      ? <Button outline color="info" onClick={() => toggleModalType('ds-condition')} block>Conditions</Button>
                      : null}
                    </Col>
                    <Col sm={6}>
                      {!_.isEmpty(overview.deploymentOverviews)
                      ? <Button outline color="info" onClick={() => toggleModalType('ds-deployment')} block>Deployments</Button>
                      : null}
                      {!_.isEmpty(overview.configMaps)
                      ? <Button outline color="info" onClick={() => toggleModalType('ds-configMap')} block>ConfigMaps</Button>
                      : null}
                    </Col>
                  </Row>
                </Col>
                : null}
              </Row>
            </small>
          </CardBody>
          <CardFooter>
            <DeploymentOverviews 
              overviews={overview.deploymentOverviews}
              replicaSets={replicaSetOverviews} />
          </CardFooter>
        </Card>

        <JsonViewModal
          title="Conditions"
          show={conditionsModalOpen}
          body={overview.conditions}
          handleClose={() => {
            toggleModalType('ds-condition');
          }} />
    
        <JsonViewModal
          title="ConfigMaps"
          show={configMapModalOpen}
          body={overview.configMaps}
          handleClose={() => {
            toggleModalType('ds-configMap');
          }} />

        <JsonViewModal
          title="Deployments"
          show={deploymentModalOpen}
          body={overview.deploymentOverviews}
          handleClose={() => {
            toggleModalType('ds-deployment');
          }} />
      </div>
    )})}
  </div>
  );
};

export default DaemonSetView;
