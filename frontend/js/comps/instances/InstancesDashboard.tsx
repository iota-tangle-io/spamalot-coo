import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {InstanceOverview} from "./InstanceOverview";

interface Props {
    instanceStore: InstanceStore;
}

@inject("instanceStore")
@observer
export class InstancesDashboard extends React.Component<Props, {}> {
    componentWillMount() {
        this.props.instanceStore.fetchInstances();
    }

    render() {
        let {instances} = this.props.instanceStore;
        let overviews = [];
        instances.keys().forEach(id => overviews.push(<InstanceOverview id={id} key={id} />));
        return (
            <div className={'instances_dashboard'}>
                <h2>Instances ({instances.size})</h2>
                <div>{overviews}</div>
            </div>
        );
    }
}