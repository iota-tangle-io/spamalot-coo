import * as React from 'react';
import {InstanceStore} from "../../stores/InstanceStore";
import dateformat from 'dateformat';
import {Link} from 'react-router-dom';
import {inject, observer} from "mobx-react";
import {StyleRulesCallback, Theme, withStyles} from 'material-ui/styles';
import Paper from 'material-ui/Paper';
import {WithStyles} from "material-ui";
import Grid from "material-ui/Grid";
import Button from "material-ui/Button";

interface Props {
    instanceStore?: InstanceStore;
    id: string;
}

const styles: StyleRulesCallback = (theme: Theme) => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    paper: {
        padding: 16,
    },
    button: {
        marginRight: theme.spacing.unit,
    }
});

@inject("instanceStore")
@observer
class instanceOverview extends React.Component<Props & WithStyles, {}> {

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
        let classes = this.props.classes;
        let key = this.props.id;
        let instance = this.props.instanceStore.instances.get(key);
        let config = instance.spammer_config;
        let lastState = instance.last_state;
        let running = instance.last_state ? instance.last_state.running : false;
        return (
            <Grid item xs={12} lg={6}>
                <Paper className={classes.paper}>
                    <h3>{instance.name} <OnlineIndicator online={instance.online}/></h3>
                    <p>Slave Address: {instance.address}</p>
                    <p>API Token: {instance.api_token}</p>
                    <p>Description: {instance.desc}</p>
                    <p>
                        Config: remote node {config.node_address}, mwm: {config.mwm}, depth: {config.depth}
                    </p>
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
                    <p>
                        Created: {dateformat(instance.created_on, 'dd.mm.yy HH:mm:ss')},
                        Last Update: {dateformat(instance.updated_on, 'dd.mm.yy HH:mm:ss')}
                    </p>

                            <Link to={`/instance/${instance.id}`}>
                                <Button raised className={classes.button}>
                                    <i className="fas fa-tachometer-alt icon_margin_right"></i>
                                    Dashboard
                                </Button>
                            </Link>
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

                            <Button
                                className={classes.button} onClick={this.restart} raised
                            >
                                <i className="fas fa-angle-double-right icon_margin_right"></i>
                                Restart
                            </Button>

                        </span>
                    }
                </Paper>
            </Grid>
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

export var InstanceOverview = withStyles(styles)(instanceOverview);