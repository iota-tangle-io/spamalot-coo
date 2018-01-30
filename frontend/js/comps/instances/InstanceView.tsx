import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {Link} from 'react-router-dom';
import {Route, withRouter} from "react-router";
import {Col, Row} from 'react-flexbox-grid';
import {
    AreaChart, Area, ReferenceLine,
    LineChart, ComposedChart, Brush, XAxis, Line, YAxis,
    CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts';
import {InstanceErrorRate} from "./InstanceErrorRate";
import {InstanceTPS} from "./InstanceTPS";
import {InstanceTXLog} from "./InstanceTXLog";


interface Props {
    instanceStore: InstanceStore;
    match: { params: { id: string } };
}

@withRouter
@inject("instanceStore")
@observer
export class InstanceView extends React.Component<Props, {}> {
    componentWillMount() {
        this.props.instanceStore.fetchInstance(this.props.match.params.id);
    }

    stop = () => {
        let id = this.props.match.params.id;
        this.props.instanceStore.stopInstance(id);
    }

    start = () => {
        let id = this.props.match.params.id;
        this.props.instanceStore.startInstance(id);
    }

    restart = () => {
        let id = this.props.match.params.id;
        this.props.instanceStore.restartInstance(id);
    }

    render() {
        let id = this.props.match.params.id;
        let instance = this.props.instanceStore.instances.get(id);
        if (!instance) {
            return <span>loading...</span>;
        }
        let config = instance.spammer_config;
        let fakeTPSDate = this.props.instanceStore.instanceTPSData;
        let fakeErrorRateData = this.props.instanceStore.instanceErrorRateData;
        let running = instance.last_state ? instance.last_state.running : false;
        return (
            <div>
                <h2><Link to={'/instances'}>Instances</Link> > Instance: "{instance.name}"</h2>
                <Row>
                    <Col xs={12} lg={4}>
                        <div className={'clearfix instance_controls'}>
                            {
                                instance.online &&
                                <div className={'instance_controls'}>

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
                                    <Link to={`/instance/editor/u/${instance.id}`}>
                                        <button className='smallButton floatRightButton' onClick={this.restart}>
                                            <i className="fas fa-wrench icon_margin_right"></i>
                                            Configure
                                        </button>
                                    </Link>
                                </div>
                            }
                        </div>
                    </Col>
                </Row>
                <Row>
                    <Col xs={12} lg={4}>
                        <div>

                            <h3 className={'bold'}>Configuration:</h3>
                            <p>Slave Address: {instance.address}</p>
                            <p>API Token: <span className={'api_token'}>{instance.api_token}</span></p>
                            <p>Description: {instance.desc}</p>
                            <p>Spammer Configuration:</p>
                            <ul className={'list'}>
                                <li>Remote Node: {config.node_address}</li>
                                <li>Minimal Weight Magnitude: {config.mwm}</li>
                                <li>Tip Selection Depth: {config.depth}</li>
                                <li>Address Security Level: {config.security_lvl}</li>
                                <li>TXs Destination Address (first 10 letters shown):
                                    {' '}
                                    {config.dest_address.substr(0, 10)}...
                                </li>
                                <li>Message: {config.message}</li>
                                <li>Tag: {config.tag}</li>
                                <li>PoW Mode: {config.pow_mode == 0 ? 'local' : 'remote'}</li>
                                <li>Filter Trunk TXs: {config.filter_branch ? 'yes' : 'no'}</li>
                                <li>Filter Branch TXs: {config.filter_branch ? 'yes' : 'no'}</li>
                            </ul>
                        </div>
                    </Col>
                    <Col xs={12} lg={8}>
                        <h3 className={'bold'}>Data:</h3>
                        <div>
                            <InstanceTPS/>
                            <InstanceErrorRate/>
                        </div>
                    </Col>
                </Row>
                <Row>
                    <Col xs={12} lg={12}>
                        <InstanceTXLog/>
                    </Col>
                </Row>
            </div>
        );
    }
}