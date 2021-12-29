import moxios from 'moxios';
import sinon from 'sinon';
import adapter from './adapter';
import { deepEqual, fail } from 'assert';

beforeEach(function () {
  // import and pass your custom axios instance to this method
  moxios.install()
})

afterEach(function () { 
  // import and pass your custom axios instance to this method
  moxios.uninstall()
})

describe('overviews should', () => { 
  test('getOverview', async () => {
    moxios.withMock(function () {
      let onFulfilled = sinon.spy()
      adapter.get("overviews", "minikube", "").then(onFulfilled);

      moxios.wait(function () { 
        let request = moxios.requests.mostRecent()
        request.respondWith({
          status: 200,
          response: [{
            name: 'app1'
          }]
        }).then(function () {
          deepEqual(onFulfilled.called, true)
          deepEqual(onFulfilled.getCall(0).args[0].data[0].name, "app1")
        }).catch(function (err) {
          fail(err)
        });
      });
    });
  });
});
