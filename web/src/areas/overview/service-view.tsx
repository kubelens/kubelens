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
import CardText from '../../components/text';
import JsonViewModal from '../../components/json-view-modal';
import { Service } from '../../types';
import _ from 'lodash';

export type ServiceOverviewProps = {
  serviceOverviews: Service[],
  toggleModalType: (type: string) => void,
  specModalOpen: boolean,
  statusModalOpen: boolean,
  configMapModalOpen: boolean
};

const ServiceOverview = ({
  serviceOverviews,
  toggleModalType,
  specModalOpen,
  statusModalOpen,
  configMapModalOpen
}: ServiceOverviewProps) => {
  // There should only ever be 1 overview for a service, kept as an array for ease.
  const overview:Service = !_.isEmpty(serviceOverviews) ? serviceOverviews[0] : {} as Service;

  return (
    <div>
      {!_.isEmpty(overview)
      ? <div>
          <h4>Service Name: {overview.name}</h4>
          <hr />
          <Card className="kind-detail-container mb-4">
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
                      {!_.isEmpty(overview.configMaps)
                      ? <Col sm={6}>
                          <Button outline color="info" onClick={() => toggleModalType('configMap')} block>Service ConfigMaps</Button>
                        </Col>
                      : null}
                      <Col sm={!_.isEmpty(overview.configMaps) ? 6 : 12}>
                        <Button outline color="info" onClick={() => toggleModalType('spec')} block>Service Spec</Button>
                        <Button outline color="info" onClick={() => toggleModalType('status')} block>Service Status</Button>
                      </Col>
                    </Row>
                  </Col>
                </Row>
              </small>
            </CardBody>
          </Card>
        </div>
      : null}

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
      title="Service ConfigMaps"
      show={configMapModalOpen}
      body={overview.configMaps}
      handleClose={() => {
        toggleModalType('configMap');
      }} />

    </div>
  );
};

export default ServiceOverview;
