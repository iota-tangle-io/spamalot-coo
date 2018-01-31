import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {InstanceOverview} from "./InstanceOverview";
import {Link} from 'react-router-dom';
import Grid from "material-ui/Grid";
import {StyleRulesCallback, Theme, withStyles} from "material-ui/styles";
import {WithStyles} from "material-ui";
import Button from "material-ui/Button";

interface Props {
    instanceStore: InstanceStore;
}

const styles: StyleRulesCallback = (theme: Theme) => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    root: {
        flexGrow: 1,
    },
    paper: {
        padding: 16,
    },
    button: {
        marginTop: theme.spacing.unit,
        marginBottom: theme.spacing.unit * 3,
    },
});

@inject("instanceStore")
@observer
class instancesDashboard extends React.Component<Props & WithStyles, {}> {
    componentWillMount() {
        this.props.instanceStore.fetchInstances();
    }

    render() {
        let classes = this.props.classes;
        let {instances} = this.props.instanceStore;
        let overviews = [];
        instances.keys().forEach(id => overviews.push(<InstanceOverview id={id} key={id}/>));
        return (
            <div className={'instances_dashboard'}>
                <h2>Instances ({instances.size})</h2>
                <Link to={'/instance/editor/create'}>
                    <Button raised color="primary" className={classes.button}>
                        New Instance
                    </Button>
                </Link>
                <Grid container className={classes.root}>
                    {overviews}
                </Grid>
            </div>
        );
    }
}

export var InstancesDashboard = withStyles(styles)(instancesDashboard);