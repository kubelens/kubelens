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
import { JobOverview, ReplicaSetOverview } from '../../types';
import DeploymentOverviews from '../../components/deployment-overviews';
import _ from 'lodash';

export type JobOverviewProps = {
  jobOverviews: JobOverview[],
  toggleModalType: (type: string) => void,
  conditionsModalOpen: boolean,
  configMapModalOpen: boolean,
  deploymentModalOpen: boolean,
  replicaSetOverviews: ReplicaSetOverview[]
};

const JobView = ({
  jobOverviews,
  toggleModalType,
  conditionsModalOpen,
  configMapModalOpen,
  deploymentModalOpen,
  replicaSetOverviews
}: JobOverviewProps) => {
  return (
    <div>
      {!_.isEmpty(jobOverviews) &&
        <span>
          <h4>Job</h4>
          <hr />
        </span>
      }
      {!_.isEmpty(jobOverviews) && jobOverviews.map((overview: JobOverview) => {
      return (
      <div key={overview.name}>
        <Card className="mb-4">
          <CardBody>
            <small>
              <Row>
                <Col sm={!_.isEmpty(overview.conditions) || !_.isEmpty(overview.deploymentOverviews) ||!_.isEmpty(overview.configMaps) ? 7 : 12}>
                  <Row>
                    <Col sm={4}>
                      <CardText label="Name" value={overview.name} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="Start Time" value={overview.startTime} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="Completion Time" value={overview.completionTime} />
                    </Col>
                  </Row>
                  <hr />
                  <Row>
                    <Col sm={4}>
                      <CardText label="Active" value={overview.active} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="Succeeded" value={overview.succeeded} />
                    </Col>
                    <Col sm={4}>
                      <CardText label="Failed" value={overview.failed} />
                    </Col>
                  </Row>
                </Col>
                {!_.isEmpty(overview.conditions) || !_.isEmpty(overview.deploymentOverviews) ||!_.isEmpty(overview.configMaps)
                ? <Col sm={5}>
                  <Row>
                    <Col sm={6}>
                      {!_.isEmpty(overview.conditions)
                      ? <Button outline color="info" onClick={() => toggleModalType('job-condition')} block>Conditions</Button>
                      : null}
                    </Col>
                    <Col sm={6}>
                      {!_.isEmpty(overview.deploymentOverviews)
                      ? <Button outline color="info" onClick={() => toggleModalType('job-deployment')} block>Deployments</Button>
                      : null}
                      {!_.isEmpty(overview.configMaps)
                      ? <Button outline color="info" onClick={() => toggleModalType('job-configMap')} block>ConfigMaps</Button>
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
            toggleModalType('job-condition');
          }} />
    
        <JsonViewModal
          title="ConfigMaps"
          show={configMapModalOpen}
          body={overview.configMaps}
          handleClose={() => {
            toggleModalType('job-configMap');
          }} />

        <JsonViewModal
          title="Deployments"
          show={deploymentModalOpen}
          body={overview.deploymentOverviews}
          handleClose={() => {
            toggleModalType('job-deployment');
          }} />
      </div>
    )})}
  </div>
  );
};

export default JobView;
