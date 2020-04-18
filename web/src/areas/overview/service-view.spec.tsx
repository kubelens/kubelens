import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import ServiceViewPage from './service-view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    serviceOverviews: [{
      appName: 'appname',
      name: 'svcname',
      namespace: 'namespace',
      type: {},
      selector: {
        "key":"value"
      },
      clusterIP: "clusterip",
      loadBalancerIP: "loadbalancerip",
      externalIPs: ["externalip"],
      ports: [{}],
      spec: {
        name: 'spec'
      },
      status: {
        name: 'status'
      },
      deployerLink: 'link',
      configMaps: [{
        name: 'configmap'
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
        conditions: []
      }]
    }],
    toggleModalType: jest.fn(),
    specModalOpen: false,
    statusModalOpen: false,
    configMapModalOpen: false,
    deploymentModalOpen: false
  }

  const wrapper = shallow(<ServiceViewPage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('service view should', () => {
  test('render with serviceOverviews', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Card').length).toBe(1);
  })

  test('have service info displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Namespace').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Cluster IP').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Type').length).toBe(1);
  })

  test('have Deployments button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Deployments').length).toBe(1);
  })

  test('have ConfigMaps button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'ConfigMaps').length).toBe(1);
  })

  test('have Spec button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Spec').length).toBe(1);
  })

  test('have Status button displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Button').findWhere(n => n.text() === 'Status').length).toBe(1);
  })

  test('have CardFooter (deployment overviews) displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('CardFooter').length).toBe(1);
  })

  test('have JsonViewModals defined', () => {
    const { wrapper } = setup();

    expect(wrapper.find('JsonViewModal').length).toBe(4);
  })
})