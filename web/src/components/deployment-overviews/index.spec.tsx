import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import DeploymentOverviews from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    keyPrefix: 'dpl',
    overviews: [{
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
  }

  const wrapper = shallow(<DeploymentOverviews {...props} />)

  return {
    props,
    wrapper
  }
}

describe('DeploymentOverviews should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div').findWhere(n => n.text() === 'Replicas')).toBeDefined();
    expect(wrapper.find('div').findWhere(n => n.text() === 'UpdatedReplicas')).toBeDefined();
    expect(wrapper.find('div').findWhere(n => n.text() === 'ReadyReplicas')).toBeDefined();
    expect(wrapper.find('div').findWhere(n => n.text() === 'UnavailableReplicas')).toBeDefined();
    expect(wrapper.find('div').findWhere(n => n.text() === 'Fully Labeled')).toBeDefined();
  })
})