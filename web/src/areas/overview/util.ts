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
import StatusRunning from '../../assets/status-running.png';
import StatusPending from '../../assets/status-pending.png';
import StatusSucceeded from '../../assets/status-succeeded.png';
import StatusFailed from '../../assets/status-failed.png';
import StatusUnknown from '../../assets/status-unknown.png';

export const PENDING = 'Pending';
export const RUNNING = 'Running';
export const SUCCEEDED = 'Succeeded';
export const FAILED = 'Failed';
export const UNKNOWN = 'Unknown';

export const PodPhaseStyle = (input: string): { img?: any, class?: string } => {
  switch (input) {
    case PENDING:
      return {
        img: StatusPending,
        class: 'pod-container pod-pending'
      };
    case RUNNING:
      return {
        img: StatusRunning,
        class: 'pod-container pod-running'
      };
    case SUCCEEDED:
      return {
        img: StatusSucceeded,
        class: 'pod-container pod-succeeded'
      };
    case FAILED:
      return {
        img: StatusFailed,
        class: 'pod-container pod-failed'
      };
    case UNKNOWN:
      return {
        img: StatusUnknown,
        class: 'pod-container pod-unknown'
      };
    default:
      return {};
  }
};
