import React from "react";
import ReactDOM from "react-dom";
import DataApp from "./apps/DataApp";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import MapApp from "./apps/MapApp";

ReactDOM.render(
  <BrowserRouter>
    <Switch>
      <Route exact path="/data" component={DataApp} />
      <Route exact path="/" component={MapApp} />
    </Switch>
  </BrowserRouter>,
  document.getElementById("root")
);
