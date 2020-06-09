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
import ReactList from 'react-list';
import { PodOverview } from '../../types';
import PodCard from '../../components/pod-card';
import _ from 'lodash';
import './styles.css';

export type PodOverviewProps = {
  podOverview: PodOverview
};

const Overview = ({
  podOverview
}: PodOverviewProps) => {
  return (
    <div>
      {!_.isEmpty(podOverview) &&
        <div>
          <h4>Pods</h4>
          <hr />
          <Row className="infinite-scroll-container">
            {!_.isEmpty(podOverview.pods)
            && <ReactList 
                itemRenderer={(index, key) => {
                  return (
                    <Col sm={6} key={key}>
                      <PodCard key={key} name={podOverview.name} pod={podOverview.pods[index]} />
                    </Col>
                  );
                }} 
                length={podOverview.pods && podOverview.pods.length || 0} 
                type="uniform"/>
            || <Col sm={12}>No Pods Returned.</Col>}
          </Row>
        </div>
      }
    </div>
  )
};

export default Overview;
