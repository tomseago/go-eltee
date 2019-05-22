
import React, { Component } from 'react';

const { ElTeeClient } = require('./api_grpc_web_pb');
const { StringMsg } = require('./api_pb.js');

var client = new ElTeeClient('http://localhost:9090', null, null);

class RpcTest extends Component {
  
  callGrpcService = () => {
    const request = new StringMsg();
    request.setVal('Hello there');

    console.log("Sending request...");
    client.ping(request, {}, (err, response) => {
      console.log("Eh???");
      if (response == null) {
        console.log(err)
      }else {
        console.log("Got a response!");
        console.log(response.getVal())
      }
    });
    console.log("request inflight");
  }

  render() {
    return (
      <div>
        <button style={{padding:10}} onClick={this.callGrpcService}>Click for grpc request</button>
      </div>
    );
  }
}

export default RpcTest;