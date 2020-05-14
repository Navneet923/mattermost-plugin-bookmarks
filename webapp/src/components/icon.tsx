// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';

type Props = {
    type?: string;
};

export default class BookmarkIcon extends React.PureComponent<Props> {
    public render() {
        let iconStyle = {};
        if (this.props.type === 'menu') {
            iconStyle = {flex: '0 0 auto', width: '20px', height: '20px', fill: '#0052CC', background: 'white', borderRadius: '50px', padding: '2px'};
        }

        return (
            <span className='MenuItem__icon'>
                <svg
                    aria-hidden='true'
                    focusable='false'
                    role='img'
                    viewBox='0 0 24 24'
                    width='14'
                    height='14'
                    style={iconStyle}
                >
                    <g transform='translate(0 -1028.4)'>
                        <path
                            d='m3 1035.4v2 1 3 1 5 1c0 1.1 0.8954 2 2 2h14c1.105 0 2-0.9 2-2v-1-5-4-3h-18z'
                            fill='#16a085'
                        />
                        <path
                            d='m3 1034.4v2 1 3 1 5 1c0 1.1 0.8954 2 2 2h14c1.105 0 2-0.9 2-2v-1-5-4-3h-18z'
                            fill='#ecf0f1'
                        />
                        <path
                            d='m3 1033.4v2 1 3 1 5 1c0 1.1 0.8954 2 2 2h14c1.105 0 2-0.9 2-2v-1-5-4-3h-18z'
                            fill='#bdc3c7'
                        />
                        <path
                            d='m3 1032.4v2 1 3 1 5 1c0 1.1 0.8954 2 2 2h14c1.105 0 2-0.9 2-2v-1-5-4-3h-18z'
                            fill='#ecf0f1'
                        />
                        <path
                            d='m5 1028.4c-1.1046 0-2 0.9-2 2v1 4 2 1 3 1 5 1c0 1.1 0.8954 2 2 2h2v-1h-1.5c-0.8284 0-1.5-0.7-1.5-1.5 0-0.9 0.6716-1.5 1.5-1.5h12.5 1c1.105 0 2-0.9 2-2v-1-5-4-3-1c0-1.1-0.895-2-2-2h-4-10z'
                            fill='#16a085'
                        />
                        <path
                            d='m8 1028.4v18h1 9 1c1.105 0 2-0.9 2-2v-1-5-4-3-1c0-1.1-0.895-2-2-2h-4-6-1z'
                            fill='#1abc9c'
                        />
                        <path
                            d='m7 1048.4v2 2l2.5-2 2.5 2v-2-2h-5z'
                            fill='#e74c3c'
                        />
                        <rect
                            height='1'
                            width='5'
                            y='1047.4'
                            x='7'
                            fill='#c0392b'
                        />
                    </g>
                </svg>
            </span>
        );
    }
}
