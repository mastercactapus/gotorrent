export const SELECT_TORRENT = "SELECT_TORRENT";
export const RECEIVE_ACTIVE = "RECEIVE_ACTIVE";
export const REQUEST_ACTIVE = "REQUEST_ACTIVE";

export function selectTorrent(id) {
  return {
    type: SELECT_TORRENT,
    id
  }
}

export function fetchActive() {
  return dispatch=> {
    dispatch(requestActive())
    return fetch("http://localhost:7080/torrents")
    .then(r=>r.json())
    .then(data=>dispatch(receiveActive(data)))
  }
}

function receiveActive(data) {
  return {
    type: RECEIVE_ACTIVE,
    data
  }
}

function requestActive() {
  return {
    type: REQUEST_ACTIVE
  }
}
