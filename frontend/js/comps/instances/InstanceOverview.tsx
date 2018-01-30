import * as React from 'react';
import {InstanceStore} from "../../stores/InstanceStore";
import dateformat from 'dateformat';
import {Link} from 'react-router-dom';
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
        let running = instance.last_state ? instance.last_state.running : false;
        return (
            <div className={'instance_overview'}>
                <h3>{instance.name} <OnlineIndicator online={instance.online}/></h3>
                <p>Slave Address: {instance.address}</p>
                <p>API Token: {instance.api_token}</p>
                <p>Description: {instance.desc}</p>
                <p>Config: remote node {config.node_address}, mwm: {config.mwm}, depth: {config.depth},</p>
                {
                    instance.online &&
                    <div>
                        <p>Spammer: {
                            running
                                ? <span className={'spammer_running_indicator'}>running</span>
                                : <span className={'spammer_stopped_indicator'}>not running</span>
                        }
                        </p>
                        {lastState && <p>Config Hash: {instance.last_state.config_hash}</p>}
                    </div>
                }
                <p>Created: {dateformat(instance.created_on, 'dd.mm.yy HH:mm:ss')}</p>
                <p>Last Update: {dateformat(instance.updated_on, 'dd.mm.yy HH:mm:ss')}</p>

                        <Link to={`/instance/${instance.id}`}>
                            <button className='smallButton' onClick={this.restart}>
                                <i className="fas fa-tachometer-alt icon_margin_right"></i>
                                Dashboard
                            </button>
                        </Link>
                {
                    instance.online &&
                    <div>

                        <button disabled={running} className='smallButton startButton' onClick={this.start}>
                            <i className="fas fa-play icon_margin_right"></i>
                            Start
                        </button>
                        <button disabled={!running} className='smallButton stopButton' onClick={this.stop}>
                            <i className="fas fa-stop-circle icon_margin_right"></i>
                            Stop
                        </button>
                        <button className='smallButton restartButton' onClick={this.restart}>
                            <i className="fas fa-angle-double-right icon_margin_right"></i>
                            Restart
                        </button>


                    </div>
                }
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