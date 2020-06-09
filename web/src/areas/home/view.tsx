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
import React, { lazy } from 'react';
import { Route } from 'react-router-dom';
import { Input } from 'reactstrap';
import { App } from '../../types';
import { RouteChildrenProps } from 'react-router';
import ReactList from 'react-list';
import AppCard from '../../components/app-card';
import _ from 'lodash';
import './home.css';

const Overview = lazy(() => import('../overview'));
const Pod = lazy(() => import('../pod'));

export type HomeViewProps =
  Partial<RouteChildrenProps<{
    appName?: string
  }>> & {
    onFilterChanged: Function,
    onViewApp(appname: string, namespace: string, labelSelector: string),
    filteredApps: App[],
    selectedAppName: string
  }

const HomePage = (props: HomeViewProps) => {
  const {
    match,
    onFilterChanged,
    onViewApp,
    filteredApps,
    selectedAppName
  } = props;

  return (
    <div className="container">
      <div className="inner-container">
        <div className="stuck">
          {/* search */}
          <Input className="search" title="Applications" type="text" placeholder="Search" onChange={onFilterChanged} />
          <p className="text-center search-title"><strong>Applications</strong></p>
          <hr className="separator" />

          {/* applications list */}
          {/* I don't understand css enough to get this to work without inline style, moving on. */}
          <div dir="rtl" style={{ maxHeight: '76vh', overflow: 'scroll', marginRight: -40, paddingRight: 40 }}>
            {filteredApps
              && <ReactList 
                  itemRenderer={(index, key) => {
                    return (
                      <AppCard 
                        key={key}
                        app={filteredApps[index]}
                        index={index}
                        match={match}
                        selectedAppName={selectedAppName}
                        onViewApp={onViewApp}/>
                    )
                  }} 
                  length={filteredApps.length} 
                  type="uniform"/>
              || <div>No Apps Returned.</div>}
          </div>
        </div>
        <div className="not-stuck">
          <Route exact path='/:appName' render={p => <Overview {...p} {...props} />} />
          <Route exact path='/:appName/pods/:podName' render={p => <Pod {...props} {...p} />} />
          <Route path='/' render={() => { return null }} />
        </div>
      </div>
    </div >
  );
};

export default HomePage;
