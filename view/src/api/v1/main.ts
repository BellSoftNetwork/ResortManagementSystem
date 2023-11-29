import { api } from "boot/axios";
import { SingleResponse } from "src/schema/response";
import { User } from "src/schema/user";
import { ServerConfig, ServerEnv } from "src/schema/server-config";

type MyParams = {
  password: string;
};

export type MyPatchParams = Partial<MyParams>;

export async function postMy() {
  const result = await api.post<SingleResponse<User>>("/api/v1/my");
  return result.data;
}

export async function patchMy(params: Partial<MyPatchParams>) {
  const result = await api.patch<SingleResponse<User>>("/api/v1/my", params);
  return result.data;
}

export async function getServerConfig() {
  const result = await api.get<SingleResponse<ServerConfig>>("/api/v1/config");
  return result.data;
}

export async function getServerEnv() {
  const result = await api.get<SingleResponse<ServerEnv>>("/api/v1/env");
  return result.data;
}
