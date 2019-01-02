import React, { Component } from "react";
import { Table } from "antd";
import moment from "moment";
import { Spin } from "antd";

const CMD_API = "http://localhost:8001/cmd";

class CmdTable extends Component {
  constructor() {
    super();
    // Set initial state
    this.state = {
      loading: true,
      cmds: []
    };
    this.fetchData();
  }

  fetchData() {
    fetch(CMD_API)
      .then(resp => resp.json())
      .then(data =>
        data.map(record => {
          var t = moment(parseInt(record.epoch, 10));

          var time = t.format("MMM Do H:mm:ss");
          record.epoch = time;
          return record;
        })
      )
      .then(data => this.setState({ cmds: data, loading: false }));
  }

  componentDidMount() {
    let secondsToWait = this.props.secondsToWait || 3;
    console.log("Cmd table refresh rate:", secondsToWait);
    this.interval = setInterval(() => this.fetchData(), secondsToWait * 1000);
  }

  render() {
    if (this.state.loading) {
      // If data has not loaded, display the spinning icon.
      console.log("Loading");
      return (
        <div>
          <Spin size="large" />
        </div>
      );
    }

    const { cmds } = this.state;
    const columns = [
      {
        title: "Time",
        dataIndex: "epoch",
        key: "time"
      },
      {
        title: "Source IP",
        dataIndex: "src_ip",
        key: "src_ip"
      }
    ];
    return (
      <Table
        pagination={{ pageSize: this.props.pagesize || 4 }}
        rowKey="input_id"
        dataSource={cmds}
        columns={columns}
        expandedRowRender={record => <p style={{ margin: 0 }}>{record.cmd}</p>}
        bordered={true}
        size="small"
      />
    );
  }
}

export default CmdTable;
