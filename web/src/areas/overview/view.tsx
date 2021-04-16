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
import { Row, Col, Card, CardBody, Button } from 'reactstrap';
import { DaemonSet, Deployment, Job, Pod, ReplicaSet, Service, ConfigMap } from '../../types';
import JsonViewModal from '../../components/json-view-modal';
import _ from 'lodash';
import './styles.css';

export type OverviewProps = {
  linkedName: string,
  namespace: string,
  serviceOverviews: Service[],
  daemonSetOverviews: DaemonSet[],
  deploymentOverviews: Deployment[],
  jobOverviews: Job[],
  podOverviews: Pod[],
  replicaSetOverviews: ReplicaSet[],
  configMapOverviews: ConfigMap[],
  toggleModalType: (type: string) => void,
  daemonSetsModalOpen: boolean,
  deploymentsModalOpen: boolean,
  jobsModalOpen: boolean,
  podsModalOpen: boolean,
  replicaSetsModalOpen: boolean,
  servicesModalOpen: boolean,
  configMapsModalOpen: boolean
};

const View = ({
  linkedName,
  namespace,
  serviceOverviews,
  daemonSetOverviews,
  deploymentOverviews,
  jobOverviews,
  podOverviews,
  replicaSetOverviews,
  configMapOverviews,
  toggleModalType,
  daemonSetsModalOpen,
  deploymentsModalOpen,
  jobsModalOpen,
  podsModalOpen,
  replicaSetsModalOpen,
  servicesModalOpen,
  configMapsModalOpen
}: OverviewProps) => {
  return (
    <div>
      <h4>{linkedName} - {namespace}</h4>
      <hr />
      <div key={`${linkedName}-${namespace}`}>
        <Card className="mb-4">
          <CardBody>
            <Row>
              {!_.isEmpty(daemonSetOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('daemonSets')} block>DaemonSets</Button>
                </Col>
              }
              {!_.isEmpty(deploymentOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('deployments')} block>Deployments</Button>
                </Col>
              }
              {!_.isEmpty(jobOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('jobs')} block>Jobs</Button>
                </Col>
              }
              {!_.isEmpty(podOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('pods')} block>Pods</Button>
                </Col>
              }
              {!_.isEmpty(replicaSetOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('replicaSets')} block>ReplicaSets</Button>
                </Col>
              }
              {!_.isEmpty(serviceOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('services')} block>Services</Button>
                </Col>
              }
              {!_.isEmpty(configMapOverviews) &&
                <Col sm={4}>
                <Button className="overview-button" outline color="info" onClick={() => toggleModalType('configMaps')} block>ConfigMaps</Button>
                </Col>
              }
            </Row>
          </CardBody>
        </Card>

        <JsonViewModal
          title="DaemonSets"
          show={daemonSetsModalOpen}
          body={daemonSetOverviews}
          handleClose={() => {
            toggleModalType('daemonSets');
          }} />
    
        <JsonViewModal
          title="Deployments"
          show={deploymentsModalOpen}
          body={deploymentOverviews}
          handleClose={() => {
            toggleModalType('deployments');
          }} />
    
        <JsonViewModal
          title="Jobs"
          show={jobsModalOpen}
          body={jobOverviews}
          handleClose={() => {
            toggleModalType('jobs');
          }} />

        <JsonViewModal
          title="Pods"
          show={podsModalOpen}
          body={podOverviews}
          handleClose={() => {
            toggleModalType('pods');
          }} />
  
        <JsonViewModal
          title="ReplicaSets"
          show={replicaSetsModalOpen}
          body={replicaSetOverviews}
          handleClose={() => {
            toggleModalType('replicaSets');
          }} />
    
        <JsonViewModal
          title="Services"
          show={servicesModalOpen}
          body={serviceOverviews}
          handleClose={() => {
            toggleModalType('services');
          }} />
        
        <JsonViewModal
          title="ConfigMaps"
          show={configMapsModalOpen}
          body={configMapOverviews}
          handleClose={() => {
            toggleModalType('configMaps');
          }} />
      </div>
  </div>
  );
};

export default View;
