import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import DaemonSetViewPage from './daemonset-view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    daemonSetOverviews: [{
      friendlyName: 'fname',
      name: 'name',
      namespace: 'namespace',
      labelSelector: {
        key: 'value'
      },
      currentNumberScheduled: 1,
      desiredNumberScheduled: 1,
      numberAvailable: 1,
      numberMisscheduled: 1,
      numberReady: 1,
      numberUnavailable: 1,
      updatedNumberScheduled: 1,
      conditions: [{
        name: 'name'
      }],
      configMaps: [{
        name: 'name'
      }],
      deploymentOverviews: [{
        friendlyName: 'fname',
        name: 'name',
        namespace: 'namespace',
        labelSelector: {
          key: 'value'
        },
        replicas: 1,
        updatedReplicas: 1,
        readyReplicas: 1,
        unavailableReplicas: 1,
        conditions: [{
          name: 'test'
        }]
      }]
    }],
    toggleModalType: jest.fn(),
    conditionsModalOpen: false,
    configMapModalOpen: false,
    deploymentModalOpen: false
  }

  const wrapper = shallow(<DaemonSetViewPage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('daemon set overview should', () => {
  test('render with a daemonsetOverviews', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Card').length).toBe(1);
  })

  test('have job info displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Ready').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Available').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Unavailable').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Misscheduled').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'CurrentNumberScheduled').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'UpdatedNumberScheduled').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'DesiredNumberScheduled').length).toBe(1);
  })

  test('have Deployments button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Deployments').length).toBe(1);
  })

  test('have ConfigMaps button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'ConfigMaps').length).toBe(1);
  })

  test('have Conditions button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Conditions').length).toBe(1);
  })

  test('have CardFooter (pod conditions) displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('CardFooter').length).toBe(1);
  })

  test('have JsonViewModals defined', () => {
    const { wrapper } = setup();

    expect(wrapper.find('JsonViewModal').length).toBe(3);
  })
})