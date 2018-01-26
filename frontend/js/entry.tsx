import * as moment from "moment";

declare let module;
declare let require;
import 'babel-polyfill';

// stylesheets
require("react-hot-loader/patch");
require("!style-loader!css-loader!sass-loader!normalize.css/normalize.css");
require("!style-loader!css-loader!sass-loader!../css/main.scss");
require("!style-loader!css-loader!sass-loader!react-datepicker/dist/react-datepicker.css");

// set moment lang
moment.locale("en-us");

import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {BrowserRouter} from 'react-router-dom';
import {useStrict} from 'mobx';
import {Provider} from 'mobx-react';
import {AppContainer} from 'react-hot-loader';
import {App} from './comps/App';

// stores
import {AppStoreInstance as appStore} from "./stores/AppStore";
import {ConfigStoreInstance as configStore} from "./stores/ConfigStore";
import {InstanceStoreInstance as instanceStore} from "./stores/InstanceStore";

// use MobX in strict mode
useStrict(true);

let stores = {appStore, configStore, instanceStore};

const render = Component => {
    ReactDOM.render(
        <AppContainer>
            <BrowserRouter>
                <Provider {...stores}>
                    <Component/>
                </Provider>
            </BrowserRouter>
        </AppContainer>,
        document.getElementById('app')
    )
};

render(App);

if (module.hot) {
    module.hot.accept()
}