import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import PodCard from './index'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    name: 'name',
    pod: {
      name: 'podname',
      namespace: 'namespace',
      hostIP: 'hostip',
      podIP: 'podip',
      startTime: '2020-04-18T13:52:39Z',
      phase: 'running',
      phaseMessage: 'started',
      images: [{
        name: 'imgname',
        containerName: 'container'
      }],
      conditions: [{
        name: 'testcondition'
      }]
    }
  }

  const wrapper = shallow(<PodCard {...props} />)

  return {
    props,
    wrapper
  }
}

describe('pod overview should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Card').length).toBe(1);
  })

  test('have pod info displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Namespace').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Start Time').length).toBe(1);
    expect(wrapper.findWhere(n => n.prop('label') === 'Status').length).toBe(1);
  })

  test('have Image displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.findWhere(n => n.prop('label') === 'Image').length).toBe(1);
  })

  test('have CardFooter (pod conditions) displayed', () => {
    const { wrapper } = setup();

    expect(wrapper.find('CardFooter').length).toBe(1);
  })

})