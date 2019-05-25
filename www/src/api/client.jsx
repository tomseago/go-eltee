import { ElTeePromiseClient } from "./api_grpc_web_pb";

const client = new ElTeePromiseClient("http://localhost:9090", null, null);

export default client;
