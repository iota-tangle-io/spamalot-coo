import {action, computed, observable, ObservableMap, runInAction} from "mobx";
import {AllInstances, Instance, RestartInstanceReq, StartInstanceReq, StopInstanceReq} from "../models/Instance";

const refetchAfter = 500;

export class InstanceStore {
    @observable instances: ObservableMap<Instance> = observable.map();

    @action
    async fetchInstances() {
        try {
            let instances = new AllInstances();
            await instances.fetch()
            runInAction('fetchInstances', () => {
                instances.models.sort((a, b) => a.name > b.name ? 1 : -1);
                instances.models.forEach(inst => this.instances.set(inst.id.toString(), inst));
            });
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async fetchInstance(id: string) {
        try {
            let instance = new Instance({id});
            await instance.fetch();
            runInAction('get instance', () => {
                this.instances.set(id, instance);
            });
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async stopInstance(id: string) {
        try {
            let stopReq = new StopInstanceReq(id);
            await stopReq.fetch();
            setTimeout(() => {this.fetchInstance(id)}, refetchAfter);
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async restartInstance(id: string) {
        try {
            let restartReq = new RestartInstanceReq(id);
            await restartReq.fetch();
            setTimeout(() => {this.fetchInstance(id)}, refetchAfter);
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async startInstance(id: string) {
        try {
            let startReq = new StartInstanceReq(id);
            await startReq.fetch();
            setTimeout(() => {this.fetchInstance(id)}, refetchAfter);
        } catch (err) {
            console.error(err)
        }
    }

    @computed
    get instancesArray(): Array<Instance> {
        let array = [];
        this.instances.forEach(inst => array.push(array));
        return array;
    }


}

// maybe should've picked another name
export let InstanceStoreInstance = new InstanceStore();