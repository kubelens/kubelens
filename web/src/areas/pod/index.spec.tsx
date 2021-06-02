import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { Pod } from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    cluster: 'minikube',
    identityToken: 'token',
    pod: {
      name: 'podname',
      namespace: 'namespace',
      hostIP: 'hostip',
      podIP: 'podip',
      startTime: '2020-04-18T13:52:39Z',
      phase: 'running',
      phaseMessage: 'started',
      containerStatus: [{
        image: 'image',
        name: 'name',
        restartCount: 0,
        ready: true
      }],
      status: {
        containerStatuses: [],
        conditions: [{
          name: 'somecondition'
        }],
        startTime: '2020-04-18T13:52:39Z',
        phase: 'running'
      },
      spec: {
        containers: [{
          env: {
            key: 'value'
          }
        }]
      },
      containerNames: [
        "container1",
        "container2"
      ]
    },
    logs: {
      pod: 'podname',
      output: 'some log text'
    },
    envBody: {
      key: 'value'
    },
    selectedOverview: 'appname',
    hasLogAccess: true,
    authClient: jest.fn(),
    setselectedOverview: jest.fn(),
    getPod: jest.fn(),
    getLogs: jest.fn(),
    toggleLogStream: jest.fn(),
    setSelectedContainerName: jest.fn(),
    selectedContainerName: 'container1',
    error: {
      apiOpen: false
    },
    match: {
      params: {
        appName: 'appname',
        podName: 'podname'
      }
    }
  }

  const wrapper = shallow(<Pod {...props} />)

  return {
    props,
    wrapper
  }
}

describe('overview should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div').length).toBe(1);
    expect(wrapper.find('PodPage').length).toBe(1);
    expect(wrapper.find('APIErrorModal').length).toBe(1);
  })
})