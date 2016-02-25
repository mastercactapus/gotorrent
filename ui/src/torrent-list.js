"use strict";
/* @flow */
import * as React from "react";

import {
  Table,
  TableHeader,
  TableRow,
  TableHeaderColumn,
  TableBody,
  TableRowColumn,
  LinearProgress
} from "material-ui"

import List from "material-ui/lib/lists/list";
import ListItem from "material-ui/lib/lists/list-item";
import {connect} from "react-redux"
import {selectTorrent, fetchActive} from "./actions"
import {formatBytes, formatSeconds} from "./util"


export class TorrentList extends React.Component {
  props: {
    selectedTorrent: string;
    active: Array<ActiveTorrent>;
    onSelectTorrent: Function;

    refreshData: Function;
  };

  updateIntervalRef: number;

  constructor() {
    super();
  }

  componentDidMount() {
    this.updateIntervalRef = setInterval(()=>this.props.refreshData(), 2000);
  }
  componentWillUnmount() {
    clearInterval(this.updateIntervalRef);
  }
  renderTorrent(t: ActiveTorrent) {
    var color = null;
    if (t.Seeding) color = "green";

    var progress;
    if (t.Seeding) {
      progress = "Seeding";
    } else if (t.Length == t.BytesCompleted) {
      progress = "Complete"
    } else {
      progress = <LinearProgress color={color} mode="determinate" max={t.Length} value={t.BytesCompleted} />;
    }

    return <TableRow key={t.ID}>
      <TableRowColumn>{t.Name}</TableRowColumn>
      <TableRowColumn>{formatBytes(t.Length)}</TableRowColumn>
      <TableRowColumn>
        {progress}
      </TableRowColumn>
      <TableRowColumn>{formatBytes(t.Downspeed)}/sec</TableRowColumn>
      <TableRowColumn>{formatSeconds(t.ETA)}</TableRowColumn>
    </TableRow>
  }
  render() {

    var activeTorrentRows = this.props.active.map(t=>this.renderTorrent(t))


    return <Table>
      <TableHeader>
        <TableRow>
          <TableHeaderColumn>Name</TableHeaderColumn>
          <TableHeaderColumn>Size</TableHeaderColumn>
          <TableHeaderColumn>Progress</TableHeaderColumn>
          <TableHeaderColumn>Down</TableHeaderColumn>
          <TableHeaderColumn>ETA</TableHeaderColumn>
        </TableRow>
      </TableHeader>
      <TableBody>
        {activeTorrentRows}
      </TableBody>
    </Table>
  }
}

const mapStateToProps = state => {
  return {
    selectedTorrent: state.torrents.selectedTorrent,
    active: state.torrents.active
  }
}
const mapDispatchToProps = dispatch => {
  return {
    onSelectTorrent: id=>{
      dispatch(selectTorrent(id))
    },
    refreshData: ()=>{
      dispatch(fetchActive())
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(TorrentList);
