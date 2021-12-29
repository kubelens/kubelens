import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import OverviewCard from './'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    match: {
      params: {
        appName: 'appname'
      }
    },
    onViewApp: jest.fn(),
    overview: {linkedName: 'name1', namespace: 'namespace'},
    selectedOverview: 'appname'
  }

  const wrapper = shallow(<OverviewCard {...props} />)

  return {
    props,
    wrapper
  }
}

describe('app card should', () => {
  test('render with filtered apps', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Card').length).toBe(1);
  })
})