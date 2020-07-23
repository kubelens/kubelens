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
import { Row, Col } from 'reactstrap';
import { DeploymentOverview, ReplicaSetOverview } from '../../types';
import CardText from '../text';
import _ from 'lodash';

export type TextItemsProps = {
  overviews?: DeploymentOverview[],
  replicaSets?: ReplicaSetOverview[]
};

const TextItems = (props: TextItemsProps) => {
  const { overviews, replicaSets } = props;
  const rso = replicaSets && replicaSets.length > 0 ? replicaSets[replicaSets.length - 1] : {} as ReplicaSetOverview;
  const unavailablers = overviews && overviews.length > 0 ? overviews[overviews.length - 1].unavailableReplicas : 0;
  return (
    !_.isEmpty(rso) ?
      <small>
        <h5>Replica Sets</h5>
        <hr />
        <Row>
          <Col>
            <CardText label="Total" value={rso.replicas} />
          </Col>
          <Col>
            <CardText label="Ready" value={rso.readyReplicas} />
          </Col>
          <Col>
            <CardText label="Available" value={rso.availableReplicas} />
          </Col>
          <Col>
            <CardText label="Unavailable" value={unavailablers} />
          </Col>
          <Col>
            <CardText label="Fully Labeled" value={rso.fullyLabeledReplicas} />
          </Col>
        </Row>
      </small>
    : null
  );
}

export default TextItems; 
