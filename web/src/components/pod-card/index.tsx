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
import { Row, Col, Card, CardHeader, CardBody, CardFooter } from 'reactstrap';
import { Link } from 'react-router-dom';
import CopyClipboard from '../copy-clipboard';
import CardText from '../text';
import { Pod } from '../../types';
import PodConditions from '../pod-conditions';
import { PodPhaseStyle } from './util';
import moment from 'moment';
import _ from 'lodash';
import './styles.css';

export type PodCardProps = {
  podName: string,
  overview: Pod
};

const PodCard = (props: PodCardProps) => {
  const { podName, overview } = props;

  return (
    <Card>
      <CardHeader className="link-card-title text-center">
        <Link
          to={{ pathname: `/${overview.linkedName}/pods/${podName}?namespace=${overview.namespace}` }}>
          <strong>{podName}</strong>
        </Link>
      </CardHeader>

      {/* body */}
      <CardBody>
        <small>
          <Row>
            <Col sm={4}>
              <CardText label="Namespace" value={overview.namespace} />
            </Col>
            <Col sm={4}>
              <CardText label="Start Time" value={moment(overview.pod.startTime).format('l LTS')} />
            </Col>
            <Col sm={4}>
              <CardText label="Status" value={<img style={{ marginTop: '5px' }} height={25} src={PodPhaseStyle(overview.pod.status.phase).img} alt="Status" />} />
            </Col>
          </Row>
          <Row >
            <Col sm={12}>
              <CardText label={overview.pod.spec.containers.length > 1 ? "Images" : "Image"} />
              {overview.pod.spec.containers.map(image => {
                return(<CopyClipboard key={image.image} labelText={`${image.name}`} value={image.name} size={16} />)
              })}
            </Col>
          </Row>
        </small>
      </CardBody>

      {/* footer */}
      {(!_.isEmpty(overview.pod.status.conditions)) ?
        <CardFooter>
          <PodConditions items={overview.pod.status.conditions} keyPrefix={overview.name} />
        </CardFooter>
        : null
      }
    </Card>
  );
}

export default PodCard; 
