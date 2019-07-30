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
import { Row, Col } from 'reactstrap';
import ReactTooltip from 'react-tooltip';
import moment from 'moment';
import './styles.css';

export type TextItemsProps = {
  keyPrefix: string,
  items: any[]
};

const TextItems = (props: TextItemsProps) => {
  const { keyPrefix, items } = props;

  return (
    <small>
      <Row>
        {items && items.map((item, index) => {
          return (
            <Col key={`${keyPrefix}-${index}`} className={`col-xs-12 col-md-3 text-center pod-condidtions-padding`} data-tip data-for={`${keyPrefix}-${index}`}>

              <div className="text-center"><strong>{item.type}</strong></div>
              <div className={item.status.toLowerCase() === 'true' ? 'text-success' : 'text-danger'}>{item.status}</div>

              <ReactTooltip id={`${keyPrefix}-tooltip-${index}`} type='info'>
                Last Transition Time: ${moment(item.lastTransitionTime).format('ll LTS')}
              </ReactTooltip>
            </Col>
          )
        })}
      </Row>
    </small>
  );
}

export default TextItems; 
