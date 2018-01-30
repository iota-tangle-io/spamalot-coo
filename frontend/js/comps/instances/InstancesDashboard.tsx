import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {InstanceOverview} from "./InstanceOverview";
import {Link} from 'react-router-dom';
import {Col, Row} from 'react-flexbox-grid';

interface Props {
    instanceStore: InstanceStore;
}

@inject("instanceStore")
@observer
export class InstancesDashboard extends React.Component<Props, {}> {
    componentWillMount() {
        this.props.instanceStore.fetchInstances();
    }

    render() {
        let {instances} = this.props.instanceStore;
        let overviews = [];
        instances.keys().forEach(id => overviews.push(<InstanceOverview id={id} key={id}/>));
        return (
            <div className={'instances_dashboard'}>
                <h2>Instances ({instances.size})</h2>
                <Row className={'box_margin_bottom'}>
                    <Col xs={12} lg={12}>
                        <Link to={'/instance/editor/new'}>
                            <div className={'default_button'}>
                                <i className="fas fa-plus icon_margin_right"></i> New Instance
                            </div>
                        </Link>
                    </Col>
                </Row>
                <Row>
                    <Col xs={12} lg={12}>
                        <div>{overviews}</div>
                    </Col>
                </Row>
            </div>
        );
    }
}