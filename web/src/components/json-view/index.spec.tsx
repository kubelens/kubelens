import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import JsonView from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    item: {
      key: 'value'
    },
    collapsed: false
  }

  const wrapper = shallow(<JsonView {...props} />)

  return {
    props,
    wrapper
  }
}

describe('JsonView should', () => {
  test('render json', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div.json-view-container').length).toBe(1);
    expect(wrapper.find('div.json-view-padding').length).toBe(1);
    expect(wrapper.find('ReactJson')).toBeDefined();
  })
})