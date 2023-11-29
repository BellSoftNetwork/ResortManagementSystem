import {
  DynamicRuleMap,
  EnumMap,
  FieldDetail,
  FieldMap,
  StaticRuleMap,
} from "src/types/map"
import { BASE_TIME_FIELD_MAP, BaseTime } from "src/schema/base"

const USER_ROLE_MAP: EnumMap = {
  NORMAL: "일반",
  ADMIN: "관리자",
  SUPER_ADMIN: "최고 관리자",
} as const;
export type UserRole = keyof typeof USER_ROLE_MAP;

export function userRoleValueToName(role: UserRole) {
  return USER_ROLE_MAP[role] || role
}

const USER_STATUS_MAP: EnumMap = {
  ACTIVE: "활성화",
  INACTIVE: "비활성화",
} as const;
export type UserStatus = keyof typeof USER_STATUS_MAP;

export type User = {
  id: number;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  profileImageUrl: string;
} & BaseTime;

const UserFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  name: { label: "이름" } as const,
  email: { label: "이메일" } as const,
  role: { label: "권한", format: userRoleValueToName } as const,
  status: { label: "상태" } as const,
  profileImageUrl: { label: "프로필 이미지 URL" } as const,
  ...BASE_TIME_FIELD_MAP,
} as const;

export function getUserFieldDetail(field: string): FieldDetail | null {
  const fieldDetail = UserFieldMap[field]

  return fieldDetail
    ? {
      field: field,
      ...fieldDetail,
    }
    : null;
}

export function formatUserFieldToLabel(field: string) {
  return getUserFieldDetail(field)?.label ?? field
}

export function formatUserValue(field: string, value: string | number | null) {
  const fieldDetail = getUserFieldDetail(field)
  return fieldDetail?.format ? fieldDetail.format(value) : value
}

export const userStaticRules: StaticRuleMap = {
  name: [
    (value: string) =>
      (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다",
  ] as const,
  email: [
    (value: string) =>
      /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value) ||
      "이메일이 유효하지 않습니다.",
  ] as const,
  password: [
    (value: string) =>
      value.length === 0 ||
      (value.length >= 8 && value.length <= 20) ||
      "비밀번호는 8~20 글자가 필요합니다.",
  ] as const,
} as const;

export const userDynamicRules: DynamicRuleMap = {
  passwordConfirm: (password: string) =>
    [
      (value: string) => password === value || "비밀번호가 일치하지 않습니다.",
    ] as const,
} as const;
