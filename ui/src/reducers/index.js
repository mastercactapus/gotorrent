import {combineReducers} from "redux"

import torrents from "./torrents"

const torrentApp = combineReducers({
  torrents,
})

export default torrentApp
