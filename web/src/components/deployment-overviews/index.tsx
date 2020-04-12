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
import { DeploymentOverview } from '../../types';

export type TextItemsProps = {
  keyPrefix: string,
  overviews?: DeploymentOverview[]
};

const TextItems = (props: TextItemsProps) => {
  const { keyPrefix, overviews } = props;

  return (
    <small>
      {overviews && overviews.map((ov, index) => {
        return(
          <Row key={`${keyPrefix}-${index}`}>
            <Col sm={3}>
              <div className="text-center"><strong>Replicas</strong></div>
              <div>{ov.replicas}</div>
            </Col>
            <Col sm={3}>
              <div className="text-center"><strong>UpdatedReplicas</strong></div>
              <div>{ov.updatedReplicas}</div>
            </Col>
            <Col sm={3}>
              <div className="text-center"><strong>ReadyReplicas</strong></div>
              <div>{ov.readyReplicas}</div>
            </Col>
            <Col sm={3}>
              <div className="text-center"><strong>UnavailableReplicas</strong></div>
              <div>{ov.unavailableReplicas}</div>
            </Col>
          </Row>
        )
      })}
    </small>
  );
}

export default TextItems; 
