import React, { Component } from "react";
import { Layout } from "antd";
import "antd/dist/antd.css";
import "../styles/App.css";
import NavBar from "../components/NavBar";

const { Content } = Layout;

class ResourceApp extends Component {
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

export default ResourceApp;
