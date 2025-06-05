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
  const result = await api.post<SingleResponse<User>>("/api/v1/auth/register", params);
  return result.data;
}

export type LoginParams = {
  username: string;
  password: string;
};

// 로그인 응답 타입 (JWT 토큰 포함)
export type LoginResponse = {
  user: User;
  accessToken: string;
  refreshToken: string;
  accessTokenExpiresIn: number;
};

export async function postLogin(params: LoginParams) {
  const result = await api.post<SingleResponse<LoginResponse>>("/api/v1/auth/login", params);
  return result.data;
}

export async function postLogout() {
  const result = await api.post("/api/v1/auth/logout");
  return result.data;
}

export type RefreshTokenParams = {
  refreshToken: string;
};

export type RefreshTokenResponse = {
  accessToken: string;
  refreshToken: string;
  accessTokenExpiresIn: number;
};

export async function postRefreshToken(params: RefreshTokenParams) {
  const result = await api.post<SingleResponse<RefreshTokenResponse>>("/api/v1/auth/refresh", params);
  return result.data;
}
