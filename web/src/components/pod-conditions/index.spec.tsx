import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import PodConditions from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    keyPrefix: 'pod',
    items: [{
      type: 'PENDING',
      status: 'pending',
      lastTransitionTime: '2020-04-18T13:52:39Z'
    }]
  }

  const wrapper = shallow(<PodConditions {...props} />)

  return {
    props,
    wrapper
  }
}

describe('PodConditions should', () => {
  test('render', () => {
    const { wrapper, props } = setup();

    expect(wrapper.find('Col').length).toBe(1);

    expect(wrapper.find('div.text-center').findWhere(n => n.text() === props.items[0].type)).toBeDefined();
    expect(wrapper.find('ReactTooltip')).toBeDefined();
  })
})