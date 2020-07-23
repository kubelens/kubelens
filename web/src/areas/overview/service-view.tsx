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
import { Service } from '../../types';
import DeploymentOverviews from '../../components/deployment-overviews';
import { ReplicaSetOverview } from '../../types';
import _ from 'lodash';

export type ServiceOverviewProps = {
  serviceOverviews: Service[],
  toggleModalType: (type: string) => void,
  specModalOpen: boolean,
  statusModalOpen: boolean,
  configMapModalOpen: boolean,
  deploymentModalOpen: boolean,
  replicaSetOverviews: ReplicaSetOverview[]
};

const ServiceOverview = ({
  serviceOverviews,
  toggleModalType,
  specModalOpen,
  statusModalOpen,
  configMapModalOpen,
  deploymentModalOpen,
  replicaSetOverviews
}: ServiceOverviewProps) => {
  return (
    <div>
      {!_.isEmpty(serviceOverviews) &&
        <span>
          <h4>Service</h4>
          <hr />
        </span>
      }
      {!_.isEmpty(serviceOverviews) && serviceOverviews.map((overview: Service) => {
      return (
      <div key={overview.name}>
        <Card className="mb-4">
          <CardBody>
            <small>
              <Row>
                <Col sm={7}>
                  <Row>
                    <Col sm={4}>
                      <CardText label="Namespace" value={overview.namespace} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="Cluster IP" value={overview.clusterIP} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="Type" value={overview.type} />
                    </Col>
                  </Row>
                </Col>
                <Col sm={5}>
                  <Row>
                    <Col sm={6}>
                      {!_.isEmpty(overview.deploymentOverviews)
                      ? <Button outline color="info" onClick={() => toggleModalType('deployment')} block>Deployments</Button>
                      : null}
                      {!_.isEmpty(overview.configMaps)
                      ? <Button outline color="info" onClick={() => toggleModalType('configMap')} block>ConfigMaps</Button>
                      : null}
                    </Col>
                    <Col sm={6}>
                      <Button outline color="info" onClick={() => toggleModalType('spec')} block>Spec</Button>
                      <Button outline color="info" onClick={() => toggleModalType('status')} block>Status</Button>
                    </Col>
                  </Row>
                </Col>
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
          title="Service Spec"
          show={specModalOpen}
          body={overview.spec}
          handleClose={() => {
            toggleModalType('spec');
          }} />
    
        <JsonViewModal
          title="Service Status"
          show={statusModalOpen}
          body={overview.status}
          handleClose={() => {
            toggleModalType('status');
          }} />
    
        <JsonViewModal
          title="ConfigMaps"
          show={configMapModalOpen}
          body={overview.configMaps}
          handleClose={() => {
            toggleModalType('configMap');
          }} />

        <JsonViewModal
          title="Deployments"
          show={deploymentModalOpen}
          body={overview.deploymentOverviews}
          handleClose={() => {
            toggleModalType('deployment');
          }} />
      </div>
    )})}
  </div>
  );
};

export default ServiceOverview;
