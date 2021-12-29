import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import PodPage from './view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    podDetail: {
      name: 'name',
      linkedName: 'name',
      namespace: 'namespace',
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
            type: 'PENDING',
            status: 'pending',
            lastTransitionTime: '2020-04-18T13:52:39Z'
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
      }
    },
    showEnvModal: true,
    showSpecModal: true,
    toggleModalType: jest.fn(),
    logs: {
      pod: 'podname',
      output: 'some log text'
    },
    openLogStream: jest.fn(),
    closeLogStream: jest.fn(),
    logStream: null,
    streamEnabled: false,
    envBody: {
      key: 'value'
    },
    hasLogAccess: true,
    toggleContainerNameSelect: jest.fn(),
    containerNameSelectOpen: false,
    setSelectedContainerName: jest.fn(),
    selectedContainerName: ''
  }

  const wrapper = shallow(<PodPage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('pod view should', () => {
  test('render a podDetail', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div.pod-container').length).toBe(1);
    expect(wrapper.find('Card').length).toBe(1);
  })

  test('have pod info displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Namespace:').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Start Time:').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'HostIP:').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'PodIP:').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Ready:').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Restarts:').length).toBe(1);
  })
  
  test('have "Seclect Container" dropdown displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('span').findWhere(n => n.text() === 'Select Container:').length).toBe(1);
  })

  test('have Image displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Image').length).toBe(1);
  })

  test('have Pod Detail button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Pod Detail').length).toBe(1);
  })

  test('have Stream Logs button displayed when "streamEnabled" is false', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Stream Logs').length).toBe(1);
  })

  test('have CardFooter (pod conditions) displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('CardFooter').length).toBe(1);
  })

  test('have JsonViewModals defined', () => {
    const { wrapper } = setup();

    expect(wrapper.find('JsonViewModal').length).toBe(2);
  })
})