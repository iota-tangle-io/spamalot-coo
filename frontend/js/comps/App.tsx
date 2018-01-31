import * as React from "react";
import {withRouter} from "react-router";
import {inject, observer} from 'mobx-react';
import {ApplicationStore} from "../stores/AppStore";
import DevTools from 'mobx-react-devtools';
import {Route, Switch} from 'react-router-dom';
import {NotFound} from "./NotFound";
import {ConfigStore} from "../stores/ConfigStore";
import {Configuration} from "./Configuration";
import {InstancesDashboard} from "./instances/InstancesDashboard";
import {Nav} from './Nav';
import {InstanceView} from "./instances/InstanceView";
import {InstanceEditor} from "./editor/InstanceEditor";

declare var __DEVELOPMENT__;

interface Props {
    appStore: ApplicationStore;
    configStore: ConfigStore;
}

@withRouter
@inject("configStore")
@inject("appStore")
@observer
export class App extends React.Component<Props, {}> {
    componentWillMount() {
        this.props.configStore.fetchConfig();
    }

    render() {
        return (
            <div>
                <Nav></Nav>
                <div className={"container"}>
                    <Switch>
                        <Route exact path={"/"} component={InstancesDashboard}/>
                        <Route exact path={"/instances"} component={InstancesDashboard}/>
                        <Route exact path={"/instance/:id"} component={InstanceView}/>
                        <Route path={"/instance/editor/create"} component={InstanceEditor} />
                        <Route path={"/instance/editor/update/:id"} component={InstanceEditor} />
                        <Route path={"/config"} component={Configuration}/>
                        <Route component={NotFound}/>
                    </Switch>
                </div>

                {__DEVELOPMENT__ ? <DevTools/> : <span></span>}
            </div>
        );
    }
}