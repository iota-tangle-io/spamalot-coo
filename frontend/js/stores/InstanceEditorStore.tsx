import {Instance} from "../models/Instance";
import {action, observable, runInAction} from "mobx";

export class InstanceEditorStore {
    @observable instance: Instance = new Instance();
    @observable created: boolean = false;
    @observable creation_error: any = null;
    @observable updated: boolean = false;
    @observable update_error: any = null;

    @action
    async loadInstance(id: string) {
        try {
            let instance = new Instance({id});
            await instance.fetch();
            runInAction('load instance in editor', () => {
                this.instance = instance;
            });
        } catch (err) {
            console.error(err)
        }
    }

    @action
    reset() {
        this.created = false;
        this.updated = false;
        this.creation_error = null;
        this.update_error = null;
        this.resetInstance();
    }

    @action
    resetInstance() {
        this.instance = new Instance();
    }

    @action
    async createInstance() {
        try {
            let copy = this.instance.clone();
            await copy.create();
            runInAction('create instance', () => {
                this.created = true;
                this.instance = copy;
            });
        } catch (err) {
            console.error(err);
            runInAction('creation error', () => {
                this.creation_error = err;
            });
        }
    }

    @action
    async updateInstance() {
        try {
            let copy = this.instance.clone();
            await copy.update();
            runInAction('update instance', () => {
                this.created = true;
                this.instance = copy;
            });
        } catch (err) {
            console.error(err);
            runInAction('update error', () => {
                this.update_error = err;
            });
        }
    }

    @action
    changeName(name: string) {
        this.instance.name = name;
    }

    @action
    changeDesc(desc: string) {
        this.instance.desc = desc;
    }

    @action
    changeAddress(addr: string) {
        this.instance.address = addr;
    }

    @action
    changeNodeAddr(addr: string) {
        this.instance.spammer_config.node_address = addr;
    }

    @action
    changeSecurityLvl(lvl: number) {
        this.instance.spammer_config.security_lvl = lvl;
    }

    @action
    changeDepth(depth: number) {
        this.instance.spammer_config.depth = depth;
    }

    @action
    changeTag(tag: string) {
        this.instance.spammer_config.tag = tag;
    }

    @action
    changeMessage(msg: string) {
        this.instance.spammer_config.message = msg;
    }

    @action
    changeDestAddr(destAddr: string) {
        this.instance.spammer_config.dest_address = destAddr;
    }

    @action
    changePoWMode(powMode: number) {
        this.instance.spammer_config.pow_mode = powMode;
    }

    @action
    changeFilterTrunk(filterTrunk: boolean) {
        this.instance.spammer_config.filter_trunk = filterTrunk;
    }

    @action
    changeFilterBranch(filterBranch: boolean) {
        this.instance.spammer_config.filter_branch = filterBranch;
    }

    @action
    changeFilterMilestone(filterMilestone: boolean) {
        this.instance.spammer_config.filter_milestone = filterMilestone;
    }

    @action
    changeCheckAddress(checkAddress: boolean) {
        this.instance.check_address = checkAddress;
    }

}

export let InstanceEditorStoreInstance = new InstanceEditorStore();