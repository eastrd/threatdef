import React, { Component } from "react";
import { Table } from "antd";
import moment from "moment";
import { Spin } from "antd";

const TUNNEL_API = "http://localhost:8001/tunnel";

class TunnelTable extends Component {
  constructor() {
    super();
    // Set initial state
    this.state = {
      loading: true,
      tunnels: []
    };
  }

  fetchData() {
    fetch(TUNNEL_API)
      .then(resp => resp.json())
      .then(data =>
        data.map(record => {
          var t = moment(parseInt(record.epoch, 10));

          var time = t.format("MMM Do H:mm:ss");
          record.epoch = time;
          return record;
        })
      )
      .then(data => this.setState({ tunnels: data }));
  }

  componentDidMount() {
    // Display Spinning Icon for 1.5 Seconds
    setTimeout(() => this.setState({ loading: false }), 1500);

    // this.fetchData();
    this.interval = setInterval(() => this.fetchData(), 3000);
  }

  render() {
    if (this.state.loading) {
      console.log("Loading");
      return (
        <div>
          <Spin size="large" />
        </div>
      );
    }

    const { tunnels } = this.state;
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
      },
      {
        title: "Destination IP",
        dataIndex: "dst_ip",
        key: "dst_ip"
      }
    ];
    return (
      <Table
        pagination={{ pageSize: 4 }}
        columnWidth="1"
        rowKey="http_id"
        dataSource={tunnels}
        columns={columns}
        expandedRowRender={record => <p style={{ margin: 0 }}>{record.data}</p>}
      />
    );
  }
}

export default TunnelTable;
