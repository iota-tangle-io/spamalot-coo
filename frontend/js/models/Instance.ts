import {Collection, Model} from "./BackOn";
import {observable} from "mobx";

export class Instance extends Model {
    address: string;
    api_token: string;
    name: string;
    desc: string;
    tags: Array<string>;
    online: boolean;
    created_on: Date;
    updated_on: Date;
    spammer_config: SpammerConfig;
    last_state: SpammerLastState;

    constructor(attrs?: any) {
        super(attrs);
        this.url = '/api/instances/id';
    }
}

class SpammerConfig extends Model {
    node_address: string;
    security_lvl: number;
    mwm: number;
    depth: number;
    tag: string;
    message: string;
    dest_address: string;
    pow_mode: number;
    filter_trunk: boolean;
    filter_branch: boolean;

    constructor(attrs?: any) {
        super(attrs);
    }
}

class SpammerLastState {
    config_hash: string;
    running: boolean;
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
