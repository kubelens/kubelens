import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import PodViewPage from './pod-view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    podOverviews: [{
      name: 'name',
      linkedName: 'name',
      namespace: 'namespace',
      pod: {
        name: 'podname',
        namespace: 'namespace',
        hostIP: 'hostip',
        podIP: 'podip',
        startTime: '2020-04-18T13:52:39Z',
        status: {
          phase: 'running',
          phaseMessage: 'started'
        },
        images: [{
          name: 'imgname',
          containerName: 'container'
        }],
        conditions: [{
          name: 'testcondition'
        }]
      }
    }]
  }

  const wrapper = shallow(<PodViewPage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('pod overview should', () => {
  test('render with a podOverview', () => {
    const { wrapper } = setup();

    expect(wrapper.find('PodCard').length).toBe(1);
  })

})