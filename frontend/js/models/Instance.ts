import {Collection, Model} from "./BackOn";
import {observable} from "mobx";

export class Instance extends Model {
    @observable address: string;
    @observable api_token: string;
    @observable name: string;
    @observable desc: string;
    @observable tags: Array<string>;
    @observable online: boolean;
    @observable check_address: boolean;
    @observable created_on: Date;
    @observable updated_on: Date;
    @observable spammer_config: SpammerConfig;
    @observable last_state: SpammerLastState;

    constructor(attrs?: any) {
        super(attrs);
        this.address = this.address || "127.0.0.1";
        this.name = this.name || "Node";
        this.desc = this.desc || "";
        this.tags = this.tags || [];
        this.online = this.online || false;
        this.check_address = this.check_address || false;
        this.spammer_config = this.spammer_config || new SpammerConfig();
        this.url = '/api/instances/id';
    }

    clone(): Instance {
        return Object.assign(new Instance(), this);
    }
}

export let POW_MODES = {
    LOCAL: 0, REMOTE: 1
}

const NirvanaAddress = "999999999999999999999999999999999999999999999999999999999999999999999999999999999"
const DefaultMessage = "GOSPAMMER9SPAMALOT"
const DefaultTag = "999SPAMALOT"

class SpammerConfig extends Model {
    @observable node_address: string;
    @observable security_lvl: number;
    @observable mwm: number;
    @observable depth: number;
    @observable tag: string;
    @observable message: string;
    @observable dest_address: string;
    @observable pow_mode: number;
    @observable filter_trunk: boolean;
    @observable filter_branch: boolean;
    @observable filter_milestone: boolean;

    constructor(attrs?: any) {
        super(attrs);
        this.node_address = this.node_address || "http://127.0.0.1:14265";
        this.security_lvl = this.security_lvl || 3;
        this.mwm = this.mwm || 14;
        this.depth = this.depth || 1;
        this.tag = this.tag || DefaultTag;
        this.message = this.message || DefaultMessage;
        this.dest_address = this.dest_address || NirvanaAddress;
        this.pow_mode = this.pow_mode || POW_MODES.LOCAL;
        this.filter_trunk = this.filter_trunk || true;
        this.filter_branch = this.filter_branch || true;
        this.filter_milestone = this.filter_milestone || true;
    }
}

class SpammerLastState {
    @observable config_hash: string;
    @observable running: boolean;
}

export class StopInstanceReq extends Instance {
    constructor(instanceID: string) {
        super({});
        this.url = `/api/instances/id/${instanceID}/stop`;
        this.noIDInURI = true;
    }
}

export class StartInstanceReq extends Instance {
    constructor(instanceID: string) {
        super({});
        this.url = `/api/instances/id/${instanceID}/start`;
        this.noIDInURI = true;
    }
}

export class RestartInstanceReq extends Instance {
    constructor(instanceID: string) {
        super({});
        this.url = `/api/instances/id/${instanceID}/restart`;
        this.noIDInURI = true;
    }
}

export class ResetInstanceConfigReq extends SpammerConfig {
    constructor(instanceID: string) {
        super({});
        this.url = `/api/instances/id/${instanceID}/reset_config`;
        this.noIDInURI = true;
    }
}

export class Instances extends Collection<Instance> {
    constructor(models?: Array<Instance>, options?: any) {
        super(Instance, models, options);
    }
}

export class AllInstances extends Instances {
    constructor(models?: Array<Instance>, options?: any) {
        super(models, options);
        this.url = `/api/instances`;
    }
}

export let MetricType = {
    INC_MILESTONE_BRANCH: 0,
    INC_MILESTONE_TRUNK: 1,
    INC_BAD_TRUNK: 2,
    INC_BAD_BRANCH: 3,
    INC_BAD_TRUNK_AND_BRANCH: 4,
    INC_FAILED_TX: 5,
    INC_SUCCESSFUL_TX: 6,
    SUMMARY: 7,
}

export class TXData {
    hash: string;
    count: number;
    created_on: Date;
}

export class MetricSummary {
    txs_succeeded: number;
    txs_failed: number;
    bad_branch: number;
    bad_trunk: number;
    bad_trunk_and_branch: number;
    milestone_trunk: number;
    milestone_branch: number;
    tps: number;
    error_rate: number;
}

export class Metric extends Model {
    @observable metric: number;
    @observable instance_id: string;
    @observable data: any;
    @observable created_on: Date;

    constructor(attrs?: any) {
        super(attrs);
        this.metric = this.metric || -1;
        this.url = '/api/instances/id';
    }
}

export class Metrics extends Collection<Metric> {
    constructor(models?: Array<Metric>, options?: any) {
        super(Metric, models, options);
    }
}