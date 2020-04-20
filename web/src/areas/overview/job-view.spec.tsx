import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import JobViewPage from './job-view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    jobOverviews: [{
      friendlyName: 'fname',
      name: 'name',
      namespace: 'namespace',
      labelSelector: {
        key: 'value'
      },
      startTime: '2020-04-18T13:52:39Z',
      completionTime: '2020-04-18T13:52:39Z',
      active: 1,
      succeeded: 1,
      failed: 1,
      conditions: [{
        name: 'test'
      }],
      configMaps: [{
        name: 'test'
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

  const wrapper = shallow(<JobViewPage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('job overview should', () => {
  test('render with a jobOverviews', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Card').length).toBe(1);
  })

  test('have job info displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Name').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Start Time').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Completion Time').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Active').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Succeeded').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Failed').length).toBe(1);
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