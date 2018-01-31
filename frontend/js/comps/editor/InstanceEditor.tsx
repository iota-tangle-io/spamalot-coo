import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import TextField from 'material-ui/TextField';
import Checkbox from 'material-ui/Checkbox';
import {FormControl, FormControlLabel, FormGroup} from 'material-ui/Form';
import Input, {InputLabel} from 'material-ui/Input'
import Select from 'material-ui/Select';
import {InstanceEditorStore} from "../../stores/InstanceEditorStore";
import {StyleRulesCallback, Theme, withStyles} from 'material-ui/styles';
import {WithStyles} from "material-ui";
import {MenuItem} from 'material-ui/Menu';
import Grid from 'material-ui/Grid';
import Paper from 'material-ui/Paper';
import Button from 'material-ui/Button';
import {Redirect, withRouter} from "react-router";
import {Link} from "react-router-dom";
import Tooltip from 'material-ui/Tooltip';

interface Props {
    instanceStore: InstanceStore;
    instanceEditorStore: InstanceEditorStore;
    match: { params: { id: string } };
}

const styles: StyleRulesCallback = (theme: Theme) => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    root: {
        flexGrow: 1,
    },
    button: {
        marginLeft: theme.spacing.unit * 2,
    },
    textField: {
        marginLeft: theme.spacing.unit * 2,
        marginRight: theme.spacing.unit * 2,
        marginTop: theme.spacing.unit * 2,
        width: 200,
    },
    divider: {
        marginTop: theme.spacing.unit * 3,
    },
    paper: {
        padding: 16,
    },
    formGroup: {
        marginLeft: theme.spacing.unit * 2,
    },
    formControl: {
        margin: theme.spacing.unit * 2,
        minWidth: 120,
    },
    info: {
        margin: theme.spacing.unit * 2,
    },
    menu: {
        width: 200,
    },
});

@withRouter
@inject("instanceEditorStore")
@inject("instanceStore")
@observer
class instanceEditor extends React.Component<Props & WithStyles, {}> {

    componentWillMount() {
        let id = this.props.match.params.id;
        if (id) {
            this.props.instanceEditorStore.loadInstance(id);
        }
    }

    componentWillUnmount() {
        this.props.instanceEditorStore.reset();
    }

    componentWillReceiveProps(newProps: Props) {
        let newID = newProps.match.params.id;
        let currentID = this.props.match.params.id;

        // new editor clean
        if (!newID && currentID) {
            this.props.instanceEditorStore.resetInstance();
            return;
        }

        // other instance to load
        if (newID !== currentID) {
            this.props.instanceEditorStore.loadInstance(newID);
        }
    }

    handleNameChange = (e: any) => {
        this.props.instanceEditorStore.changeName(e.target.value);
    }

    handleAddrChange = (e: any) => {
        this.props.instanceEditorStore.changeAddress(e.target.value);
    }

    handleDescChange = (e: any) => {
        this.props.instanceEditorStore.changeDesc(e.target.value);
    }

    handleNodeAddrChange = (e: any) => {
        this.props.instanceEditorStore.changeNodeAddr(e.target.value);
    }

    handleSecurityLvlChange = (e: any) => {
        let lvl = parseInt(e.target.value) || 3;
        this.props.instanceEditorStore.changeSecurityLvl(lvl);
    }

    handleDepthChange = (e: any) => {
        let depth = parseInt(e.target.value) || 1;
        this.props.instanceEditorStore.changeDepth(depth);
    }

    handleTagChange = (e: any) => {
        this.props.instanceEditorStore.changeTag(e.target.value);
    }

    handleMessageChange = (e: any) => {
        this.props.instanceEditorStore.changeMessage(e.target.value);
    }

    handleDestAddressChange = (e: any) => {
        this.props.instanceEditorStore.changeDestAddr(e.target.value);
    }

    handlePoWModeChange = (e: any) => {
        this.props.instanceEditorStore.changePoWMode(e.target.value);
    }

    handleFilterTrunk = (e: any) => {
        this.props.instanceEditorStore.changeFilterTrunk(e.target.checked);
    }

    handleFilterBranch = (e: any) => {
        this.props.instanceEditorStore.changeFilterBranch(e.target.checked);
    }

    handleFilterMilestone = (e: any) => {
        this.props.instanceEditorStore.changeFilterMilestone(e.target.checked);
    }

    handleCheckAddress = (e: any) => {
        this.props.instanceEditorStore.changeCheckAddress(e.target.checked);
    }

    createInstance = () => {
        this.props.instanceEditorStore.createInstance();
    }

    updateInstance = () => {
        this.props.instanceEditorStore.updateInstance();
    }

    render() {
        let {instance, updated, created} = this.props.instanceEditorStore;
        let spammer = instance.spammer_config;
        const {classes} = this.props;

        if (updated || created) {
            return <Redirect to={`/instance/${instance.id}`}/>;
        }

        return (
            <div>
                <h2>
                    <Link to={'/instances'}>Instances</Link>
                    {
                        instance.name.length > 0 &&
                        <span>
                                <i className="fas fa-angle-right icon_margin_both"></i> {instance.name}
                            {
                                instance.address.length > 0 &&
                                <span>{' '} @ {' '} {instance.address}</span>
                            }
                            </span>
                    }
                </h2>

                <Grid container className={classes.root}>
                    <Grid item xs={12} lg={12}>
                        <Paper className={classes.paper}>

                            <Tooltip id="tooltip-icon" title="The address of the slave" placement="top-start">
                                <TextField label="Name" margin="normal"
                                           className={classes.textField}
                                           value={instance.name} onChange={this.handleNameChange}
                                />
                            </Tooltip>


                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="If 'Check Address' is checked only calls from the given address are allowed"
                            >
                                <TextField label="Instance Address" margin="normal"
                                           className={classes.textField}
                                           value={instance.address} onChange={this.handleAddrChange}
                                />
                            </Tooltip>

                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="the IRI node address the spammer will use"
                            >
                                <TextField label="IRI Node Address" margin="normal"
                                           className={classes.textField}
                                           value={spammer.node_address} onChange={this.handleNodeAddrChange}
                                />
                            </Tooltip>


                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="a description to better identify the slave"
                            >
                                <TextField label="Description" margin="normal"
                                           multiline={true} fullWidth={true}
                                           className={classes.textField}
                                           value={instance.desc} onChange={this.handleDescChange}
                                />
                            </Tooltip>


                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="the security level of the generated addresses by the spammer"
                            >
                                <TextField label="Address Security Level" margin="normal" type="number"
                                           className={classes.textField}
                                           value={spammer.security_lvl} onChange={this.handleSecurityLvlChange}
                                />
                            </Tooltip>

                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="the depth of the Markov chain Monte Carlo algorithm used by the spammer"
                            >
                                <TextField label="MCMC Depth" margin="normal" type="number"
                                           className={classes.textField}
                                           value={spammer.depth} onChange={this.handleDepthChange}
                                />
                            </Tooltip>

                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="the tag the spammer will put on each created TX"
                            >
                                <TextField label="TX Tag" margin="normal"
                                           className={classes.textField}
                                           value={spammer.tag} onChange={this.handleTagChange}
                                />
                            </Tooltip>

                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="the message the spammer will put on each created TX"
                            >
                                <TextField label="TX Message" margin="normal"
                                           className={classes.textField}
                                           value={spammer.message} onChange={this.handleMessageChange}
                                />
                            </Tooltip>

                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="the target address for spammer created TXs"
                            >
                                <TextField label="TX Destination Address" margin="normal"
                                           className={classes.textField}
                                           value={spammer.dest_address} onChange={this.handleDestAddressChange}
                                />
                            </Tooltip>

                            <Tooltip id="tooltip-icon" placement="top-start"
                                     title="proof of work: either done locally by the spammer or by the specified IRI node"
                            >
                                <FormControl className={classes.formControl}>
                                    <InputLabel htmlFor="uncontrolled-native">Proof of Work</InputLabel>
                                    <Select
                                        value={spammer.pow_mode} input={<Input id="uncontrolled-native"/>}
                                        onChange={this.handlePoWModeChange}
                                    >
                                        <MenuItem value={0}>Local</MenuItem>
                                        <MenuItem value={1}>IRI Node</MenuItem>
                                    </Select>
                                    {/*<FormHelperText>Uncontrolled</FormHelperText>*/}
                                </FormControl>
                            </Tooltip>

                            <FormGroup row className={classes.formGroup}>
                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={spammer.filter_trunk} onChange={this.handleFilterTrunk}
                                        />
                                    }
                                    label="Filter Trunk TXs"
                                />

                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={spammer.filter_branch} onChange={this.handleFilterBranch}
                                        />
                                    }
                                    label="Filter Branch TXs"
                                />

                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={spammer.filter_milestone} onChange={this.handleFilterMilestone}
                                        />
                                    }
                                    label="Filter Milestones"
                                />

                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={instance.check_address} onChange={this.handleCheckAddress}
                                        />
                                    }
                                    label="Check Address"
                                />
                            </FormGroup>

                            {
                                !instance.id &&
                                <p className={classes.info}>
                                    <i className="fas fa-info-circle"></i>{' '}
                                    API token will be displayed once the instance is created.
                                </p>
                            }

                            {
                                instance.id ?
                                    <Button raised color="primary" className={classes.button}
                                            onClick={this.updateInstance}
                                    >
                                        Update Instance
                                    </Button>
                                    :
                                    <Button raised color="primary" className={classes.button}
                                            onClick={this.createInstance}
                                    >
                                        Create Instance
                                    </Button>
                            }

                        </Paper>
                    </Grid>
                </Grid>


            </div>
        );
    }
}

export var InstanceEditor = withStyles(styles)<Props>(instanceEditor);
