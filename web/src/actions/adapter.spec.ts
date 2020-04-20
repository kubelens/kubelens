import moxios from 'moxios';
import sinon from 'sinon';
import adapter from './adapter';
import { equal } from 'assert';

beforeEach(function () {
  // import and pass your custom axios instance to this method
  moxios.install()
})

afterEach(function () { 
  // import and pass your custom axios instance to this method
  moxios.uninstall()
})

describe('apps should', () => { 
  test('getAppOverview', async (done) => {
    moxios.withMock(function () {
      let onFulfilled = sinon.spy()
      adapter.get("apps", "minikube", "").then(onFulfilled);

      moxios.wait(function () { 
        let request = moxios.requests.mostRecent()
        request.respondWith({
          status: 200,
          response: [{
            name: 'app1'
          }]
        }).then(function () {
          equal(onFulfilled.called, true)
          equal(onFulfilled.getCall(0).args[0].data[0].name, "app1")
          done()
        });
      });
    });
  });
});
