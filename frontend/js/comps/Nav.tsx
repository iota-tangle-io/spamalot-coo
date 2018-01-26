import * as React from 'react';


interface Props {

}

export class Nav extends React.Component<Props, {}> {
    render() {
        return (
            <div className={'nav'}>
                <div className={'site_title'}>spamalot-coo</div>
                <ul className={'nav_menu'}>
                    <li>Instances</li>
                    <li>Configuration</li>
                </ul>
            </div>
        );
    }
}