import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import ReplicaSetViewPage from './replicaset-view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    replicaSetOverviews: [{
      friendlyName: 'fname',
      name: 'name',
      namespace: 'namespace',
      labelSelector: {
        key: 'value'
      },
      availableReplicas: 1,
      fullyLabeledReplicas: 1,
      readyReplicas: 1,
      replicas: 1,
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

  const wrapper = shallow(<ReplicaSetViewPage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('replica set overview should', () => {
  test('render with a replicaOverviews', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Card').length).toBe(1);
  })

  test('have job info displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Replicas').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Ready Replicas').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Available').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Fully Labeled').length).toBe(1);
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