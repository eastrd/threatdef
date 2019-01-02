import React, { Component } from "react";
import { Table } from "antd";
import { Spin } from "antd";
import "../../styles/tables.css";

const CRED_API = "http://localhost:8001/login";

class CredTable extends Component {
  constructor() {
    super();
    // Set initial state
    this.state = {
      loading: true,
      creds: []
    };
    this.fetchData();
  }

  fetchData() {
    fetch(CRED_API)
      .then(resp => resp.json())
      .then(data => this.setState({ creds: data, loading: false }));
  }

  componentDidMount() {
    let secondsToWait = this.props.secondsToWait || 3;
    console.log("Credential table refresh rate:", secondsToWait);
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

    const { creds } = this.state;
    const columns = [
      { title: "Username", dataIndex: "username", key: "username" },
      { title: "Password", dataIndex: "password", key: "password" },
      { title: "Attempts", dataIndex: "num_attempts", key: "num_attempts" }
    ];
    return (
      <Table
        pagination={{ pageSize: this.props.pagesize || 4 }}
        rowKey={record => record.username + ":" + record.password}
        dataSource={creds}
        columns={columns}
        bordered={true}
        size="small"
      />
    );
  }
}

export default CredTable;
