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
import { ReplicaSetOverview } from '../../types';
import DeploymentOverviews from '../../components/deployment-overviews';
import _ from 'lodash';

export type ReplicaSetOverviewProps = {
  replicaSetOverviews: ReplicaSetOverview[],
  toggleModalType: (type: string) => void,
  conditionsModalOpen: boolean,
  configMapModalOpen: boolean,
  deploymentModalOpen: boolean
};

const ReplicaSetView = ({
  replicaSetOverviews,
  toggleModalType,
  conditionsModalOpen,
  configMapModalOpen,
  deploymentModalOpen
}: ReplicaSetOverviewProps) => {
  return (
    <div>
      {!_.isEmpty(replicaSetOverviews) &&
        <span>
          <h4>ReplicaSet</h4>
          <hr />
        </span>
      }
      {!_.isEmpty(replicaSetOverviews) && replicaSetOverviews.map((overview: ReplicaSetOverview) => {
      return (
      <div key={overview.name}>
        <Card className="mb-4">
          <CardBody>
            <small>
              <Row>
                <Col sm={!_.isEmpty(overview.conditions) || !_.isEmpty(overview.deploymentOverviews) ||!_.isEmpty(overview.configMaps) ? 7 : 12}>
                  <Row>
                    <Col sm={3}>
                      <CardText label="Replicas" value={overview.replicas} />
                    </Col>
                    <Col sm={3}>
                      <CardText label="Ready Replicas" value={overview.readyReplicas} />
                    </Col>
                    <Col sm={3}>
                      <CardText label="Available" value={overview.availableReplicas} />
                    </Col>
                    <Col sm={3}>
                      <CardText label="Fully Labeled" value={overview.fullyLabeledReplicas} />
                    </Col>
                  </Row>
                </Col>
              </Row>
            </small>
          </CardBody>
          {!_.isEmpty(overview.deploymentOverviews) ?
            <CardFooter>
              <DeploymentOverviews overviews={overview.deploymentOverviews} keyPrefix={`ds-${overview.name}`} />
            </CardFooter>
          : null
          }
        </Card>
      </div>
    )})}
  </div>
  );
};

export default ReplicaSetView;
