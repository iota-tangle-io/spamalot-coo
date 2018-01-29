import {action, computed, observable, ObservableMap, runInAction} from "mobx";
import {AllInstances, Instance, RestartInstanceReq, StartInstanceReq, StopInstanceReq} from "../models/Instance";

const refetchAfter = 500;

class FakeTPS {
    ts: number;
    value: number;
}

export class InstanceStore {
    @observable instances: ObservableMap<Instance> = observable.map();
    @observable fake_tps_data: ObservableMap<FakeTPS> = observable.map();
    @observable fake_error_data: ObservableMap<FakeTPS> = observable.map();

    constructor() {
        this.generateRandomTPSData();
    }

    // TODO: remove later
    generateRandomTPSData() {
        let counter = 0;
        let id = setInterval(() => {
            counter++;
            if (counter == 100) {
                clearInterval(id);
                return;
            }
            let r = Math.random() * 10;
            runInAction('add tps data', () => {
                let fakeTPS = new FakeTPS();
                fakeTPS.value = r;
                fakeTPS.ts = Date.now();
                this.fake_tps_data.set(fakeTPS.ts.toString(), fakeTPS);
            });
            runInAction('add error rate data', () => {
                let fakeTPS = new FakeTPS();
                fakeTPS.value = r;
                fakeTPS.ts = Date.now();
                if (Math.floor(r) % 2 == 0) {
                    fakeTPS.value -= r;
                } else {
                    fakeTPS.value -= Math.floor(Math.random() * (4 - 2 + 1)) + 2;
                }
                fakeTPS.value = fakeTPS.value < 0 ? 0 : fakeTPS.value;
                this.fake_error_data.set(fakeTPS.ts.toString(), fakeTPS);
            });
        }, 250);
    }

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

    @computed
    get instanceTPSData(): Array<any> {
        let a = [];
        this.fake_tps_data.forEach(tps => {
            a.push({name: tps.ts, value: tps.value});
        });
        return a;
    }

    @computed
    get instanceErrorRateData(): Array<any> {
        let a = [];
        this.fake_error_data.forEach(tps => {
            a.push({name: tps.ts, value: tps.value});
        });
        return a;
    }

    @action
    async stopInstance(id: string) {
        try {
            let stopReq = new StopInstanceReq(id);
            await stopReq.fetch();
            setTimeout(() => {
                this.fetchInstance(id)
            }, refetchAfter);
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async restartInstance(id: string) {
        try {
            let restartReq = new RestartInstanceReq(id);
            await restartReq.fetch();
            setTimeout(() => {
                this.fetchInstance(id)
            }, refetchAfter);
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async startInstance(id: string) {
        try {
            let startReq = new StartInstanceReq(id);
            await startReq.fetch();
            setTimeout(() => {
                this.fetchInstance(id)
            }, refetchAfter);
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