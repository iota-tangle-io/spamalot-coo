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

    stop = () => {
        this.props.instanceStore.stopInstance(this.props.id);
    }

    start = () => {
        this.props.instanceStore.startInstance(this.props.id);
    }

    restart = () => {
        this.props.instanceStore.restartInstance(this.props.id);
    }

    render() {
        let key = this.props.id;
        let instance = this.props.instanceStore.instances.get(key);
        let config = instance.spammer_config;
        let lastState = instance.last_state;
        return (
            <div className={'instance_overview'}>
                <h3>{instance.name} <OnlineIndicator online={instance.online}/></h3>
                <p>Slave address: {instance.address}</p>
                <p>API token: {instance.api_token}</p>
                <p>Description: {instance.desc}</p>
                <p>Config: remote node {config.node_address}, mwm: {config.mwm}, depth: {config.depth},</p>
                <p>Spammer: {
                    instance.last_state.running
                        ? <span className={'spammer_running_indicator'}>running</span>
                        : <span className={'spammer_stopped_indicator'}>not running</span>
                    }
                </p>
                {lastState && <p>Config hash: {instance.last_state.config_hash}</p>}
                <p>Created: {dateformat(instance.created_on, 'dd.mm.yy HH:mm:ss')}</p>

                <div className='smallButton' onClick={this.start}>Start</div>
                <div className='smallButton' onClick={this.stop}>Stop</div>
                <div className='smallButton' onClick={this.restart}>Restart</div>
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