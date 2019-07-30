/*
MIT License

Copyright (c) 2019 The KubeLens Authors

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
import { Card, CardHeader, CardBody, CardFooter, Row, Col } from 'reactstrap';
import JsonViewModal from '../../components/json-view-modal';
import CopyClipboard from '../../components/copy-clipboard';
import CardText from '../../components/text';
import { Button } from 'reactstrap';
import PodConditions from '../../components/pod-conditions';
import { PodDetail, Log } from '../../types';
import moment from 'moment';
import _ from 'lodash';
import './styles.css';

export type PodPageProps = {
  podDetail: PodDetail,
  showEnvModal: boolean,
  showSpecModal: boolean,
  showStatusModal: boolean,
  toggleModalType: (type: string) => void,
  logs: Log,
  openLogStream: () => void,
  closeLogStream: () => void,
  logStream?: string,
  streamEnabled: boolean,
  envBody: {},
  hasLogAccess: boolean
};

const PodPage = ({
  podDetail,
  showEnvModal,
  showSpecModal,
  showStatusModal,
  toggleModalType,
  logs,
  openLogStream,
  closeLogStream,
  logStream,
  streamEnabled,
  envBody,
  hasLogAccess
}: PodPageProps) => {
  // is this right? not sure if there would ever be more than 1 container status
  let ready: boolean = false;
  let restartCount: number = 0;
  let image: string = '';

  if (!_.isEmpty(podDetail) && podDetail.status && !_.isEmpty(podDetail.status.containerStatuses)) {
    const podDetailStatus = podDetail.status.containerStatuses[0];
    ready = podDetailStatus.ready;
    restartCount = podDetailStatus.restartCount;
    image = podDetailStatus.image;
  }

  return (
    <div className="pod-container">
      {!_.isEmpty(podDetail) ?
        <div>
          <Card className="kind-detail-container mb-4">
            <CardHeader className="kind-detail-title text-center">
              {podDetail.name}
            </CardHeader>

            {/* body */}
            <CardBody>
              <Col xs={12}>
                <Row>
                  <Col xs={12} md={9}>
                    <Row>
                      <Col sm={6}>
                        <CardText label="Start Time:" value={podDetail.status && moment(podDetail.status.startTime).format('ll LTS')} />
                        <br />
                        <CardText label="Namespace:" value={podDetail.namespace} />
                      </Col>
                      <Col sm={6}>
                        <div className="row">
                          <div className="col-sm-6">
                            <CardText label="HostIP:" value={podDetail.hostIP} />
                          </div>
                          <div className="col-sm-6">
                            <CardText label="PodIP:" value={podDetail.podIP} />
                          </div>
                        </div>
                        <br />
                        <div className="row">
                          <div className="col-sm-6">
                            <CardText label="Ready:" value={<span className={`${ready ? 'text-success' : 'text-danger'}`}>{`${ready}`}</span>} />
                          </div>
                          <div className="col-sm-6">
                            <CardText label="Restarts:" value={`${restartCount}`} />
                          </div>
                        </div>
                      </Col>
                      <br />
                    </Row>
                    <br />
                    <CardText label="Image" />
                    <CopyClipboard labelText={image} value={image} size={16} />
                  </Col>
                  <Col xs={12} md={3}>
                    {streamEnabled
                      ? <Button block color="warning" onClick={closeLogStream}>Close Stream</Button>
                      : <Button block color="info" onClick={openLogStream}>Stream Logs</Button>
                    }
                    <br />
                    {envBody && <Button outline color="info" onClick={() => toggleModalType('env')} block>Environment Variables</Button>}
                    <Button outline color="info" onClick={() => toggleModalType('spec')} block>Pod Spec</Button>
                    <Button outline color="info" onClick={() => toggleModalType('status')} block>Pod Status</Button>
                  </Col>
                </Row>
                <br />
                {hasLogAccess
                  ? <div>
                    <h4><CopyClipboard labelText={`Log Stdout ${logStream ? 'Stream' : ''}`} value={logStream ? logStream : (logs ? logs.output : '')} size={22} /></h4>
                    <hr />
                    {logStream || logs
                      ? <div style={{ backgroundColor: 'rgb(39, 40, 34)', padding: '10px', overflow: 'auto', maxHeight: '300px' }}>
                        <pre style={{ whiteSpace: 'pre-line', color: 'white', fontSize: '12px' }}>
                          {logStream && logStream}
                          {logs && !logStream && logs.output}
                        </pre>
                      </div>
                      : null
                    }
                  </div>
                  : <h4>You do not have access to view logs for this pod.</h4>
                }
              </Col>
            </CardBody>

            {/* footer */}
            {(podDetail.status && !_.isEmpty(podDetail.status.conditions)) ?
              <CardFooter className="kind-detail-footer text-center">
                <PodConditions items={podDetail.status && podDetail.status.conditions} keyPrefix={podDetail.name} />
              </CardFooter>
              : null
            }
          </Card>

          <JsonViewModal
            title="Pod Environment Variables"
            show={showEnvModal}
            body={envBody}
            handleClose={() => {
              toggleModalType('env');
            }} />

          <JsonViewModal
            title="Pod Spec"
            show={showSpecModal}
            body={podDetail.spec}
            handleClose={() => {
              toggleModalType('spec');
            }} />

          <JsonViewModal
            title="Pod Status"
            show={showStatusModal}
            body={podDetail.status}
            handleClose={() => {
              toggleModalType('status');
            }} />
        </div>
        : null
      }
    </div >
  )
}

export default PodPage;
