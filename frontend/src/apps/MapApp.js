import React, { Component } from "react";
import { Layout } from "antd";
import "antd/dist/antd.css";
import NavBar from "../components/NavBar";

const { Content } = Layout;

class MapApp extends Component {
  render() {
    return (
      <div>
        <Layout>
          <NavBar />
          <Content />
        </Layout>
      </div>
    );
  }
}

export default MapApp;
