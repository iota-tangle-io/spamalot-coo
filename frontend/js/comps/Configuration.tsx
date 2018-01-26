import * as React from 'react';
import {inject, observer} from "mobx-react";
import {ConfigStore} from "../stores/ConfigStore";

interface Props {
    configStore: ConfigStore;
}

@inject("configStore")
@observer
export class Configuration extends React.Component<Props, {}> {

    render() {
        return <div></div>;
    }
}