import {Config} from "../models/Config";
import {action, observable, runInAction} from "mobx";

export class ConfigStore {
    @observable config: Config = new Config();
    @observable saved: boolean = false;

    @action
    async fetchConfig() {
        try {
            let config = new Config();
            await config.fetch()
            runInAction('setConfig', () => {
                this.config = config;
                this.saved = true;
            });
        } catch (err) {
            console.error(err)
        }
    }

    @action
    async saveConfig() {
        try {
            let config = Object.assign(new Config(), this.config);
            await config.update();
            runInAction('saveConfig', () => {
                this.saved = true;
            });
        } catch (err) {
            console.error(err)
        }
    }

}

export let ConfigStoreInstance = new ConfigStore();