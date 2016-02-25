import * as React from "react"
import {render} from "react-dom"
import thunkMiddleware from "redux-thunk"
import {compose, applyMiddleware, createStore} from "redux"
import {Provider} from "react-redux"
import torrentApp from "./reducers"

import TorrentList from "./torrent-list";
let store = createStore(torrentApp, compose(applyMiddleware(thunkMiddleware),  window.devToolsExtension ? window.devToolsExtension() : f => f))


import injectTapEventPlugin from 'react-tap-event-plugin';

// Needed for onTouchTap
// Can go away when react 1.0 release
// Check this repo:
// https://github.com/zilverline/react-tap-event-plugin
injectTapEventPlugin();

render(
    <Provider store={store}>
      <TorrentList />
    </Provider>,
  document.getElementById("app")
)
