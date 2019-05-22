import React from "react";
import ReactDOM from "react-dom";

import RpcTest from "./rpc_test";

class HelloMessage extends React.Component {
  render() {
    return <div>Hello {this.props.name}</div>;
  }
}

var mountNode = document.getElementById("app");
ReactDOM.render(
    <div>
        <HelloMessage name="Fred" />
        <RpcTest />
    </div>
    , mountNode);
