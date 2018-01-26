import * as React from 'react';
import {InstanceStore} from "../../stores/InstanceStore";
import dateformat from 'dateformat';
import {inject, observer} from "mobx-react";

interface Props {
    instanceStore?: InstanceStore;
    id: string;
}

@inject("instanceStore")
@observer
export class InstanceOverview extends React.Component<Props, {}> {
    render() {
        let key = this.props.id;
        let instance = this.props.instanceStore.instances.get(key);
        return (
            <div className={'instance_overview'}>
                <h3>{instance.name} <OnlineIndicator online={instance.online}/></h3>
                <p>Address: {instance.address}</p>
                <p>API token: {instance.api_token}</p>
                <p>Description: {instance.desc}</p>
                <p>Created: {dateformat(instance.created_on, 'dd.mm.yy HH:mm:ss')}</p>
            </div>
        );
    }
}

class OnlineIndicator extends React.Component<{ online: boolean }, {}> {
    render() {
        if (this.props.online) {
            return <div className={'online_indicator'}></div>;
        }
        return <div className={'offline_indicator'}></div>;
    }
}