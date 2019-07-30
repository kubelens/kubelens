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
import { Row, Col, Button } from 'reactstrap';
import ReactTooltip from 'react-tooltip';
import { Service } from '../../types';
import _ from 'lodash';

export type ServiceOverviewProps = {
  serviceOverviews: Service[]
};

const ServiceOverview = ({
  serviceOverviews
}: ServiceOverviewProps) => {
  return (
    <div>
      {serviceOverviews && !_.isEmpty(serviceOverviews) &&
        <div>
          <h4>Services</h4>
          <hr />
          <Row>
            {serviceOverviews.map(svc => {
              return (
                <Col sm={6} key={svc.name} style={{ marginBottom: '10px' }}>

                  <div key={`${svc.name}-btn`} style={{ padding: 0 }} data-tip data-for={`${svc.name}-btn`}>
                    <Button id={`${svc.name}-btn`} disabled={true} block key={svc.name} color="info">{svc.name}</Button>
                    <ReactTooltip id={`${svc.name}-btn`} type='info'>
                      Coming Soon.
                    </ReactTooltip>
                  </div>
                </Col>
              )
            })
            }
          </Row>
          <br />
        </div>
      }
    </div>
  );
};

export default ServiceOverview;
