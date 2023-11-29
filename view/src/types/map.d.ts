import { ValidationRule } from "quasar";

export type EnumMap = {
  [field: string]: string;
};

export type FieldDetail = {
  label: string;
  field?: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  format?: (value: any) => any;
};

export type FieldMap = {
  [field: string]: FieldDetail;
};

export type StaticRuleMap = {
  [field: string]: ValidationRule[];
};

export type DynamicRuleMap = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  [field: string]: (value: any) => ValidationRule[];
};
