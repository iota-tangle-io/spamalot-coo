import {Config} from "../models/Config";
import {action, computed, observable, ObservableMap, runInAction} from "mobx";
import {Instance, Metric, MetricSummary, MetricType, TXData} from "../models/Instance";

export class MetricsStore {
    @observable metrics: ObservableMap<Metric> = observable.map();
    @observable txs: ObservableMap<Metric> = observable.map();
    @observable last_metric: Metric = new Metric();
    ws: WebSocket = null;

    @action
    async pullMetrics(id: string) {

        // pull metrics from server
        this.ws = new WebSocket(`ws://${location.host}/api/instances/id/${id}/metrics`);
        this.ws.onmessage = (ev: MessageEvent) => {
            let metric = <Metric> JSON.parse(ev.data);
            switch (metric.metric) {
                case MetricType.SUMMARY:
                    runInAction('add metric', () => {
                        this.metrics.set(metric.id, metric);
                        this.last_metric = metric;
                    });
                    break;
                case MetricType.INC_SUCCESSFUL_TX:
                    runInAction('add tx', () => {
                        this.txs.set(metric.id, metric);
                    });
            }
        };

        this.ws.onopen = (ev: Event) => {
            console.log('websocket open');
        };

        this.ws.onerror = (ev: Event) => {
            console.log('received error from websocket');
        }
    }

    stopPullingMetrics() {
        if (this.ws) {
            this.ws.close();
        }
    }

    @computed
    get transactions(): Array<any> {
        let a = [];
        this.txs.forEach(metricTxs => {
            let data = <TXData> metricTxs.data;
            data.created_on = metricTxs.created_on;
            a.push(data);
        });
        a.sort((a, b) => a.created_on > b.created_on ? -1 : 1);
        return a;
    }

    @computed
    get tps(): Array<any> {
        let a = [];
        this.metrics.forEach(metric => {
            a.push({
                name: metric.created_on,
                value: metric.data.tps,
            });
        });
        return a;
    }

    @computed
    get errorRate(): Array<any> {
        let a = [];
        this.metrics.forEach(metric => {
            let d = <MetricSummary> metric.data;
            a.push({
                name: metric.created_on,
                value: d.error_rate,
            });
        });
        return a;
    }

}

export let MetricsStoreInstance = new MetricsStore();