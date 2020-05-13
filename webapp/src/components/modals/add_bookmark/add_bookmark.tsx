// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, { PureComponent } from 'react';
import { Modal } from 'react-bootstrap';

import AddBookmarkForm from './add_bookmark_form';

export type Props = {
    visible: boolean;
}

export default class AddBookmarkModal extends PureComponent<Props, State> {
    render() {
        let content;
        if (this.props.visible) {
            content = (
                <AddBookmarkForm
                    {...this.props}
                />
            );
        }

        const style = getStyle();
        return (
            <Modal
                dialogClassName='modal--scroll'
                style={style.modal}
                show={this.props.visible}
                bsSize='large'
                backdrop='static'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Edit Bookmark'}
                    </Modal.Title>
                </Modal.Header>
                {content}
            </Modal>
        );
    }
}
const getStyle = () => ({
    textarea: {
        resize: 'none',
    },
    modal: {
        height: '100%',
    },
});
