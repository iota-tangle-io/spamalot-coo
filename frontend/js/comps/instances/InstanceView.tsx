import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {Link} from 'react-router-dom';
import {withRouter} from "react-router";
import {
    Area, AreaChart, Brush, CartesianGrid, ComposedChart, Legend, Line, LineChart, ReferenceLine,
    ResponsiveContainer, Tooltip, XAxis, YAxis
} from 'recharts';
import {InstanceErrorRate} from "./InstanceErrorRate";
import {InstanceTPS} from "./InstanceTPS";
import {InstanceTXLog} from "./InstanceTXLog";
import {WithStyles} from "material-ui";
import {StyleRulesCallback, Theme, withStyles} from 'material-ui/styles';
import Grid from 'material-ui/Grid';
import Button from "material-ui/Button";
import Paper from 'material-ui/Paper';
import Divider from "material-ui/Divider";

interface Props {
    instanceStore: InstanceStore;
    match: { params: { id: string } };
}

const styles: StyleRulesCallback = (theme: Theme) => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    root: {
        flexGrow: 1,
        marginTop: theme.spacing.unit * 2,
    },
    divider: {
        marginTop: theme.spacing.unit * 3,
        marginBottom: theme.spacing.unit * 3,
    },
    paper: {
        padding: 16,
    },
    button: {
        marginRight: theme.spacing.unit,
    },
    configureButton: {
        float: 'right',
        marginRight: theme.spacing.unit,
    }
});

@withRouter
@inject("instanceStore")
@observer
class instanceView extends React.Component<Props & WithStyles, {}> {
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
        let classes = this.props.classes;
        let config = instance.spammer_config;
        let fakeTPSDate = this.props.instanceStore.instanceTPSData;
        let fakeErrorRateData = this.props.instanceStore.instanceErrorRateData;
        let running = instance.last_state ? instance.last_state.running : false;
        return (
            <div>
                <h2>
                    <Link to={'/instances'}>Instances</Link>
                    <i className="fas fa-angle-right icon_margin_both"></i>
                    Instance: "{instance.name}"
                </h2>

                <Grid container className={classes.root}>
                    <Grid item xs={12} lg={12}>
                        {
                            instance.online &&
                            <span>

                                    <Button
                                        className={classes.button} onClick={this.start}
                                        disabled={running} raised
                                    >
                                        <i className="fas fa-play icon_margin_right"></i>
                                        Start
                                    </Button>

                                    <Button
                                        className={classes.button} onClick={this.stop}
                                        disabled={!running} raised
                                    >
                                        <i className="fas fa-stop-circle icon_margin_right"></i>
                                        Stop
                                    </Button>

                                    <Button className={classes.button} onClick={this.restart} raised>
                                        <i className="fas fa-angle-double-right icon_margin_right"></i>
                                        Restart
                                    </Button>

                                </span>
                        }
                        <Link to={`/instance/editor/update/${instance.id}`}>
                            <Button className={classes.configureButton} raised>
                                <i className="fas fa-wrench icon_margin_right"></i>
                                Configure
                            </Button>
                        </Link>
                    </Grid>
                </Grid>

                <Grid container spacing={40} className={classes.root}>
                    <Grid item xs={12} lg={4}>
                        <Paper className={classes.paper}>
                            <h3>Configuration</h3>
                            <Divider className={classes.divider}/>
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
                                <li>Filter Milestones: {config.filter_milestone ? 'yes' : 'no'}</li>
                            </ul>
                        </Paper>
                    </Grid>
                    <Grid item xs={12} lg={8}>
                        <Paper className={classes.paper}>
                            <h3>Data</h3>
                            <Divider className={classes.divider}/>
                            <InstanceTPS/>
                            <InstanceErrorRate/>
                        </Paper>
                    </Grid>
                </Grid>

                <Grid item xs={12} lg={12} className={classes.root}>
                    <Paper className={classes.paper}>
                        <InstanceTXLog/>
                    </Paper>
                </Grid>
            </div>
        );
    }
}

export var InstanceView = withStyles(styles)<Props>(instanceView);