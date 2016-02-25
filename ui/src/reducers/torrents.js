import {
  SELECT_TORRENT,
  REQUEST_ACTIVE,
  RECEIVE_ACTIVE
} from "../actions"

const initialState = {
  active: [],
  isFetching: false,
  lastUpdate: Date.now()
}

export default function torrents(state=initialState, action) {
  switch (action.type) {
    case RECEIVE_ACTIVE:
      var old = {};
      var lastUpdate = Date.now();
      var secondsElapsed = (lastUpdate - state.lastUpdate) / 1000;
      state.active.forEach(t=>old[t.ID]=t);
      var active = action.data.map((d: ActiveTorrent)=>{
        if (d.Length == d.BytesCompleted) {
          d.Downspeed = 0;
          d.ETA = 0;
          return d;
        }
        if (!old[d.ID]) {
          d.Downspeed = 0;
          d.ETA = Infinity;
          return d;
        }

        d.Downspeed = (d.BytesCompleted - old[d.ID].BytesCompleted) / secondsElapsed;
        d.ETA = (d.Length - d.BytesCompleted) / d.Downspeed;
        return d;
      })
      return Object.assign({}, state, {isFetching: false, lastUpdate, active})
    case REQUEST_ACTIVE:
      return Object.assign({}, state, {isFetching: true})
    case SELECT_TORRENT:
      return Object.assign({}, state, {selectedTorrent: action.id})
    default:
      return state;
  }
}
