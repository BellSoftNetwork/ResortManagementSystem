import { api } from "boot/axios";
import {
  ListResponse,
  PageRequestParams,
  SingleResponse,
  SortRequestParams,
} from "src/schema/response";
import { User, UserRole } from "src/schema/user";

export type FetchAdminAccountsRequestParams = Partial<
  PageRequestParams & SortRequestParams
>;

export async function fetchAdminAccounts(
  params: FetchAdminAccountsRequestParams,
) {
  const result = await api.get<ListResponse<User>>("/api/v1/admin/accounts", {
    params,
  });
  return result.data;
}

type AdminAccountParams = {
  name: string;
  email: string;
  password: string;
  role: UserRole;
};

export type AdminAccountCreateParams = AdminAccountParams;

export async function createAdminAccount(params: AdminAccountCreateParams) {
  const result = await api.post<SingleResponse<User>>(
    "/api/v1/admin/accounts",
    params,
  );
  return result.data;
}

export type AdminAccountPatchParams = Partial<AdminAccountParams>;

export async function patchAdminAccount(
  id: number,
  params: Partial<AdminAccountPatchParams>,
) {
  const result = await api.patch<SingleResponse<User>>(
    `/api/v1/admin/accounts/${id}`,
    params,
  );
  return result.data;
}

export async function deleteAdminAccount(id: number) {
  const result = await api.delete(`/api/v1/admin/accounts/${id}`);
  return result.data;
}
