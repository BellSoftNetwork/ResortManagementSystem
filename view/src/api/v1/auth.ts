import { api } from "boot/axios";
import { SingleResponse } from "src/schema/response";
import { User } from "src/schema/user";

export type RegisterParams = {
  userId: string;
  email: string;
  name: string;
  password: string;
};

export async function postRegister(params: RegisterParams) {
  const result = await api.post<SingleResponse<User>>(
    "/api/v1/auth/register",
    params,
  );
  return result.data;
}

export type LoginParams = {
  username: string;
  password: string;
};

export async function postLogin(params: LoginParams) {
  const result = await api.post<SingleResponse<User>>(
    "/api/v1/auth/login",
    params,
  );
  return result.data;
}

export async function postLogout() {
  const result = await api.post("/api/v1/auth/logout");
  return result.data;
}
