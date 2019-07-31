# kubelens-web

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Build & Deploy

__`DOCKER_ID=id DOCKER_USER=user GIT_BRANCH=master npm run docker-build-push`__

Build and push the docker image

__`npm run build-server`__

Within the ./public folder contains both the uncompiled & compiled code to serve this application. server-x86-x64 was compiled for linux to be used in the docker image. If you wish to customize/alter & build, you'll need to install [Golang](https://golang.org/doc/install). Other than that case, you shouldn't need go installed to work on this piece of Kubelens.

__`INGRESS_HOST=kubelens.minikube-local npm run helm-upgrade`__

Deploy via [Helm](https://helm.sh/)

## Default Available Scripts

In the project directory, you can run:

__`npm start`__

Runs the app in the development mode.<br>
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br>
You will also see any lint errors in the console.

__`npm test`__

Launches the test runner in the interactive watch mode.<br>
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

__`npm run build`__

Builds the app for production to the `build` folder.<br>
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

__`npm run eject`__

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (Webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).
