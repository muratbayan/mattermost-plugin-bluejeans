import React from 'react';

import {id as pluginId} from './manifest';

import Icon from './components/icon.jsx';
import PostTypeBluejeans from './components/post_type_bluejeans';
import {startMeeting} from './actions';

class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerChannelHeaderButtonAction(
            <Icon/>,
            (channel) => {
                startMeeting(channel.id)(store.dispatch, store.getState);
            },
            'Start Bluejeans Meeting'
        );
        registry.registerPostTypeComponent('custom_bluejeans', PostTypeBluejeans);
    }
}

window.registerPlugin(pluginId, new Plugin());
