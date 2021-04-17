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
import { Button, Modal, ModalHeader, ModalBody, ModalFooter, Row, Col, Card, CardHeader, CardBody, Input } from 'reactstrap';
import { AvailableCluster } from 'types';
import './styles.css';

export type ClusterModalProps = {
  handleClose(),
  onSelect(cluster: AvailableCluster),
  onFilterChanged: any,
  open: boolean,
  availableClusters: AvailableCluster[]
}

const ClusterModal = (props: ClusterModalProps) => {
  return (
    <Modal isOpen={props.open} toggle={props.handleClose}>
      <ModalHeader toggle={props.handleClose}>Available Clusters</ModalHeader>
      <ModalBody>
      <Input className="search" title="Available Clusters" type="text" placeholder="Search" onChange={props.onFilterChanged} />
      <hr />
      { props.availableClusters.map(c => {
        return (
          <div key={c.name} id="anti-shadow-div">
            <div id="shadow-div" >
              <a style={{ cursor: 'pointer' }} onClick={() => props.onSelect(c)}>
                <Card dir="ltr" style={{ marginBottom: '10px' }}>
                  <CardHeader className="link-card-title text-center">
                    <strong>
                      {c.name}
                    </strong>
                  </CardHeader>
                  <CardBody>
                    <Row>
                      <Col sm={12}>
                        <div>
                          <div className="app-list-text-root">
                            <small>Cluster: <strong>{c.cluster}</strong></small>
                          </div>
                        </div>
                      </Col>
                    </Row>
                  </CardBody>
                </Card>
              </a>
            </div>
          </div> 
        )
      })}
      </ModalBody>
      <ModalFooter>
        <Button color="primary" onClick={props.handleClose}>Close</Button>
      </ModalFooter>
    </Modal>
  );
}

export default ClusterModal;
