import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import PodCard from './index'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    podName: 'podname',
    overview: {
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