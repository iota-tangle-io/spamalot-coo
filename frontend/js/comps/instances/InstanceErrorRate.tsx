import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {withRouter} from "react-router";
import {
    AreaChart, Area, ReferenceLine,
    LineChart, ComposedChart, Brush, XAxis, Line, YAxis,
    CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts';
import {MetricsStore} from "../../stores/MetricsStore";


interface Props {
    instanceStore?: InstanceStore;
    metricsStore?: MetricsStore;
    match?: { params: { id: string } };
}

@withRouter
@inject("metricsStore")
@inject("instanceStore")
@observer
export class InstanceErrorRate extends React.Component<Props, {}> {
    render() {
        let id = this.props.match.params.id;
        let instance = this.props.instanceStore.instances.get(id);
        let config = instance.spammer_config;
        let errorRate = this.props.metricsStore.errorRate;
        return (

            <div>
                <ResponsiveContainer width="100%" height={200}>
                    <ComposedChart data={errorRate} syncId="error_rate">
                        <XAxis dataKey="name"/>
                        <YAxis/>
                        <CartesianGrid strokeDasharray="2 2"/>
                        <Tooltip/>
                        <Legend/>
                        <Line type="linear" isAnimationActive={false} name="Error Rate" dataKey="value"
                              fill="#db172a" stroke="#db172a" dot={false} activeDot={{r: 8}}/>
                    </ComposedChart>
                </ResponsiveContainer>
            </div>
        );
    }
}